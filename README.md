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
2. node.js 18 버전 설치
3. yarn 설치 (`npm install -g yarn`)
4. 프로젝트 라이브러리 설치
    ```cmd
    cd view
    yarn
    ```

### 개발환경 실행 방법

1. 백엔드 앱 실행 (프로필: `local` (기본값))
2. 프론트엔드 앱 실행
    ```cmd
    cd view
    yarn dev
    ```
3. 브라우저 접속 (http://localhost:9000)  
   프론트엔드 앱에서 백엔드 API 호출 시 8080포트에 프록시를 태우도록 되어 있으므로 8080 포트로 접속할 필요 없음  
   (`view\quasar.config.js` 파일 내 `devServer.proxy` 항목 참고)

#### 로컬 환경 실행 팁

인텔리제이에서 백엔드 실행 구성과 프론트엔드 실행 구성 생성 후 복합 실행 구성으로 두 구성을 동시에 실행하게 설정하면 편함

### IntelliJ 필수 플러그인 설치
- `Kotest`
- `JPA Buddy`
  이외 추가 플러그인은 BSN 위키 내 [JetBrains IDE 추천 플러그인](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/recommended-plugins) 문서 참고

### IntelliJ 설정
BSN 위키 내 [IntelliJ 기본 설정](https://wiki.bellsoft.net/ko/dev-1-team/development-environment-settings/tip/jetbrains-ide/intellij-basic-settings) 문서 참고
