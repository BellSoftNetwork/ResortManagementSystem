package net.bellsoft.rms.user.repository.impl

import com.querydsl.jpa.impl.JPAQueryFactory
import net.bellsoft.rms.user.entity.QUser
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.repository.UserCustomRepository
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
