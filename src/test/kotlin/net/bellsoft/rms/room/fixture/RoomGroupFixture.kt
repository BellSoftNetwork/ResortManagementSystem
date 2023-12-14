package net.bellsoft.rms.room.fixture

import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.room.entity.RoomGroup
import net.datafaker.Faker
import java.util.*

class RoomGroupFixture {
    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(RoomGroup::name) { FAKER.random().hex(20) }
            property(RoomGroup::description) { FAKER.random().hex(100) }
            property(RoomGroup::rooms) { mutableListOf() }
        }
    }
}
