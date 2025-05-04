package net.bellsoft.rms.authentication.fixture

import net.bellsoft.rms.authentication.entity.LoginAttempt
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.time.LocalDateTime
import java.util.*

class LoginAttemptFixture {
    enum class Feature : FixtureFeature {
        SUCCESSFUL {
            override fun config() = fixtureConfig {
                property(LoginAttempt::successful) { true }
            }
        },
        FAILED {
            override fun config() = fixtureConfig {
                property(LoginAttempt::successful) { false }
            }
        },
        PAST {
            override fun config() = fixtureConfig {
                property(LoginAttempt::attemptAt) { LocalDateTime.now().minusHours(1) }
            }
        },
        RECENT {
            override fun config() = fixtureConfig {
                property(LoginAttempt::attemptAt) { LocalDateTime.now().minusMinutes(5) }
            }
        },
        WITH_DEVICE_FINGERPRINT {
            override fun config() = fixtureConfig {
                property(LoginAttempt::deviceFingerprint) { "device-fingerprint-${FAKER.random().hex(10)}" }
            }
        },
        WITHOUT_DEVICE_FINGERPRINT {
            override fun config() = fixtureConfig {
                property(LoginAttempt::deviceFingerprint) { null }
            }
        },
        PAST_FAILED {
            override fun config() = fixtureConfig {
                property(LoginAttempt::successful) { false }
                property(LoginAttempt::attemptAt) { LocalDateTime.now().minusHours(1) }
            }
        },
        CURRENT_FAILED {
            override fun config() = fixtureConfig {
                property(LoginAttempt::successful) { false }
                property(LoginAttempt::attemptAt) { LocalDateTime.now() }
            }
        },
        CURRENT_SUCCESSFUL {
            override fun config() = fixtureConfig {
                property(LoginAttempt::successful) { true }
                property(LoginAttempt::attemptAt) { LocalDateTime.now() }
            }
        },
        WITH_SPECIFIC_USERNAME {
            override fun config() = fixtureConfig {
                property(LoginAttempt::username) { "specific-user" }
            }
        },
        WITH_SPECIFIC_IP {
            override fun config() = fixtureConfig {
                property(LoginAttempt::ipAddress) { "127.0.0.1" }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(LoginAttempt::username) { "user-${FAKER.random().hex(5)}" }
            property(LoginAttempt::ipAddress) { FAKER.internet().ipV4Address() }
            property(LoginAttempt::successful) { FAKER.random().nextBoolean() }
            property(LoginAttempt::attemptAt) { LocalDateTime.now() }
            property(LoginAttempt::deviceFingerprint) {
                if (FAKER.random().nextBoolean()) "device-${
                    FAKER.random().hex(8)
                }" else null
            }
        }
    }
}
