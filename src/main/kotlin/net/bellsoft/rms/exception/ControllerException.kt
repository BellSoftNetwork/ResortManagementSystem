package net.bellsoft.rms.exception

import org.springframework.dao.DataIntegrityViolationException

class BadRequestException(message: String) : Exception(message)
class UnprocessableEntityException(message: String) : DataIntegrityViolationException(message)
