package net.bellsoft.rms.domain.paymentmethod

import jakarta.persistence.AttributeConverter
import jakarta.persistence.Converter

@Converter(autoApply = true)
class PaymentMethodStatusConverter : AttributeConverter<PaymentMethodStatus, Int> {
    override fun convertToDatabaseColumn(attribute: PaymentMethodStatus): Int = attribute.value
    override fun convertToEntityAttribute(dbData: Int): PaymentMethodStatus =
        PaymentMethodStatus.values().first { it.value == dbData }
}
