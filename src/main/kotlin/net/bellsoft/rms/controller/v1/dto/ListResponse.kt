package net.bellsoft.rms.controller.v1.dto

abstract class ListResponse(
    val page: Int,
    val perPage: Int,
    val totalPages: Int?,
    val totalEntries: Int?,
)
