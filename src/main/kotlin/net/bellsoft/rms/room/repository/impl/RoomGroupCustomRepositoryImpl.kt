package net.bellsoft.rms.room.repository.impl

import com.querydsl.core.types.dsl.BooleanExpression
import com.querydsl.jpa.JPAExpressions
import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.common.exception.DataNotFoundException
import net.bellsoft.rms.reservation.entity.QReservation.reservation
import net.bellsoft.rms.reservation.entity.QReservationRoom.reservationRoom
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.dto.projection.QRoomLastReservationProjection
import net.bellsoft.rms.room.dto.projection.RoomLastReservationProjection
import net.bellsoft.rms.room.entity.QRoom.room
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.entity.RoomGroup
import net.bellsoft.rms.room.repository.RoomGroupCustomRepository
import net.bellsoft.rms.room.type.RoomStatus
import org.springframework.stereotype.Repository
import java.time.LocalDate

@Repository
class RoomGroupCustomRepositoryImpl(
    private val jpaQueryFactory: JPAQueryFactory,
) : RoomGroupCustomRepository {
    override fun getFilteredRoomsOrderByLastStayAt(
        roomGroup: RoomGroup,
        filter: RoomFilterDto,
    ): List<RoomLastReservationProjection> {
        return availableFilteredRooms(roomGroup, filter)
            .map { getRoomWithLastReservation(it.id, filter) }
            .sortedBy { it.lastReservation?.stayEndAt }
    }

    private fun getRoomWithLastReservation(roomId: Long, filter: RoomFilterDto): RoomLastReservationProjection {
        val reservationJoinExpressions = listOfNotNull(
            reservationRoom.reservation.eq(reservation),
            loeStayEndAt(filter.stayStartAt),
        )

        return jpaQueryFactory
            .select(QRoomLastReservationProjection(room, reservation))
            .from(room)
            .leftJoin(reservationRoom).on(room.eq(reservationRoom.room))
            .leftJoin(reservation).on(*reservationJoinExpressions.toTypedArray())
            .where(room.id.eq(roomId))
            .orderBy(reservation.stayEndAt.desc())
            .limit(1)
            .fetchOne() ?: throw DataNotFoundException("객실 정보 조회 실패")
    }

    private fun availableFilteredRooms(roomGroup: RoomGroup, filter: RoomFilterDto): List<Room> {
        return jpaQueryFactory
            .select(room)
            .from(room)
            .where(
                eqRoomGroupId(roomGroup.id),
                eqStatus(filter.status),
                reservedRooms(filter),
            )
            .fetch()
    }

    private fun reservedRooms(filter: RoomFilterDto): BooleanExpression? {
        if (filter.stayStartAt == null || filter.stayEndAt == null)
            return null

        return room.id.notIn(
            JPAExpressions
                .select(reservationRoom.room.id)
                .from(reservationRoom)
                .where(
                    filterReservationRooms(filter),
                )
                .distinct(),
        )
    }

    private fun filterReservationRooms(filter: RoomFilterDto): BooleanExpression? =
        reservationRoom.reservation.id.`in`(
            JPAExpressions
                .select(reservation.id)
                .from(reservation)
                .where(
                    neReservationId(filter.excludeReservationId),
                    beforeDateFilterExpressions(filter)
                        ?.or(afterDateFilterExpressions(filter))
                        ?.or(wrapDateFilterExpressions(filter)),
                ),
        )

    /**
     * 기존 예약 기간: ###=
     * 희망 예약 기간: =###
     */
    private fun beforeDateFilterExpressions(filter: RoomFilterDto) =
        loeStayStartAt(filter.stayStartAt)?.and(gtStayEndAt(filter.stayStartAt))

    /**
     * 기존 예약 기간: =###
     * 희망 예약 기간: ###=
     */
    private fun afterDateFilterExpressions(filter: RoomFilterDto) =
        ltStayStartAt(filter.stayEndAt)?.and(goeStayEndAt(filter.stayEndAt))

    /**
     * 기존 예약 기간: ####
     * 희망 예약 기간: =##=
     */
    private fun wrapDateFilterExpressions(filter: RoomFilterDto) =
        goeStayStartAt(filter.stayStartAt)?.and(loeStayEndAt(filter.stayEndAt))

    private fun goeStayStartAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayStartAt.goe(it) }

    private fun loeStayStartAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayStartAt.loe(it) }

    private fun ltStayStartAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayStartAt.lt(it) }

    private fun goeStayEndAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayEndAt.goe(it) }

    private fun gtStayEndAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayEndAt.gt(it) }

    private fun loeStayEndAt(localDate: LocalDate?) =
        localDate?.let { reservation.stayEndAt.loe(it) }

    private fun eqStatus(status: RoomStatus?) =
        status?.let { room.status.eq(status) }

    private fun neReservationId(reservationId: Long?) =
        reservationId?.let { reservation.id.ne(reservationId) }

    private fun eqRoomGroupId(roomId: Long?) =
        roomId?.let { room.roomGroup.id.eq(roomId) }
}
