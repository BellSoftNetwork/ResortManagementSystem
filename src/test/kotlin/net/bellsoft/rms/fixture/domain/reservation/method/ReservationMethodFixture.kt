package net.bellsoft.rms.fixture.domain.reservation.method

import net.bellsoft.rms.domain.reservation.method.ReservationMethod
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class ReservationMethodFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(ReservationMethod::name) { "name-${FAKER.random().hex(10)}" }
        }
    }
}
