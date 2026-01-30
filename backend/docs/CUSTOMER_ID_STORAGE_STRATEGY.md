# Gateway Customer ID Storage Strategy

## The Question

Where should we store gateway customer IDs (e.g., `stripe_customer_id`)?

**Option A**: In `users` table  
**Option B**: In `payment_methods` table  
**Our Choice**: **Hybrid approach** (both, with precedence rules)

---

## Final Design (Hybrid)

### Users Table (Default Customer IDs)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255),
    
    -- Default customer IDs (one per gateway per user)
    stripe_customer_id VARCHAR(255) UNIQUE,
    razorpay_customer_id VARCHAR(255) UNIQUE,
    juspay_customer_id VARCHAR(255) UNIQUE,
    ...
);
```

**Use case**: 99% of users (personal accounts)

### Payment Methods Table (Override Customer IDs)
```sql
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    gateway_provider VARCHAR(50),
    gateway_payment_method_id VARCHAR(255),
    
    -- Optional: Override customer ID for this specific payment method
    gateway_customer_id VARCHAR(255),  -- Usually NULL
    ...
);
```

**Use case**: 1% of users (business accounts, multiple billing entities)

---

## Customer ID Resolution Logic

```go
func (s *PaymentService) GetCustomerID(user *User, paymentMethod *PaymentMethod) string {
    // 1. Check if payment method has override customer ID
    if paymentMethod != nil && paymentMethod.GatewayCustomerID != "" {
        return paymentMethod.GatewayCustomerID  // Business account
    }
    
    // 2. Fall back to user's default customer ID
    switch paymentMethod.GatewayProvider {
    case "stripe":
        return user.StripeCustomerID
    case "razorpay":
        return user.RazorpayCustomerID
    case "juspay":
        return user.JuspayCustomerID
    default:
        return ""
    }
}
```

---

## Use Cases

### Case 1: Personal User (99% of users)

**Alice** has a personal account:

```
users table:
├── id = alice-uuid
├── email = alice@example.com
└── stripe_customer_id = "cus_alice_personal"

payment_methods table:
├── Visa 4242 → gateway_customer_id = NULL  (uses cus_alice_personal)
├── Mastercard 5555 → gateway_customer_id = NULL  (uses cus_alice_personal)
└── Amex 3782 → gateway_customer_id = NULL  (uses cus_alice_personal)
```

**All cards use the same customer ID from users table.**

---

### Case 2: Business User (1% of users)

**Bob** works at Acme Corp and has both personal and company cards:

```
users table:
├── id = bob-uuid
├── email = bob@acme.com
└── stripe_customer_id = "cus_bob_personal"  (default)

payment_methods table:
├── Personal Visa → gateway_customer_id = NULL  (uses cus_bob_personal)
├── Company Amex → gateway_customer_id = "cus_acme_corp"  (override!)
└── Company Visa → gateway_customer_id = "cus_acme_corp"  (override!)
```

**Personal cards use user's customer ID, company cards use override.**

---

### Case 3: Multi-Gateway User

**Carol** uses different gateways in different regions:

```
users table:
├── id = carol-uuid
├── stripe_customer_id = "cus_carol_stripe"  (for international)
└── razorpay_customer_id = "cust_carol_razorpay"  (for India)

