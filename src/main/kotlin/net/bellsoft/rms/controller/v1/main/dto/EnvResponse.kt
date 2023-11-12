package net.bellsoft.rms.controller.v1.main.dto

import java.time.Instant
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter
import java.time.format.DateTimeParseException
import java.util.*

data class EnvResponse(
    val applicationFullName: String,
    val applicationShortName: String,
    val commitSha: String,
    val commitShortSha: String,
    val commitTitle: String,
    val commitTimestamp: LocalDateTime,
) {
    companion object {
        fun of(
            applicationFullName: String,
            applicationShortName: String,
            commitSha: String,
            commitShortSha: String,
            commitTitle: String,
            commitTimestamp: String,
        ) = EnvResponse(
            applicationFullName = applicationFullName,
            applicationShortName = applicationShortName,
            commitSha = commitSha,
            commitShortSha = commitShortSha,
            commitTitle = commitTitle,
            commitTimestamp = convertTimestamp(commitTimestamp),
        )

        private fun convertTimestamp(timestamp: String) = try {
            LocalDateTime.from(
                Instant.from(
                    DateTimeFormatter.ISO_DATE_TIME.parse(timestamp),
                ).atZone(TimeZone.getDefault().toZoneId()),
            )
        } catch (_: DateTimeParseException) {
            LocalDateTime.of(1970, 1, 3, 0, 0)
        }
    }
}
