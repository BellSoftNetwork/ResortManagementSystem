package net.bellsoft.rms.domain.paymentmethod

import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.stereotype.Repository

@Repository
interface PaymentMethodRepository : JpaRepository<PaymentMethod, Long>
