package models

import "database/sql"

// User : struct of User
type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
	IsActive bool
}

// List of users
func (u *User) List(db *sql.DB) ([]User, error) {
	var list []User
	const q = `SELECT id, username, password, email, is_active FROM users`

	rows, err := db.Query(q)
	if err != nil {
		return list, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsActive); err != nil {
			return list, err
		}
		list = append(list, user)
	}

	return list, rows.Err()
}
