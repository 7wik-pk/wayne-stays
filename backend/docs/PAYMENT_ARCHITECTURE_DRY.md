# DRY Payment Architecture - Wayne Stays

> **Status: partly superseded.** Two decisions here changed: gateway customer IDs moved from `users` columns to the `gateway_customers` table ([GATEWAY_CUSTOMERS_TABLE.md](GATEWAY_CUSTOMERS_TABLE.md)), and payments stay in separate `booking_payments` / `order_payments` tables ([WHY_SEPARATE_PAYMENT_TABLES.md](WHY_SEPARATE_PAYMENT_TABLES.md)). `schema.dbml` is the source of truth.


## Overview

The payment system has been refactored to be maximally DRY (Don't Repeat Yourself) while supporting multiple payment gateways.

---

## Key Changes

### ❌ Before (Repetitive)
```
booking_payments table
├── gateway_provider
├── gateway_payment_id
├── gateway_charge_id
├── gateway_customer_id  ← Repeated!
└── gateway_order_id

order_payments table
├── gateway_provider
├── gateway_payment_id
├── gateway_charge_id
├── gateway_customer_id  ← Repeated!
└── gateway_order_id
```

### ✅ After (DRY)
```
users table
├── stripe_customer_id    ← Stored once per user
├── razorpay_customer_id  ← Stored once per user
└── juspay_customer_id    ← Stored once per user

payment_methods table (NEW)
├── gateway_provider
├── gateway_payment_method_id
├── card_last4, card_brand, etc.
└── nickname ("My Visa", "Work Card")

payments table (UNIFIED)
├── payable_type (booking/order)  ← Polymorphic
├── payable_id
├── gateway_provider
├── gateway_payment_id
├── payment_method_id  ← Links to saved card
└── NO gateway_customer_id (it's in users table!)
```

---

## Database Schema

### 1. Users Table (Gateway Customer IDs)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    
    -- Gateway customer IDs (for saved payment methods)
    stripe_customer_id VARCHAR(255) UNIQUE,
    razorpay_customer_id VARCHAR(255) UNIQUE,
    juspay_customer_id VARCHAR(255) UNIQUE,
    
    ...
);
```

**Why here?**
- ✅ One customer ID per user per gateway
- ✅ Used for saved payment methods
- ✅ No duplication across payments

### 2. Payment Methods Table (Saved Cards)
```sql
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    
    -- Gateway info
    gateway_provider VARCHAR(50) NOT NULL,  -- stripe, razorpay, juspay
    gateway_payment_method_id VARCHAR(255) NOT NULL,
    
    -- Card metadata (safe to store)
    type VARCHAR(20) NOT NULL,  -- card, bank_account, upi, wallet
    card_brand VARCHAR(20),     -- visa, mastercard, amex
    card_last4 VARCHAR(4),      -- Last 4 digits
    card_exp_month INTEGER,
    card_exp_year INTEGER,
    
    -- User preferences
    nickname VARCHAR(100),      -- "My Visa", "Work Card"
    is_default BOOLEAN DEFAULT false,
    
    ...
);
```

**Why separate table?**
- ✅ Users can have multiple saved cards
- ✅ Each card has metadata (last4, brand, expiry)
- ✅ User can nickname cards
- ✅ Mark default payment method

### 3. Payments Table (Unified)
```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    
    -- Polymorphic relationship
    payable_type VARCHAR(20) NOT NULL,  -- 'booking' or 'order'
    payable_id UUID NOT NULL,
    
    -- Gateway
    gateway_provider VARCHAR(50) NOT NULL,
    
    -- Payment details
    amount INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'AUD',
    payment_method_id UUID REFERENCES payment_methods(id),  -- If using saved card
    payment_method_type VARCHAR(50),  -- card, upi, wallet, etc.
    
    -- Gateway transaction IDs
    gateway_payment_id VARCHAR(255) UNIQUE,
    gateway_charge_id VARCHAR(255),
    gateway_order_id VARCHAR(255),
    
    -- NO gateway_customer_id (it's in users table!)
    
    status VARCHAR(20) DEFAULT 'pending',
    ...
);
```

**Why polymorphic?**
- ✅ Single table for all payment types
- ✅ No duplication of payment logic
- ✅ Easy to add new payable types (subscriptions, donations, etc.)

---

## How It Works

### Scenario 1: User Saves a Card (First Time)

```go
// 1. Create Stripe customer (if not exists)
if user.StripeCustomerID == "" {
    customer, _ := stripe.Customer.New(&stripe.CustomerParams{
        Email: user.Email,
    })
    user.StripeCustomerID = customer.ID
    userRepo.Update(user)
}

