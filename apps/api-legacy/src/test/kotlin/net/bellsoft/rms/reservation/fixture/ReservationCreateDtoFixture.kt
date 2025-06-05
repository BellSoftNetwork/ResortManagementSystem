package net.bellsoft.rms.reservation.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.reservation.dto.service.ReservationCreateDto
import net.bellsoft.rms.reservation.type.ReservationStatus
import net.datafaker.Faker
import java.util.*

class ReservationCreateDtoFixture {
    enum class Feature : FixtureFeature {
        CANCEL {
            override fun config() = fixtureConfig {
                property(ReservationCreateDto::status) { ReservationStatus.CANCEL }
            }
        },
        NORMAL {
            override fun config() = fixtureConfig {
                property(ReservationCreateDto::status) { ReservationStatus.NORMAL }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(ReservationCreateDto::name) { "name-${FAKER.random().hex(10)}" }
            property(ReservationCreateDto::rooms) { emptySet() }
            property(ReservationCreateDto::phone) { FAKER.phoneNumber().phoneNumber() }
        }
    }
}
