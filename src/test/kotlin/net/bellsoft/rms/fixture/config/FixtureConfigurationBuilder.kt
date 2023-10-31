package net.bellsoft.rms.fixture.config

import com.appmattus.kotlinfixture.Fixture
import com.appmattus.kotlinfixture.config.ConfigurationBuilder

fun integratedFixture(configuration: ConfigurationBuilder.() -> Unit = {}): Fixture {
    val integratedConfigBuilder = ConfigurationBuilder().apply(configuration)

    integratedFixtureConfigurations.forEach(integratedConfigBuilder::apply)

    return Fixture(integratedConfigBuilder.build())
}
