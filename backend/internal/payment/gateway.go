// Payment Gateway Abstraction Layer
// Supports multiple payment providers (Stripe, Juspay, Razorpay, etc.)

package payment

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ============================================
// Core Payment Gateway Interface
// ============================================

// PaymentGateway defines the contract all payment providers must implement
type PaymentGateway interface {
	// GetProvider returns the provider name (stripe, juspay, razorpay)
	GetProvider() string

	// CreatePaymentIntent creates a new payment intent/order
	CreatePaymentIntent(ctx context.Context, req *PaymentIntentRequest) (*PaymentIntentResponse, error)

	// ConfirmPayment confirms/captures a payment
	ConfirmPayment(ctx context.Context, paymentID string) (*PaymentResponse, error)

	// GetPaymentStatus retrieves the current status of a payment
	GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentResponse, error)

	// RefundPayment processes a full or partial refund
	RefundPayment(ctx context.Context, req *RefundRequest) (*RefundResponse, error)

	// CreateCustomer creates a customer profile for saved payment methods
	CreateCustomer(ctx context.Context, req *CustomerRequest) (*CustomerResponse, error)

	// AttachPaymentMethod attaches a payment method to a customer
	AttachPaymentMethod(ctx context.Context, customerID, paymentMethodID string) error
}

// ============================================
// Request/Response Types
// ============================================

type PaymentIntentRequest struct {
	Amount        int64             // Amount in cents
	Currency      string            // AUD, USD, INR, etc.
	CustomerID    string            // Optional: existing customer ID
	PaymentMethod string            // card, upi, wallet, etc.
	Description   string            // Payment description
	Metadata      map[string]string // Custom metadata
	ReturnURL     string            // For redirect-based flows
	WebhookURL    string            // For async notifications
}

type PaymentIntentResponse struct {
	PaymentID    string            // Gateway-specific payment ID
	OrderID      string            // Gateway-specific order ID (if applicable)
	ChargeID     string            // Gateway-specific charge ID (if applicable)
	ClientSecret string            // For frontend confirmation
	Status       string            // pending, requires_action, succeeded, etc.
	Metadata     map[string]string // Gateway-specific data
}

type PaymentResponse struct {
	PaymentID     string
	ChargeID      string
	CustomerID    string
	OrderID       string
	Amount        int64
	Currency      string
	Status        string // pending, processing, succeeded, failed
	FailureReason string
	Metadata      map[string]string
}

type RefundRequest struct {
	PaymentID string
	Amount    int64  // Amount to refund in cents (0 = full refund)
	Reason    string // Optional refund reason
}

type RefundResponse struct {
	RefundID string
	Amount   int64
	Status   string // pending, succeeded, failed
}

type CustomerRequest struct {
	Email       string
	Name        string
	Phone       string
	Description string
	Metadata    map[string]string
}

type CustomerResponse struct {
	CustomerID string
	Metadata   map[string]string
}

// ============================================
// Payment Gateway Factory
// ============================================

type GatewayFactory struct {
	gateways map[string]PaymentGateway
}

func NewGatewayFactory() *GatewayFactory {
	return &GatewayFactory{
		gateways: make(map[string]PaymentGateway),
	}
}

// Register registers a payment gateway
func (f *GatewayFactory) Register(gateway PaymentGateway) {
	f.gateways[gateway.GetProvider()] = gateway
}

// Get retrieves a payment gateway by provider name
func (f *GatewayFactory) Get(provider string) (PaymentGateway, error) {
	gateway, exists := f.gateways[provider]
	if !exists {
		return nil, errors.New("payment gateway not found: " + provider)
	}
	return gateway, nil
}

// ListProviders returns all registered provider names
func (f *GatewayFactory) ListProviders() []string {
	providers := make([]string, 0, len(f.gateways))
	for provider := range f.gateways {
		providers = append(providers, provider)
	}
	return providers
}

// ============================================
// Payment Service (Uses Gateway Abstraction)
// ============================================

type PaymentService struct {
	factory        *GatewayFactory
	paymentRepo    PaymentRepository
	defaultGateway string
}

