package usercfg_test

import (
	"testing"
	"time"

	"github.com/vogtp/rag/pkg/usercfg"
)

func TestCollection_CalculateNextDBUpdate(t *testing.T) {
	tests := []struct {
		name            string
		now             string
		last            string
		updateIntervall time.Duration
	}{
		{"early", "2015-02-26 04:25:35", "2015-02-25 04:25:35", 24 * time.Hour},
		{"moring", "2015-02-26 09:13:12", "2015-02-25 09:13:12", 24 * time.Hour},
		{"noon", "2015-02-26 12:00:00", "2015-02-25 12:00:00", 24 * time.Hour},
		{"afternoon", "2015-02-26 16:24:32", "2015-02-25 16:24:32", 24 * time.Hour},
		{"night", "2015-02-26 23:45:15", "2015-02-25 23:45:15", 24 * time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now, err := time.Parse(time.DateTime, tt.now)
			if err != nil {
				t.Fatalf("cannot parse time: %v", err)
			}
			var c usercfg.Collection
			c.DBUpdateIntervall = tt.updateIntervall
			c.CalculateNextDBUpdate(now)
			delta := c.NextDBUpdate.Sub(now)
			if delta > tt.updateIntervall+12*time.Hour {
				t.Errorf("Timediff to big: %s last: %s next: %s", delta, now, c.NextDBUpdate)
			}
			if delta < tt.updateIntervall-12*time.Hour {
				t.Errorf("Timediff to small: %s last: %s next: %s", delta, now, c.NextDBUpdate)
			}
			h := c.NextDBUpdate.Hour()
			if h < 20 && h > 6 {
				t.Errorf("Time during the day: last: %s next: %s", now, c.NextDBUpdate)
			}
		})
	}
}
