---
id: api-legacy-reservation-management
title: "api-legacy 예약 관리"
status: completed
type: product
version: 1.0.0
created: 2026-01-07
updated: 2026-01-07

impact: backend
risk: high
effort: medium
---

# api-legacy 예약 관리

> 예약 CRUD, 통계, Envers 히스토리

---

## 1. 엔드포인트

| Method | Path | 설명 | 권한 |
|--------|------|------|------|
| GET | `/api/v1/reservations` | 목록 | USER |
| GET | `/api/v1/reservations/{id}` | 상세 | USER |
| POST | `/api/v1/reservations` | 생성 | ADMIN |
| PATCH | `/api/v1/reservations/{id}` | 수정 | ADMIN |
| DELETE | `/api/v1/reservations/{id}` | 삭제 | ADMIN |
| GET | `/api/v1/reservations/{id}/histories` | 이력 | ADMIN |
| GET | `/api/v1/reservation-statistics` | 통계 | USER |

---

## 2. Reservation 엔티티

```kotlin
@Entity
@Audited(withModifiedFlag = true)
@AuditTable("reservation_history")
@SQLDelete(sql = "UPDATE reservation SET deleted_at = NOW() WHERE id = ?")
@Where(clause = "deleted_at = '1970-01-01 00:00:00'")
@Comment("예약 정보")
class Reservation(
    @Audited(withModifiedFlag = true, targetAuditMode = NOT_AUDITED)
    @ManyToOne(fetch = LAZY, optional = false)
    @JoinColumn(nullable = false)
    var paymentMethod: PaymentMethod,
    
    @OneToMany(mappedBy = "reservation", cascade = [ALL], orphanRemoval = true)
    @OrderBy("id ASC")
    val rooms: MutableList<ReservationRoom> = mutableListOf(),
    
    @Column(nullable = false, length = 30)
    var name: String,
    
    @Column(nullable = false, length = 20)
    var phone: String,
    
    @Column(name = "people_count", nullable = false)
    var peopleCount: Int = 0,
    
    @Column(name = "stay_start_at", nullable = false)
    var stayStartAt: LocalDate,
    
    @Column(name = "stay_end_at", nullable = false)
    var stayEndAt: LocalDate,
    
    @Column(name = "check_in_at")
    var checkInAt: LocalDateTime? = null,
    
    @Column(name = "check_out_at")
    var checkOutAt: LocalDateTime? = null,
    
    @Column(nullable = false)
    var price: Int,
    
    @Column(nullable = false)
    var deposit: Int = 0,
    
    @Column(name = "payment_amount", nullable = false)
    var paymentAmount: Int = 0,
    
    @Column(name = "refund_amount", nullable = false)
    var refundAmount: Int = 0,
    
    @Column(name = "broker_fee", nullable = false)
    var brokerFee: Int = 0,
    
    @Column(nullable = false, length = 200)
    var note: String = "",
    
    @Column(name = "canceled_at")
    var canceledAt: LocalDateTime? = null,
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var status: ReservationStatus,
    
    @Column(nullable = false, columnDefinition = "TINYINT")
    var type: ReservationType
) : BaseMustAuditEntity()
```

---

## 3. 상태 및 유형

```kotlin
enum class ReservationStatus(val value: Int) {
    REFUND(-10),
    CANCEL(-1),
    PENDING(0),
    NORMAL(1)
}

enum class ReservationType(val value: Int) {
    STAY(0),
    MONTHLY_RENT(10)
}
```

---

## 4. 서비스

```kotlin
interface ReservationService {
    fun getAll(pageable: Pageable, filter: ReservationFilterRequest): EntityListDto<Reservation, ReservationDetailDto>
    fun getById(id: Long): ReservationDetailDto
    fun create(request: ReservationCreateRequest): ReservationDetailDto
    fun update(id: Long, request: ReservationPatchRequest): ReservationDetailDto
    fun delete(id: Long)
    fun getHistories(id: Long, pageable: Pageable): EntityListDto<..., EntityRevisionDto<ReservationDetailDto>>
    fun getStatistics(query: ReservationStatisticsQuery): ReservationStatisticsDto
}
```

---

## 5. 필터

```kotlin
data class ReservationFilterRequest(
    val status: ReservationStatus?,
    val type: ReservationType?,
    val roomId: Long?,
    val stayStartAt: LocalDate?,
    val stayEndAt: LocalDate?,
    val search: String?
)
```

---

## 6. 통계

### 6.1 쿼리

```kotlin
data class ReservationStatisticsQuery(
    val startDate: LocalDate,
    val endDate: LocalDate,
    val periodType: PeriodType = MONTHLY
)

enum class PeriodType {
    DAILY, MONTHLY, YEARLY
}
```

### 6.2 응답

```kotlin
data class ReservationStatisticsDto(
    val periodType: PeriodType,
    val stats: List<StatisticsData>,
    val monthlyStats: List<MonthlyStats>
)

data class StatisticsData(
    val period: String,
    val totalSales: Long,
    val totalReservations: Int,
    val totalGuests: Int
)
```

---

## 7. 중개료 계산

```kotlin
@Transactional
override fun create(request: ReservationCreateRequest): ReservationDetailDto {
    val paymentMethod = paymentMethodRepository.findByIdOrThrow(request.paymentMethodId)
    
    val brokerFee = (request.paymentAmount * paymentMethod.commissionRate).toInt()
    
    val reservation = Reservation(
        paymentMethod = paymentMethod,
        brokerFee = brokerFee,
        // ...
    )
    
    return reservationMapper.toDetailDto(reservationRepository.save(reservation))
}
```
