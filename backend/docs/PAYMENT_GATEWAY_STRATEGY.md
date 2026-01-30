# Multi-Gateway Payment Integration Strategy

## Overview

Wayne Stays supports multiple payment gateways through an **Adapter Pattern** abstraction layer. This allows seamless switching between providers (Stripe, Juspay, Razorpay) without changing business logic.

---

## Database Schema (DRY Approach)

### Single Payment Table Structure

```sql
-- Works for ALL payment gateways
CREATE TABLE booking_payments (
    id UUID PRIMARY KEY,
    booking_id UUID NOT NULL,
    user_id UUID NOT NULL,
    
    -- Gateway identification
    gateway_provider VARCHAR(50) NOT NULL DEFAULT 'stripe',
    
    -- Generic gateway fields (work for all providers)
    gateway_payment_id VARCHAR(255) UNIQUE,  -- Main payment ID
    gateway_charge_id VARCHAR(255),          -- Charge/capture ID
    gateway_customer_id VARCHAR(255),        -- Customer ID
    gateway_order_id VARCHAR(255),           -- Order ID (some providers)
    
    -- Payment details
    amount INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'AUD',
    payment_method VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    
    -- Provider-specific data
    metadata JSONB,  -- Flexible storage for provider-specific fields
    
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Field Mapping by Provider

### Stripe
```
gateway_provider     = "stripe"
gateway_payment_id   = payment_intent_id  (e.g., "pi_abc123")
gateway_charge_id    = charge_id          (e.g., "ch_xyz789")
gateway_customer_id  = customer_id        (e.g., "cus_def456")
gateway_order_id     = NULL               (not used)
metadata             = {
    "payment_method_id": "pm_card_123",
    "receipt_url": "https://stripe.com/receipt/...",
    ...
}
```

### Razorpay
```
gateway_provider     = "razorpay"
gateway_payment_id   = payment_id         (e.g., "pay_abc123")
gateway_charge_id    = NULL               (not separate from payment_id)
gateway_customer_id  = customer_id        (e.g., "cust_xyz789")
gateway_order_id     = order_id           (e.g., "order_def456")
metadata             = {
    "method": "card",
    "card_id": "card_123",
    "bank": "HDFC",
    ...
}
```

### Juspay
```
gateway_provider     = "juspay"
gateway_payment_id   = payment_id         (e.g., "jpy_abc123")
gateway_charge_id    = transaction_id     (e.g., "txn_xyz789")
gateway_customer_id  = customer_id        (e.g., "cust_def456")
gateway_order_id     = order_id           (e.g., "ord_ghi789")
metadata             = {
    "payment_method_type": "CARD",
    "card_network": "VISA",
    ...
}
```

---

## Architecture

### Layer 1: Payment Service (Business Logic)
```go
// Handles business logic, validation, database operations
type PaymentService struct {
    factory     *GatewayFactory
    paymentRepo PaymentRepository
}

func (s *PaymentService) CreateBookingPayment(
    ctx context.Context,
    bookingID uuid.UUID,
    amount int64,
    provider string,  // "stripe", "razorpay", "juspay"
) (*Payment, error) {
    // Business logic here
    // Gateway-agnostic
}
```

### Layer 2: Gateway Interface (Abstraction)
```go
// All gateways implement this interface
type PaymentGateway interface {
    GetProvider() string
    CreatePaymentIntent(ctx context.Context, req *PaymentIntentRequest) (*PaymentIntentResponse, error)
    ConfirmPayment(ctx context.Context, paymentID string) (*PaymentResponse, error)
    RefundPayment(ctx context.Context, req *RefundRequest) (*RefundResponse, error)
    // ... other methods
}
```

### Layer 3: Gateway Adapters (Provider-Specific)
```go
// Stripe adapter
type StripeGateway struct { ... }
func (g *StripeGateway) CreatePaymentIntent(...) { /* Stripe SDK calls */ }

// Razorpay adapter
type RazorpayGateway struct { ... }
func (g *RazorpayGateway) CreatePaymentIntent(...) { /* Razorpay SDK calls */ }

