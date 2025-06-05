rootProject.name = "ResortManagementSystem"

include("apps:api-legacy")

dependencyResolutionManagement {
    versionCatalogs {
        create("paketobuildpacks") {
            from(files("./gradle/paketobuildpacks.versions.toml"))
        }
    }
}
