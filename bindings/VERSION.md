# paketo-buildpacks
- BSN 내부 인프라에서 캐시 사용이 가능한 빌드팩 버전 정의
- `build.gradle.kts` 파일에서 빌더 버전 업데이트 시 사용하는 버전에 맞게 아래 항목을 오브젝트 스토리지에 캐시용 파일 업로드 필요

## 현재 버전
- base-builder: [0.3.44](https://github.com/paketo-buildpacks/base-builder/blob/v0.3.44/builder.toml)
  - java: [6.36.0](https://github.com/paketo-buildpacks/java/blob/v6.36.0/buildpack.toml)
    - bellsoft-liberica: [9.5.0](https://github.com/paketo-buildpacks/bellsoft-liberica/blob/v9.5.0/buildpack.toml)
      - jre:17.0.4 (`04d5a8540a5ed57af493137be85ce996ebb3afc86390f685c0f922a9c6fa1d62`)
    - spring-boot: [5.16.0](https://github.com/paketo-buildpacks/spring-boot/blob/v5.16.0/buildpack.toml)
      - spring-cloud-bindings:1.10.0 (`1f5b781f8bd0d6b85ab2462e4b98d36782a2227fef5b168db174b3959a0ebebe`)
    - syft: [1.17.0](https://github.com/paketo-buildpacks/syft/blob/v1.17.0/buildpack.toml)
      - syft:0.54.0 (`00794a0eafe3e86cf10e096fffd8f983551fc8585937deb4fc44203635f0b2ca`)
