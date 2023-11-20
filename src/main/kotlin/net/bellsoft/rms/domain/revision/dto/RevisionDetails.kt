package net.bellsoft.rms.domain.revision.dto

import net.bellsoft.rms.domain.revision.RevisionInfo
import org.hibernate.envers.RevisionType

data class RevisionDetails<T>(
    val entity: T,
    val revisionInfo: RevisionInfo,
    val revisionType: RevisionType,
    val modifiedFields: HashSet<String>,
) {
    companion object {
        @Suppress("UNCHECKED_CAST")
        fun <T : Any> of(result: Any?): RevisionDetails<T> {
            check(result is Array<*> && result.isArrayOf<Any>()) { "result must be Array<Any>" }
            require(result.size == 4) { "result must be Array<Any> of size 4 (size: ${result.size})" }

            return RevisionDetails(
                entity = result[0] as T,
                revisionInfo = result[1] as RevisionInfo,
                revisionType = result[2] as RevisionType,
                modifiedFields = result[3] as HashSet<String>,
            )
        }
    }
}
