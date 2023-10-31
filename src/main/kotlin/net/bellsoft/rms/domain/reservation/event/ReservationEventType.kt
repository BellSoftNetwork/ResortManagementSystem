package net.bellsoft.rms.domain.reservation.event

enum class ReservationEventType(val value: Int) {
    REFUND(-10),
    CANCEL(-1),
    PENDING(0),
    NORMAL(1),
}
