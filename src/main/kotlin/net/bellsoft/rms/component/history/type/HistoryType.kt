package net.bellsoft.rms.component.history.type

import org.hibernate.envers.RevisionType

enum class HistoryType {
    CREATED,
    UPDATED,
    DELETED,
    ;

    companion object {
        fun of(revisionType: RevisionType) = when (revisionType) {
            RevisionType.ADD -> CREATED
            RevisionType.MOD -> UPDATED
            RevisionType.DEL -> DELETED
        }
    }
}
