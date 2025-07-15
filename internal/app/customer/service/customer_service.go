package service

import (
	"context"
	"fmt"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/utils"
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

	expiresAt := utils.GetCurrentTime().AddDate(0, 0, 30)

	tier := entities.ValueOfCustomerTier(req.CustomerTier)
	if tier == entities.CustomerTierUnknown {
		cs.logger.Warnf("Invalid customer tier : %v", req.CustomerTier)
		return dto.GenerateAPIKeyResponse{}, errorz.ErrInvalidCustomerTier
	}
	log.NewLogger().Println(tier)
	customer := entities.Customer{
		ID:           req.UserID,
		CustomerTier: tier,
		ApiKey:       utils.HashAPIKey(req.Prefix, apiKey),
		Prefix:       req.Prefix,
		MonthlyLimit: entities.TierLimit(tier),
		ExpiresAt:    expiresAt,
	}
	
	customer.OrderID = "free"
	if tier != entities.CustomerTierFree {
		var order dto.PaymentResponse
		if tier != entities.CustomerTierFree {
			order, err = cs.ps.LatestPaymentStatus(ctx, dto.GetPaymentStatusRequest{UserID: req.UserID})
			if err != nil {
				cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GenerateAPIKey] failed to fetch latest payment")
				return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToGetLatestPaymentStatus.WithTraceID(traceID)
			}

			if order.Status != "success" {
				return dto.GenerateAPIKeyResponse{}, fmt.Errorf("payment still %v, try again later until success", order.Status)
			}
			customer.OrderID = order.OrderID
		}
	}

	if err := cs.cr.CreateCustomer(ctx, &customer); err != nil {
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
		return dto.GenerateAPIKeyResponse{}, errorz.ErrFailedToUpdateAPIKey.WithTraceID(traceID)
	}

	return dto.GenerateAPIKeyResponse{
		APIKey:    req.Prefix + "_" + apiKey,
		ExpiresAt: respCustomer.ExpiresAt,
	}, nil
}

func (cs *customerService) GetCustomerAPIStatus(ctx context.Context, req dto.GetCustomerRequest) (dto.GetCustomerAPIStatusResponse, error) {
	traceID := utils.GetTraceID(ctx)
	customer, err := cs.cr.GetCustomerByID(ctx, req.CustomerID)
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GetCustomerAPIStatus] user not found")
		return dto.GetCustomerAPIStatusResponse{}, errorz.ErrUserNotFound.WithTraceID(traceID)
	}
	var cust entities.Customer
	if customer == cust {
		return dto.GetCustomerAPIStatusResponse{}, errorz.ErrCustomerNotRegistered.WithTraceID(traceID)
	}

	return dto.GetCustomerAPIStatusResponse{
		IsCustomer: true,
		ExpiresAt:  customer.ExpiresAt,
	}, nil
}

func (cs *customerService) GetCustomerTotalCalls(ctx context.Context, req dto.GetCustomerRequest) (dto.GetCustomerTotalCallsResponse, error) {
	traceID := utils.GetTraceID(ctx)
	customer, err := cs.cr.GetCustomerByID(ctx, req.CustomerID)
	if err != nil {
		cs.logger.WithFields(log.WithTraceID(traceID, err)).Error("[CustomerService][GetCustomerAPIStatus] user not found")
		return dto.GetCustomerTotalCallsResponse{}, errorz.ErrUserNotFound.WithTraceID(traceID)
	}
	var cust entities.Customer
	if customer == cust {
		return dto.GetCustomerTotalCallsResponse{}, errorz.ErrCustomerNotRegistered.WithTraceID(traceID)
	}

	return dto.GetCustomerTotalCallsResponse{
		TotalCalls: customer.CurrentUsage,
		ExpiresAt:  customer.ExpiresAt,
	}, nil
}
