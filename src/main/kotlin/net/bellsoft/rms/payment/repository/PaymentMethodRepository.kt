package net.bellsoft.rms.payment.repository

import net.bellsoft.rms.payment.entity.PaymentMethod
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Modifying
import org.springframework.data.jpa.repository.Query
import org.springframework.stereotype.Repository
import org.springframework.transaction.annotation.Transactional

@Repository
interface PaymentMethodRepository : JpaRepository<PaymentMethod, Long> {
    @Transactional
    @Modifying
    @Query("UPDATE PaymentMethod p SET p.isDefaultSelect = false WHERE p.isDefaultSelect = true")
    fun updateAllIsDefaultSelectToFalse(): Int
}
