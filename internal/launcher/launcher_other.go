//go:build !linux

package launcher

// newPlatformUpdater is a no-op on non-Linux platforms. The Unity LauncherEntry
// API is a Free-desktop (Linux) mechanism; macOS/Windows badge support, if any,
// is handled by their respective native notification paths.
func newPlatformUpdater(_ string) Updater {
	return &noopUpdater{}
}
