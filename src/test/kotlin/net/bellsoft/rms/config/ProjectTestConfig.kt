package net.bellsoft.rms.config

import io.kotest.core.config.AbstractProjectConfig
import io.kotest.core.spec.IsolationMode
import io.kotest.core.test.AssertionMode
import io.kotest.extensions.spring.SpringTestExtension
import io.kotest.extensions.spring.SpringTestLifecycleMode
import io.mockk.clearAllMocks

object ProjectTestConfig : AbstractProjectConfig() {
    override val parallelism = 1
    override val assertionMode = AssertionMode.Error
    override val failOnIgnoredTests = true
    override val isolationMode = IsolationMode.InstancePerLeaf
    override val testNameRemoveWhitespace = true
    override var displayFullTestPath: Boolean? = true
    override val failOnEmptyTestSuite = true

    override suspend fun afterProject() {
        super.afterProject()
        clearAllMocks()
    }

    override fun extensions() = listOf(SpringTestExtension(SpringTestLifecycleMode.Root))
}
