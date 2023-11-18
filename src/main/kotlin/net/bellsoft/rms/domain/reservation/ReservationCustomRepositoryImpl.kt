package net.bellsoft.rms.domain.reservation

import com.querydsl.core.types.Order
import com.querydsl.core.types.OrderSpecifier
import com.querydsl.core.types.dsl.Expressions
import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.service.reservation.dto.ReservationFilterDto
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
            goeStayStartAt(filter.stayStartAt),
            loeStayStartAt(filter.stayEndAt),
            likeSearchText(filter.searchText),
        )

    private fun goeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.goe(it) }

    private fun loeStayStartAt(stayStartAt: LocalDate?) =
        stayStartAt?.let { QReservation.reservation.stayStartAt.loe(it) }

    private fun likeSearchText(searchText: String?) =
        searchText?.let {
            QReservation.reservation.name.like("%$it%")
                .or(QReservation.reservation.phone.like("%$it%"))
        }
}
