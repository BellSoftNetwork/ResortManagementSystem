package net.bellsoft.rms.mapper.model

import net.bellsoft.rms.domain.paymentmethod.PaymentMethod
import net.bellsoft.rms.mapper.config.BaseMapperConfig
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodCreateDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodDetailDto
import net.bellsoft.rms.service.paymentmethod.dto.PaymentMethodPatchDto
import org.mapstruct.BeanMapping
import org.mapstruct.Mapper
import org.mapstruct.MappingTarget
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(config = BaseMapperConfig::class)
interface PaymentMethodMapper {
    fun toDto(entity: PaymentMethod): PaymentMethodDetailDto

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.SET_TO_DEFAULT)
    fun toEntity(dto: PaymentMethodCreateDto): PaymentMethod

    fun updateEntity(dto: PaymentMethodPatchDto, @MappingTarget entity: PaymentMethod)
}
