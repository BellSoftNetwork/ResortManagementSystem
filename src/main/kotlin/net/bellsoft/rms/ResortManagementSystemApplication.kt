package net.bellsoft.rms

import net.bellsoft.rms.annotation.ExcludeFromJacocoGeneratedReport
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class ResortManagementSystemApplication

@ExcludeFromJacocoGeneratedReport
fun main(args: Array<String>) {
    runApplication<ResortManagementSystemApplication>(*args)
}
