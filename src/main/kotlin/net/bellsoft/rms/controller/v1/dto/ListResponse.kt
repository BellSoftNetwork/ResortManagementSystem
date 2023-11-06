package net.bellsoft.rms.controller.v1.dto

import org.springframework.data.domain.Pageable

abstract class ListResponse<T>(
    open val pageable: Pageable,
    open val values: Collection<T>,
)
