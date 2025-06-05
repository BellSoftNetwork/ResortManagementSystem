package net.bellsoft.rms.migration.change

import liquibase.database.jvm.JdbcConnection
import net.bellsoft.rms.migration.AbstractCustomTaskChange
import java.sql.ResultSet
import java.sql.Timestamp

@Suppress("unused")
class ReservationRoomMigration : AbstractCustomTaskChange() {
    override fun migrate(connection: JdbcConnection) {
        val fetchQuery = """
            SELECT id, room_id, updated_at, deleted_at, updated_by
            FROM reservation
            WHERE room_id IS NOT NULL
        """.trimIndent()
        val resultSet = connection.prepareStatement(fetchQuery).executeQuery()
        val reservations = mutableListOf<ReservationResultSet>()

        while (resultSet.next())
            reservations.add(ReservationResultSet.of(resultSet))

        logger.info("Target ${reservations.size} reservation ids: ${reservations.map { it.id }.joinToString()}")
        insertReservationRooms(connection, reservations)
    }

    private fun insertReservationRooms(connection: JdbcConnection, reservations: List<ReservationResultSet>) {
        val insertQuery = """
            INSERT INTO reservation_room (reservation_id, room_id, created_at, updated_at, deleted_at, created_by, updated_by)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        val result = connection.prepareStatement(insertQuery).apply {
            reservations.forEach { reservation ->
                setLong(1, reservation.id)
                setLong(2, reservation.roomId)
                setTimestamp(3, reservation.updatedAt)
                setTimestamp(4, reservation.updatedAt)
                setTimestamp(5, reservation.deletedAt)
                setLong(6, reservation.updatedBy)
                setLong(7, reservation.updatedBy)

                addBatch()
                clearParameters()
            }
        }.executeBatch()

        logger.info("${result.sum()} rows inserted")
    }

    private data class ReservationResultSet(
        val id: Long,
        val roomId: Long,
        val updatedAt: Timestamp,
        val deletedAt: Timestamp,
        val updatedBy: Long,
    ) {
        companion object {
            fun of(resultSet: ResultSet) = ReservationResultSet(
                id = resultSet.getLong("id"),
                roomId = resultSet.getLong("room_id"),
                updatedAt = resultSet.getTimestamp("updated_at"),
                deletedAt = resultSet.getTimestamp("deleted_at"),
                updatedBy = resultSet.getLong("updated_by"),
            )
        }
    }
}
