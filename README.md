# Resort Management System
리조트 통합 관리 시스템

## 🚀 Docker 기반 개발 환경 (권장)

모든 개발 환경을 Docker로 통합하여 어느 PC에서든 동일한 환경을 즉시 구축할 수 있습니다.

### 빠른 시작

#### 🚀 개발 환경 시작
```shell
# 커스텀 설정이 필요한 경우 (.env 파일 생성)
cp .env.example .env
```

```shell
# Docker 개발 환경 시작 (.env 파일 없이도 기본값으로 동작)
docker compose up -d
```

#### 🛑 개발 환경 중지 (CPU 절약)
```shell
# 모든 컨테이너 중지 (데이터는 유지됨)
docker compose stop
```

#### ▶️ 개발 환경 재시작
```shell
# 중지된 컨테이너 다시 시작
docker compose start
```

#### 🧹 개발 환경 완전 제거
```shell
# 컨테이너 및 네트워크 제거 (볼륨은 유지)
docker compose down
```

```shell
# 데이터 포함 완전 제거 (주의!)
docker compose down -v
```

#### 📊 상태 확인
```shell
# 실행 중인 컨테이너 확인
docker compose ps
```

#### 📋 로그 확인
```shell
# 모든 서비스 로그
docker compose logs -f
```

```shell
# 특정 서비스 로그
docker compose logs -f api-core
```

접속 URL:
- Frontend: http://localhost:9000 (api-core 연동)
- API Core (Go): http://localhost:8080 (메인 API)
- API Legacy (Spring Boot): http://localhost:8081

### 포트 구성

| 서비스 | 설명 | 호스트 포트 | 컨테이너 포트 | 접속 URL |
|--------|------|-------------|----------------|----------|
| MySQL | 데이터베이스 | 3306 | 3306 | `localhost:3306` |
| Redis | 캐시 서버 | 6379 | 6379 | `localhost:6379` |
| API Core (Go) | 새 API (메인) | 8080 | 8080 | `http://localhost:8080` |
| API Legacy (Spring) | 기존 API | 8081 | 8080 | `http://localhost:8081` |
| Frontend | Vue.js 앱 | 9000 | 9000 | `http://localhost:9000` |

### 테스트 실행

#### 🧪 전체 테스트 실행
```shell
# 모든 서비스의 유닛 테스트 실행
./scripts/dev-test.sh
```

#### 🎯 특정 서비스 테스트
```shell
# Go API 테스트
./scripts/dev-test.sh --service api-core

# Spring Boot API 테스트
./scripts/dev-test.sh --service api-legacy

# Frontend 테스트
./scripts/dev-test.sh --service frontend
```

#### 🔍 통합 테스트 실행
```shell
# 모든 통합 테스트
./scripts/dev-test.sh --type integration

# 특정 서비스의 모든 테스트 (유닛 + 통합)
./scripts/dev-test.sh --service api-core --type all
```

### 개별 서비스 접속

#### 🔧 컨테이너 내부 접속
```shell
# Go API 컨테이너 접속
docker compose exec api-core bash

# Spring Boot API 컨테이너 접속
docker compose exec api-legacy bash

# Frontend 컨테이너 접속
docker compose exec frontend sh
```

### 기본 접속 정보 (로컬 개발)

- **MySQL**: 
  - Host: localhost
  - Port: 3306
  - Database: rms-legacy (api-legacy), rms-core (api-core)
  - User: rms
  - Password: rms123
  - Root Password: root

- **Redis**:
  - Host: localhost
  - Port: 6379

- **JWT Secret**: local-development-secret

> 💡 **참고**: 모든 로컬 개발 환경 설정은 기본값이 내장되어 있어 .env 파일 없이도 즉시 실행 가능합니다. 커스텀 설정이 필요한 경우 .env.example을 참고하여 .env 파일을 생성하세요.

---

## 💻 로컬 개발 환경 설정 (Docker 없이)

Docker를 사용하지 않고 로컬에서 직접 실행하는 경우의 설정입니다.

### 백엔드 설정

1. IntelliJ 에서 jdk 17 다운로드 후 `gradle refresh` 실행
2. `gradle ktlintApplyToIdea` 실행
3. `gradle addKtlintCheckGitPreCommitHook` 실행
4. MySQL 서버 개별 설치
5. Redis 서버 개별 설치
6. 로컬 실행 구성에 필수 환경 변수 설정
    - `DATABASE_MYSQL_HOST`: MySQL 서버 IP
    - `DATABASE_MYSQL_USER`: MySQL 계정 ID
    - `DATABASE_MYSQL_PASSWORD`: MySQL 계정 비밀번호
    - `REDIS_HOST`: Redis 서버 IP

### 프론트엔드 설정

1. nvm 설치
2. node.js 20 버전 설치 (기존 문서는 22버전이었으나 프로젝트는 20버전 사용)
3. yarn 설치 (`npm install -g yarn`)
4. 프로젝트 라이브러리 설치
    ```cmd
    cd apps/frontend-web
    yarn
    ```

### 로컬 실행 방법

