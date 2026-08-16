package people

import (
	"context"
	"fmt"
	"sync"
)

type Repository interface {
	List(context.Context) ([]Person, error)
	Create(context.Context, Person) (Person, error)
	UpdateStatus(context.Context, string, Status) (Person, error)
	Delete(context.Context, string) error
}

type MemoryRepository struct {
	mu     sync.RWMutex
	people map[string]Person
	order  []string
	nextID int
}

func NewMemoryRepository(initial []Person) *MemoryRepository {
	repo := &MemoryRepository{people: make(map[string]Person), nextID: 1}
	for _, person := range initial {
		if person.ID == "" {
			person.ID = repo.allocateID()
		}
		repo.people[person.ID] = clonePerson(person)
		repo.order = append(repo.order, person.ID)
	}
	return repo
}

func (r *MemoryRepository) List(ctx context.Context) ([]Person, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]Person, 0, len(r.order))
	for _, id := range r.order {
		items = append(items, clonePerson(r.people[id]))
	}
	return items, nil
}

func (r *MemoryRepository) Create(ctx context.Context, person Person) (Person, error) {
	if err := ctx.Err(); err != nil {
		return Person{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.people {
		if existing.Phone != "" && existing.Phone == person.Phone {
			return Person{}, ErrContactConflict
		}
	}
	person.ID = r.allocateID()
	r.people[person.ID] = clonePerson(person)
	r.order = append(r.order, person.ID)
	return clonePerson(person), nil
}

func (r *MemoryRepository) UpdateStatus(ctx context.Context, id string, status Status) (Person, error) {
	if err := ctx.Err(); err != nil {
		return Person{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	person, ok := r.people[id]
	if !ok {
		return Person{}, ErrNotFound
	}
	person.Status = status
	r.people[id] = person
	return clonePerson(person), nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.people[id]; !ok {
		return ErrNotFound
	}
	delete(r.people, id)
	for index, item := range r.order {
		if item == id {
			r.order = append(r.order[:index], r.order[index+1:]...)
			break
		}
	}
	return nil
}

func (r *MemoryRepository) allocateID() string {
	id := "person-" + formatID(r.nextID)
	r.nextID++
	return id
}

func formatID(value int) string {
	return fmt.Sprintf("%03d", value)
}

func clonePerson(person Person) Person {
	return person
}
