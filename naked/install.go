package naked

import (
	"path/filepath"

	"github.com/itchio/hush"
	"github.com/itchio/hush/bfs"
	"github.com/itchio/hush/download"
	"github.com/pkg/errors"
)

func (m *Manager) Install(params hush.InstallParams) (*hush.InstallResult, error) {
	consumer := params.Consumer

	stats, err := params.File.Stat()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	destName := filepath.Base(stats.Name())
	destAbsolutePath := filepath.Join(params.InstallFolderPath, destName)

	err = download.DownloadInstallSource(download.DownloadInstallSourceParams{
		Context:       params.Context,
		Consumer:      params.Consumer,
		StageFolder:   params.StageFolderPath,
		OperationName: "naked-installer",
		File:          params.File,
		DestPath:      destAbsolutePath,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var res = hush.InstallResult{
		Files: []string{
			destName,
		},
	}

	consumer.Opf("Busting ghosts...")
	var bustGhostStats bfs.BustGhostStats
	err = bfs.BustGhosts(bfs.BustGhostsParams{
		Folder:   params.InstallFolderPath,
		NewFiles: res.Files,
		Receipt:  params.ReceiptIn,
		Consumer: params.Consumer,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = params.EventSink.PostGhostBusting("install::naked", bustGhostStats)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
