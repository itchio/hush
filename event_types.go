package hush

import "time"

type InstallEvent struct {
	Type      InstallEventType `json:"type"`
	Timestamp time.Time        `json:"timestamp"`

	Heal         *HealInstallEvent         `json:"heal,omitempty"`
	Install      *InstallInstallEvent      `json:"install,omitempty"`
	Upgrade      *UpgradeInstallEvent      `json:"upgrade,omitempty"`
	GhostBusting *GhostBustingInstallEvent `json:"ghostBusting,omitempty"`
	Patching     *PatchingInstallEvent     `json:"patching,omitempty"`
	Problem      *ProblemInstallEvent      `json:"problem,omitempty"`
	Fallback     *FallbackInstallEvent     `json:"fallback,omitempty"`
}

type InstallEventType string

const (
	// Started for the first time or resumed after a pause
	// or exit or whatever
	InstallEventResume InstallEventType = "resume"

	// Stopped explicitly (pausing downloads), can't rely
	// on this being present because BRÜTAL PÖWER LÖSS will
	// not announce itself 🔥
	InstallEventStop InstallEventType = "stop"

	// Regular install from archive or naked file
	InstallEventInstall InstallEventType = "install"

	// Reverting to previous version or re-installing
	// wharf-powered upload
	InstallEventHeal InstallEventType = "heal"

	// Applying one or more wharf patches
	InstallEventUpgrade InstallEventType = "upgrade"

	// Applying a single wharf patch
	InstallEventPatching InstallEventType = "patching"

	// Cleaning up ghost files
	InstallEventGhostBusting InstallEventType = "ghostBusting"

	// Any kind of step failing
	InstallEventProblem InstallEventType = "problem"

	// Any operation we do as a result of another one failing,
	// but in a case where we're still expecting a favorable
	// outcome eventually.
	InstallEventFallback InstallEventType = "fallback"
)

type InstallInstallEvent struct {
	Manager string `json:"manager"`
}

type HealInstallEvent struct {
	TotalCorrupted int64 `json:"totalCorrupted"`

	AppliedCaseFixes bool `json:"appliedCaseFixes"`
}

type UpgradeInstallEvent struct {
	NumPatches int `json:"numPatches"`
}

type ProblemInstallEvent struct {
	// Short error
	Error string `json:"error"`
	// Longer error
	ErrorStack string `json:"errorStack"`
}

type FallbackInstallEvent struct {
	// Name of the operation we were trying to do
	Attempted string `json:"attempted"`

	// Problem encountered while trying "attempted"
	Problem ProblemInstallEvent `json:"problem"`

	// Name of the operation we're falling back to
	NowTrying string `json:"nowTrying"`
}

type PatchingInstallEvent struct {
	// Build we patched to
	BuildID int64 `json:"buildID"`

	// "default" or "optimized" (for the +bsdiff variant)
	Subtype string `json:"subtype"`
}

type GhostBustingInstallEvent struct {
	// Operation that requested the ghost busting (install, upgrade, heal)
	Operation string `json:"operation"`

	// Number of ghost files found
	Found int64 `json:"found"`

	// Number of ghost files removed
	Removed int64 `json:"removed"`
}
