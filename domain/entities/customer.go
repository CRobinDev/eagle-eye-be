package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID           uuid.UUID    `db:"id"`
	OrderID      string       `db:"order_id"`
	CustomerTier CustomerTier `db:"customer_tier"`
	ApiKey       string       `db:"hashed_key"`
	Prefix       string       `db:"prefix"`
	CurrentUsage uint64       `db:"current_usage"`
	MonthlyLimit uint32       `db:"monthly_limit"`
	ExpiresAt    time.Time    `db:"expires_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	RevokedAt    time.Time    `db:"revoked_at"`
	LastUsed     time.Time    `db:"last_used"`

	Client User    `db:"-"`
	Order  Payment `db:"-"`
}

type CustomerTier uint8

const (
	CustomerTierUnknown CustomerTier = 0
	CustomerTierFree    CustomerTier = 1
	CustomerTierBasic   CustomerTier = 2
	CustomerTierPremium CustomerTier = 3
)

var (
	CustomerTierMap = map[CustomerTier]string{
		CustomerTierFree:    "Free",
		CustomerTierBasic:   "Basic",
		CustomerTierPremium: "Premium",
	}

	CustomerTierLimit = map[CustomerTier]uint32{
		CustomerTierFree:    20,
		CustomerTierBasic:   2500,
		CustomerTierPremium: 15000,
	}
)

func (t CustomerTier) String() string {
	if val, ok := CustomerTierMap[t]; ok {
		return val
	}
	return "Unknown"
}

func (t CustomerTier) IsValid() bool {
	_, ok := CustomerTierMap[t]
	return ok
}

func ValueOfCustomerTier(value string) CustomerTier {
	value = strings.ToLower(value)
	realValue := strings.ToUpper(string(value[0])) + value[1:]

	for k, v := range CustomerTierMap {
		if v == realValue {
			return k
		}
	}
	return CustomerTierUnknown
}

func TierLimit(t CustomerTier) uint32 {
	return CustomerTierLimit[t]
}
