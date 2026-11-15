package score

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScoringStopsAndReleasesBothTickers(t *testing.T) {
	// Header-only fixtures make every quest a no-op: no network or database.
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "data"), 0o700); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 5; i++ {
		if err := os.WriteFile(filepath.Join(root, "data", fmt.Sprintf("quest%d.csv", i)), []byte("ip,hostname\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Error(err)
		}
	})

	stop := make(chan bool)
	// Reuse the same stop channel just like Menu does on stop/start cycles.
	for cycle := 0; cycle < 3; cycle++ {
		scoring := time.NewTicker(10 * time.Millisecond)
		health := time.NewTicker(10 * time.Millisecond)
		defer scoring.Stop()
		defer health.Stop()
		done := make(chan struct{})
		go func() {
			Scoring(nil, scoring, health, stop)
			close(done)
		}()
		select {
		case stop <- true:
		case <-time.After(time.Second):
			t.Fatal("scoring did not receive stop")
		}
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("scoring did not return after stop")
		}
		// Old Go timer channels may retain one tick queued before Stop.
		for _, ticks := range []<-chan time.Time{scoring.C, health.C} {
			select {
			case <-ticks:
			default:
			}
		}
		select {
		case <-scoring.C:
			t.Fatal("scoring ticker still running after stop")
		case <-health.C:
			t.Fatal("health ticker still running after stop")
		case <-time.After(40 * time.Millisecond):
		}
	}
}
