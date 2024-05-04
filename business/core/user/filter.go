package user

import (
	"github/islamghany/blog/foundation/validate"
	"time"

	"github.com/google/uuid"
)

type QueryFilter struct {
	ID        *uuid.UUID `validate:"omitempty"`
	Username  *string    `validate:"omitempty,min=3,alphanum"`
	Email     *string    `validate:"omitempty"`
	FirstName *string    `validate:"omitempty"`
	LastName  *string    `validate:"omitempty"`
	CreatedAt *time.Time `validate:"omitempty"`
	UpdatedAt *time.Time `validate:"omitempty"`
}

func (qf *QueryFilter) Validate() error {
	return validate.Check(qf)
}

func (qf *QueryFilter) WithUserID(id uuid.UUID) {
	qf.ID = &id
}

func (qf *QueryFilter) WithUsername(username string) {
	qf.Username = &username
}

func (qf *QueryFilter) WithEmail(email string) {
	qf.Email = &email
}

func (qf *QueryFilter) WithFirstName(firstName string) {
	qf.FirstName = &firstName
}
func (qf *QueryFilter) WithLastName(lastName string) {
	qf.LastName = &lastName
}

func (qf *QueryFilter) WithCreatedAt(createdAt time.Time) {
	d := createdAt.UTC()
	qf.CreatedAt = &d
}

func (qf *QueryFilter) WithUpdatedAt(updatedAt time.Time) {
	d := updatedAt.UTC()
	qf.UpdatedAt = &d
}
