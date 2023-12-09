package net.bellsoft.rms.payment.mapper

import net.bellsoft.rms.common.config.BaseMapperConfig
import net.bellsoft.rms.payment.dto.response.PaymentMethodDetailDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodCreateDto
import net.bellsoft.rms.payment.dto.service.PaymentMethodPatchDto
import net.bellsoft.rms.payment.entity.PaymentMethod
import org.mapstruct.BeanMapping
import org.mapstruct.Mapper
import org.mapstruct.Mapping
import org.mapstruct.MappingTarget
import org.mapstruct.Mappings
import org.mapstruct.NullValuePropertyMappingStrategy

@Mapper(config = BaseMapperConfig::class)
interface PaymentMethodMapper {
    fun toDto(entity: PaymentMethod): PaymentMethodDetailDto

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.SET_TO_DEFAULT)
    fun toEntity(dto: PaymentMethodCreateDto): PaymentMethod

    @Mappings(
        Mapping(target = "status", ignore = true),
    )
    fun updateEntity(dto: PaymentMethodPatchDto, @MappingTarget entity: PaymentMethod)
}
