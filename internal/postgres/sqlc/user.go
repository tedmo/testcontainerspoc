package sqlc

import "github.com/tedmo/testcontainerspoc/internal/app"

func (u *User) DomainModel() *app.User {
	return &app.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
