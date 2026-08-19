package hush

import (
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/itchio/boar"
	"github.com/itchio/headway/state"
	"github.com/itchio/httpkit/eos"
	"github.com/itchio/savior"
	"github.com/pkg/errors"
)

type GetInstallerInfoParams struct {
	Consumer *state.Consumer
	File     eos.File

	// see boar.ProbeParams
	NormalizeZipBackslashes bool
}

func GetInstallerInfo(consumer *state.Consumer, file eos.File) (*InstallerInfo, error) {
	return GetInstallerInfoWithParams(GetInstallerInfoParams{
		Consumer: consumer,
		File:     file,
	})
}

func GetInstallerInfoWithParams(params GetInstallerInfoParams) (*InstallerInfo, error) {
	consumer := params.Consumer
	file := params.File

	stat, err := file.Stat()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	target := stat.Name()
	ext := strings.ToLower(filepath.Ext(target))
	name := filepath.Base(target)

	consumer.Infof("↝ For source (%s)", name)

	var installerType = InstallerTypeUnknown

	if extType, ok := installerForExt[ext]; ok {
		consumer.Infof("✓ Using file extension registry (%s) => (%s)", ext, extType)
		installerType = extType
	} else {
		consumer.Warnf("  No mapping for file extension (%s)", ext)
	}

	if installerType == InstallerTypeArchive {
		beforeArchiveProbe := time.Now()
		consumer.Infof("  Probing with boar (because installer type is archive)...")

		var entries []*savior.Entry
		archiveInfo, err := boar.Probe(boar.ProbeParams{
			File:     file,
			Consumer: consumer,
			OnEntries: func(es []*savior.Entry) {
				entries = es
			},
			NormalizeZipBackslashes: params.NormalizeZipBackslashes,
		})
		consumer.Debugf("  (archive probe took %s)", time.Since(beforeArchiveProbe))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if archiveInfo == nil {
			consumer.Infof("✗ Source is not a supported archive format")
		} else {
			consumer.Infof("✓ Source is a supported archive format (%s)", archiveInfo.Format)
			consumer.Infof("  Features: %s", archiveInfo.Features)
			return &InstallerInfo{
				Type:        InstallerTypeArchive,
				ArchiveInfo: archiveInfo,
				Entries:     entries,
			}, nil
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &InstallerInfo{
		Type: installerType,
	}, nil
}
