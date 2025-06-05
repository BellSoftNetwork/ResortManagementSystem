package net.bellsoft.rms.user.repository

import net.bellsoft.rms.user.entity.User

interface UserCustomRepository {
    fun getYearCreatedUsers(year: Int): List<User>
}
