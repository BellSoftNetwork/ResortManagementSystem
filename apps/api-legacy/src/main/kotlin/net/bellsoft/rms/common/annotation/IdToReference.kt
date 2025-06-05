package net.bellsoft.rms.common.annotation

import org.mapstruct.Qualifier

@Qualifier
@Target(AnnotationTarget.FUNCTION)
@Retention(AnnotationRetention.BINARY)
annotation class IdToReference
