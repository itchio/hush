package download

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/itchio/headway/state"
	"github.com/itchio/httpkit/eos"
	"github.com/itchio/httpkit/htfs"
	"github.com/itchio/httpkit/retrycontext"
	"github.com/itchio/hush/download/downloadextractor"
	"github.com/itchio/hush/intervalsaveconsumer"
	"github.com/itchio/intact"
	"github.com/itchio/savior"

	"github.com/pkg/errors"
)

type DownloadInstallSourceParams struct {
	// State consumer
	Consumer *state.Consumer

	// For cancellation
	Context context.Context

	// Staging folder we can use to store the download state file
	StageFolder string

	// Used to determine the state file name - must be filesystem-safe
	OperationName string

	// File to download from (usually an HTTP file)
	File eos.File

	// Where to download the file
	DestPath string
}

func DownloadInstallSource(params DownloadInstallSourceParams) error {
	consumer := params.Consumer
	ctx := params.Context
	stageFolder := params.StageFolder
	file := params.File
	destPath := params.DestPath

	stateName := fmt.Sprintf("downsource-%s-state.dat", params.OperationName)
	statePath := filepath.Join(stageFolder, stateName)
	sc := intervalsaveconsumer.New(statePath, intervalsaveconsumer.DefaultInterval, consumer, ctx)

	checkpoint, err := sc.Load()
	if err != nil {
		consumer.Warnf("Could not load checkpoint: %s", err.Error())
	}

	destName := filepath.Base(destPath)
	sink := &savior.FolderSink{
		Directory: filepath.Dir(destPath),
		Consumer:  consumer,
	}

	retryCtx := retrycontext.NewDefault()
	retryCtx.Settings.Consumer = consumer

	tryDownload := func() error {
		ex := downloadextractor.New(file, destName)
		ex.SetConsumer(consumer)
		ex.SetSaveConsumer(sc)
		_, err := ex.Resume(checkpoint, sink)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	}

	for retryCtx.ShouldTry() {
		err := tryDownload()
		if err != nil {
			if errors.Cause(err) == savior.ErrStop {
				return err
			}

			if intact.IsIntegrityError(err) {
				consumer.Warnf("Had integrity errors, we have to start over")
				checkpoint = nil
				retryCtx.Retry(err)
				continue
			}

			if se, ok := asServerError(err); ok {
				if se.Code == htfs.ServerErrorCodeNoRangeSupport {
					consumer.Warnf("%s does not support range requests (boo, hiss), we have to start over", se.Host)
					checkpoint = nil
					retryCtx.Retry(err)
					continue
				}
			}

			// if it's not an integrity error, just bubble it up
			return err
		}

		// N.B: we don't remove the state file in case we retry
		// so that we know the operation is done already.
		// That means state file names should be unique per operation,
		// cf. https://github.com/itchio/itch/issues/2311

		return nil
	}

	return errors.WithMessage(retryCtx.LastError, "download")
}

type causer interface {
	Cause() error
}

func asServerError(err error) (*htfs.ServerError, bool) {
	if err == nil {
		return nil, false
	}

	if se, ok := err.(causer); ok {
		return asServerError(se.Cause())
	}

	if se, ok := err.(*htfs.ServerError); ok {
		return se, true
	}

	return nil, false
}
