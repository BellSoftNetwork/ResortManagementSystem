package net.bellsoft.rms.reservation.type

enum class ReservationStatus(val value: Int) {
    REFUND(-10),
    CANCEL(-1),
    PENDING(0),
    NORMAL(1),
}
