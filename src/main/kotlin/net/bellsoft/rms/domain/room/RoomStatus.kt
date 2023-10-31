package net.bellsoft.rms.domain.room

enum class RoomStatus(val value: Int) {
    DAMAGED(-10),
    CONSTRUCTION(-1),
    INACTIVE(0),
    NORMAL(1),
}
