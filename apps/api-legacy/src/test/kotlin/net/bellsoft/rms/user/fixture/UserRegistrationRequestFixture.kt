package net.bellsoft.rms.user.fixture

import net.bellsoft.rms.authentication.dto.request.UserRegistrationRequest
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class UserRegistrationRequestFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(UserRegistrationRequest::userId) { "userId-${FAKER.random().hex(5)}" }
            property(UserRegistrationRequest::email) { "${FAKER.random().hex(5)}-${FAKER.internet().emailAddress()}" }
            property(UserRegistrationRequest::name) { "name-${FAKER.random().hex(10)}" }
        }
    }
}
