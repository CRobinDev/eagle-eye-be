package errorz

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Errors struct {
	Code    int        `json:"code"`
	Err     error      `json:"error"`
	TraceID *uuid.UUID `json:"trace_id,omitempty"`
}

func (e *Errors) Error() string {
	return e.Err.Error()
}

func NewError(code int, err string) *Errors {
	return &Errors{
		Code: code,
		Err:  errors.New(err),
	}
}

func (e *Errors) WithTraceID(traceID uuid.UUID) *Errors {
	e.TraceID = &traceID
	return e
}

var (
	// User Errors
	ErrFailedToCreateUser       = NewError(fiber.StatusInternalServerError, "Failed to save user. Please try again later")
	ErrUserNotFound             = NewError(fiber.StatusNotFound, "User not found. Please check your details and try again.")
	ErrUserAlreadyExists        = NewError(fiber.StatusConflict, "This email is already registered. Please try logging in or use a different email.")
	ErrHashPassword             = NewError(fiber.StatusInternalServerError, "Something went wrong while processing your password. Please try again.")
	ErrFailedToGenerateJWT      = NewError(fiber.StatusInternalServerError, "There was an error generating your authentication token. Please try again.")
	ErrFailedToDecodeJWT        = NewError(fiber.StatusInternalServerError, "There was an error decoding the token. Please try again.")
	ErrInvalidEmail             = NewError(fiber.StatusBadRequest, "Please provide a valid email address.")
	ErrInvalidPassword          = NewError(fiber.StatusBadRequest, "Your password must be at least 8 characters long and contain a mix of letters, numbers, and special characters.")
	ErrInvalidToken             = NewError(fiber.StatusUnauthorized, "The token is invalid. Please try again.")
	ErrFailedToSendNotification = NewError(fiber.StatusInternalServerError, "There was an error sending the notification. Please try again later.")
	ErrUnauthorized             = NewError(fiber.StatusUnauthorized, "Invalid login credentials. Please check your username and password.")
	ErrCredentialMismatch       = NewError(fiber.StatusUnauthorized, "The credentials provided do not match our records.")
	ErrForbiddenRole            = NewError(fiber.StatusForbidden, "You do not have permission to access this resource.")
	ErrUserNotVerified          = NewError(fiber.StatusForbidden, "Your account is not verified. Please check your email to verify your account.")
	ErrUserAlreadyVerified      = NewError(fiber.StatusForbidden, "Your account is already verified. You can log in now.")
	ErrFailedToDeleteUser       = NewError(fiber.StatusInternalServerError, "Failed to delete user. Please try again later.")
	ErrFailedToUpdateUser       = NewError(fiber.StatusInternalServerError, "Failed to update user. Please try again later.")
	ErrFailedToUpdatePassword   = NewError(fiber.StatusInternalServerError, "Failed to update password. Please try again later.")
	ErrPasswordMismatch         = NewError(fiber.StatusBadRequest, "Password Mismatch. Please input the correct password validation")
	ErrSetHTMLTemplate          = NewError(fiber.StatusInternalServerError, "There was an error setting the HTML template. Please try again.")
	ErrExecuteHTML              = NewError(fiber.StatusInternalServerError, "We encountered an issue while processing the page. Please try again.")

	// Google Errors
	ErrStateNoMatch        = NewError(fiber.StatusBadRequest, "The state does not match. Please try again.")
	ErrFailedFetchUserInfo = NewError(fiber.StatusInternalServerError, "Failed to fetch your Google user data. Please try again later.")
	ErrReadResponseBody    = NewError(fiber.StatusInternalServerError, "There was an issue reading the response. Please try again.")
	ErrUnmarshal           = NewError(fiber.StatusBadRequest, "There was a problem with the data we received. Please try again.")
	ErrRegisterFailed      = NewError(fiber.StatusInternalServerError, "An error occurred while registering. Please try again later.")
	ErrInvalidCode         = NewError(fiber.StatusBadRequest, "The code provided is invalid. Please check and try again.")

	// Gemini Errors
	ErrFailedToAnalyzeFood            = NewError(fiber.StatusInternalServerError, "Failed to analyze food. Please try again later.")
	ErrFailedToOpenFile               = NewError(fiber.StatusInternalServerError, "Failed to open the file. Please check the file and try again.")
	ErrFailedToReadFile               = NewError(fiber.StatusInternalServerError, "Failed to read the file. Please check the file and try again.")
	ErrFailedToAnalyzeImage           = NewError(fiber.StatusInternalServerError, "Failed to analyze the image. Please try again later.")
	ErrFailedToSaveAnalyzation        = NewError(fiber.StatusInternalServerError, "Failed to save analyzation history. Please try again later.")
	ErrFailedToGenerateRecommendation = NewError(fiber.StatusInternalServerError, "Failed to generate recommendation. Please try again later.")

	// Payment Errors
	ErrFailedToCreatePayment          = NewError(fiber.StatusInternalServerError, "Failed to create payment. Please try again later.")
	ErrPaymentNotFound                = NewError(fiber.StatusNotFound, "Payment not found. Please check your details and try again.")
	ErrSavePayment                    = NewError(fiber.StatusInternalServerError, "Payment Internal Server Error")
	ErrUpdatePaymentStatus            = NewError(fiber.StatusInternalServerError, "Update Status Internal Server Error")
	ErrFetchPaymentStatus             = NewError(fiber.StatusInternalServerError, "Fetch Status Internal Server Error")
	ErrInvalidSignatureKey            = NewError(fiber.StatusBadRequest, "Invalid signature key. Please check your request and try again.")
	ErrFailedToGetLatestPaymentStatus = NewError(fiber.StatusInternalServerError, "Failed to get the latest payment status. Please try again later.")

	// Customer Errors
	ErrFailedToGenerateKey   = NewError(fiber.StatusInternalServerError, "Something went wrong when generating api key. Please try again later.")
	ErrInvalidCustomerTier   = NewError(fiber.StatusBadRequest, "Invalid customer tier. Please choose one of the list !")
	ErrPaymentUnsuccess      = NewError(fiber.StatusPaymentRequired, "Payment still pending. Try again later")
	ErrFailedToSaveCustomer  = NewError(fiber.StatusInternalServerError, "Failed to save customer information. Please try again later.")
	ErrFailedToUpdateAPIKey  = NewError(fiber.StatusInternalServerError, "Failed to update api key. Please try again later.")
	ErrCustomerNotRegistered = NewError(fiber.StatusNotFound, "There is no customer for your account. Please try to subcribe first !")

	// Detection Errors
	ErrMissingAPIKey              = NewError(fiber.StatusUnauthorized, "Key Empty. Please try again later!")
	ErrInvalidAPIKey              = NewError(fiber.StatusUnauthorized, "Invalid Key. Don't try to brute force!")
	ErrMismatchAPIKey             = NewError(fiber.StatusUnauthorized, "Mismatch Key. Please input the correct one !")
	ErrFailedToCreateFormFile     = NewError(fiber.StatusInternalServerError, "Failed to create form file. Please try again later.")
	ErrFailedToWriteBytes         = NewError(fiber.StatusInternalServerError, "Failed to write file bytes. Please try again later.")
	ErrFailedToCreateHTTPRequest  = NewError(fiber.StatusInternalServerError, "Failed to create HTTP request. Please try again later.")
	ErrFailedToInitiateHTTPClient = NewError(fiber.StatusInternalServerError, "Failed to make HTTP request. Please try again later.")
	ErrFailedToCreateDetection    = NewError(fiber.StatusInternalServerError, "Failed to create detection. Please try again later.")
	ErrSaveDetection              = NewError(fiber.StatusInternalServerError, "Failed to save detection history. Please try again later.")
	ErrBanIP                      = NewError(fiber.StatusInternalServerError, "Failed to ban ip. Please try again later.")
	ErrFailedToUpdateUsage        = NewError(fiber.StatusInternalServerError, "Failed to update customer usage. Please try again later.")
	ErrMonthlyLimitReached        = NewError(fiber.StatusPaymentRequired, "Monthly limit reached. Please subscribe to a higher tier.")

	// Other Errors
	ErrFailedGenerateUUIDV7 = NewError(fiber.StatusInternalServerError, "Failed to generate UUIDV7 for User ID. Please try again later !")
)
