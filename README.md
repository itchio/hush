# hush

A Go library for installing files from itch.io. Used by the [itch](https://itch.io/app)
desktop app to handle game and application installation.

## Features

- **Multi-format support** - Handles both archives and naked (single-file) downloads
- **Resumable installation** - Interruptible and can be resumed, even after a crash
- **Ghost busting** - Automatically removes obsolete files from previous versions
- **Angel saving** - Preserves user-created files (saves, mods, configs) during updates
- **Receipt tracking** - Maintains installation records for clean updates

## Supported Formats

**Archives:** zip, 7z, tar, rar, gz, bz2, xz

**Naked files:** exe, dmg, deb, rpm, pkg, msi, jar, pdf, epub, mp3, mp4, png, html, and more

## Terminology

- **Ghosts** - Files from a previous installation that are no longer part of the
  current version. These obsolete files are removed during updates.

- **Angels** - User-created or runtime-generated files (saves, mods, configs) that
  should be preserved during updates even though they weren't part of the original
  installation.

- **Ghost Busting** - The process of detecting and removing ghost files after a
  new installation completes.

- **Angel Saving** - The process of preserving angel files by temporarily backing
  them up, performing a fresh install, then merging them back.

- **Receipt** - A compressed JSON file (`.itch/receipt.json.gz`) that tracks what
  was installed. Used to detect ghosts and angels during updates.

## How It Works

1. **Installer detection** - File extension determines the installer type (archive
   or naked). Archives are probed to verify format support.

2. **Two installer types**:
   - `archive` - Extracts contents to the destination folder
   - `naked` - Copies the single file directly to the destination

3. **Installation flow** - Install files, bust ghosts (if receipt exists), then
   write a new receipt for future updates.

## Usage

```go
package main

import (
    "context"
    "os"

    "github.com/itchio/headway/state"
    "github.com/itchio/hush"
    "github.com/itchio/hush/bfs"
    "github.com/itchio/hush/installers"
)

func main() {
    consumer := &state.Consumer{} // progress/logging listener

    // Open the source file
    file, _ := os.Open("/path/to/game.zip")
    defer file.Close()

    // Detect installer type
    info, _ := hush.GetInstallerInfo(consumer, file)

    // Get the appropriate manager (archive or naked)
    manager := installers.GetManager(info.Type)

    // Read existing receipt (if any) for ghost busting
    receipt, _ := bfs.ReadReceipt("/path/to/install")

    // Perform installation
    result, _ := manager.Install(hush.InstallParams{
        File:              file,
        ReceiptIn:         receipt,
        StageFolderPath:   "/path/to/staging",
        InstallFolderPath: "/path/to/install",
        Consumer:          consumer,
        InstallerInfo:     info,
        Context:           context.Background(),
    })

    // Write receipt for future updates
    newReceipt := &bfs.Receipt{
        Files:         result.Files,
        InstallerName: manager.Name(),
    }
    newReceipt.WriteReceipt("/path/to/install")
}
```

