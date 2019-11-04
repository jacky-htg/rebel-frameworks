package response

import "essentials/models"

// UserResponse is response format of user
type UserResponse struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

// Transform models.User to UserResponse
func (u *UserResponse) Transform(user *models.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.IsActive = user.IsActive
}
