package net.bellsoft.rms.payment.repository

import net.bellsoft.rms.payment.entity.PaymentMethod
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface PaymentMethodRepository : JpaRepository<PaymentMethod, Long>
