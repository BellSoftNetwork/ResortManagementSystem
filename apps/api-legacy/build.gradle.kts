import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootBuildImage
import org.springframework.boot.gradle.tasks.run.BootRun

@Suppress("DSL_SCOPE_VIOLATION")
plugins {
    val kotlinVersion = libs.versions.kotlin.get()

    alias(libs.plugins.springBoot)
    alias(libs.plugins.springDependencyManagement)
    alias(libs.plugins.ktlint)
    alias(libs.plugins.liquibaseGradle)
    alias(libs.plugins.jacocoLog)

    kotlin("jvm") version kotlinVersion
    kotlin("plugin.spring") version kotlinVersion
    kotlin("plugin.jpa") version kotlinVersion
    kotlin("kapt") version kotlinVersion

    jacoco
    idea
}

jacoco {
    toolVersion = "0.8.8"
}

group = "net.bellsoft"
version = "0.0.1-SNAPSHOT"

java {
    sourceCompatibility = JavaVersion.VERSION_17
}

repositories {
    mavenCentral()
}

dependencies {
    // NOTE: Kotlin Support
    implementation(libs.kotlinReflect)
    implementation(libs.kotlinStdlibJdk8)

    // NOTE: Web
    implementation(libs.springBootStarterWeb)

    // NOTE: Database
    implementation(libs.springBootStarterDataJpa)
    implementation(libs.springBootStarterRedis)
    implementation(libs.springSessionDataRedis)
    implementation(libs.springDataEnvers)
    testRuntimeOnly(libs.h2database)
    runtimeOnly(libs.mysqlConnector)
    implementation(variantOf(libs.queryDslJpa) { classifier("jakarta") })
    implementation(libs.queryDslSql)
    kapt(variantOf(libs.queryDslApt) { classifier("jakarta") })
    kaptTest(variantOf(libs.queryDslApt) { classifier("jakarta") })

    implementation(libs.liquibase)
    liquibaseRuntime(libs.liquibase)
    liquibaseRuntime(libs.mysqlConnector)
    liquibaseRuntime("info.picocli:picocli:4.7.5")
    liquibaseRuntime("org.yaml:snakeyaml:1.33")
    liquibaseRuntime(sourceSets["main"].output)
    liquibaseRuntime(libs.kotlinLogging)

    // NOTE: Security
    implementation(libs.springBootStarterSecurity)
    implementation(libs.jjwtApi)
    runtimeOnly(libs.jjwtImpl)
    runtimeOnly(libs.jjwtJackson)
    testImplementation(libs.springSecurityTest)

    // NOTE: Validation
    implementation(libs.springBootStarterValidation)

    // NOTE: Data Process
    implementation(libs.jacksonKotlin)
    implementation(libs.jacksonDatabindNullable)
    implementation(libs.mapstruct)
    kapt(libs.mapstructProcessor)
    kaptTest(libs.mapstructProcessor)

    // NOTE: Data Support
    implementation(libs.dataFaker)
    implementation(libs.ulidCreator)
    implementation(libs.findbugsJsr305)

    // NOTE: Code Support
    implementation(libs.springBootStarterAop)

    // NOTE: Logging
    implementation(libs.kotlinLogging)
    implementation(libs.p6spy)

    // NOTE: Monitoring
    implementation(libs.springBootStarterActuator)

    // NOTE: Docs
    implementation(libs.springdocUi)

    // NOTE: Development Tools
    developmentOnly(libs.springBootDevTools)

    // NOTE: Test Support
    testImplementation(libs.springBootStarterTest) {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
        exclude(module = "mockito-core")
    }
    testImplementation(libs.bundles.kotest)
    testImplementation(libs.bundles.mock)
    testImplementation(libs.bundles.fixture)
}

liquibase {
    activities {
        val DB_JDBC_URL = System.getenv("DB_JDBC_URL")
        val DB_USERNAME = System.getenv("DB_USERNAME")
        val DB_PASSWORD = System.getenv("DB_PASSWORD")

        register("main") {
            arguments = mapOf(
                "logLevel" to "info",
                "changelogFile" to "db/changelog-master.yaml",
                "databaseChangeLogTableName" to "database_changelog",
                "databaseChangeLogLockTableName" to "database_changelog_lock",
                "url" to DB_JDBC_URL,
                "username" to DB_USERNAME,
                "password" to DB_PASSWORD,
            )
        }
    }
}

tasks.withType<KotlinCompile> {
    kotlinOptions {
        freeCompilerArgs += "-Xjsr305=strict"
        jvmTarget = "17"
    }
}

tasks.named<BootRun>("bootRun") {
    setupEnvironment()
}

tasks.named<org.springframework.boot.gradle.tasks.bundling.BootJar>("bootJar") {
    archiveFileName.set("app.jar")
    destinationDirectory.set(layout.buildDirectory.dir("libs"))
}

fun BootRun.setupEnvironment() {
    // Debugger configuration is handled by JAVA_OPTS environment variable
    // jvmArgs = listOf("-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005")
    environment("SPRING_PROFILES_ACTIVE", "local")
}

// https://docs.spring.io/spring-boot/docs/current/gradle-plugin/reference/htmlsingle/#build-image
tasks.named<BootBuildImage>("bootBuildImage") {
    setupDocker(this)
    setupBuildProperty(this)
    setupImageProperty(this)
    setupEnvironment(this)
}

