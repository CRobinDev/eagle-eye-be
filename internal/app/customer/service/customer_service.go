package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
)

type customerService struct {
	cr     interfaces.ICustomerRepository
	ur     interfaces.IUserRepository
	ps     interfaces.IPaymentService
	logger *logrus.Logger
}

func NewCustomerService(cr interfaces.ICustomerRepository, ur interfaces.IUserRepository, ps interfaces.IPaymentService, logger *logrus.Logger) interfaces.ICustomerService {
	return &customerService{
		cr:     cr,
		ur:     ur,
		ps:     ps,
		logger: logger,
	}
}

func (cs *customerService) GenerateAPIKey(ctx context.Context, req dto.GenerateAPIKeyRequest) (dto.GenerateAPIKeyResponse, error) {
	traceID := utils.GetTraceID(ctx)
	_, err := cs.ur.GetUserByID(ctx, req.UserID)
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GenerateAPIKey] can't find user by id")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrUserNotFound
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GenerateAPIKey] failed to generate api key")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToGenerateKey.WithTraceID(traceID)
	}

	var orderID string
	var tier string
	order, err := cs.ps.LatestPaymentStatus(ctx, dto.GetPaymentStatusRequest{UserID: req.UserID})
	if err != nil && order == (dto.PaymentResponse{}) {
		orderID = "free"
		tier = "free"
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Warn("[CustomerService][GenerateAPIKey] latest payment status not found, assuming free tier")
	} else {
		if order.Status != "success" {
			cs.logger.WithFields(log.WithTraceID(traceID, err)).Warnf("[CustomerService][GenerateAPIKey] latest payment status is %v, waiting for payment to complete", order.Status)
			return dto.GenerateAPIKeyResponse{}, fmt.Errorf("payment still %v, try again later until success", order.Status)
		}
		tier = order.Tier
		orderID = order.OrderID
	}

	custTier := entities.ValueOfCustomerTier(tier)

	log.NewLogger().Println(tier)

	expiresAt := utils.GetCurrentTime().AddDate(0, 0, 30)
	customer := entities.Customer{
		ID:           req.UserID,
		OrderID:      orderID,
		CustomerTier: custTier,
		ApiKey:       utils.HashAPIKey(req.Prefix, apiKey),
		Prefix:       req.Prefix,
		MonthlyLimit: entities.TierLimit(custTier),
		ExpiresAt:    expiresAt,
	}

	if err := cs.cr.CreateCustomer(ctx, &customer); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			cs.logger.WithFields(log.WithTraceID(traceID, err)).Warn("[CustomerService][GenerateAPIKey] duplicate key, possibly prefix already exists")
			return dto.GenerateAPIKeyResponse{}, errorz.ErrDuplicatePrefix.WithTraceID(traceID)
		}
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GenerateAPIKey] failed to save customer information")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToSaveCustomer.WithTraceID(traceID)
	}

	go func() {
		var user entities.User
		user.ID = req.UserID
		user.IsCustomer = true

		if err := cs.ur.UpdateUser(context.Background(), &user); err != nil {
			cs.logger.Errorf("[CustomerService][GenerateAPIKey] failed to update user status : %v", err)
		}
	}()

	return dto.GenerateAPIKeyResponse{
		APIKey:    req.Prefix + "_" + apiKey,
		ExpiresAt: expiresAt,
	}, nil
}

func (cs *customerService) UpdateAPIKey(ctx context.Context, req dto.GenerateAPIKeyRequest) (dto.GenerateAPIKeyResponse, error) {
	traceID := utils.GetTraceID(ctx)
	_, err := cs.ur.GetUserByID(ctx, req.UserID)
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][UpdateAPIKey] can't find user by id")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrUserNotFound.WithTraceID(traceID)
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][UpdateAPIKey] failed to generate api key")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToGenerateKey.WithTraceID(traceID)
	}

	var customer entities.Customer
	customer.ID = req.UserID
	customer.ApiKey = utils.HashAPIKey(req.Prefix, apiKey)
	customer.Prefix = req.Prefix

	respCustomer, err := cs.cr.UpdateAPIKey(ctx, &customer)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			if pqErr.ConstraintName == "idx_customers_prefix" {
				cs.logger.WithFields(log.WithTraceID(traceID, err)).Warn("[CustomerService][GenerateAPIKey] duplicate key, possibly prefix already exists")
				return dto.GenerateAPIKeyResponse{}, errorz.ErrDuplicatePrefix.WithTraceID(traceID)
			}
		}
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][UpdateAPIKey] failed to update api key on db")
		return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToUpdateAPIKey.WithTraceID(traceID)
	}

	return dto.GenerateAPIKeyResponse{
		APIKey:    req.Prefix + "_" + apiKey,
		ExpiresAt: respCustomer.ExpiresAt,
	}, nil
}

func (cs *customerService) UpdateUsage(ctx context.Context, req dto.UpdateUsageRequest) (dto.UpdateUsageResponse, error) {
	customer, err := cs.cr.UpdateUsage(ctx, req.Prefix)
	if err != nil {
		return dto.UpdateUsageResponse{}, err
	}

	return dto.UpdateUsageResponse{
		ID:           customer.ID,
		CustomerTier: customer.CustomerTier.String(),
		HashedKey:    customer.ApiKey,
		Prefix:       customer.Prefix,
		CurrentUsage: customer.CurrentUsage,
		MonthlyLimit: customer.MonthlyLimit,
		LastUsed:     customer.LastUsed.Time,
	}, err
}

func (cs *customerService) GetCustomerAPIStatus(ctx context.Context, req dto.GetCustomerRequest) (dto.GetCustomerAPIStatusResponse, error) {
	customer, _ := cs.cr.GetCustomerByID(ctx, req.CustomerID)

	var cust entities.Customer
	var isCustomer bool
	if customer == cust {
		isCustomer = false
	} else {
		isCustomer = true
	}

	return dto.GetCustomerAPIStatusResponse{
		IsCustomer: isCustomer,
		ExpiresAt:  customer.ExpiresAt,
	}, nil
}

func (cs *customerService) GetCustomerTotalCalls(ctx context.Context, req dto.GetCustomerRequest) (dto.GetCustomerTotalCallsResponse, error) {
	customer, _ := cs.cr.GetCustomerByID(ctx, req.CustomerID)

	var cust entities.Customer
	var totalCalls uint64
	if customer == cust {
		totalCalls = 0
	} else {
		totalCalls = customer.CurrentUsage
	}

	return dto.GetCustomerTotalCallsResponse{
		TotalCalls: totalCalls,
		ExpiresAt:  customer.ExpiresAt,
	}, nil
}

func (cs *customerService) GetCustomerTier(ctx context.Context, req dto.GetCustomerRequest) (dto.GetCustomerTierResponse, error) {
	traceID := utils.GetTraceID(ctx)
	var userTier string

	tier, err := cs.cr.GetCustomerTier(ctx, req.CustomerID)
	userTier = strings.ToLower(tier.String())
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Warn("[CustomerService][GetCustomerTier] customer tier not found, assuming free tier")
		userTier = "free"
	}

	return dto.GetCustomerTierResponse{
		Tier: userTier,
	}, nil
}
