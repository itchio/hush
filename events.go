package hush

import (
	"fmt"
	"time"

	"github.com/itchio/hush/bfs"
)

type InstallEventSink struct {
	Append func(ev InstallEvent) error
}

func (ies *InstallEventSink) PostEvent(event InstallEvent) error {
	if ies == nil {
		return nil
	}

	event.Timestamp = time.Now()
	switch true {
	case event.Install != nil:
		event.Type = InstallEventInstall
	case event.Heal != nil:
		event.Type = InstallEventHeal
	case event.Upgrade != nil:
		event.Type = InstallEventUpgrade
	case event.Problem != nil:
		event.Type = InstallEventProblem
	case event.GhostBusting != nil:
		event.Type = InstallEventGhostBusting
	case event.Patching != nil:
		event.Type = InstallEventPatching
	case event.Fallback != nil:
		event.Type = InstallEventFallback
	}

	if event.Type == "" {
		// wee runtime checks
		panic("InstallEventSink events should always have Type set")
	}

	return ies.Append(event)
}

func (ies *InstallEventSink) PostProblem(err error) error {
	prob := ies.MakeProblem(err)
	return ies.PostEvent(InstallEvent{
		Type:    InstallEventProblem,
		Problem: &prob,
	})
}

func (ies *InstallEventSink) MakeProblem(err error) ProblemInstallEvent {
	return ProblemInstallEvent{
		Error:      fmt.Sprintf("%v", err),
		ErrorStack: fmt.Sprintf("%+v", err),
	}
}

func (ies *InstallEventSink) PostGhostBusting(operation string, stats bfs.BustGhostStats) error {
	return ies.PostEvent(InstallEvent{
		GhostBusting: &GhostBustingInstallEvent{
			Operation: "heal",
			Found:     stats.Found,
			Removed:   stats.Removed,
		},
	})
}
