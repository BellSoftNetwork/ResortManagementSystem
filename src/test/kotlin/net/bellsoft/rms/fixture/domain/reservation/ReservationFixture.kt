package net.bellsoft.rms.fixture.domain.reservation

import net.bellsoft.rms.domain.reservation.Reservation
import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class ReservationFixture {
    enum class Feature : FixtureFeature {
        CANCEL {
            override fun config() = fixtureConfig {
                property(Reservation::status) { ReservationStatus.CANCEL }
            }
        },
        NORMAL {
            override fun config() = fixtureConfig {
                property(Reservation::status) { ReservationStatus.NORMAL }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(Reservation::user) { baseFixture() }
            property(Reservation::reservationMethod) { baseFixture() }
            property(Reservation::room) { baseFixture() }
            property(Reservation::name) { "name-${FAKER.random().hex(10)}" }
            property(Reservation::phone) { FAKER.phoneNumber().phoneNumber() }
        }
    }
}
