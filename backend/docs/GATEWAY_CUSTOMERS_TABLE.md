# Gateway Customers Table - The Right Approach

## The Problem

**Bad Approach** (adding columns for each gateway):
```sql
CREATE TABLE users (
    id UUID,
    stripe_customer_id VARCHAR(255),
    razorpay_customer_id VARCHAR(255),
    juspay_customer_id VARCHAR(255),
    -- Add new column for every new gateway! ❌
);
```

**Problems:**
- ❌ Schema change for every new gateway
- ❌ Migration required for every new gateway
- ❌ Sparse columns (most users don't use all gateways)
- ❌ Not scalable

---

## The Solution: `gateway_customers` Table ✅

```sql
CREATE TABLE gateway_customers (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    gateway_provider VARCHAR(50) NOT NULL,  -- 'stripe', 'razorpay', 'juspay'
    gateway_customer_id VARCHAR(255) NOT NULL,
    
    UNIQUE(user_id, gateway_provider)  -- One customer ID per user per gateway
);
```

**Benefits:**
- ✅ **No schema changes** when adding gateways
- ✅ **Just insert rows** for new gateways
- ✅ **Normalized** (no sparse columns)
- ✅ **Scalable** (supports unlimited gateways)

---

## How It Works

### Data Structure

```
users table:
├── id = alice-uuid
├── email = alice@example.com
└── (no gateway-specific columns)

gateway_customers table:
├── user_id = alice-uuid, gateway_provider = 'stripe', gateway_customer_id = 'cus_alice_stripe'
├── user_id = alice-uuid, gateway_provider = 'razorpay', gateway_customer_id = 'cust_alice_razorpay'
└── user_id = bob-uuid, gateway_provider = 'stripe', gateway_customer_id = 'cus_bob_stripe'
```

### Relationship

```
1 user : N gateway_customers (one per gateway they use)
```

---

## Adding a New Gateway

### Old Approach (Column-based) ❌
```sql
-- Step 1: Add column
ALTER TABLE users ADD COLUMN paypal_customer_id VARCHAR(255);

-- Step 2: Add index
CREATE INDEX idx_users_paypal_customer_id ON users(paypal_customer_id);

-- Step 3: Update application code
-- Step 4: Deploy migration
```

**Downtime**: Yes (schema change)  
**Risk**: High (migration can fail)

### New Approach (Row-based) ✅
```sql
-- Just insert rows!
INSERT INTO gateway_customers (user_id, gateway_provider, gateway_customer_id)
VALUES ('alice-uuid', 'paypal', 'paypal_alice_123');
```

**Downtime**: No  
**Risk**: None (no schema change)  
**Code change**: Update enum/constant only

---

## Code Examples

### Get Customer ID for a Gateway

```go
func (r *GatewayCustomerRepo) GetCustomerID(userID uuid.UUID, provider string) (string, error) {
    var gc GatewayCustomer
    err := r.db.Where("user_id = ? AND gateway_provider = ?", userID, provider).
        First(&gc).Error
    
    if err != nil {
        return "", err  // Customer ID doesn't exist yet
    }
    
    return gc.GatewayCustomerID, nil
}
```

### Create Customer ID

```go
func (s *PaymentService) GetOrCreateCustomerID(userID uuid.UUID, provider string) (string, error) {
    // Try to get existing
    customerID, err := s.gatewayCustomerRepo.GetCustomerID(userID, provider)
    if err == nil {
        return customerID, nil  // Already exists
    }
    
    // Create new customer in gateway
    user, _ := s.userRepo.FindByID(userID)
    
    var gatewayCustomerID string
    switch provider {
    case "stripe":
        customer, _ := stripe.Customer.New(&stripe.CustomerParams{
            Email: user.Email,
        })
        gatewayCustomerID = customer.ID
        
    case "razorpay":
        customer, _ := razorpay.Customer.Create(map[string]interface{}{
            "email": user.Email,
        })
        gatewayCustomerID = customer.ID
        
    case "juspay":
        // Juspay customer creation
        gatewayCustomerID = "juspay_customer_id"
    }
    
    // Save to database
    gc := &GatewayCustomer{
        UserID:            userID,
        GatewayProvider:   provider,
        GatewayCustomerID: gatewayCustomerID,
    }
    s.gatewayCustomerRepo.Create(gc)
    
    return gatewayCustomerID, nil
}
```

### Use in Payment Flow

```go
func (s *PaymentService) CreatePayment(userID uuid.UUID, amount int64, provider string) error {
    // Get or create customer ID for this gateway
    customerID, _ := s.GetOrCreateCustomerID(userID, provider)
    
    // Create payment intent
    gateway, _ := s.gatewayFactory.Get(provider)
    intent, _ := gateway.CreatePaymentIntent(&PaymentIntentRequest{
        Amount:     amount,
        Currency:   "AUD",
        CustomerID: customerID,  // ✅ Works for any gateway
    })
    
    return nil
}
```

---

## Database Queries

### Get all gateways a user has customer IDs for
```sql
SELECT gateway_provider, gateway_customer_id
FROM gateway_customers
WHERE user_id = 'alice-uuid';
```

Result:
```
gateway_provider | gateway_customer_id
-----------------|--------------------
stripe           | cus_alice_stripe
razorpay         | cust_alice_razorpay
```

### Get Stripe customer ID for a user
```sql
SELECT gateway_customer_id
FROM gateway_customers
WHERE user_id = 'alice-uuid' AND gateway_provider = 'stripe';
```

### Get all users with Stripe customer IDs
```sql
SELECT DISTINCT user_id
FROM gateway_customers
WHERE gateway_provider = 'stripe';
```

### Join with users
```sql
SELECT u.email, gc.gateway_provider, gc.gateway_customer_id
FROM users u
JOIN gateway_customers gc ON u.id = gc.user_id
WHERE u.id = 'alice-uuid';
```

---

## GORM Model

```go
type GatewayCustomer struct {
    ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID            uuid.UUID `gorm:"type:uuid;not null"`
    GatewayProvider   string    `gorm:"type:varchar(50);not null"`
    GatewayCustomerID string    `gorm:"type:varchar(255);not null"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    
    // Relationships
    User User `gorm:"foreignKey:UserID"`
}

