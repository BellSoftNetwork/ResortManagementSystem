package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.domain.paymentmethod.PaymentMethod
import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.mapper.common.ReferenceMapper
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
import org.springframework.beans.factory.annotation.Autowired

abstract class ReservationMapperDecorator : ReservationMapper {
    @Autowired
    private lateinit var delegate: ReservationMapper

    @Autowired
    private lateinit var referenceMapper: ReferenceMapper

    override fun toEntity(dto: ReservationCreateDto): Reservation {
        return Reservation(
            paymentMethod = referenceMapper.longIdToReference(
                dto.paymentMethodId,
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
        )
    }
}
