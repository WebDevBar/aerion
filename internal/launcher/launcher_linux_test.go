//go:build linux

package launcher

import (
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

// TestLauncherEntryEmit verifies that SetCount broadcasts a well-formed
// com.canonical.Unity.LauncherEntry "Update" signal on the session bus. It
// requires a running session bus and is skipped otherwise (e.g. headless CI).
func TestLauncherEntryEmit(t *testing.T) {
	recv, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Skipf("no session bus available: %v", err)
	}
	defer recv.Close()

	if err := recv.AddMatchSignal(
		dbus.WithMatchInterface(launcherEntryInterface),
		dbus.WithMatchMember("Update"),
	); err != nil {
		t.Fatalf("AddMatchSignal: %v", err)
	}
	signals := make(chan *dbus.Signal, 8)
	recv.Signal(signals)

	u := New("io.github.hkdb.Aerion.desktop")
	defer u.Close()
	u.SetCount(7)

	for {
		select {
		case sig := <-signals:
			if len(sig.Body) != 2 {
				continue
			}
			appURI, _ := sig.Body[0].(string)
			if appURI != "application://io.github.hkdb.Aerion.desktop" {
				continue
			}
			props, ok := sig.Body[1].(map[string]dbus.Variant)
			if !ok {
				t.Fatalf("props not a map: %T", sig.Body[1])
			}
			if got := props["count"].Value(); got != int64(7) {
				t.Fatalf("count = %v (%T), want int64(7)", got, got)
			}
			if got := props["count-visible"].Value(); got != true {
				t.Fatalf("count-visible = %v, want true", got)
			}
			return
		case <-time.After(2 * time.Second):
			t.Fatal("no matching LauncherEntry Update signal received within 2s")
		}
	}
}