// 2. Attach payment method to customer
stripe.PaymentMethod.Attach(paymentMethodID, &stripe.PaymentMethodAttachParams{
    Customer: user.StripeCustomerID,
})

// 3. Save payment method metadata in our DB
paymentMethod := &PaymentMethod{
    UserID:                 user.ID,
    GatewayProvider:        "stripe",
    GatewayPaymentMethodID: paymentMethodID,
    Type:                   "card",
    CardBrand:              "visa",
    CardLast4:              "4242",
    CardExpMonth:           12,
    CardExpYear:            2026,
    Nickname:               "My Visa",
    IsDefault:              true,
}
paymentMethodRepo.Create(paymentMethod)
```

### Scenario 2: User Makes a Payment (Using Saved Card)

```go
// 1. Get user's saved payment method
paymentMethod, _ := paymentMethodRepo.FindDefaultByUserID(user.ID)

// 2. Get user's gateway customer ID
var customerID string
switch paymentMethod.GatewayProvider {
case "stripe":
    customerID = user.StripeCustomerID
case "razorpay":
    customerID = user.RazorpayCustomerID
case "juspay":
    customerID = user.JuspayCustomerID
}

// 3. Create payment intent with saved card
intent, _ := stripe.PaymentIntent.New(&stripe.PaymentIntentParams{
    Amount:        10000,
    Currency:      "aud",
    Customer:      customerID,
    PaymentMethod: paymentMethod.GatewayPaymentMethodID,
})

// 4. Save payment record
payment := &Payment{
    UserID:            user.ID,
    PayableType:       "booking",
    PayableID:         booking.ID,
    GatewayProvider:   "stripe",
    GatewayPaymentID:  intent.ID,
    Amount:            10000,
    Currency:          "AUD",
    PaymentMethodID:   &paymentMethod.ID,  // Link to saved card
    PaymentMethodType: "card",
    Status:            "pending",
}
paymentRepo.Create(payment)
```

### Scenario 3: User Makes a Payment (New Card, Don't Save)

```go
// 1. Create payment intent (no customer ID)
intent, _ := stripe.PaymentIntent.New(&stripe.PaymentIntentParams{
    Amount:   10000,
    Currency: "aud",
    // No customer or payment method - user will enter card on frontend
})

// 2. Save payment record
payment := &Payment{
    UserID:            user.ID,
    PayableType:       "order",
    PayableID:         order.ID,
    GatewayProvider:   "stripe",
    GatewayPaymentID:  intent.ID,
    Amount:            10000,
    Currency:          "AUD",
    PaymentMethodID:   nil,  // No saved card used
    PaymentMethodType: "card",
    Status:            "pending",
}
paymentRepo.Create(payment)
```

---

## Query Examples

### Get all payments for a booking
```sql
SELECT * FROM payments 
WHERE payable_type = 'booking' AND payable_id = 'booking-123';
```

### Get all payments for a user
```sql
SELECT * FROM payments 
WHERE user_id = 'user-123'
ORDER BY created_at DESC;
```

### Get user's saved cards
```sql
SELECT * FROM payment_methods 
WHERE user_id = 'user-123' AND deleted_at IS NULL
ORDER BY is_default DESC, created_at DESC;
```

### Get payment with saved card details
```sql
SELECT 
    p.*,
    pm.card_brand,
    pm.card_last4,
    pm.nickname
