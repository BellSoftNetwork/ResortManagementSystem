package net.bellsoft.rms.fixture.controller.v1.admin

import net.bellsoft.rms.controller.v1.admin.dto.AccountCreateRequest
import net.bellsoft.rms.domain.user.UserRole
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class AccountCreateRequestFixture {
    enum class Feature : FixtureFeature {
        NORMAL {
            override fun config() = fixtureConfig {
                property(AccountCreateRequest::role) { UserRole.NORMAL }
            }
        },
        ADMIN {
            override fun config() = fixtureConfig {
                property(AccountCreateRequest::role) { UserRole.ADMIN }
            }
        },
        SUPER_ADMIN {
            override fun config() = fixtureConfig {
                property(AccountCreateRequest::role) { UserRole.SUPER_ADMIN }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(AccountCreateRequest::name) { "name-${FAKER.random().hex(10)}" }
            property(AccountCreateRequest::email) { "${FAKER.random().hex(5)}-${FAKER.internet().emailAddress()}" }
            property(AccountCreateRequest::password) { FAKER.random().hex(10).toString() }
        }
    }
}
