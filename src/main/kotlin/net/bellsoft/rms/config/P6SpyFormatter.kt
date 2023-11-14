package net.bellsoft.rms.config

import com.p6spy.engine.logging.Category
import com.p6spy.engine.spy.P6SpyOptions
import com.p6spy.engine.spy.appender.MessageFormattingStrategy
import jakarta.annotation.PostConstruct
import org.hibernate.engine.jdbc.internal.FormatStyle
import org.springframework.context.annotation.Configuration

@Configuration
class P6SpyFormatter : MessageFormattingStrategy {
    @PostConstruct
    fun setLogMessageFormat() {
        P6SpyOptions.getActiveInstance().logMessageFormat = this.javaClass.getName()
    }

    override fun formatMessage(
        connectionId: Int,
        now: String,
        elapsed: Long,
        category: String,
        prepared: String,
        sql: String,
        url: String,
    ) = String.format("[%s] | %d ms | %s", category, elapsed, formatSql(category, sql))

    private fun formatSql(category: String, sql: String?): String? {
        if (sql == null || sql.trim { it <= ' ' }.isEmpty() || Category.STATEMENT.name != category)
            return sql

        val trimmedSQL = sql.trim { it <= ' ' }.lowercase()

        return if (
            trimmedSQL.startsWith("create") ||
            trimmedSQL.startsWith("alter") ||
            trimmedSQL.startsWith("comment")
        ) {
            FormatStyle.DDL.formatter.format(sql)
        } else {
            FormatStyle.BASIC.formatter.format(sql)
        }
    }
}
