package net.bellsoft.rms.reservation.repository.impl

import com.querydsl.core.types.Order
import com.querydsl.core.types.OrderSpecifier
import com.querydsl.core.types.dsl.Expressions
import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.reservation.dto.filter.ReservationFilterDto
import net.bellsoft.rms.reservation.entity.QReservation
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.repository.ReservationCustomRepository
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.bellsoft.rms.reservation.type.ReservationType
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.data.support.PageableExecutionUtils
import org.springframework.stereotype.Repository
import java.time.LocalDate

@Repository
class ReservationCustomRepositoryImpl(
    private val jpaQueryFactory: JPAQueryFactory,
) : ReservationCustomRepository {
    override fun getFilteredReservations(pageable: Pageable, filter: ReservationFilterDto): Page<Reservation> {
        val result = getFilteredReservationsBaseQuery(filter)
            .select(QReservation.reservation)
            .orderBy(
                *pageable.sort.map {
                    OrderSpecifier(
                        if (it.isAscending) Order.ASC else Order.DESC,
                        Expressions.path(String::class.java, QReservation.reservation, it.property),
                    )
                }.toList().toTypedArray<OrderSpecifier<String>?>(),
            )
            .offset(pageable.offset)
            .limit(pageable.pageSize.toLong())
            .fetch()

        return PageableExecutionUtils.getPage(result, pageable) {
            getFilteredReservationsBaseQuery(filter)
                .select(QReservation.reservation.count())
                .fetchOne()!!
        }
    }

    private fun getFilteredReservationsBaseQuery(filter: ReservationFilterDto) = jpaQueryFactory
        .from(QReservation.reservation)
        .where(
            goeStayStartAt(filter.stayStartAt)?.or(goeStayEndAt(filter.stayStartAt)),
            loeStayStartAt(filter.stayEndAt)?.or(loeStayEndAt(filter.stayEndAt)),
            likeSearchText(filter.searchText),
            eqStatus(filter.status),
            eqType(filter.type),
        )

    private fun goeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.goe(it) }

    private fun loeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.loe(it) }

    private fun goeStayEndAt(stayEndAt: LocalDate?) =
        stayEndAt?.let { QReservation.reservation.stayEndAt.goe(it) }

    private fun loeStayEndAt(stayEndAt: LocalDate?) =
        stayEndAt?.let { QReservation.reservation.stayEndAt.loe(it) }

    private fun likeSearchText(searchText: String?) =
        searchText?.let {
            QReservation.reservation.name.like("%$it%")
                .or(QReservation.reservation.phone.like("%$it%"))
        }

    private fun eqStatus(status: ReservationStatus?) =
        status?.let { QReservation.reservation.status.eq(it) }

    private fun eqType(status: ReservationType?) =
        status?.let { QReservation.reservation.type.eq(it) }
}
