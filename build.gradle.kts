// Root project build file
// This is a simple wrapper that delegates to the apps/api-legacy subproject

plugins {
    base
    idea
}

// Define the root project group and version
group = "net.bellsoft"
version = "0.0.1-SNAPSHOT"

// Define repositories for all projects
repositories {
    mavenCentral()
}

// Make the root project depend on the subprojects
dependencies {
    // No dependencies at the root level
}

// Define tasks that delegate to the backend subproject
tasks.register("bootRun") {
    dependsOn(":apps:api-legacy:bootRun")
    description = "Runs the Spring Boot application (delegated to apps:api-legacy)"
}

tasks.register("test") {
    dependsOn(":apps:api-legacy:test")
    description = "Runs the tests (delegated to apps:api-legacy)"
}

tasks.register("bootBuildImage") {
    dependsOn(":apps:api-legacy:bootBuildImage")
    description = "Builds the Docker image (delegated to apps:api-legacy)"
}

tasks.register("ktlintCheck") {
    dependsOn(":apps:api-legacy:ktlintCheck")
    description = "Runs ktlint check (delegated to apps:api-legacy)"
}

tasks.register("jacocoTestReport") {
    dependsOn(":apps:api-legacy:jacocoTestReport")
    description = "Generates JaCoCo test report (delegated to apps:api-legacy)"
}