1. 백엔드 앱 실행 (프로필: `local` (기본값))
2. 프론트엔드 앱 실행
    ```cmd
    cd apps/frontend-web
    yarn dev
    ```
3. 브라우저 접속 (http://localhost:9000)  
   프론트엔드 앱에서 백엔드 API 호출 시 8080포트에 프록시를 태우도록 되어 있으므로 8080 포트로 접속할 필요 없음  
   (`apps\frontend-web\quasar.config.js` 파일 내 `devServer.proxy` 항목 참고)

---

## 📱 모바일 앱 테스트 환경 설정

### Android 앱 테스트

1. **필수 설치 항목**
    - Android Studio 설치
    - Android SDK 설치 (Android Studio를 통해 설치 가능)
    - Java Development Kit (JDK) 설치

2. **환경 변수 설정**
    - Android Studio 경로 설정 (Windows에서 필수)
   1. Windows 키 + R 키를 눌러 "실행" 창 열기
   2. `sysdm.cpl` 입력 후 엔터
   3. "고급" 탭 선택
   4. "환경 변수" 버튼 클릭
   5. "시스템 변수" 섹션에서 "새로 만들기" 클릭
   6. 변수 이름: `CAPACITOR_ANDROID_STUDIO_PATH`
   7. 변수 값: Android Studio 실행 파일의 전체 경로 입력 (JetBrains ToolBox 로 설치 시 `C:\Users\[사용자이름]\AppData\Local\Programs\Android Studio\bin\studio64.exe`)

3. **프로젝트 설정**
   ```cmd
   cd apps/frontend-web
   # 이미 package.json에 @capacitor/android가 있으므로 다음 명령은 생략 가능
   # yarn add @capacitor/android     

   # Android 프로젝트 생성 (처음 한 번만)
   yarn cap:add-android

   # 웹 앱 빌드 후 Android 프로젝트에 동기화
   yarn android:build
   ```

4. **Android 앱 실행 방법**
   ```cmd
   cd apps/frontend-web
   # Android Studio에서 프로젝트 열기
   yarn cap:open-android

   # 또는 한 번에 빌드하고 Android Studio 열기 (에뮬레이터용)
   yarn android:run

   # 개발 환경 API 서버를 사용하여 실제 기기에서 테스트
   yarn android:run:dev

   # 프로덕션 환경 API 서버를 사용하여 실제 기기에서 테스트
   yarn android:run:prod
   ```

5. **환경별 API 서버 설정**
   - 로컬 개발 환경 (에뮬레이터): `http://10.0.2.2:8080` (에뮬레이터에서 호스트 PC의 localhost를 가리키는 특수 IP)
   - 로컬 개발 환경 (실제 기기): 개발 환경 API 서버 주소 사용 (기본값: `http://dev-api.example.com`)
   - 개발 환경 배포: 개발 환경 API 서버 주소 사용
   - 프로덕션 환경 배포: 프로덕션 API 서버 주소 사용 (기본값: `https://api.example.com`)

   환경 변수를 통해 API 서버 주소 설정:
   ```cmd
   # 개발 환경 API 서버 주소 설정
   set VITE_API_URL_DEV=http://dev-api.example.com

   # 프로덕션 환경 API 서버 주소 설정
   set VITE_API_URL_PROD=https://api.example.com

   # 특정 API 서버 주소 직접 설정 (모든 환경에 적용)
   set VITE_API_URL=http://specific-api.example.com
   ```

6. **문제 해결**
    - `Unable to launch Android Studio. Is it installed?` 오류 발생 시:
- CAPACITOR_ANDROID_STUDIO_PATH 환경 변수가 올바르게 설정되었는지 확인
- 새 명령 프롬프트/터미널을 열어 다시 시도 (환경 변수 변경 후 필수)
- 직접 Android Studio를 실행하고 `[프로젝트 경로]/apps/frontend-web/android` 폴더를 열기

7. **실제 기기에서 테스트하기**
    - USB 디버깅 활성화 (Android 기기의 개발자 옵션에서 설정)
    - 기기를 컴퓨터에 USB로 연결
    - Android Studio에서 연결된 기기 선택 후 실행
    - 또는 위의 환경별 명령어 사용 (`yarn android:run:dev` 또는 `yarn android:run:prod`)

---

## 🛠 개발 도구

### IntelliJ 필수 플러그인
- `Kotest`
- `JPA Buddy`
  이외 추가 플러그인은 BSN 위키 내 [JetBrains IDE 추천 플러그인](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/recommended-plugins) 문서 참고

### IntelliJ 설정
BSN 위키 내 [IntelliJ 기본 설정](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/intellij-basic-settings) 문서 참고

### 로컬 환경 실행 팁
인텔리제이에서 백엔드 실행 구성과 프론트엔드 실행 구성 생성 후 복합 실행 구성으로 두 구성을 동시에 실행하게 설정하면 편함

---

## 📚 추가 문서

- [프로젝트 구조 및 개발 가이드](./CLAUDE.md)
- [API 문서](./docs/api/)