// Unique constraint
func (GatewayCustomer) TableName() string {
    return "gateway_customers"
}

// In migration
db.Exec(`
    CREATE UNIQUE INDEX idx_user_gateway_unique 
    ON gateway_customers(user_id, gateway_provider)
`)
```

---

## Comparison

| Aspect | Column-based | Row-based (gateway_customers) |
|--------|--------------|-------------------------------|
| **Add new gateway** | ALTER TABLE | INSERT rows |
| **Schema changes** | Yes (migration) | No |
| **Downtime** | Possible | No |
| **Sparse data** | Yes (NULL columns) | No |
| **Scalability** | Limited | Unlimited |
| **Query complexity** | Simple | Simple (indexed) |
| **Industry standard** | No | Yes |

---

## Industry Examples

### Companies Using Row-based Approach:
- **Stripe** - Separate `customer_payment_methods` table
- **Shopify** - Separate `payment_gateway_customers` table
- **AWS** - Separate tables for different service integrations
- **Multi-tenant SaaS** - Standard pattern for extensibility

### Companies Using Column-based:
- Small apps with 1-2 payment gateways (not scalable)

---

## Migration Path

### Phase 1: Create Table
```sql
CREATE TABLE gateway_customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    gateway_provider VARCHAR(50) NOT NULL,
    gateway_customer_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, gateway_provider)
);

CREATE INDEX idx_gateway_customers_user_id ON gateway_customers(user_id);
CREATE INDEX idx_gateway_customers_provider ON gateway_customers(gateway_provider);
CREATE INDEX idx_gateway_customers_customer_id ON gateway_customers(gateway_customer_id);
```

### Phase 2: Start with Stripe
```go
// When user saves first payment method
customerID, _ := s.GetOrCreateCustomerID(userID, "stripe")
```

### Phase 3: Add Razorpay (No Schema Change!)
```go
// Just use different provider
customerID, _ := s.GetOrCreateCustomerID(userID, "razorpay")
```

### Phase 4: Add Juspay (No Schema Change!)
```go
// Just use different provider
customerID, _ := s.GetOrCreateCustomerID(userID, "juspay")
```

---

## Summary

### Your Instinct Was Correct ✅

**Problem**: Adding columns for each gateway doesn't scale  
**Solution**: Separate `gateway_customers` table

**Benefits:**
- ✅ No schema changes when adding gateways
- ✅ Normalized (no sparse columns)
- ✅ Scalable (unlimited gateways)
- ✅ Industry standard pattern
- ✅ Zero downtime for new gateways

**Relationship:**
```
1 user : N gateway_customers (one per gateway)
```

**To add new gateway:**
```
1. Update gateway_provider enum/constant
2. Insert rows (no migration needed)
```

This is the **correct, scalable approach**! 🎯
