---
id: api-legacy-payment-method
title: "api-legacy 결제 수단 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: medium
effort: small
---

# api-legacy 결제 수단 관리

> 결제 수단 CRUD 및 수수료율 관리

---

## 1. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/payment-methods` | 목록 | USER |
| GET | `/api/v1/payment-methods/{id}` | 상세 | USER |
| POST | `/api/v1/payment-methods` | 생성 | ADMIN |
| PATCH | `/api/v1/payment-methods/{id}` | 수정 | ADMIN |
| DELETE | `/api/v1/payment-methods/{id}` | 삭제 | ADMIN |

---

## 2. PaymentMethod 엔티티

```kotlin
@Entity
@Table(
    uniqueConstraints = [
        UniqueConstraint(columnNames = ["name", "deleted_at"])
    ]
)
@SQLDelete(sql = "UPDATE payment_method SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at = '1970-01-01 00:00:00'")
@Comment("결제 수단")
class PaymentMethod(
    @Column(nullable = false, length = 20)
    var name: String,
    
    @Column(name = "commission_rate", nullable = false)
    var commissionRate: Double,
    
    @Column(name = "require_unpaid_amount_check", nullable = false)
    var requireUnpaidAmountCheck: Boolean = false,
    
    @Column(name = "is_default_select", nullable = false)
    var isDefaultSelect: Boolean = false,
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var status: PaymentMethodStatus = PaymentMethodStatus.ACTIVE
) : BaseTimeEntity()
```

---

## 3. 상태

```kotlin
enum class PaymentMethodStatus(val value: Int) {
    INACTIVE(-1),
    ACTIVE(1)
}
```

---

## 4. 서비스

```kotlin
interface PaymentMethodService {
    fun getAll(pageable: Pageable): EntityListDto<PaymentMethod, PaymentMethodDetailDto>
    fun getById(id: Long): PaymentMethodDetailDto
    fun create(request: PaymentMethodCreateRequest): PaymentMethodDetailDto
    fun update(id: Long, request: PaymentMethodPatchRequest): PaymentMethodDetailDto
    fun delete(id: Long)
}
```

---

## 5. 기본 선택 로직

```kotlin
@Transactional
override fun update(id: Long, request: PaymentMethodPatchRequest): PaymentMethodDetailDto {
    val paymentMethod = paymentMethodRepository.findByIdOrThrow(id)
    
    // 기본 선택 설정 시 다른 것들 해제
    if (request.isDefaultSelect.isPresent && request.isDefaultSelect.get() == true) {
        paymentMethodRepository.resetAllDefaultSelects()
    }
    
    paymentMethodMapper.updateEntity(paymentMethod, request)
    return paymentMethodMapper.toDetailDto(paymentMethodRepository.save(paymentMethod))
}
```

---

## 6. DTO

```kotlin
data class PaymentMethodDetailDto(
    val id: Long,
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean,
    val isDefaultSelect: Boolean,
    val status: PaymentMethodStatus,
    val createdAt: LocalDateTime,
    val updatedAt: LocalDateTime
)

data class PaymentMethodCreateRequest(
    val name: String,
    val commissionRate: Double,
    val requireUnpaidAmountCheck: Boolean = false
)

data class PaymentMethodPatchRequest(
    val name: JsonNullable<String>,
    val commissionRate: JsonNullable<Double>,
    val requireUnpaidAmountCheck: JsonNullable<Boolean>,
    val isDefaultSelect: JsonNullable<Boolean>,
    val status: JsonNullable<PaymentMethodStatus>
)
```
