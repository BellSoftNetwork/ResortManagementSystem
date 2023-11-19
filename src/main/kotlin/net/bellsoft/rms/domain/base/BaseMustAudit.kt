package net.bellsoft.rms.domain.base

import jakarta.persistence.EntityListeners
import jakarta.persistence.FetchType
import jakarta.persistence.JoinColumn
import jakarta.persistence.ManyToOne
import jakarta.persistence.MappedSuperclass
import net.bellsoft.rms.domain.user.User
import org.hibernate.envers.Audited
import org.hibernate.envers.RelationTargetAuditMode
import org.springframework.data.annotation.CreatedBy
import org.springframework.data.annotation.LastModifiedBy
import org.springframework.data.jpa.domain.support.AuditingEntityListener
import java.io.Serializable

@MappedSuperclass
@EntityListeners(AuditingEntityListener::class)
abstract class BaseMustAudit : Serializable, BaseTime() {
    @CreatedBy
    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "created_by_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "created_by", nullable = false, updatable = false)
    lateinit var createdBy: User
        private set

    @LastModifiedBy
    @Audited(
        withModifiedFlag = true,
        modifiedColumnName = "updated_by_mod",
        targetAuditMode = RelationTargetAuditMode.NOT_AUDITED,
    )
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "updated_by", nullable = false)
    lateinit var updatedBy: User
}
