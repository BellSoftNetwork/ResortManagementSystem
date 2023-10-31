package net.bellsoft.rms.fixture.domain.room

import net.bellsoft.rms.domain.room.Room
import net.bellsoft.rms.domain.room.RoomStatus
import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.datafaker.Faker
import java.util.*

class RoomFixture {
    enum class Feature : FixtureFeature {
        INACTIVE {
            override fun config() = fixtureConfig {
                property(Room::status) { RoomStatus.INACTIVE }
            }
        },
        NORMAL {
            override fun config() = fixtureConfig {
                property(Room::status) { RoomStatus.NORMAL }
            }
        },
    }

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(Room::number) { FAKER.random().hex(10) }
        }
    }
}
