package net.bellsoft.rms.main.service

import net.bellsoft.rms.main.dto.response.AppConfigDto

interface ConfigService {
    fun getAppConfig(): AppConfigDto
}
