# paketo-buildpacks
- BSN 내부 인프라에서 캐시 사용이 가능한 빌드팩 버전 정의
- `build.gradle.kts` 파일에서 빌더 버전 업데이트 시 사용하는 버전에 맞게 아래 항목을 오브젝트 스토리지에 캐시용 파일 업로드 필요

## 현재 버전

- builder: [0.4.254](https://github.com/paketo-buildpacks/builder-jammy-base/blob/v0.4.254/builder.toml)
    - java: [10.3.3](https://github.com/paketo-buildpacks/java/blob/v10.3.3/buildpack.toml)
        -
        bellsoft-liberica: [10.4.2](https://github.com/paketo-buildpacks/bellsoft-liberica/blob/v10.4.2/buildpack.toml)
            - jre:17.0.9 (`8129150fa39f1fdf6fdef6b05cf74ff570dae35bf540d2bdd9bf915532e12d55`)
        - spring-boot: [5.27.5](https://github.com/paketo-buildpacks/spring-boot/blob/v5.27.5/buildpack.toml)
            - spring-cloud-bindings:2.0.2 (`dab47967bffb29c5b4c41653e9741b6ca15cde483926d7f9d0a04956247fb680`)
        - syft: [1.39.0](https://github.com/paketo-buildpacks/syft/blob/v1.39.0/buildpack.toml)
            - syft:0.94.0 (`a18f10ba6add219b2680687450869db3c6a8b71e68ca6ae3925f9e53964cacbd`)
