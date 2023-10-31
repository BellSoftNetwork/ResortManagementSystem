package net.bellsoft.rms.fixture

import com.appmattus.kotlinfixture.decorator.nullability.AlwaysNullStrategy
import com.appmattus.kotlinfixture.decorator.nullability.NeverNullStrategy
import com.appmattus.kotlinfixture.decorator.nullability.RandomlyNullStrategy
import com.appmattus.kotlinfixture.decorator.nullability.nullabilityStrategy
import com.appmattus.kotlinfixture.decorator.optional.AlwaysOptionalStrategy
import com.appmattus.kotlinfixture.decorator.optional.NeverOptionalStrategy
import com.appmattus.kotlinfixture.decorator.optional.RandomlyOptionalStrategy
import com.appmattus.kotlinfixture.decorator.optional.optionalStrategy
import net.bellsoft.rms.fixture.config.integratedFixture

@Suppress("ktlint:experimental:property-naming")
val baseFixture = integratedFixture {
    // NOTE: https://github.com/appmattus/kotlinfixture/blob/main/fixture/configuration-options.adoc#overriding-nullability-with-nullabilitystrategy
    nullabilityStrategy(RandomlyNullStrategy)

    // NOTE: https://github.com/appmattus/kotlinfixture/blob/main/fixture/configuration-options.adoc#overriding-the-use-of-default-values-with-optionalstrategy
    optionalStrategy(RandomlyOptionalStrategy)
}

@Suppress("ktlint:experimental:property-naming")
val baseNullFixture = baseFixture.new {
    nullabilityStrategy(AlwaysNullStrategy)
    optionalStrategy(AlwaysOptionalStrategy)
}

@Suppress("ktlint:experimental:property-naming")
val baseNotNullFixture = baseFixture.new {
    nullabilityStrategy(NeverNullStrategy)
    optionalStrategy(NeverOptionalStrategy)
}
