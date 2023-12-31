package net.bellsoft.rms.user.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.user.dto.service.UserCreateDto
import net.bellsoft.rms.user.type.UserRole
import net.datafaker.Faker
import java.util.*

class UserCreateDtoFixture {
    enum class Feature : FixtureFeature {
        NORMAL {
            override fun config() = fixtureConfig {
                property(UserCreateDto::role) { UserRole.NORMAL }
            }
        },
        ADMIN {
            override fun config() = fixtureConfig {
                property(UserCreateDto::role) { UserRole.ADMIN }
            }
        },
        SUPER_ADMIN {
            override fun config() = fixtureConfig {
                property(UserCreateDto::role) { UserRole.SUPER_ADMIN }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(UserCreateDto::name) { "name-${FAKER.random().hex(10)}" }
            property(UserCreateDto::userId) { "userId-${FAKER.random().hex(5)}" }
            property(UserCreateDto::email) { "${FAKER.random().hex(5)}-${FAKER.internet().emailAddress()}" }
            property(UserCreateDto::password) { FAKER.random().hex(10).toString() }
        }
    }
}
