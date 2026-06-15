package app

// refreshLauncherBadge recomputes the unified-inbox unread total and updates the
// desktop taskbar unread-count badge (Unity LauncherEntry on Linux, no-op
// elsewhere). It is safe to call before the badge updater is initialized (no-op)
// and from any goroutine.
func (a *App) refreshLauncherBadge() {
	if a.launcherBadge == nil {
		return
	}
	count, err := a.GetUnifiedInboxUnreadCount()
	if err != nil {
		return
	}
	a.launcherBadge.SetCount(count)
}