payment_methods table:
├── Visa (Stripe) → gateway_customer_id = NULL  (uses cus_carol_stripe)
├── UPI (Razorpay) → gateway_customer_id = NULL  (uses cust_carol_razorpay)
└── Wallet (Juspay) → gateway_customer_id = "cust_carol_juspay_business"  (override)
```

---

## Relationship Cardinality

### Users ↔ Customer IDs

**Option A (users table):**
```
1 user : 1 stripe_customer_id
1 user : 1 razorpay_customer_id
1 user : 1 juspay_customer_id
```

**Option B (payment_methods table):**
```
1 user : N payment_methods : N customer_ids
```

**Our Hybrid:**
```
1 user : 1 default_customer_id (per gateway)
1 user : N payment_methods : 0-N override_customer_ids
```

---

## Why Hybrid is Best

### ✅ Advantages

1. **Simple for 99% of users**
   - Default customer ID in users table
   - No NULL handling in payment_methods
   - Fast queries

2. **Flexible for 1% of users**
   - Business accounts can override
   - Multiple billing entities supported
   - No schema changes needed

3. **Industry standard**
   - Follows Stripe's customer model
   - Matches Shopify, WooCommerce patterns
   - Familiar to developers

4. **Performance**
   - Most queries hit users table (indexed)
   - Override is optional (sparse column)
   - No joins for common case

---

## Implementation Examples

### Creating a Payment Method (Personal User)

```go
func (s *PaymentService) SavePaymentMethod(userID uuid.UUID, provider string, paymentMethodID string) error {
    user, _ := s.userRepo.FindByID(userID)
    
    // Get or create customer ID in users table
    var customerID string
    switch provider {
    case "stripe":
        if user.StripeCustomerID == "" {
            // Create Stripe customer
            customer, _ := stripe.Customer.New(&stripe.CustomerParams{
                Email: user.Email,
            })
            user.StripeCustomerID = customer.ID
            s.userRepo.Update(user)
        }
        customerID = user.StripeCustomerID
    }
    
    // Attach payment method to customer
    stripe.PaymentMethod.Attach(paymentMethodID, &stripe.PaymentMethodAttachParams{
        Customer: customerID,
    })
    
    // Save payment method (no customer ID override)
    paymentMethod := &PaymentMethod{
        UserID:                 userID,
        GatewayProvider:        provider,
        GatewayPaymentMethodID: paymentMethodID,
        GatewayCustomerID:      nil,  // Use default from users table
        CardLast4:              "4242",
        CardBrand:              "visa",
    }
    
    return s.paymentMethodRepo.Create(paymentMethod)
}
```

### Creating a Payment Method (Business User)

```go
func (s *PaymentService) SaveBusinessPaymentMethod(
    userID uuid.UUID,
    provider string,
    paymentMethodID string,
    companyCustomerID string,  // e.g., "cus_acme_corp"
) error {
    // Attach to company's customer ID (not user's)
    stripe.PaymentMethod.Attach(paymentMethodID, &stripe.PaymentMethodAttachParams{
        Customer: companyCustomerID,
    })
    
    // Save with customer ID override
    paymentMethod := &PaymentMethod{
        UserID:                 userID,
        GatewayProvider:        provider,
        GatewayPaymentMethodID: paymentMethodID,
        GatewayCustomerID:      &companyCustomerID,  // Override!
        CardLast4:              "5555",
        CardBrand:              "amex",
        Nickname:               "Company Card",
    }
    
    return s.paymentMethodRepo.Create(paymentMethod)
}
```

### Using a Payment Method

```go
func (s *PaymentService) CreatePayment(paymentMethodID uuid.UUID, amount int64) error {
    // Get payment method
    pm, _ := s.paymentMethodRepo.FindByID(paymentMethodID)
    
    // Get user
    user, _ := s.userRepo.FindByID(pm.UserID)
    
    // Resolve customer ID (with precedence)
    customerID := s.GetCustomerID(user, pm)
    
    // Create payment intent
    intent, _ := stripe.PaymentIntent.New(&stripe.PaymentIntentParams{
        Amount:        amount,
        Currency:      "aud",
        Customer:      customerID,  // Resolved customer ID
        PaymentMethod: pm.GatewayPaymentMethodID,
    })
    
    return nil
}
```

---

## Database Queries

### Get user's default customer ID
```sql
SELECT stripe_customer_id 
FROM users 
WHERE id = 'user-123';
```

### Get customer ID for a specific payment method
```sql
SELECT 
    COALESCE(pm.gateway_customer_id, u.stripe_customer_id) as customer_id
FROM payment_methods pm
JOIN users u ON pm.user_id = u.id
WHERE pm.id = 'pm-123';
```

### Get all payment methods with resolved customer IDs
```sql
SELECT 
    pm.id,
    pm.nickname,
    pm.card_last4,
    CASE pm.gateway_provider
        WHEN 'stripe' THEN COALESCE(pm.gateway_customer_id, u.stripe_customer_id)
        WHEN 'razorpay' THEN COALESCE(pm.gateway_customer_id, u.razorpay_customer_id)
        WHEN 'juspay' THEN COALESCE(pm.gateway_customer_id, u.juspay_customer_id)
    END as resolved_customer_id
FROM payment_methods pm
JOIN users u ON pm.user_id = u.id
WHERE pm.user_id = 'user-123';
```

---

## Migration Path

### Phase 1: Start Simple (Current)
```sql
-- Most users: customer ID in users table
users.stripe_customer_id = "cus_alice"

-- Payment methods: no override
payment_methods.gateway_customer_id = NULL
```

### Phase 2: Add Business Accounts (Future)
```sql
-- Business users: override for company cards
payment_methods.gateway_customer_id = "cus_acme_corp"
```

**No schema changes needed!** Just start populating the optional field.

---

## Summary

### Decision: Hybrid Approach ✅

**Default**: Customer IDs in `users` table (99% of cases)  
**Override**: Customer IDs in `payment_methods` table (1% of cases)

### Benefits:
- ✅ Simple for personal users
- ✅ Flexible for business users
- ✅ Industry standard pattern
- ✅ No schema changes for future features
- ✅ Optimal performance

### Relationship:
```
1 user : 1 default_customer_id (per gateway)
1 user : N payment_methods
1 payment_method : 0-1 override_customer_id
```

This gives you the **best of both worlds**! 🎯
