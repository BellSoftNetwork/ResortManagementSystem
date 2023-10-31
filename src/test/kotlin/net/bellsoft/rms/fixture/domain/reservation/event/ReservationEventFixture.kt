package net.bellsoft.rms.fixture.domain.reservation.event

import net.bellsoft.rms.domain.reservation.event.ReservationEvent
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class ReservationEventFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(ReservationEvent::user) { baseFixture() }
            property(ReservationEvent::reservation) { baseFixture() }
        }
    }
}
