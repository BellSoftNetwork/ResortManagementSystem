package net.bellsoft.rms.reservation.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.reservation.entity.ReservationRoom

class ReservationRoomFixture {
    enum class Feature : FixtureFeature

    companion object {
        val BASE_CONFIGURATION = fixtureConfig {
            property(ReservationRoom::reservation) { baseFixture() }
            property(ReservationRoom::room) { baseFixture() }
        }
    }
}