func NewPaymentService(factory *GatewayFactory, repo PaymentRepository, defaultGateway string) *PaymentService {
	return &PaymentService{
		factory:        factory,
		paymentRepo:    repo,
		defaultGateway: defaultGateway,
	}
}

// CreatePayment creates a payment using the specified gateway
func (s *PaymentService) CreatePayment(ctx context.Context, provider string, req *PaymentIntentRequest) (*PaymentIntentResponse, error) {
	// Use default if provider not specified
	if provider == "" {
		provider = s.defaultGateway
	}

	// Get the appropriate gateway
	gateway, err := s.factory.Get(provider)
	if err != nil {
		return nil, err
	}

	// Create payment intent via gateway
	response, err := gateway.CreatePaymentIntent(ctx, req)
	if err != nil {
		return nil, err
	}

	// Save to database
	payment := &BookingPayment{
		ID:               uuid.New(),
		GatewayProvider:  provider,
		GatewayPaymentID: response.PaymentID,
		GatewayChargeID:  response.ChargeID,
		GatewayOrderID:   response.OrderID,
		Amount:           req.Amount,
		Currency:         req.Currency,
		PaymentMethod:    req.PaymentMethod,
		Status:           "pending",
		Metadata:         response.Metadata,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return response, nil
}

// ConfirmPayment confirms a payment
func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID string) error {
	// Get payment from database
	payment, err := s.paymentRepo.FindByGatewayPaymentID(paymentID)
	if err != nil {
		return err
	}

	// Get the gateway
	gateway, err := s.factory.Get(payment.GatewayProvider)
	if err != nil {
		return err
	}

	// Confirm via gateway
	response, err := gateway.ConfirmPayment(ctx, paymentID)
	if err != nil {
		payment.Status = "failed"
		payment.FailureReason = err.Error()
		payment.UpdatedAt = time.Now()
		s.paymentRepo.Update(payment)
		return err
	}

	// Update database
	payment.Status = response.Status
	payment.GatewayChargeID = response.ChargeID
	payment.UpdatedAt = time.Now()
	return s.paymentRepo.Update(payment)
}

// RefundPayment processes a refund
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID string, amount int64, reason string) error {
	// Get payment from database
	payment, err := s.paymentRepo.FindByGatewayPaymentID(paymentID)
	if err != nil {
		return err
	}

	// Get the gateway
	gateway, err := s.factory.Get(payment.GatewayProvider)
	if err != nil {
		return err
	}

	// Process refund via gateway
	refundReq := &RefundRequest{
		PaymentID: paymentID,
		Amount:    amount,
		Reason:    reason,
	}

	refundResp, err := gateway.RefundPayment(ctx, refundReq)
	if err != nil {
		return err
	}

	// Update database
	payment.RefundAmount = refundResp.Amount
	payment.RefundedAt = timePtr(time.Now())
	if refundResp.Amount >= payment.Amount {
		payment.Status = "refunded"
	}
	payment.UpdatedAt = time.Now()

	return s.paymentRepo.Update(payment)
}

// ============================================
// Repository Interface
// ============================================

type PaymentRepository interface {
	Create(payment *BookingPayment) error
	Update(payment *BookingPayment) error
	FindByID(id uuid.UUID) (*BookingPayment, error)
	FindByGatewayPaymentID(paymentID string) (*BookingPayment, error)
	FindByBookingID(bookingID uuid.UUID) ([]*BookingPayment, error)
	FindByUserID(userID uuid.UUID) ([]*BookingPayment, error)
}

// ============================================
// Model (matches your schema)
// ============================================

type BookingPayment struct {
	ID                uuid.UUID
	BookingID         uuid.UUID
	UserID            uuid.UUID
	GatewayProvider   string // stripe, juspay, razorpay
	GatewayPaymentID  string
	GatewayChargeID   string
	GatewayCustomerID string
	GatewayOrderID    string
	Amount            int64
	Currency          string
	PaymentMethod     string
	Status            string
	FailureReason     string
	RefundAmount      int64
	RefundedAt        *time.Time
	Metadata          map[string]string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
