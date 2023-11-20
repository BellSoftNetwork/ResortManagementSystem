package net.bellsoft.rms.service.reservation.dto

import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.domain.reservation.method.ReservationMethodRepository
import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomRepository
import java.time.LocalDate
import java.time.LocalDateTime

data class ReservationCreateDto(
    val reservationMethodId: Long,
    val roomId: Long? = null,
    val name: String,
    val phone: String = "",
    val peopleCount: Int = 0,
    val stayStartAt: LocalDate,
    val stayEndAt: LocalDate,
    val checkInAt: LocalDateTime? = null,
    val checkOutAt: LocalDateTime? = null,
    val price: Int = 0,
    val paymentAmount: Int = 0,
    val refundAmount: Int = 0,
    val brokerFee: Int = 0,
    val note: String = "",
    val canceledAt: LocalDateTime? = null,
    val status: ReservationStatus = ReservationStatus.PENDING,
) {
    private var isLoadedEntities = false

    private lateinit var reservationMethod: ReservationMethod
    private var room: Room? = null

    fun loadProxyEntities(reservationMethodRepository: ReservationMethodRepository, roomRepository: RoomRepository) {
        reservationMethodId.let { reservationMethod = reservationMethodRepository.getReferenceById(it) }
        roomId?.let { room = roomRepository.getReferenceById(it) }

        isLoadedEntities = true
    }

    fun toEntity(): Reservation {
        require(isLoadedEntities) { "loadProxyEntities 함수 실행 필요" }

        return Reservation(
            reservationMethod = reservationMethod,
            room = room,
            name = name,
            phone = phone,
            peopleCount = peopleCount,
            stayStartAt = stayStartAt,
            stayEndAt = stayEndAt,
            checkInAt = checkInAt,
            checkOutAt = checkOutAt,
            price = price,
            paymentAmount = paymentAmount,
            refundAmount = refundAmount,
            brokerFee = brokerFee,
            note = note,
            canceledAt = canceledAt,
            status = status,
        )
    }
}
