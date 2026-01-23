package hush

var installerForExt = map[string]InstallerType{

	///////////////////////////////////////////////////////////
	// Generic archives
	///////////////////////////////////////////////////////////

	".zip": InstallerTypeArchive,
	".gz":  InstallerTypeArchive,
	".bz2": InstallerTypeArchive,
	".7z":  InstallerTypeArchive,
	".tar": InstallerTypeArchive,
	".xz":  InstallerTypeArchive,
	".rar": InstallerTypeArchive,

	///////////////////////////////////////////////////////////
	// Known naked (because we no longer support them)
	///////////////////////////////////////////////////////////

	".dmg":          InstallerTypeNaked,
	".exe":          InstallerTypeNaked,
	".x86_64":       InstallerTypeNaked,
	".appimage":     InstallerTypeNaked,
	".deb":          InstallerTypeNaked,
	".rpm":          InstallerTypeNaked,
	".pkg":          InstallerTypeNaked,
	".msi":          InstallerTypeNaked,
	".jar":          InstallerTypeNaked,
	".air":          InstallerTypeNaked,
	".love":         InstallerTypeNaked,
	".unitypackage": InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// Books!
	///////////////////////////////////////////////////////////

	".pdf":    InstallerTypeNaked,
	".ps":     InstallerTypeNaked,
	".djvu":   InstallerTypeNaked,
	".cbr":    InstallerTypeNaked,
	".cbz":    InstallerTypeNaked,
	".cb7":    InstallerTypeNaked,
	".cbt":    InstallerTypeNaked,
	".cba":    InstallerTypeNaked,
	".doc":    InstallerTypeNaked,
	".docx":   InstallerTypeNaked,
	".epub":   InstallerTypeNaked,
	".mobi":   InstallerTypeNaked,
	".pdb":    InstallerTypeNaked,
	".fb2":    InstallerTypeNaked,
	".xeb":    InstallerTypeNaked,
	".ceb":    InstallerTypeNaked,
	".ibooks": InstallerTypeNaked,
	".txt":    InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// Media
	///////////////////////////////////////////////////////////

	".ogg": InstallerTypeNaked,
	".mp3": InstallerTypeNaked,
	".wav": InstallerTypeNaked,
	".mp4": InstallerTypeNaked,
	".avi": InstallerTypeNaked,
	".mkv": InstallerTypeNaked,
	".flac": InstallerTypeNaked,
	".opus": InstallerTypeNaked,
	".webm": InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// Images
	///////////////////////////////////////////////////////////

	".png": InstallerTypeNaked,
	".jpg": InstallerTypeNaked,
	".gif": InstallerTypeNaked,
	".bmp": InstallerTypeNaked,
	".tga": InstallerTypeNaked,
	".webp": InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// Game Maker assets
	///////////////////////////////////////////////////////////

	".gmez": InstallerTypeNaked,
	".gmz":  InstallerTypeNaked,
	".yyz":  InstallerTypeNaked,
	".yymp": InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// ROMs
	///////////////////////////////////////////////////////////

	".gb":  InstallerTypeNaked,
	".gbc": InstallerTypeNaked,
	".sfc": InstallerTypeNaked,
	".smc": InstallerTypeNaked,
	".swc": InstallerTypeNaked,
	".gen": InstallerTypeNaked,
	".sg":  InstallerTypeNaked,
	".smd": InstallerTypeNaked,
	".md":  InstallerTypeNaked,

	///////////////////////////////////////////////////////////
	// Miscellaneous other things
	///////////////////////////////////////////////////////////

	// Some html games provide a single .html file
	// Now that's dedication.
	".html": InstallerTypeNaked,
}
