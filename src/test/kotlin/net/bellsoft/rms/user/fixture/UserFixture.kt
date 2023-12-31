package net.bellsoft.rms.user.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.user.entity.User
import net.bellsoft.rms.user.type.UserRole
import net.bellsoft.rms.user.type.UserStatus
import net.datafaker.Faker
import java.util.*

class UserFixture {
    enum class Feature : FixtureFeature {
        NORMAL {
            override fun config() = fixtureConfig {
                property(User::name) { "normal-${FAKER.random().hex(10)}" }
                property(User::role) { UserRole.NORMAL }
                property(User::status) { UserStatus.ACTIVE }
            }
        },
        ADMIN {
            override fun config() = fixtureConfig {
                property(User::name) { "admin-${FAKER.random().hex(10)}" }
                property(User::role) { UserRole.ADMIN }
                property(User::status) { UserStatus.ACTIVE }
            }
        },
        SUPER_ADMIN {
            override fun config() = fixtureConfig {
                property(User::name) { "super_admin-${FAKER.random().hex(5)}" }
                property(User::role) { UserRole.SUPER_ADMIN }
                property(User::status) { UserStatus.ACTIVE }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(User::userId) { "userId-${FAKER.random().hex(5)}" }
            property(User::email) { "${FAKER.random().hex(5)}-${FAKER.internet().emailAddress()}" }
            property(User::name) { "name-${FAKER.random().hex(10)}" }
        }
    }
}
