package service

import (
	"context"
	"errors"

	"github.com/CRobinDev/karsa/config/env"
	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/internal/app/user/mapper"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/jwt"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	ar     interfaces.IAuthRepository
	logger *logrus.Logger
}

func NewAuthService(ar interfaces.IAuthRepository, logger *logrus.Logger) interfaces.IAuthService {
	return &authService{
		ar:     ar,
		logger: logger,
	}
}

func (as *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	traceID := utils.GetTraceID(ctx)
	userID, err := uuid.NewV7()
	if err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Register] failed to generate uuid")
		return errorz.ErrFailedGenerateUUIDV7.WithTraceID(traceID)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), env.GetEnv().BcryptCost)
	if err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Register] failed to hash password")
		return errorz.ErrHashPassword.WithTraceID(traceID)
	}

	user := entities.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := as.ar.CreateUser(ctx, &user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Register] user exists.")
			return errorz.ErrUserAlreadyExists.WithTraceID(traceID)
		}
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Register] failed to create user")
		return errorz.ErrFailedToCreateUser.WithTraceID(traceID)
	}

	return nil
}

func (as *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	traceID := utils.GetTraceID(ctx)
	user, err := as.ar.GetUserByEmail(ctx, req.Email)
	if err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Login] user not found by email")
		return dto.LoginResponse{}, errorz.ErrUserNotFound.WithTraceID(traceID)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Register] failed to compare password")
		return dto.LoginResponse{}, errorz.ErrCredentialMismatch.WithTraceID(traceID)
	}

	token, err := jwt.CreateToken(&user)
	if err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][Login] failed to generate jwt")
		return dto.LoginResponse{}, errorz.ErrFailedToGenerateJWT.WithTraceID(traceID)
	}

	return dto.LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}

func (as *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (dto.UserResponse, error) {
	traceID := utils.GetTraceID(ctx)
	user, err := as.ar.GetUserByID(ctx, userID)
	if err != nil {
		as.logger.WithFields(log.WithTraceID(traceID, err)).Error("[AuthService][GetUserByID] user not found by ID")
		return dto.UserResponse{}, errorz.ErrUserNotFound.WithTraceID(traceID)
	}

	return mapper.ToUserResponse(&user), nil
}
