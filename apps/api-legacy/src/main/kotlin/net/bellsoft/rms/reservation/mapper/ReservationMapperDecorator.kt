package net.bellsoft.rms.reservation.mapper

import net.bellsoft.rms.common.mapper.ReferenceMapper
import net.bellsoft.rms.payment.entity.PaymentMethod
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.entity.Reservation
import org.springframework.beans.factory.annotation.Autowired

abstract class ReservationMapperDecorator : ReservationMapper {
    @Autowired
    private lateinit var delegate: ReservationMapper

    @Autowired
    private lateinit var referenceMapper: ReferenceMapper

    override fun toEntity(dto: ReservationCreateDto): Reservation {
        return Reservation(
            paymentMethod = referenceMapper.refIdToReference(
                dto.paymentMethod,
                PaymentMethod::class.java,
            ),
            name = dto.name,
            phone = dto.phone,
            peopleCount = dto.peopleCount,
            stayStartAt = dto.stayStartAt,
            stayEndAt = dto.stayEndAt,
            checkInAt = dto.checkInAt,
            checkOutAt = dto.checkOutAt,
            price = dto.price,
            paymentAmount = dto.paymentAmount,
            refundAmount = dto.refundAmount,
            brokerFee = dto.brokerFee,
            note = dto.note,
            canceledAt = dto.canceledAt,
            status = dto.status,
            type = dto.type,
        )
    }
}
