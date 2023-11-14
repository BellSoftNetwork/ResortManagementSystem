package net.bellsoft.rms.fixture.service.room

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.service.room.dto.RoomCreateDto
import net.datafaker.Faker
import java.util.*

class RoomCreateDtoFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(RoomCreateDto::number) { FAKER.random().hex(10) }
        }
    }
}
