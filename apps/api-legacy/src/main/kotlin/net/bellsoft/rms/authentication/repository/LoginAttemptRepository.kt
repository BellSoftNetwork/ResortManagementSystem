package net.bellsoft.rms.authentication.repository

import net.bellsoft.rms.authentication.entity.LoginAttempt
import org.springframework.data.jpa.repository.JpaRepository
import org.springframework.data.jpa.repository.Query
import org.springframework.data.repository.query.Param
import org.springframework.stereotype.Repository
import java.time.LocalDateTime

/**
 * 로그인 시도 저장소
 */
@Repository
interface LoginAttemptRepository : JpaRepository<LoginAttempt, Long> {
    /**
     * 특정 사용자의 최근 로그인 시도 횟수를 조회한다.
     *
     * @param username 사용자 ID
     * @param since 조회 시작 시간
     * @return 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.username = :username AND la.attemptAt > :since
        """,
    )
    fun countRecentAttempts(
        @Param("username") username: String,
        @Param("since") since: LocalDateTime,
    ): Long

    /**
     * 특정 사용자의 최근 실패한 로그인 시도 횟수를 조회한다.
     *
     * @param username 사용자 ID
     * @param since 조회 시작 시간
     * @return 실패한 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.username = :username AND la.successful = false AND la.attemptAt > :since
        """,
    )
    fun countRecentFailedAttempts(
        @Param("username") username: String,
        @Param("since") since: LocalDateTime,
    ): Long

    /**
     * 특정 IP 주소의 최근 실패한 로그인 시도 횟수를 조회한다.
     *
     * @param ipAddress IP 주소
     * @param since 조회 시작 시간
     * @return 실패한 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.ipAddress = :ipAddress AND la.successful = false AND la.attemptAt > :since
        """,
    )
    fun countRecentFailedAttemptsByIp(
        @Param("ipAddress") ipAddress: String,
        @Param("since") since: LocalDateTime,
    ): Long

    /**
     * 특정 사용자와 IP 주소의 최근 실패한 로그인 시도 횟수를 조회한다.
     *
     * @param username 사용자 ID
     * @param ipAddress IP 주소
     * @param since 조회 시작 시간
     * @return 실패한 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.username = :username AND la.ipAddress = :ipAddress AND la.successful = false AND la.attemptAt > :since
        """,
    )
    fun countRecentFailedAttemptsByUsernameAndIp(
        @Param("username") username: String,
        @Param("ipAddress") ipAddress: String,
        @Param("since") since: LocalDateTime,
    ): Long

    /**
     * 특정 사용자의 마지막 성공한 로그인 시도를 조회한다.
     *
     * @param username 사용자 ID
     * @return 마지막 성공한 로그인 시도
     */
    fun findTopByUsernameAndSuccessfulTrueOrderByAttemptAtDesc(username: String): LoginAttempt?

    /**
     * 특정 사용자와 IP 주소의 마지막 성공한 로그인 시도를 조회한다.
     *
     * @param username 사용자 ID
     * @param ipAddress IP 주소
     * @return 마지막 성공한 로그인 시도
     */
    fun findTopByUsernameAndIpAddressAndSuccessfulTrueOrderByAttemptAtDesc(
        username: String,
        ipAddress: String,
    ): LoginAttempt?

    /**
     * 특정 IP 주소의 마지막 성공한 로그인 시도를 조회한다.
     *
     * @param ipAddress IP 주소
     * @return 마지막 성공한 로그인 시도
     */
    fun findTopByIpAddressAndSuccessfulTrueOrderByAttemptAtDesc(ipAddress: String): LoginAttempt?

    /**
     * 특정 사용자와 IP 주소의 마지막 성공한 로그인 이후의 실패한 로그인 시도 횟수를 조회한다.
     * 성공한 로그인이 없는 경우 주어진 시간 이후의 모든 실패한 로그인 시도를 조회한다.
     *
     * @param username 사용자 ID
     * @param ipAddress IP 주소
     * @param since 조회 시작 시간
     * @return 실패한 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.username = :username AND la.ipAddress = :ipAddress AND la.successful = false
        AND la.attemptAt > COALESCE(
            (SELECT MAX(la2.attemptAt) FROM LoginAttempt la2
             WHERE la2.username = :username AND la2.ipAddress = :ipAddress AND la2.successful = true),
            :since
        )
        """,
    )
    fun countFailedAttemptsAfterLastSuccessfulByUsernameAndIp(
        @Param("username") username: String,
        @Param("ipAddress") ipAddress: String,
        @Param("since") since: LocalDateTime,
    ): Long

    /**
     * 특정 IP 주소의 마지막 성공한 로그인 이후의 실패한 로그인 시도 횟수를 조회한다.
     * 성공한 로그인이 없는 경우 주어진 시간 이후의 모든 실패한 로그인 시도를 조회한다.
     *
     * @param ipAddress IP 주소
     * @param since 조회 시작 시간
     * @return 실패한 로그인 시도 횟수
     */
    @Query(
        """
        SELECT COUNT(la) FROM LoginAttempt la
        WHERE la.ipAddress = :ipAddress AND la.successful = false
        AND la.attemptAt > COALESCE(
            (SELECT MAX(la2.attemptAt) FROM LoginAttempt la2
             WHERE la2.ipAddress = :ipAddress AND la2.successful = true),
            :since
        )
        """,
    )
    fun countFailedAttemptsAfterLastSuccessfulByIp(
        @Param("ipAddress") ipAddress: String,
        @Param("since") since: LocalDateTime,
    ): Long
}