fun setupDocker(bootBuildImage: BootBuildImage) {
    bootBuildImage.run {
        val dockerHost: String by project
        val isDockerTlsVerify: String by project
        val dockerCertPath: String by project

        val proxyRegistryUrl: String by project
        val projectRegistryUrl: String by project
        val registryUser: String by project
        val registryPassword: String by project
        val registryEmail: String by project

        docker {
            if (project.hasProperty("dockerHost")) host.set(dockerHost)
            if (project.hasProperty("isDockerTlsVerify")) tlsVerify.set(isDockerTlsVerify.toBoolean())
            if (project.hasProperty("dockerCertPath")) certPath.set(dockerCertPath)

            builderRegistry {
                if (project.hasProperty("proxyRegistryUrl")) url.set(proxyRegistryUrl)
                if (project.hasProperty("registryUser")) username.set(registryUser)
                if (project.hasProperty("registryPassword")) password.set(registryPassword)
                if (project.hasProperty("registryEmail")) email.set(registryEmail)
            }

            publishRegistry {
                if (project.hasProperty("projectRegistryUrl")) url.set(projectRegistryUrl)
                if (project.hasProperty("registryUser")) username.set(registryUser)
                if (project.hasProperty("registryPassword")) password.set(registryPassword)
                if (project.hasProperty("registryEmail")) email.set(registryEmail)
            }
        }
    }
}

fun setupBuildProperty(bootBuildImage: BootBuildImage) {
    bootBuildImage.run {
        val proxyRegistryUrl: String by project
        val bindingsDir: String by project
        val gradleDir: String by project

        val bindingVolumes = mutableListOf<String>()

        if (project.hasProperty("bindingsDir")) bindingVolumes.add("$bindingsDir:/platform/bindings:rw")
        if (project.hasProperty("gradleDir")) bindingVolumes.add("$gradleDir:/home/cnb/.gradle:rw")

        bindings.set(bindingVolumes)

        val buildPackBaseUrl = if (project.hasProperty("proxyRegistryUrl"))
            "$proxyRegistryUrl/paketobuildpacks"
        else
            "paketobuildpacks"

        runImage.set("$buildPackBaseUrl/run-jammy-base:latest")
        builder.set("$buildPackBaseUrl/builder-jammy-base:${paketobuildpacks.versions.builder.get()}")
    }
}

fun setupImageProperty(bootBuildImage: BootBuildImage) {
    bootBuildImage.run {
        val imagePath: String by project
        val imageBaseName: String by project
        val imageTag: String by project

        if (project.hasProperty("imagePath")) imageName.set(imagePath)
        if (project.hasProperty("imageBaseName")) imageName.set("${imageName.get()}/$imageBaseName")
        if (project.hasProperty("imageTag")) tags.set(mutableListOf("${imageName.get()}:$imageTag"))
    }
}

fun setupEnvironment(bootBuildImage: BootBuildImage) {
    bootBuildImage.run {
        environment.set(environment.get() + mapOf("SPRING_PROFILES_ACTIVE" to "production"))
    }
}

ktlint {
    version.set("0.48.2")
    verbose.set(true)
    relative.set(true)
    outputColorName.set("RED")
    enableExperimentalRules.set(true)
}

tasks.named<Test>("test") {
    useJUnitPlatform()

    jvmArgs(
        "--add-opens",
        "java.base/java.time=ALL-UNNAMED",
        "--add-opens",
        "java.base/java.lang.reflect=ALL-UNNAMED",
    )
}

tasks.jacocoTestReport {
    reports {
        xml.required.set(true)
        xml.outputLocation.set(file("${layout.buildDirectory.get()}/jacoco/jacoco.xml"))
    }

    finalizedBy("jacocoTestCoverageVerification") // NOTE: 활성화시 violationRules 통과 실패할경우 테스트도 실패처리 됨
}

private object JacocoViolationRuleSet {
    object Default {
        private val QUERY_DSL_DOMAINS = ('A'..'Z').map { "*.Q$it*" }

        val GENERATED_FILES = listOf(
            "*.mapper.*MapperImpl",
        ) + QUERY_DSL_DOMAINS
        val EXCLUDE_FILES = listOf(
            "*ApplicationKt",
            "*.config.*Config",
            "*.domain.*Converter",
            "*.exception.*Exception",
            "*.dto.*",
            "*Dto",
            "*DTO",
        )
    }

    object Business {
        val INCLUDE_FILES = listOf(
            "*.service.*",
        )
    }
}

tasks.jacocoTestCoverageVerification {
    violationRules {
        rule {
            // NOTE: element 가 없으면 프로젝트의 전체 파일을 합친 값 기준

            limit {
                // NOTE: counter 를 지정하지 않으면 default 는 INSTRUCTION
                // NOTE: value 를 지정하지 않으면 default 는 COVEREDRATIO
                minimum = "0.30".toBigDecimal()
            }
        }

        rule {
            enabled = true
            element = "CLASS"
            excludes = JacocoViolationRuleSet.Default.GENERATED_FILES

            // NOTE: 빈 줄을 제외한 코드의 라인수를 최대 200라인으로 제한
            limit {
                counter = "LINE"
                value = "TOTALCOUNT"
                maximum = "200".toBigDecimal()
            }
        }

        rule {
            enabled = true
            element = "CLASS"
            includes = JacocoViolationRuleSet.Business.INCLUDE_FILES
            excludes = JacocoViolationRuleSet.Default.GENERATED_FILES + JacocoViolationRuleSet.Default.EXCLUDE_FILES

            // NOTE: 브랜치 커버리지 최소 90% 만족
            limit {
                counter = "BRANCH"
                value = "COVEREDRATIO"
                minimum = "0.90".toBigDecimal()
            }

            // NOTE: 라인 커버리지 최소 80% 만족
            limit {
                counter = "LINE"
                value = "COVEREDRATIO"
                minimum = "0.80".toBigDecimal()
            }
        }
    }
}
