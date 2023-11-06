package net.bellsoft.rms.exception

import org.springframework.dao.DataAccessException

class DataNotFoundException(message: String = "존재하지 않는 데이터") : DataAccessException(message)
class PermissionRequiredDataException(message: String = "데이터 접근 권한 부족") : DataAccessException(message)
