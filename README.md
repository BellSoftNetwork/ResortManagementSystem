# Resort Management System
리조트 통합 관리 시스템

### 개발 환경 설정

#### 백엔드

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

#### 프론트엔드

1. nvm 설치
2. node.js 22 버전 설치
3. yarn 설치 (`npm install -g yarn`)
4. 프로젝트 라이브러리 설치
    ```cmd
    cd apps/frontend-web
    yarn
    ```

### 개발환경 실행 방법

1. 백엔드 앱 실행 (프로필: `local` (기본값))
2. 프론트엔드 앱 실행
    ```cmd
    cd apps/frontend-web
    yarn dev
    ```
3. 브라우저 접속 (http://localhost:9000)  
   프론트엔드 앱에서 백엔드 API 호출 시 8080포트에 프록시를 태우도록 되어 있으므로 8080 포트로 접속할 필요 없음  
   (`apps\frontend-web\quasar.config.js` 파일 내 `devServer.proxy` 항목 참고)

### 모바일 앱 테스트 환경 설정

#### Android 앱 테스트

1. 필수 설치 항목
    - Android Studio 설치
    - Android SDK 설치 (Android Studio를 통해 설치 가능)
    - Java Development Kit (JDK) 설치

2. 환경 변수 설정
    - Android Studio 경로 설정 (Windows에서 필수)
   1. Windows 키 + R 키를 눌러 "실행" 창 열기
   2. `sysdm.cpl` 입력 후 엔터
   3. "고급" 탭 선택
   4. "환경 변수" 버튼 클릭
   5. "시스템 변수" 섹션에서 "새로 만들기" 클릭
   6. 변수 이름: `CAPACITOR_ANDROID_STUDIO_PATH`
   7. 변수 값: Android Studio 실행 파일의 전체 경로 입력 (JetBrains ToolBox 로 설치 시 `C:\Users\[사용자이름]\AppData\Local\Programs\Android Studio\bin\studio64.exe`)

3. 프로젝트 설정
   ```cmd
   cd apps/frontend-web
   # 이미 package.json에 @capacitor/android가 있으므로 다음 명령은 생략 가능
   # yarn add @capacitor/android     

   # Android 프로젝트 생성 (처음 한 번만)
   yarn cap:add-android

   # 웹 앱 빌드 후 Android 프로젝트에 동기화
   yarn android:build
   ```

4. Android 앱 실행 방법
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

5. 환경별 API 서버 설정
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

6. 문제 해결
    - `Unable to launch Android Studio. Is it installed?` 오류 발생 시:
- CAPACITOR_ANDROID_STUDIO_PATH 환경 변수가 올바르게 설정되었는지 확인
- 새 명령 프롬프트/터미널을 열어 다시 시도 (환경 변수 변경 후 필수)
- 직접 Android Studio를 실행하고 `[프로젝트 경로]/apps/frontend-web/android` 폴더를 열기

7. 실제 기기에서 테스트하기
    - USB 디버깅 활성화 (Android 기기의 개발자 옵션에서 설정)
    - 기기를 컴퓨터에 USB로 연결
    - Android Studio에서 연결된 기기 선택 후 실행
    - 또는 위의 환경별 명령어 사용 (`yarn android:run:dev` 또는 `yarn android:run:prod`)

#### 로컬 환경 실행 팁

인텔리제이에서 백엔드 실행 구성과 프론트엔드 실행 구성 생성 후 복합 실행 구성으로 두 구성을 동시에 실행하게 설정하면 편함

### IntelliJ 필수 플러그인 설치
- `Kotest`
- `JPA Buddy`
  이외 추가 플러그인은 BSN 위키 내 [JetBrains IDE 추천 플러그인](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/recommended-plugins) 문서 참고

### IntelliJ 설정
BSN 위키 내 [IntelliJ 기본 설정](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/intellij-basic-settings) 문서 참고
