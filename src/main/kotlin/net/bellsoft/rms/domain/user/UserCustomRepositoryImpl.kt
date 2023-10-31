package net.bellsoft.rms.domain.user

import com.querydsl.jpa.impl.JPAQueryFactory
import org.springframework.stereotype.Repository
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime

@Repository
class UserCustomRepositoryImpl(
    private val jpaQueryFactory: JPAQueryFactory,
) : UserCustomRepository {
    override fun getYearCreatedUsers(year: Int): List<User> {
        val startDate = LocalDateTime.of(LocalDate.of(year, 1, 1), LocalTime.MIN)
        val endDate = startDate.plusYears(1)

        return jpaQueryFactory
            .selectFrom(QUser.user)
            .where(QUser.user.createdAt.goe(startDate).and(QUser.user.createdAt.lt(endDate)))
            .fetch()
    }
}
