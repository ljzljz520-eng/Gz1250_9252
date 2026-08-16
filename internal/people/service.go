package people

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Clock func() time.Time

type Service struct {
	repo  Repository
	clock Clock
}

func NewService(repo Repository, clock Clock) *Service {
	if clock == nil {
		clock = func() time.Time { return time.Unix(0, 0).UTC() }
	}
	return &Service{repo: repo, clock: clock}
}

func FixedClock(value time.Time) Clock {
	value = value.UTC()
	return func() time.Time { return value }
}

func (s *Service) List(ctx context.Context) ([]Person, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (Person, error) {
	if err := validateCreate(input); err != nil {
		return Person{}, err
	}
	person := Person{
		Name:      strings.TrimSpace(input.Name),
		Phone:     strings.TrimSpace(input.Phone),
		Email:     strings.TrimSpace(input.Email),
		Role:      input.Role,
		Status:    input.Status,
		CreatedAt: s.clock().UTC(),
	}
	if person.Status == "" {
		person.Status = StatusActive
	}
	created, err := s.repo.Create(ctx, person)
	if err != nil {
		if errors.Is(err, ErrContactConflict) {
			return Person{}, fmt.Errorf("create person failed: %s", err)
		}
		return Person{}, fmt.Errorf("create person: %w", err)
	}
	return created, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status Status) (Person, error) {
	if strings.TrimSpace(id) == "" {
		return Person{}, ErrInvalidRequest
	}
	if !validStatus(status) {
		return Person{}, ErrInvalidStatus
	}
	person, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return Person{}, fmt.Errorf("update person status: %w", err)
	}
	return person, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidRequest
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return nil
}

func validateCreate(input CreateInput) error {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Phone) == "" {
		return ErrInvalidRequest
	}
	if !validRole(input.Role) {
		return ErrInvalidRole
	}
	if input.Status != "" && !validStatus(input.Status) {
		return ErrInvalidStatus
	}
	return nil
}

func validRole(role Role) bool {
	switch role {
	case RolePhotographer, RolePhotoEditor, RoleMakeupArtist, RoleCustomerService:
		return true
	default:
		return false
	}
}

func validStatus(status Status) bool {
	return status == StatusActive || status == StatusInactive
}
