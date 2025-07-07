package mapper

import (
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
)

func ToUserResponse(user *entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsCustomer: user.IsCustomer,
	}
}
