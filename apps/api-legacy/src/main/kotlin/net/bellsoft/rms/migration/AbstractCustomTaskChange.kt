package net.bellsoft.rms.migration

import liquibase.change.custom.CustomTaskChange
import liquibase.database.Database
import liquibase.database.jvm.JdbcConnection
import liquibase.exception.CustomChangeException
import liquibase.exception.ValidationErrors
import liquibase.resource.ResourceAccessor
import mu.KLogging

/** ## liquibase CustomTaskChange 실행 절차
 * 1. set Variables (Optional)
 * 2. call `setFileOpener()`
 * 3. call `setUp()`
 * 4. call `validate()`
 * 5. call `execute()`
 * 6. call `getConfirmationMessage()`
 */
abstract class AbstractCustomTaskChange : CustomTaskChange {
    private var resourceAccessor: ResourceAccessor? = null
    private var isExecuted = false

    override fun setFileOpener(resourceAccessor: ResourceAccessor) {
        this.resourceAccessor = resourceAccessor
    }

    override fun setUp() {
        isExecuted = false
    }

    override fun execute(database: Database?) {
        if (isExecuted)
            return

        isExecuted = true
        val connection = getJdbcConnection(database ?: throw CustomChangeException("Database is null"))

        try {
            migrate(connection)
            connection.commit()
        } catch (ex: Exception) {
            connection.rollback()
            throw CustomChangeException(ex)
        }
    }

    abstract fun migrate(connection: JdbcConnection)

    override fun validate(database: Database?) = ValidationErrors()

    override fun getConfirmationMessage() = "[${this::class.simpleName}] Successfully custom change executed"

    private fun getJdbcConnection(database: Database): JdbcConnection {
        return database.connection as JdbcConnection
    }

    companion object : KLogging()
}
