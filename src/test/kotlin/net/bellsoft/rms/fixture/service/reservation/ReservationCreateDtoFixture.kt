package net.bellsoft.rms.fixture.service.reservation

import net.bellsoft.rms.domain.reservation.ReservationStatus
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.service.reservation.dto.ReservationCreateDto
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
