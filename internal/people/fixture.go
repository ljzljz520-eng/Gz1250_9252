package people

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

type Fixture struct {
	People []CreateInput `yaml:"people"`
}

func LoadFixture(reader io.Reader, clock Clock) ([]Person, error) {
	var fixture Fixture
	if err := yaml.NewDecoder(reader).Decode(&fixture); err != nil {
		return nil, fmt.Errorf("decode people fixture: %w", err)
	}
	if clock == nil {
		clock = FixedClock(zeroTime)
	}
	people := make([]Person, 0, len(fixture.People))
	for _, input := range fixture.People {
		if err := validateCreate(input); err != nil {
			return nil, fmt.Errorf("invalid people fixture: %w", err)
		}
		people = append(people, Person{
			Name:      input.Name,
			Phone:     input.Phone,
			Email:     input.Email,
			Role:      input.Role,
			Status:    input.Status,
			CreatedAt: clock().UTC(),
		})
		if people[len(people)-1].Status == "" {
			people[len(people)-1].Status = StatusActive
		}
	}
	return people, nil
}

var zeroTime = timeValue()

func timeValue() (value time.Time) {
	return time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
}
