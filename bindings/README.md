## bootBuildImage 태스크에 의존성 캐시 등록 방법

1. 각 빌드팩 의존성의 `buildpack.toml` 파일에 있는 `id` 값으로 디렉토리를 생성
2. `type` 파일을 생성하여 내용에 `dependency-mapping` 작성
3. 해당 의존성의 `buildpack.toml` 파일에 정의된 `sha256` 값을 위에서 생성한 디렉토리 내에 파일명으로 생성하고 내용에는 `uri` 값을 등록

### example

#### 1. 베이스 빌더의 버전 정의 확인
- 베이스 빌더 `builder.toml` 파일: https://github.com/paketo-buildpacks/base-builder/blob/main/builder.toml

```toml
[[buildpacks]]
  uri = "docker://gcr.io/paketo-buildpacks/java:6.36.0"
  version = "6.36.0"
```
- 위 예제에서 자바 빌드팩 버전이 `6.36.0`


#### 2. 해당 빌드팩 버전의 정의 확인
- Java 빌드팩의 `buildpack.toml` 파일 (`6.36.0`): https://github.com/paketo-buildpacks/java/blob/v6.36.0/buildpack.toml

```toml
[[order]]

  [[order.group]]
    id = "paketo-buildpacks/bellsoft-liberica"
    version = "9.5.0"
```
- `6.36.0` 버전 자바 빌드팩에서 사용하는 의존성 중 `bellsoft-liberica` 버전이 `9.5.0`


#### 3. 해당 의존성 내의 파일 확인
- bellsoft liberica `buildpack.toml` 파일 (`9.5.0`): https://github.com/paketo-buildpacks/bellsoft-liberica/blob/v9.5.0/buildpack.toml

```toml
[[metadata.dependencies]]
  cpes = ["cpe:2.3:a:oracle:jre:17.0.4:*:*:*:*:*:*:*"]
  id = "jre"
  name = "BellSoft Liberica JRE"
  purl = "pkg:generic/bellsoft-jre@17.0.4?arch=amd64"
  sha256 = "04d5a8540a5ed57af493137be85ce996ebb3afc86390f685c0f922a9c6fa1d62"
  stacks = ["io.buildpacks.stacks.bionic", "io.paketo.stacks.tiny", "*"]
  uri = "https://github.com/bell-sw/Liberica/releases/download/17.0.4.1+1/bellsoft-jre17.0.4.1+1-linux-amd64.tar.gz"
  version = "17.0.4"
```
- `9.5.0` 버전 bellsoft liberica 에서 사용하는 JRE 파일 정의 확인

1. 파일 내용에 `id` 값인 `jre`을 디렉토리명으로 생성
2. `jre/type` 파일에 `dependency-mapping` 작성
3. `jre/04d5a8540a5ed57af493137be85ce996ebb3afc86390f685c0f922a9c6fa1d62` 파일에 대체 다운로드할 URI 작성
