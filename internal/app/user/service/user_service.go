package service

import (
	"context"
	"errors"
	"strings"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	ur     interfaces.IUserRepository
	logger *logrus.Logger
}

func NewUserService(ur interfaces.IUserRepository, logger *logrus.Logger) interfaces.IUserService {
	return &userService{
		ur:     ur,
		logger: logger,
	}
}

func (us *userService) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) error {
	traceID := utils.GetTraceID(ctx)
	user, err := us.ur.GetUserByID(ctx, req.UserID)
	if err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdateUser] user not found by id")
		return errorz.ErrUserNotFound.WithTraceID(traceID)
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if err := us.ur.UpdateUser(ctx, &user); err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdateUser] failed to update user")
		return errorz.ErrFailedToUpdateUser.WithTraceID(traceID)
	}

	return nil
}

func (us *userService) UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest) error {
	traceID := utils.GetTraceID(ctx)
	user, err := us.ur.GetUserByID(ctx, req.UserID)
	if err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdatePassword] user not found by id")
		return errorz.ErrUserNotFound.WithTraceID(traceID)
	}

	result := strings.Compare(req.NewPassword, req.ConfirmPassword)
	if result == 1 {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdatePassword] password mismatch")
		return errorz.ErrPasswordMismatch.WithTraceID(traceID)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), env.GetEnv().BcryptCost)
	if err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdatePassword] failed to hash password")
		return errorz.ErrHashPassword.WithTraceID(traceID)
	}

	user.Password = string(hashedPassword)

	if err := us.ur.UpdatePassword(ctx, &user); err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][UpdatePassword] failed to update password")
		return errorz.ErrFailedToUpdatePassword.WithTraceID(traceID)
	}

	return nil
}

func (us *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	traceID := utils.GetTraceID(ctx)
	_, err := us.ur.GetUserByID(ctx, id)
	if err != nil {
		us.logger.WithFields(log.WithTraceID(traceID, err)).Error("[UserService][DeleteUser] user not found by id")
		return errorz.ErrUserNotFound
	}

	if err := us.ur.DeleteUser(ctx, id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return errorz.ErrFailedToDeleteUser.WithTraceID(traceID)
		} else {
			return err
		}
	}

	return nil
}
