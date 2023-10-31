package net.bellsoft.rms.config

import io.kotest.core.config.AbstractProjectConfig
import io.kotest.core.spec.IsolationMode
import io.kotest.core.test.AssertionMode
import io.kotest.extensions.spring.SpringExtension
import io.mockk.clearAllMocks

object ProjectTestConfig : AbstractProjectConfig() {
    override val parallelism = 1
    override val assertionMode = AssertionMode.Warn
    override val failOnIgnoredTests = false
    override val isolationMode = IsolationMode.SingleInstance
    override val testNameRemoveWhitespace = true
    override var displayFullTestPath: Boolean? = true

    override suspend fun afterProject() {
        super.afterProject()
        clearAllMocks()
    }

    override fun extensions() = listOf(SpringExtension)
}
