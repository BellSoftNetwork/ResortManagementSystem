package net.bellsoft.rms.user.service

import net.bellsoft.rms.common.dto.response.EntityListDto
import net.bellsoft.rms.user.dto.response.UserDetailDto
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.dto.service.UserPatchDto
import net.bellsoft.rms.user.entity.User
import org.springframework.data.domain.Pageable

interface UserService {
    fun register(userCreateDto: UserCreateDto): UserDetailDto

    fun findAll(pageable: Pageable): EntityListDto<UserDetailDto>

    fun updateAccount(updateUserId: Long, userPatchDto: UserPatchDto): UserDetailDto

    fun isUpdatableAccount(requestUser: User, targetUserId: Long): Boolean
}
