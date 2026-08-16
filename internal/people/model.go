package people

import "time"

type Role string

const (
	RolePhotographer    Role = "photographer"
	RolePhotoEditor     Role = "photo_editor"
	RoleMakeupArtist    Role = "makeup_artist"
	RoleCustomerService Role = "customer_service"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Person struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
	Role      Role      `json:"role"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateInput struct {
	Name   string `json:"name" yaml:"name"`
	Phone  string `json:"phone" yaml:"phone"`
	Email  string `json:"email" yaml:"email"`
	Role   Role   `json:"role" yaml:"role"`
	Status Status `json:"status" yaml:"status"`
}

type StatusInput struct {
	Status Status `json:"status"`
}
