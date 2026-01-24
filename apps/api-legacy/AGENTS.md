# api-legacy 개발 가이드

Kotlin + Spring Boot 기반 API 서버 (레거시, 마이그레이션 예정)

## 아키텍처

도메인 주도 설계 (DDD) 기반:

```
src/main/kotlin/
├── config/         # Spring 설정
├── controller/     # REST 컨트롤러
├── domain/         # 도메인 모델
├── dto/            # 데이터 전송 객체
├── repository/     # JPA 리포지토리
├── service/        # 비즈니스 로직
└── security/       # 인증/인가
```

## 개발 명령어

```bash
# Docker 컨테이너 접속
docker compose exec api-legacy bash

# 컨테이너 내부에서 실행
./gradlew bootRun         # 애플리케이션 실행 (프로필: local)
./gradlew test           # 테스트 실행
./gradlew jacocoTestReport  # 테스트 커버리지 리포트
./gradlew ktlintCheck    # 린트 체크
./gradlew ktlintFormat   # 린트 포맷
./gradlew bootBuildImage # Docker 이미지 빌드
```

## DB 마이그레이션 (Liquibase)

변경셋 위치: `src/main/resources/db/changelog/`

Spring Boot 앱 시작 시 자동으로 Liquibase가 실행됩니다.

## 코드 품질 표준

### JaCoCo 커버리지

- 전체: 최소 30%
- 비즈니스 로직 서비스: 80% 라인 커버리지, 90% 브랜치 커버리지
- 클래스당 최대 200줄

### 린팅

- Ktlint 준수 필수
- pre-commit 훅으로 자동 검사

## Hibernate Envers

변경 이력 추적 테이블:

- `revision_info`: 리비전 정보
- `*_history`: 엔티티별 히스토리

RevisionType:
- 0 = ADD (CREATED)
- 1 = MOD (UPDATED)
- 2 = DEL (DELETED)

상세: [Hibernate Envers 참조](../../docs/references/hibernate-envers.md)

## 환경 변수

```bash
# 데이터베이스
DATABASE_MYSQL_HOST=mysql
DATABASE_MYSQL_PORT=3306
DATABASE_MYSQL_USER=root
DATABASE_MYSQL_PASSWORD=root
DATABASE_MYSQL_DATABASE=rms-legacy

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
```

## 마이그레이션 상태

이 앱은 api-core로 마이그레이션 예정입니다.
마이그레이션 스펙: [specs/migration/](../../specs/migration/)

### 주의사항

- 새로운 기능은 가급적 api-core에 구현
- api-legacy 변경 시 api-core 호환성 확인 필요
