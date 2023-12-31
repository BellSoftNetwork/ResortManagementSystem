package net.bellsoft.rms.revision.entity

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.Table
import org.hibernate.envers.RevisionEntity
import org.hibernate.envers.RevisionNumber
import org.hibernate.envers.RevisionTimestamp
import java.io.Serial
import java.io.Serializable
import java.time.LocalDateTime

@Entity
@RevisionEntity
@Table(name = "revision_info")
class RevisionInfo : Serializable {
    @Id
    @RevisionNumber
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set

    @RevisionTimestamp
    @Column(name = "created_at")
    var createdAt: LocalDateTime = LocalDateTime.MIN
        private set

    companion object {
        @Serial
        private const val serialVersionUID: Long = 2472706047115136983L
    }
}
