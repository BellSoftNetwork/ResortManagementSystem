package net.bellsoft.rms.fixture.domain.room.event

import net.bellsoft.rms.domain.room.event.RoomEvent
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class RoomEventFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(RoomEvent::user) { baseFixture() }
            property(RoomEvent::room) { baseFixture() }
        }
    }
}
