# Resort Management System 프로젝트 구성

## 목차
1. [기술 스택](#기술-스택)
2. [프로젝트 아키텍처](#프로젝트-아키텍처)
3. [개발 원칙](#개발-원칙)
4. [코드 스타일](#코드-스타일)
5. [명명 규칙](#명명-규칙)
6. [API 패턴](#api-패턴)
7. [보안 가이드라인](#보안-가이드라인)
8. [에러 처리](#에러-처리)
9. [환경 설정 및 배포](#환경-설정-및-배포)
10. [모니터링 및 로깅](#모니터링-및-로깅)
11. [문서화](#문서화)
12. [다국어 지원](#다국어-지원)

## 기술 스택

### 인프라
- 온프레미스 환경에 Kubernetes 클러스터에서 운영
- 고가용성 및 확장성을 위한 컨테이너 기반 아키텍처
- 무중단 배포를 위한 Blue-Green 배포 전략 적용
- Helm 차트를 통한 애플리케이션 배포 관리

### 데이터베이스
- MySQL 8.0 데이터베이스 사용 (UTF-8mb4 인코딩)
- Liquibase를 통한 데이터베이스 스키마 마이그레이션 관리
  - 변경 이력 추적 및 롤백 지원
  - 환경별 데이터베이스 스키마 일관성 유지
- Redis를 사용한 세션 관리 및 캐싱
  - 분산 환경에서의 세션 공유
  - 성능 향상을 위한 데이터 캐싱

### 백엔드
- Kotlin 1.8+ 및 Spring Boot 3.x 프레임워크 기반
- Spring Data JPA를 통한 데이터 액세스
- Spring Security를 활용한 인증 및 권한 관리
- JWT 기반 인증 시스템 사용
- Kotest를 이용한 테스트 주도 개발(TDD) 적용
- Fixture 패턴을 활용한 테스트 데이터 생성
- Spring Actuator를 통한 애플리케이션 상태 모니터링

### 프론트엔드
- TypeScript 5.x + Vue.js 3.x + Quasar 2.x 프레임워크 기반
- Vue 3 Composition API 사용
- Pinia 상태 관리 라이브러리 사용
- dayjs 날짜 처리 라이브러리 사용
- lodash 유틸리티 라이브러리 사용
- 컴포넌트 기반 UI 아키텍처
- Capacitor를 통한 모바일 앱 지원

## 프로젝트 아키텍처

### 백엔드 아키텍처
- 도메인 주도 설계(DDD) 패턴 적용
- 계층형 아키텍처 사용:
  - Controller: API 엔드포인트 정의
  - Service: 비즈니스 로직 처리
  - Repository: 데이터 액세스 로직
  - Entity: 도메인 모델
  - DTO: 데이터 전송 객체
- 인터페이스와 구현체 분리 패턴 사용 (예: AuthController 인터페이스와 AuthControllerImpl 구현체)
- 의존성 주입을 통한 컴포넌트 결합도 낮춤
- Envers를 활용한 엔티티 변경 이력 추적

### 프론트엔드 아키텍처
- 컴포넌트 기반 아키텍처
- API 클라이언트 모듈을 통한 백엔드 통신
- 스키마 기반 데이터 모델링
- 재사용 가능한 컴포넌트 분리
- 라우터 기반 페이지 구성
- 반응형 디자인으로 모바일 및 데스크톱 지원

## 개발 원칙

### 코드 품질
- 클린 코드와 클린 아키텍처 원칙 준수
- SOLID 원칙 적용
- 코드 중복 최소화 (DRY 원칙)
- 단일 책임 원칙(SRP) 준수
- 가독성과 유지보수성 우선

### 테스트 전략
- 테스트 주도 개발(TDD) 지향
- 중요 기능은 반드시 유닛 테스트 작성
- 통합 테스트를 통한 컴포넌트 간 상호작용 검증
- kotlinfixture 라이브러리를 활용한 테스트 데이터 생성
- 각 도메인별 전용 Fixture 클래스 사용

### 성능 최적화
- 데이터베이스 쿼리 최적화
- N+1 문제 방지를 위한 적절한 페치 전략 사용
- 캐싱 전략 적용
- 불필요한 API 호출 최소화
- 대용량 데이터 처리 시 페이징 적용

## 코드 스타일
- 각 언어별 코드 린터를 통과하도록 작성
- 3번 이상 반복되는 항목은 추출하여 재사용 가능하도록 구성
- 주석 없이도 이해할 수 있는 명확한 코드 작성
- 주석보다는 명확한 함수명과 변수명 사용
- 복잡한 로직은 작은 단위의 함수로 분리

## 명명 규칙

### 백엔드
- 패키지명: 소문자 사용, 도메인 기반 구조 (예: net.bellsoft.rms.authentication)
- 클래스명: PascalCase 사용 (예: AuthController, JwtTokenProvider)
- 함수명: camelCase 사용 (예: registerUser, login)
- 변수명: camelCase 사용 (예: userService, jwtTokenProvider)
- 상수명: UPPER_SNAKE_CASE 사용
- 테스트 클래스: 테스트 대상 클래스명 + Test (예: UserServiceTest)

### 프론트엔드
- 컴포넌트명: PascalCase 사용 (예: ReservationEditor, RoomGroupSelector)
- 함수명: camelCase 사용 (예: loadPaymentMethods, recalculateStayEndAt)
- 변수명: camelCase 사용 (예: formModel, selectedRooms)
- 상수명: UPPER_SNAKE_CASE 사용
- 파일명: 컴포넌트는 PascalCase, 유틸리티는 kebab-case 사용

## API 패턴

### 백엔드
- RESTful API 설계 원칙 준수
- 응답 형식 통일: SingleResponse 패턴 사용
- 적절한 HTTP 상태 코드 사용
- 예외 처리 및 에러 응답 일관성 유지
- JWT 기반 인증 적용
- API 버전 관리 (v1, v2 등)
- OpenAPI(Swagger)를 통한 API 문서화

### 프론트엔드
- API 클라이언트 함수 모듈화
- 비동기 요청에 Promise 패턴 사용
- 에러 핸들링 일관성 유지
- 로딩 상태 관리
- 재시도 메커니즘 구현
- 토큰 만료 시 자동 갱신

## 보안 가이드라인
- 모든 API 엔드포인트에 적절한 인증 및 권한 검사 적용
- 사용자 입력 데이터 검증 및 이스케이프 처리
- CSRF 방어 메커니즘 적용
- 민감한 정보는 암호화하여 저장
- 로그인 시도 제한을 통한 무차별 대입 공격 방어
- 디바이스 변경 감지 및 추가 인증 요구

## 에러 처리

### 백엔드
- 커스텀 예외 클래스 사용
- 예외 발생 시 적절한 로깅
- 클라이언트에게 의미 있는 에러 메시지 제공
- 민감한 오류 정보는 클라이언트에 노출하지 않음

### 프론트엔드
- try-catch 블록을 통한 에러 처리
- 사용자 친화적인 에러 메시지 표시
- Quasar Notify 컴포넌트를 통한 알림 표시
- 네트워크 오류 시 적절한 재시도 메커니즘 구현

## 환경 설정 및 배포

### 환경 변수 관리
- 환경별 설정은 환경 변수를 통해 관리 (개발, 테스트, 운영)
- 민감한 정보(DB 비밀번호, API 키 등)는 환경 변수로 주입
- 기본값은 application.yaml에 정의하되, 실제 값은 환경 변수로 오버라이드
- 환경 변수 목록:
  - DATABASE_MYSQL_HOST: MySQL 데이터베이스 호스트
  - DATABASE_MYSQL_USER: MySQL 사용자 이름
  - DATABASE_MYSQL_PASSWORD: MySQL 비밀번호
  - DATABASE_MYSQL_SCHEMA: 데이터베이스 스키마 (기본값: rms)
  - REDIS_HOST: Redis 호스트 (기본값: 127.0.0.1)
  - REDIS_PORT: Redis 포트 (기본값: 6379)
  - REDIS_PASSWORD: Redis 비밀번호
  - REDIS_NAMESPACE: Redis 네임스페이스 (기본값: production.v1)
  - JWT_SECRET: JWT 서명 비밀키
  - JWT_ACCESS_TOKEN_VALIDITY: 액세스 토큰 유효 시간(시간) (기본값: 1)
  - JWT_REFRESH_TOKEN_VALIDITY: 리프레시 토큰 유효 시간(시간) (기본값: 720)

### CI/CD 파이프라인
- GitLab CI/CD를 통한 자동화된 빌드 및 배포
- 브랜치 기반 배포 전략:
  - develop 브랜치: 개발 환경 자동 배포
  - staging 브랜치: 테스트 환경 자동 배포
  - main 브랜치: 운영 환경 수동 승인 후 배포
- 배포 전 자동화된 테스트 실행
- 배포 정보 환경 변수 주입:
  - DEPLOY_COMMIT_SHA: 배포된 커밋 해시
  - DEPLOY_COMMIT_SHORT_SHA: 짧은 커밋 해시
  - DEPLOY_COMMIT_TITLE: 커밋 제목
  - DEPLOY_COMMIT_TIMESTAMP: 커밋 타임스탬프

### 컨테이너화
- Dockerfile을 통한 애플리케이션 컨테이너화
- 멀티스테이지 빌드를 통한 이미지 크기 최적화
- 컨테이너 헬스 체크 엔드포인트 구성
- 무상태(Stateless) 설계로 수평적 확장 지원

## 모니터링 및 로깅

### 애플리케이션 모니터링
- Spring Actuator를 통한 애플리케이션 상태 모니터링
- 헬스 체크 엔드포인트: /actuator/health
- 리브니스 프로브: /actuator/health/liveness
- 레디니스 프로브: /actuator/health/readiness
- 프로메테우스 메트릭 수집 및 Grafana 대시보드 구성

### 로깅 전략
- 구조화된 JSON 형식 로그 사용
- 로그 레벨 관리:
  - TRACE: 상세한 디버깅 정보 (트랜잭션 인터셉터 등)
  - DEBUG: 개발 환경에서의 디버깅 정보
  - INFO: 일반적인 애플리케이션 이벤트 (기본 레벨)
  - WARN: 잠재적 문제 상황
  - ERROR: 오류 및 예외 상황
- 중요 비즈니스 이벤트는 INFO 레벨로 로깅
- 예외 발생 시 ERROR 레벨로 로깅하며 스택 트레이스 포함
- 민감한 정보(개인정보, 비밀번호 등)는 로그에 포함하지 않음

### 성능 모니터링
- 데이터베이스 쿼리 성능 모니터링
- API 응답 시간 측정 및 모니터링
- 리소스 사용량(CPU, 메모리) 모니터링
- 병목 현상 식별 및 최적화

## 문서화

### API 문서화
- OpenAPI(Swagger)를 통한 API 문서 자동화
- API 문서 접근 경로:
  - API 스키마: /docs/schema
  - Swagger UI: /docs/swagger-ui
- 각 API 엔드포인트에 대한 상세 설명 제공
- 요청/응답 예시 포함
- 인증 방식 및 권한 요구사항 명시

### 코드 문서화
- 주요 클래스 및 함수에 KDoc/JavaDoc 스타일 주석 작성
- 복잡한 비즈니스 로직에 대한 설명 주석 추가
- README 파일을 통한 프로젝트 개요 및 설정 방법 제공
- 아키텍처 다이어그램 및 시퀀스 다이어그램 제공

## 다국어 지원
- 한국어 기본 지원
- UI 텍스트는 한국어로 작성
- 향후 다국어 지원을 위한 구조 준비
- 국제화(i18n) 리소스 번들 사용 준비
- 날짜, 시간, 통화 등의 로케일 기반 포맷팅 지원
