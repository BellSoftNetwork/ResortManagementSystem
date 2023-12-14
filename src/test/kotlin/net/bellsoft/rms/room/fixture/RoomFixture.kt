package net.bellsoft.rms.room.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.baseFixture
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.room.entity.Room
import net.bellsoft.rms.room.type.RoomStatus
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
            property(Room::roomGroup) { baseFixture() }
        }
    }
}
