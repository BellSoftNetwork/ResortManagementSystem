package net.bellsoft.rms.room.repository.impl

import com.querydsl.core.types.Order
import com.querydsl.core.types.OrderSpecifier
import com.querydsl.core.types.dsl.BooleanExpression
import com.querydsl.core.types.dsl.Expressions
import com.querydsl.jpa.JPAExpressions
import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.reservation.entity.QReservation.reservation
import net.bellsoft.rms.reservation.entity.QReservationRoom.reservationRoom
import net.bellsoft.rms.room.dto.filter.RoomFilterDto
import net.bellsoft.rms.room.entity.QRoom.room
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.repository.RoomCustomRepository
import net.bellsoft.rms.room.type.RoomStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.support.PageableExecutionUtils
import org.springframework.stereotype.Repository
import java.time.LocalDate

@Repository
class RoomCustomRepositoryImpl(
    private val jpaQueryFactory: JPAQueryFactory,
) : RoomCustomRepository {
    override fun getFilteredRooms(pageable: Pageable, filter: RoomFilterDto): Page<Room> {
        val result = getFilteredRoomsBaseQuery(filter)
            .select(room)
            .orderBy(
                *pageable.sort.map {
                    OrderSpecifier(
                        if (it.isAscending) Order.ASC else Order.DESC,
                        Expressions.path(String::class.java, room, it.property),
                    )
                }.toList().toTypedArray<OrderSpecifier<String>?>(),
            )
            .offset(pageable.offset)
            .limit(pageable.pageSize.toLong())
            .fetch()

        return PageableExecutionUtils.getPage(result, pageable) {
            getFilteredRoomsBaseQuery(filter)
                .select(room.count())
                .fetchOne()!!
        }
    }

    override fun getReservedRooms(filter: RoomFilterDto, roomIds: Set<Long>): Set<Room> {
        return jpaQueryFactory
            .from(room)
            .where(
                eqStatus(filter.status),
                duplicatedRooms(filter, roomIds),
            )
            .select(room)
            .fetch()
            .toSet()
    }

    private fun getFilteredRoomsBaseQuery(filter: RoomFilterDto) = jpaQueryFactory
        .from(room)
        .where(
            eqStatus(filter.status),
            reservedRooms(filter),
        )

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

    private fun duplicatedRooms(filter: RoomFilterDto, roomIds: Set<Long>): BooleanExpression? {
        if (filter.stayStartAt == null || filter.stayEndAt == null)
            return null

        return room.id.`in`(
            JPAExpressions
                .select(reservationRoom.room.id)
                .from(reservationRoom)
                .where(
                    reservationRoom.room.id.`in`(roomIds),
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
}
