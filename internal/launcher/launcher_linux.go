//go:build linux

package launcher

import (
	"github.com/godbus/dbus/v5"
	"github.com/hkdb/aerion/internal/logging"
	"github.com/rs/zerolog"
)

const (
	launcherEntryInterface = "com.canonical.Unity.LauncherEntry"
	// launcherEntryPath is the object path the Update signal is emitted on.
	// Receivers (KDE Plasma's task manager, compatible widgets) match on the
	// interface + app URI rather than the path, so any valid path works.
	launcherEntryPath = dbus.ObjectPath("/com/canonical/unity/launcherentry/aerion")
)

// linuxUpdater broadcasts the Unity LauncherEntry "Update" signal over the
// session bus so the desktop draws an unread-count badge on the taskbar icon.
type linuxUpdater struct {
	conn   *dbus.Conn
	appURI string
	log    zerolog.Logger
}

func newPlatformUpdater(desktopID string) Updater {
	log := logging.WithComponent("launcher")
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Warn().Err(err).Msg("LauncherEntry: no session bus; taskbar unread badge disabled")
		return &noopUpdater{}
	}
	return &linuxUpdater{
		conn:   conn,
		appURI: "application://" + desktopID,
		log:    log,
	}
}

// SetCount broadcasts the unread count as a Unity LauncherEntry Update signal.
func (u *linuxUpdater) SetCount(count int) {
	if u.conn == nil {
		return
	}
	if count < 0 {
		count = 0
	}
	props := map[string]dbus.Variant{
		"count":         dbus.MakeVariant(int64(count)),
		"count-visible": dbus.MakeVariant(count > 0),
	}
	if err := u.conn.Emit(launcherEntryPath, launcherEntryInterface+".Update", u.appURI, props); err != nil {
		u.log.Debug().Err(err).Msg("LauncherEntry: failed to emit Update signal")
	}
}

func (u *linuxUpdater) Close() {
	if u.conn != nil {
		_ = u.conn.Close()
		u.conn = nil
	}
}
