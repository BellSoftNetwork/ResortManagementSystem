package net.bellsoft.rms.main.service

import net.bellsoft.rms.main.dto.response.AppConfigDto
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ConfigService(
    private val userRepository: UserRepository,
) {
    fun getAppConfig() = AppConfigDto(
        isAvailableRegistration = userRepository.count() < 1,
    )
}
