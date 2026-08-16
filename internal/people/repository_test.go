package people

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestMemoryRepositoryConcurrentDuplicatePhoneHasSingleWinner(t *testing.T) {
	repo := NewMemoryRepository(nil)
	ready := make(chan struct{}, 2)
	start := make(chan struct{})
	results := make(chan error, 2)
	var group sync.WaitGroup
	for index := 0; index < 2; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			ready <- struct{}{}
			<-start
			_, err := repo.Create(context.Background(), Person{Name: "person", Phone: "13800000888", Role: RolePhotographer, Status: StatusActive})
			results <- err
		}()
	}
	<-ready
	<-ready
	close(start)
	group.Wait()
	close(results)

	var successes int
	var conflicts int
	for err := range results {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrContactConflict) {
			conflicts++
		} else {
			t.Fatalf("unexpected error = %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes = %d, conflicts = %d", successes, conflicts)
	}
	people, err := repo.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(people) != 1 || people[0].Phone != "13800000888" {
		t.Fatalf("stored people = %#v", people)
	}
}
