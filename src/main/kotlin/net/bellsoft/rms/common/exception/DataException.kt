package net.bellsoft.rms.common.exception

import org.springframework.dao.DataAccessException
import org.springframework.dao.DataIntegrityViolationException

class DataNotFoundException(message: String = "존재하지 않는 데이터") : DataAccessException(message)
class PermissionRequiredDataException(message: String = "데이터 접근 권한 부족") : DataAccessException(message)
open class DuplicateDataException(message: String = "중복된 데이터") : DataIntegrityViolationException(message)
