package people

import (
	"strings"
	"testing"
	"time"
)

func TestLoadFixtureUsesDeterministicCreatedAt(t *testing.T) {
	fixture := `people:
  - name: fixture photographer
    phone: "13800000777"
    role: photographer
    status: active
`
	created, err := LoadFixture(strings.NewReader(fixture), FixedClock(time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)))
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0].CreatedAt.Format(time.RFC3339) != "2026-02-03T04:05:06Z" {
		t.Fatalf("fixture people = %#v", created)
	}
}
