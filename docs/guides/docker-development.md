# Docker 기반 개발 환경 가이드

> 모든 개발 작업은 반드시 Docker 기반으로 실행. 로컬 환경 직접 실행 금지.

---

## 빠른 시작

### 개발 환경 시작

```bash
# 환경 변수 설정 (필요시)
cp .env.example .env

# Docker 개발 환경 시작 (.env 없이도 기본값으로 동작)
docker compose up -d
```

### 개발 환경 중지

```bash
# 모든 컨테이너 중지 (데이터 유지)
docker compose stop

# 컨테이너 재시작
docker compose start
```

### 완전 제거

```bash
# 컨테이너 및 네트워크 제거 (볼륨 유지)
docker compose down

# 데이터 포함 완전 제거 (주의!)
docker compose down -v
```

### 상태 및 로그 확인

```bash
docker compose ps              # 실행 중인 컨테이너 확인
docker compose logs -f         # 모든 서비스 로그
docker compose logs -f api-core  # 특정 서비스 로그
```

---

## 서비스 구성

### 포트 구성

| 서비스 | 호스트 포트 | 컨테이너 포트 | 접속 URL |
|--------|-------------|----------------|----------|
| MySQL | 3306 | 3306 | `localhost:3306` |
| Redis | 6379 | 6379 | `localhost:6379` |
| API Core (Go) | 8080 | 8080 | `http://localhost:8080` |
| API Legacy (Spring) | 8081 | 8080 | `http://localhost:8081` |
| Frontend | 9000 | 9000 | `http://localhost:9000` |

### 기본 접속 정보

**MySQL**:
- Host: localhost (컨테이너 내부: `mysql`)
- Port: 3306
- Database: `rms-legacy` (api-legacy), `rms-core` (api-core)
- User: rms / Password: rms123
- Root Password: root

**Redis**:
- Host: localhost (컨테이너 내부: `redis`)
- Port: 6379

---

## 컨테이너 내부 작업

### Go API (api-core)

```bash
# 컨테이너 접속
docker compose exec api-core bash

# 컨테이너 내부에서 실행
make dev              # 개발 모드
make test            # 테스트
make lint            # 린트
make build           # 빌드
```

### Spring Boot API (api-legacy)

```bash
# 컨테이너 접속
docker compose exec api-legacy bash

# 컨테이너 내부에서 실행
./gradlew bootRun         # 실행
./gradlew test           # 테스트
./gradlew ktlintCheck    # 린트
```

### Frontend (Vue.js)

```bash
# 컨테이너 접속
docker compose exec frontend sh

# 컨테이너 내부에서 실행
yarn dev            # 개발 서버
yarn build          # 프로덕션 빌드
yarn lint           # 린트
```

---

## 개별 서비스 개발

```bash
# DB만 실행하고 API는 로컬에서 개발 (권장하지 않음)
docker compose up -d mysql redis

# API만 실행
docker compose up -d mysql redis api-core

# Frontend만 실행 (API 필요)
docker compose up -d api-core
docker compose up frontend
```

---

## 테스트 실행

### 스크립트 사용

```bash
# 모든 서비스 테스트
./scripts/dev-test.sh

# 특정 서비스 테스트
./scripts/dev-test.sh --service api-core
./scripts/dev-test.sh --service api-legacy
./scripts/dev-test.sh --service frontend

# 통합 테스트
./scripts/dev-test.sh --type integration
```

### 직접 실행

```bash
docker compose exec api-core make test
docker compose exec api-legacy ./gradlew test
docker compose exec frontend yarn test
```

---

## Docker Compose 구성 원칙

### 완전한 독립성

- 호스트 OS에 의존하지 않음
- macOS, Windows, Linux에서 동일하게 작동

### 개발 편의 기능

- 핫 리로드 (Air for Go, Vite for Vue)
- 볼륨 마운트로 코드 변경 즉시 반영

### 12-Factor App 준수

- 환경 변수로 설정 관리
- 상태를 저장하지 않는 프로세스
- 수평 확장 가능한 구조

---

## 문제 해결

### 포트 충돌

```bash
# 사용 중인 포트 확인
lsof -i :8080
lsof -i :3306

# 해당 프로세스 종료 후 다시 시작
docker compose down && docker compose up -d
```

### 볼륨 문제

```bash
# 볼륨 초기화
docker compose down -v
docker compose up -d
```

### 이미지 재빌드

```bash
# 캐시 없이 재빌드
docker compose build --no-cache
docker compose up -d
```
