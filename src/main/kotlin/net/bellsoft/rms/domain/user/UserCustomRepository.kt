package net.bellsoft.rms.domain.user

interface UserCustomRepository {
    fun getYearCreatedUsers(year: Int): List<User>
}
