package net.bellsoft.rms.fixture.domain.reservation

import net.bellsoft.rms.domain.reservation.ReservationRoom
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig

class ReservationRoomFixture {
    enum class Feature : FixtureFeature

    companion object {
        val BASE_CONFIGURATION = fixtureConfig {
            property(ReservationRoom::reservation) { baseFixture() }
            property(ReservationRoom::room) { baseFixture() }
        }
    }
}