FROM payments p
LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
WHERE p.id = 'payment-123';
```

### User payment history with details
```sql
SELECT 
    p.id,
    p.amount,
    p.currency,
    p.status,
    p.payable_type,
    CASE 
        WHEN p.payable_type = 'booking' THEN b.booking_reference
        WHEN p.payable_type = 'order' THEN o.order_number
    END as reference,
    pm.card_brand,
    pm.card_last4,
    p.created_at
FROM payments p
LEFT JOIN bookings b ON p.payable_type = 'booking' AND p.payable_id = b.id
LEFT JOIN orders o ON p.payable_type = 'order' AND p.payable_id = o.id
LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
WHERE p.user_id = 'user-123'
ORDER BY p.created_at DESC;
```

---

## Benefits of This Architecture

### ✅ Maximum DRY
- **One** payments table (not two)
- **One** place for customer IDs (users table)
- **One** place for card metadata (payment_methods table)

### ✅ Saved Payment Methods
- Users can save multiple cards
- Each card has a nickname ("My Visa", "Work Card")
- Mark default payment method
- Secure (only metadata stored, not actual card numbers)

### ✅ Multi-Gateway Support
- Each user can have customer IDs for multiple gateways
- Each saved card linked to its gateway
- Easy to switch gateways per transaction

### ✅ Extensible
- Add new payable types (subscriptions, donations) without new tables
- Add new gateways by adding customer ID column to users table
- Add new payment method types (UPI, wallets) easily

### ✅ User-Friendly
- Users see friendly card names ("My Visa" instead of "pm_abc123")
- Last 4 digits displayed for identification
- Expiry dates shown for card management

---

## Frontend Implementation

### Display Saved Cards
```typescript
const SavedCards = ({ userId }) => {
    const [cards, setCards] = useState([]);
    
    useEffect(() => {
        fetch(`/api/users/${userId}/payment-methods`)
            .then(r => r.json())
            .then(data => setCards(data));
    }, [userId]);
    
    return (
        <div>
            {cards.map(card => (
                <div key={card.id} className="card-item">
                    <span className="brand">{card.card_brand}</span>
                    <span className="last4">•••• {card.card_last4}</span>
                    <span className="nickname">{card.nickname}</span>
                    <span className="expiry">{card.card_exp_month}/{card.card_exp_year}</span>
                    {card.is_default && <span className="badge">Default</span>}
                </div>
            ))}
        </div>
    );
};
```

### Checkout with Saved Card
```typescript
const Checkout = ({ amount, bookingId }) => {
    const [selectedCard, setSelectedCard] = useState(null);
    
    const handlePayment = async () => {
        const response = await fetch('/api/payments/create', {
            method: 'POST',
            body: JSON.stringify({
                payable_type: 'booking',
                payable_id: bookingId,
                amount: amount,
                payment_method_id: selectedCard?.id,  // Use saved card
                gateway_provider: selectedCard?.gateway_provider || 'stripe',
            })
        });
        
        const { payment_id, client_secret } = await response.json();
        
        // Confirm payment with Stripe
        await stripe.confirmCardPayment(client_secret);
    };
    
    return (
        <div>
            <SavedCards onSelect={setSelectedCard} />
            <button onClick={handlePayment}>Pay ${amount/100}</button>
        </div>
    );
};
```

---

## Summary

**Before**: 2 payment tables + repeated customer IDs  
**After**: 1 payment table + customer IDs in users + saved cards table

**Result**: 
- ✅ More DRY
- ✅ Better UX (saved cards with nicknames)
- ✅ More flexible (polymorphic payments)
- ✅ Easier to maintain

This is the **most DRY architecture** while still supporting all your requirements! 🚀
