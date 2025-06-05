package net.bellsoft.rms.authentication.fixture

import net.bellsoft.rms.authentication.dto.DeviceInfoDto
import net.bellsoft.rms.authentication.util.AuthTestConstants
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class DeviceInfoFixture {
    enum class Feature : FixtureFeature {
        WITH_LOCAL_IP {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::ipAddress) { "127.0.0.1" }
            }
        },
        WITH_INTERNAL_IP {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::ipAddress) { "192.168.1.1" }
            }
        },
        WITH_HACKER_IP {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::ipAddress) { "18.18.18.18" }
            }
        },
        WITH_WINDOWS_DEVICE_FINGERPRINT {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::deviceFingerprint) { "win" }
            }
        },
        WITH_ANDROID_DEVICE_FINGERPRINT {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::deviceFingerprint) { "and" }
            }
        },
        WITHOUT_DEVICE_FINGERPRINT {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::deviceFingerprint) { "" }
            }
        },
        WITH_WINDOWS_OS {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::userAgent) { AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT }
                property(DeviceInfoDto::osInfo) { "Windows" }
            }
        },
        WITH_KOREAN_LANGUAGE {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::languageInfo) { AuthTestConstants.DEFAULT_KOREAN_LANGUAGE }
            }
        },
        WITH_ANDROID_USER_AGENT {
            override fun config() = fixtureConfig {
                property(DeviceInfoDto::userAgent) { AuthTestConstants.DEFAULT_ANDROID_USER_AGENT }
                property(DeviceInfoDto::osInfo) { "Android" }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(DeviceInfoDto::ipAddress) { FAKER.internet().ipV4Address() }
            property(DeviceInfoDto::deviceFingerprint) { "device-${FAKER.random().hex(8)}" }
            property(DeviceInfoDto::osInfo) { "Windows" }
            property(DeviceInfoDto::languageInfo) { AuthTestConstants.DEFAULT_KOREAN_LANGUAGE }
            property(DeviceInfoDto::userAgent) { AuthTestConstants.DEFAULT_WINDOWS_USER_AGENT }
        }
    }
}
