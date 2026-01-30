# Database Design Decision: user_id in Payment Tables

## Question
Should we keep `user_id` as a foreign key in both `bookings` and `booking_payments` tables?

## Decision: YES ✅

Keep `user_id` in payment tables despite it being derivable from the parent table.

---

## Rationale

### 1. **Query Performance**
```sql
-- Direct lookup (fast - uses index)
SELECT * FROM booking_payments WHERE user_id = 'user-123';

-- vs. Join required (slower)
SELECT bp.* 
FROM booking_payments bp
JOIN bookings b ON bp.booking_id = b.id
WHERE b.user_id = 'user-123';
```

**Impact**: Payment queries are frequent (user dashboards, analytics, reports).

### 2. **Audit Trail**
Financial records should be self-contained. If a booking is deleted (even soft delete), the payment record still knows who paid.

```sql
-- Payment survives booking deletion
SELECT user_id, amount, created_at 
FROM booking_payments 
WHERE id = 'payment-123';
```

### 3. **Future Flexibility**
Allows for scenarios where the payer differs from the booker:
- Gift purchases
- Corporate bookings (assistant books, company pays)
- Family bookings (parent pays for adult child's booking)

```
booking.user_id = 'alice'   (who booked)
payment.user_id = 'bob'     (who paid)
```

### 4. **Industry Standard**
Major payment platforms (Stripe, PayPal, Square) all track the payer separately from the transaction entity.

### 5. **Reporting & Analytics**
```sql
-- User lifetime value (simple query)
SELECT user_id, SUM(amount) as total_spent 
FROM booking_payments 
GROUP BY user_id;

-- User payment history (no joins needed)
SELECT * FROM booking_payments 
WHERE user_id = 'user-123' 
ORDER BY created_at DESC;
```

---

## Trade-offs

### Cons
- ❌ **Data redundancy**: Violates 3NF (Third Normal Form)
- ❌ **Update anomalies**: Must update multiple tables if user changes
- ❌ **Storage cost**: Extra UUID per payment row (~16 bytes)
- ❌ **Consistency risk**: Could have mismatched user_ids

### Mitigation Strategies

#### 1. **Application-Level Validation** (Primary)
Validate in the service layer before creating payments:

```go
func (s *BookingPaymentService) CreatePayment(bookingID, userID uuid.UUID, amount int) error {
    // Validate user owns the booking
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        return err
    }
    
    if booking.UserID != userID {
        return errors.New("user does not own this booking")
    }
    
    // Create payment with validated user_id
    payment := &BookingPayment{
        BookingID: bookingID,
        UserID:    userID,
        Amount:    amount,
    }
    
    return s.repo.Create(payment)
}
```

#### 2. **GORM Hooks** (Secondary)
Add validation at the model level:

```go
func (bp *BookingPayment) BeforeCreate(tx *gorm.DB) error {
    // Verify user_id matches booking.user_id
    var booking Booking
    if err := tx.First(&booking, bp.BookingID).Error; err != nil {
        return err
    }
    
    if booking.UserID != bp.UserID {
        return errors.New("payment user_id must match booking user_id")
    }
    
    return nil
}
```

#### 3. **Database Triggers** (Optional - Future)
For additional safety, PostgreSQL triggers can be added later if needed.

---

## Implementation Guidelines

### 1. **Service Layer Validation** (Primary)
```go
// In BookingPaymentService
func (s *Service) CreatePayment(bookingID, userID uuid.UUID, amount int) error {
    // Validate user owns the booking
    booking, err := s.bookingRepo.FindByID(bookingID)
    if err != nil {
        return err
    }
    
    if booking.UserID != userID {
        return errors.New("user does not own this booking")
    }
    
    // Create payment
    payment := &BookingPayment{
        BookingID: bookingID,
        UserID:    userID,  // Explicitly set (matches booking.UserID)
        Amount:    amount,
    }
    
    return s.repo.Create(payment)
}
```

### 2. **GORM Model Hooks** (Secondary)
```go
// In models/booking_payment.go
func (bp *BookingPayment) BeforeCreate(tx *gorm.DB) error {
    var booking Booking
    if err := tx.First(&booking, bp.BookingID).Error; err != nil {
        return fmt.Errorf("booking not found: %w", err)
    }
    
    if booking.UserID != bp.UserID {
        return errors.New("payment user_id must match booking user_id")
    }
    
    return nil
}
```

### 3. **Testing**
```go
// Test: Payment user must match booking user
func TestCreatePayment_UserMismatch(t *testing.T) {
    booking := createTestBooking(userAlice)
    
    // Try to create payment with different user
    err := paymentService.CreatePayment(booking.ID, userBob.ID, 10000)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "does not own this booking")
}

// Test: Valid payment succeeds
func TestCreatePayment_Success(t *testing.T) {
    booking := createTestBooking(userAlice)
    
    // Create payment with correct user
    err := paymentService.CreatePayment(booking.ID, userAlice.ID, 10000)
    
    assert.NoError(t, err)
    
    // Verify payment was created
    payment, err := paymentRepo.FindByBookingID(booking.ID)
    assert.NoError(t, err)
    assert.Equal(t, userAlice.ID, payment.UserID)
}
```

---

## Alternative Considered: Remove user_id from Payments

### Pros
- ✅ Normalized (3NF compliant)
- ✅ No redundancy
- ✅ Single source of truth

### Cons
- ❌ Slower queries (always need join)
- ❌ More complex reporting
- ❌ Less flexible for future features
- ❌ Payment record not self-contained

**Decision**: Rejected due to performance and audit requirements.

---

## Conclusion

**Keep `user_id` in payment tables** with database triggers to ensure consistency.

This is a **controlled denormalization** - we accept the redundancy for:
- Better query performance
- Self-contained audit records
- Future flexibility
- Industry-standard pattern

The consistency risk is mitigated by triggers and application-level validation.

---

## References

- PostgreSQL Triggers: https://www.postgresql.org/docs/current/triggers.html
- Database Normalization: https://en.wikipedia.org/wiki/Database_normalization
- Stripe API Design: https://stripe.com/docs/api/charges
- Denormalization for Performance: https://en.wikipedia.org/wiki/Denormalization
