package postgres

import (
	"fmt"

	"github.com/abayomipopoola/reddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) User(id uuid.UUID) (reddit.User, error) {
	var u reddit.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return reddit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) UserByUsername(username string) (reddit.User, error) {
	var u reddit.User
	if err := s.Get(&u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return reddit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]reddit.User, error) {
	var uu []reddit.User
	if err := s.Select(&uu, `SELECT * FROM users`); err != nil {
		return []reddit.User{}, fmt.Errorf("error getting users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(u *reddit.User) error {
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		u.ID,
		u.Username,
		u.Password); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(u *reddit.User) error {
	if err := s.Get(u, `UPDATE users SET username = $1, password = $2 WHERE id = $3 RETURNING *`,
		u.Username,
		u.Password,
		u.ID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
