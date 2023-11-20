package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomRepository
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationUpdateDto(
    val reservationMethodId: Long? = null,
    val roomId: Long? = null,
    val name: String? = null,
    val phone: String? = null,
    val peopleCount: Int? = null,
    val stayStartAt: LocalDate? = null,
    val stayEndAt: LocalDate? = null,
    val checkInAt: LocalDateTime? = null,
    val checkOutAt: LocalDateTime? = null,
    val price: Int? = null,
    val paymentAmount: Int? = null,
    val refundAmount: Int? = null,
    val brokerFee: Int? = null,
    val note: String? = null,
    val canceledAt: LocalDateTime? = null,
    val status: ReservationStatus? = null,
) {
    private var isLoadedEntities = false

    private var reservationMethod: ReservationMethod? = null
    private var room: Room? = null

    fun loadProxyEntities(reservationMethodRepository: ReservationMethodRepository, roomRepository: RoomRepository) {
        reservationMethodId?.let { reservationMethod = reservationMethodRepository.getReferenceById(it) }
        roomId?.let { room = roomRepository.getReferenceById(it) }

        isLoadedEntities = true
    }

    fun updateEntity(entity: Reservation) {
        require(isLoadedEntities) { "loadProxyEntities 함수 실행 필요" }

        // FIXME: 현재 이 코드로는 실제 nullable 한 칼럼에 null 을 설정하려는 케이스는 대응할 수 없음
        reservationMethod?.let { entity.reservationMethod = it }
        room?.let { entity.room = it }
        name?.let { entity.name = it }
        phone?.let { entity.phone = it }
        peopleCount?.let { entity.peopleCount = it }
        stayStartAt?.let { entity.stayStartAt = it }
        stayEndAt?.let { entity.stayEndAt = it }
        checkInAt?.let { entity.checkInAt = it }
        checkOutAt?.let { entity.checkOutAt = it }
        price?.let { entity.price = it }
        paymentAmount?.let { entity.paymentAmount = it }
        refundAmount?.let { entity.refundAmount = it }
        brokerFee?.let { entity.brokerFee = it }
        note?.let { entity.note = it }
        canceledAt?.let { entity.canceledAt = it }
        status?.let { entity.status = it }
    }
}
