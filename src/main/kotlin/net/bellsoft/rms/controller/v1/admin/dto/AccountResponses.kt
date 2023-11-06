package net.bellsoft.rms.controller.v1.admin.dto

import net.bellsoft.rms.controller.v1.dto.ListResponse
import net.bellsoft.rms.domain.user.User
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

data class AccountResponses(
    override val pageable: Pageable,
    val totalPages: Int,
    val totalElements: Long,
    override val values: List<AccountResponse>,
) : ListResponse<AccountResponse>(pageable = pageable, values = values) {
    companion object {
        fun of(userPage: Page<User>): AccountResponses {
            val map = userPage.get().map { AccountResponse.of(it) }
            userPage.totalElements
            return AccountResponses(
                pageable = userPage.pageable,
                totalPages = userPage.totalPages,
                totalElements = userPage.totalElements,
                values = map.toList(),
            )
        }
    }
}
