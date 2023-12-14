package net.bellsoft.rms.room.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.room.dto.service.RoomGroupCreateDto
import net.datafaker.Faker
import java.util.*

class RoomGroupCreateDtoFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(RoomGroupCreateDto::name) { FAKER.random().hex(20) }
        }
    }
}
