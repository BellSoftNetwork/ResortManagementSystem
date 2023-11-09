package net.bellsoft.rms.domain.room

import jakarta.persistence.Column
import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.GenerationType
import jakarta.persistence.Id
import jakarta.persistence.Table
import net.bellsoft.rms.domain.base.BaseTime
import org.hibernate.annotations.SQLDelete
import org.hibernate.annotations.Where

@Entity
@Table(name = "room")
@SQLDelete(sql = "UPDATE room SET deleted_at = NOW() WHERE id = ?")
@Where(clause = BaseTime.SOFT_DELETE_CONDITION)
class Room(
    number: String,
    peekPrice: Int? = null,
    offPeekPrice: Int? = null,
    description: String = "",
    note: String = "",

    @Column(
        name = "status",
        nullable = false,
        columnDefinition = "TINYINT",
    )
    var status: RoomStatus = RoomStatus.INACTIVE,
) : BaseTime() {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id", nullable = false, unique = true, updatable = false)
    var id: Long = 0
        private set

    @Column(name = "number", nullable = false, length = 10)
    var number: String = number
        private set

    @Column(name = "peek_price")
    var peekPrice: Int? = peekPrice
        private set

    @Column(name = "off_peek_price")
    var offPeekPrice: Int? = offPeekPrice
        private set

    @Column(name = "desciption", nullable = false, length = 200)
    var description: String = description
        private set

    @Column(name = "note", nullable = false, length = 200)
    var note: String = note
        private set
}
