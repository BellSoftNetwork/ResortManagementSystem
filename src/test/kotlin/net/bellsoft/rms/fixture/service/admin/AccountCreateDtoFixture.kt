package net.bellsoft.rms.fixture.service.admin

import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.service.admin.dto.AccountCreateDto
import net.datafaker.Faker
import java.util.*

class AccountCreateDtoFixture {
    enum class Feature : FixtureFeature {
        NORMAL {
            override fun config() = fixtureConfig {
                property(AccountCreateDto::role) { UserRole.NORMAL }
            }
        },
        ADMIN {
            override fun config() = fixtureConfig {
                property(AccountCreateDto::role) { UserRole.ADMIN }
            }
        },
        SUPER_ADMIN {
            override fun config() = fixtureConfig {
                property(AccountCreateDto::role) { UserRole.SUPER_ADMIN }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(AccountCreateDto::name) { "name-${FAKER.random().hex(10)}" }
            property(AccountCreateDto::email) { "${FAKER.random().hex(5)}-${FAKER.internet().emailAddress()}" }
            property(AccountCreateDto::password) { FAKER.random().hex(10).toString() }
        }
    }
}
