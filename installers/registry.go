package installers

import (
	"github.com/itchio/hush"
	"github.com/itchio/hush/archive"
	"github.com/itchio/hush/naked"
)

func GetManager(typ hush.InstallerType) hush.Manager {
	switch typ {
	case hush.InstallerTypeArchive:
		return &archive.Manager{}
	case hush.InstallerTypeNaked:
		return &naked.Manager{}
	default:
		return nil
	}
}
