// Package launcher emits the desktop taskbar unread-count badge.
//
// On Linux it broadcasts the com.canonical.Unity.LauncherEntry "Update" D-Bus
// signal (the "smart launcher" / Unity LauncherEntry API), which KDE Plasma's
// task manager and compatible widgets read to draw an unread-count badge on the
// application's taskbar icon. On other platforms it is a no-op.
package launcher

// Updater sets the unread-count badge shown on the app's taskbar/launcher icon.
type Updater interface {
	// SetCount updates the badge to count. The badge is hidden when count <= 0.
	SetCount(count int)
	// Close releases any resources held by the updater (e.g. the D-Bus connection).
	Close()
}

// New returns a platform-appropriate Updater for the given desktop file id,
// e.g. "io.github.hkdb.Aerion.desktop".
func New(desktopID string) Updater {
	return newPlatformUpdater(desktopID)
}

// noopUpdater is used on non-Linux platforms and as a fallback when no session
// bus is available.
type noopUpdater struct{}

func (noopUpdater) SetCount(int) {}
func (noopUpdater) Close()       {}
