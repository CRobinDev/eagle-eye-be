package service

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"github.com/CRobinDev/karsa/domain/dto"
	"github.com/CRobinDev/karsa/domain/entities"
	"github.com/CRobinDev/karsa/domain/interfaces"
	"github.com/CRobinDev/karsa/pkg/errorz"
	"github.com/CRobinDev/karsa/pkg/log"
	"github.com/CRobinDev/karsa/pkg/midtrans"
	"github.com/CRobinDev/karsa/pkg/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type paymentService struct {
	pr       interfaces.IPaymentRepository
	midtrans midtrans.IMidtrans
	logger   *logrus.Logger
}

func NewPaymentService(pr interfaces.IPaymentRepository, midtrans midtrans.IMidtrans, logger *logrus.Logger) *paymentService {
	return &paymentService{
		pr:       pr,
		midtrans: midtrans,
		logger:   logger,
	}
}

func (ps *paymentService) CreatePayment(ctx context.Context, req dto.PaymentRequest) (dto.PaymentResponse, error) {
	traceID := utils.GetTraceID(ctx)
	orderID := uuid.NewString()

	req.OrderID = orderID
	snapResp, err := ps.midtrans.NewTransactionToken(req)
	if err != nil {
		ps.logger.WithFields(log.WithTraceID(traceID, err)).Error("[PaymentService][CreatePayment] failed to create snap midtrans payment")
		return dto.PaymentResponse{}, errorz.ErrFailedToCreatePayment.WithTraceID(traceID)
	}

	payment := entities.Payment{
		OrderID: orderID,
		UserID:  req.UserID,
		Amount:  float64(req.Amount),
		Tier:    req.TierOrder,
	}

	if err := ps.pr.CreatePayment(ctx, &payment); err != nil {
		ps.logger.WithFields(log.WithTraceID(traceID, err)).Error("[PaymentService][CreatePayment] failed to save payment")
		return dto.PaymentResponse{}, errorz.ErrFailedToCreatePayment.WithTraceID(traceID)
	}

	return dto.PaymentResponse{
		SnapURL: snapResp.RedirectURL,
		OrderID: orderID,
		Status:  "pending",
	}, nil

}

func (ps *paymentService) UpdatePaymentStatus(ctx context.Context, req dto.UpdatePaymentStatusRequest) (dto.PaymentResponse, error) {
	var status string

	traceID := utils.GetTraceID(ctx)
	midtransKey := ps.midtrans.GetServerKey()

	hash := sha512.New()
	hash.Write([]byte(req.OrderID + req.Code + req.Amount + midtransKey))

	hashedSignature := hex.EncodeToString(hash.Sum(nil))
	if equals := strings.Compare(hashedSignature, req.SignatureKey); equals != 0 {
		ps.logger.WithFields(log.WithTraceID(traceID, errorz.ErrInvalidSignatureKey)).Error("[PaymentService][UpdatePaymentStatus] invalid signature key")
		return dto.PaymentResponse{}, errorz.ErrInvalidSignatureKey.WithTraceID(traceID)
	}

	switch req.TransactionStatus {
	case "capture":
		switch req.FraudStatus {
		case "challenge":
			status = "challenge"
		case "accept":
			status = "success"
		default:
			status = "unknown"
		}

	case "settlement":
		status = "success"

	case "cancel", "expire":
		status = "failure"

	case "pending":
		status = "pending"

	case "deny":
		status = "denied"

	default:
		status = "unknown"
	}

	payment := entities.Payment{
		Status:  status,
		OrderID: req.OrderID,
	}

	if err := ps.pr.UpdatePaymentStatus(ctx, &payment); err != nil {
		ps.logger.WithFields(log.WithTraceID(traceID, err)).Error("[PaymentService][UpdatePaymentStatus] failed to create snap midtrans payment")
		return dto.PaymentResponse{}, errorz.ErrUpdatePaymentStatus.WithTraceID(traceID)
	}

	// if status == "success" {
	// 	go func() {
	// 	if err := ps.gomail.SendConsultationEmail(emailReq); err != nil {
	// 		ps.logger.WithFields(map[string]interface{}{
	// 			"error": err.Error(),
	// 		}).Error("[userService.VerificationEmail] failed to send consultation email (async)")
	// 	}
	// }()

	// }

	return dto.PaymentResponse{
		OrderID: req.OrderID,
		Status:  status,
	}, nil
}

func (ps *paymentService) LatestPaymentStatus(ctx context.Context, req dto.GetPaymentStatusRequest) (dto.PaymentResponse, error) {
	traceID := utils.GetTraceID(ctx)
	resp, err := ps.pr.GetStatusByUserID(ctx, req.UserID)
	if err != nil {
		return dto.PaymentResponse{}, errorz.ErrFetchPaymentStatus.WithTraceID(traceID)
	}

	return dto.PaymentResponse{
		OrderID: resp.OrderID,
		Status:  resp.Status,
	}, nil
}
