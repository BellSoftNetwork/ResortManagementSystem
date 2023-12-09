package net.bellsoft.rms.payment.fixture

import net.bellsoft.rms.fixture.FixtureFeature
import net.bellsoft.rms.fixture.util.fixtureConfig
import net.bellsoft.rms.payment.entity.PaymentMethod
import net.datafaker.Faker
import java.util.*

class PaymentMethodFixture {
    enum class Feature : FixtureFeature

    companion object {
        private val FAKER = Faker(Locale.KOREA)

        val BASE_CONFIGURATION = fixtureConfig {
            property(PaymentMethod::name) { "name-${FAKER.random().hex(10)}" }
        }
    }
}
