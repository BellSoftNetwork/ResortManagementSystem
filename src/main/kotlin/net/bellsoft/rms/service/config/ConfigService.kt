package net.bellsoft.rms.service.config

import net.bellsoft.rms.domain.user.UserRepository
import net.bellsoft.rms.service.config.dto.AppConfigDto
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
