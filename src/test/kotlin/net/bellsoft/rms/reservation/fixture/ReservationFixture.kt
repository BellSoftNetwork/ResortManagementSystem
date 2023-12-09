package net.bellsoft.rms.reservation.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.reservation.entity.Reservation
import net.bellsoft.rms.reservation.type.ReservationStatus
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
            property(Reservation::paymentMethod) { baseFixture() }
            property(Reservation::rooms) { mutableListOf() }
            property(Reservation::name) { "name-${FAKER.random().hex(10)}" }
            property(Reservation::phone) { FAKER.phoneNumber().phoneNumber() }
        }
    }
}