// Juspay adapter
type JuspayGateway struct { ... }
func (g *JuspayGateway) CreatePaymentIntent(...) { /* Juspay SDK calls */ }
```

---

## Benefits of This Approach

### ✅ DRY (Don't Repeat Yourself)
- Single payment table structure
- Single payment service
- Single set of business logic

### ✅ Extensible
- Add new gateway = implement interface + register
- No database changes needed
- No business logic changes needed

### ✅ Flexible
- Switch providers per transaction
- Support multiple providers simultaneously
- A/B test different providers

### ✅ Maintainable
- Provider-specific code isolated in adapters
- Easy to test (mock the interface)
- Clear separation of concerns

---

## Implementation Example

### 1. Register Gateways (main.go)
```go
func main() {
    // Create gateway factory
    factory := payment.NewGatewayFactory()
    
    // Register all available gateways
    factory.Register(payment.NewStripeGateway(cfg.StripeAPIKey))
    factory.Register(payment.NewRazorpayGateway(cfg.RazorpayKey, cfg.RazorpaySecret))
    factory.Register(payment.NewJuspayGateway(cfg.JuspayAPIKey))
    
    // Create payment service
    paymentService := payment.NewPaymentService(factory, paymentRepo, "stripe")
    
    // Use in handlers
    paymentHandler := payment.NewHandler(paymentService)
}
```

### 2. Create Payment (Service Layer)
```go
func (s *PaymentService) CreatePayment(
    ctx context.Context,
    provider string,
    req *PaymentIntentRequest,
) (*PaymentIntentResponse, error) {
    // Get the appropriate gateway
    gateway, err := s.factory.Get(provider)
    if err != nil {
        return nil, err
    }
    
    // Create payment via gateway (provider-agnostic call)
    response, err := gateway.CreatePaymentIntent(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Save to database (same structure for all providers)
    payment := &BookingPayment{
        GatewayProvider:  provider,
        GatewayPaymentID: response.PaymentID,
        GatewayOrderID:   response.OrderID,
        Amount:           req.Amount,
        Currency:         req.Currency,
        Status:           "pending",
        Metadata:         response.Metadata,  // Provider-specific data
    }
    
    if err := s.paymentRepo.Create(payment); err != nil {
        return nil, err
    }
    
    return response, nil
}
```

### 3. Frontend (Provider Selection)
```typescript
// User selects payment method
const handlePayment = async (provider: 'stripe' | 'razorpay' | 'juspay') => {
    // Create payment intent
    const response = await fetch('/api/payments/create', {
        method: 'POST',
        body: JSON.stringify({
            provider: provider,  // User's choice
            amount: 10000,
            currency: 'AUD',
        })
    });
    
    const { payment_id, client_secret } = await response.json();
    
    // Provider-specific frontend handling
    if (provider === 'stripe') {
        // Use Stripe.js
        await stripe.confirmCardPayment(client_secret);
    } else if (provider === 'razorpay') {
        // Use Razorpay checkout
        const rzp = new Razorpay({ key: 'rzp_key', order_id: payment_id });
        rzp.open();
    } else if (provider === 'juspay') {
        // Use Juspay SDK
        Juspay.startPayment({ payment_id });
    }
};
```

---

## Query Examples

### Get all Stripe payments
```sql
SELECT * FROM booking_payments 
WHERE gateway_provider = 'stripe';
```

### Get all payments for a user (any provider)
```sql
SELECT * FROM booking_payments 
WHERE user_id = 'user-123'
ORDER BY created_at DESC;
```

### Provider-specific reporting
```sql
SELECT 
    gateway_provider,
    COUNT(*) as total_payments,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount
FROM booking_payments
WHERE status = 'succeeded'
GROUP BY gateway_provider;
```

---

## Migration Strategy

### Phase 1: Start with Stripe (Current)
```go
factory.Register(payment.NewStripeGateway(cfg.StripeAPIKey))
defaultGateway := "stripe"
```

### Phase 2: Add Razorpay (India market)
```go
factory.Register(payment.NewStripeGateway(cfg.StripeAPIKey))
factory.Register(payment.NewRazorpayGateway(cfg.RazorpayKey, cfg.RazorpaySecret))

// Route based on currency
if currency == "INR" {
    provider = "razorpay"
} else {
    provider = "stripe"
}
```

### Phase 3: Add Juspay (Advanced routing)
```go
// All gateways registered
factory.Register(payment.NewStripeGateway(...))
factory.Register(payment.NewRazorpayGateway(...))
factory.Register(payment.NewJuspayGateway(...))

// Smart routing based on:
// - User location
// - Currency
// - Payment method
// - Success rates
provider := s.routingService.SelectBestGateway(req)
```

---

## Testing Strategy

### Mock the Interface
```go
type MockGateway struct {
    mock.Mock
}

func (m *MockGateway) CreatePaymentIntent(ctx context.Context, req *PaymentIntentRequest) (*PaymentIntentResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*PaymentIntentResponse), args.Error(1)
}

// Test
func TestPaymentService_CreatePayment(t *testing.T) {
    mockGateway := new(MockGateway)
    mockGateway.On("CreatePaymentIntent", mock.Anything, mock.Anything).
        Return(&PaymentIntentResponse{PaymentID: "test_123"}, nil)
    
    factory := NewGatewayFactory()
    factory.Register(mockGateway)
    
    service := NewPaymentService(factory, mockRepo, "mock")
    
    // Test business logic without hitting real payment gateways
    result, err := service.CreatePayment(ctx, "mock", req)
    assert.NoError(t, err)
    assert.Equal(t, "test_123", result.PaymentID)
}
```

---

## Summary

**Database**: Single DRY schema with `gateway_provider` + generic `gateway_*` fields  
**Code**: Adapter pattern with `PaymentGateway` interface  
**Flexibility**: Add new providers without changing existing code  
**Maintainability**: Provider-specific logic isolated in adapters  

This approach gives you the best of both worlds: **DRY database schema** + **extensible architecture**! 🚀
