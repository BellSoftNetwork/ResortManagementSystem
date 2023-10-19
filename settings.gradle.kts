rootProject.name = "ResortManagementSystem"

dependencyResolutionManagement {
    versionCatalogs {
        create("paketobuildpacks") {
            from(files("./gradle/paketobuildpacks.versions.toml"))
        }
    }
}
