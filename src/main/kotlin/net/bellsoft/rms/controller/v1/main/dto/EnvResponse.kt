package net.bellsoft.rms.controller.v1.main.dto

data class EnvResponse(
    val applicationFullName: String,
    val applicationShortName: String,
    val version: String,
) {
    companion object {
        fun of(applicationFullName: String, applicationShortName: String, version: String) =
            EnvResponse(applicationFullName, applicationShortName, version)
    }
}
