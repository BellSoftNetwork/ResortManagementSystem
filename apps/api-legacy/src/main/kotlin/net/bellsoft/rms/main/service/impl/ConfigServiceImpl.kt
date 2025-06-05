package net.bellsoft.rms.main.service.impl

import net.bellsoft.rms.main.dto.response.AppConfigDto
import net.bellsoft.rms.main.service.ConfigService
import net.bellsoft.rms.user.repository.UserRepository
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
@Transactional(readOnly = true)
class ConfigServiceImpl(
    private val userRepository: UserRepository,
) : ConfigService {
    override fun getAppConfig() = AppConfigDto(
        isAvailableRegistration = userRepository.count() < 1,
    )
}
