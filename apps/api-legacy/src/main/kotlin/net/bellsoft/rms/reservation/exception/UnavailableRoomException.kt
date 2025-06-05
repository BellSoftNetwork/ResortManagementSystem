package net.bellsoft.rms.reservation.exception

import net.bellsoft.rms.common.exception.DuplicateDataException

class UnavailableRoomException(message: String = "배정 불가 객실") : DuplicateDataException(message)
