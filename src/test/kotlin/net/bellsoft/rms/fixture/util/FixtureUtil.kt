package net.bellsoft.rms.fixture.util

import com.appmattus.kotlinfixture.Fixture
import com.appmattus.kotlinfixture.config.ConfigurationBuilder
import net.bellsoft.rms.fixture.FixtureFeature

fun fixtureConfig(
    vararg baseConfigurations: ConfigurationBuilder.() -> Unit,
    overrideConfiguration: ConfigurationBuilder.() -> Unit = {},
): ConfigurationBuilder.() -> Unit = {
    baseConfigurations.forEach(::apply)

    apply(overrideConfiguration)
}

infix fun (Fixture).feature(feature: FixtureFeature): Fixture {
    return this.new(feature.config())
}

fun (Fixture).features(vararg features: FixtureFeature): Fixture {
    return this.new {
        features.forEach { apply(it.config()) }
    }
}
