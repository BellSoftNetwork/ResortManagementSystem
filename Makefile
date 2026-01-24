# Resort Management System - Development Makefile
# 개발에 필요한 모든 명령어를 모아둔 Makefile
#
# Usage: make <target>
# Help:  make help

.PHONY: help
.DEFAULT_GOAL := help

# Colors for terminal output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

##@ 도움말
help: ## 사용 가능한 명령어 목록
	@awk 'BEGIN {FS = ":.*##"; printf "\n$(CYAN)사용법:$(RESET) make $(GREEN)<target>$(RESET)\n"} \
		/^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Docker Compose - 전체 환경
up: ## 전체 개발 환경 시작 (mysql, redis, api-core, api-legacy, frontend)
	docker compose up -d

down: ## 전체 환경 중지 (데이터 유지)
	docker compose down

stop: ## 컨테이너 중지 (CPU 절약, 데이터 유지)
	docker compose stop

start: ## 중지된 컨테이너 재시작
	docker compose start

restart: ## 전체 환경 재시작
	docker compose restart

ps: ## 실행 중인 컨테이너 상태 확인
	docker compose ps -a

logs: ## 전체 로그 확인 (실시간)
	docker compose logs -f

clean: ## 컨테이너 및 네트워크 제거 (볼륨 유지)
	docker compose down

clean-all: ## 완전 제거 (볼륨 포함, 주의!)
	docker compose down -v
	docker system prune -f

rebuild: ## 이미지 재빌드 후 시작
	docker compose build --no-cache
	docker compose up -d

##@ 개별 서비스 시작/중지
up-db: ## DB 서비스만 시작 (mysql, redis)
	docker compose up -d mysql redis

up-api: ## API 서비스만 시작 (api-core, api-legacy + DB)
	docker compose up -d mysql redis api-core api-legacy

up-core: ## api-core만 시작 (+ DB)
	docker compose up -d mysql redis api-core

up-legacy: ## api-legacy만 시작 (+ DB)
	docker compose up -d mysql redis api-legacy

up-frontend: ## frontend만 시작 (+ 전체 백엔드)
	docker compose up -d

##@ 개별 서비스 로그
logs-core: ## api-core 로그
	docker compose logs -f api-core

logs-legacy: ## api-legacy 로그
	docker compose logs -f api-legacy

logs-frontend: ## frontend 로그
	docker compose logs -f frontend

logs-db: ## MySQL 로그
	docker compose logs -f mysql

logs-redis: ## Redis 로그
	docker compose logs -f redis

##@ 컨테이너 접속
sh-core: ## api-core 컨테이너 쉘 접속
	docker compose exec api-core sh

sh-legacy: ## api-legacy 컨테이너 쉘 접속
	docker compose exec api-legacy bash

sh-frontend: ## frontend 컨테이너 쉘 접속
	docker compose exec frontend sh

sh-db: ## MySQL 컨테이너 쉘 접속
	docker compose exec mysql bash

##@ api-core (Go)
core-test: ## api-core 테스트 실행
	docker compose exec api-core go test ./...

core-lint: ## api-core 린트 실행
	docker compose exec api-core golangci-lint run

core-build: ## api-core 빌드
	docker compose exec api-core go build -o /tmp/main ./cmd/server

core-migrate: ## api-core DB 마이그레이션 실행
	docker compose exec api-core go run cmd/migrate/main.go -action=migrate

core-migrate-status: ## api-core 마이그레이션 상태 확인
	docker compose exec api-core go run cmd/migrate/main.go -action=status

##@ api-legacy (Kotlin/Spring)
legacy-test: ## api-legacy 테스트 실행
	docker compose exec api-legacy ./gradlew test

legacy-lint: ## api-legacy 린트 체크
	docker compose exec api-legacy ./gradlew ktlintCheck

legacy-format: ## api-legacy 코드 포맷팅
	docker compose exec api-legacy ./gradlew ktlintFormat

legacy-build: ## api-legacy 빌드
	docker compose exec api-legacy ./gradlew build -x test

##@ frontend (Vue.js/Quasar)
frontend-dev: ## frontend 개발 서버 (이미 실행 중이면 불필요)
	docker compose exec frontend yarn dev

frontend-build: ## frontend 프로덕션 빌드
	docker compose exec frontend yarn build

frontend-lint: ## frontend 린트 실행
	docker compose exec frontend yarn lint

frontend-test: ## frontend 테스트 실행
	docker compose exec frontend yarn test

##@ 데이터베이스
db-shell: ## MySQL 쉘 접속
	docker compose exec mysql mysql -uroot -proot

db-core: ## rms-core DB 접속
	docker compose exec mysql mysql -uroot -proot rms-core

db-legacy: ## rms-legacy DB 접속
	docker compose exec mysql mysql -uroot -proot rms-legacy

db-reset-core: ## rms-core DB 초기화 (주의!)
	docker compose exec mysql mysql -uroot -proot -e "DROP DATABASE IF EXISTS \`rms-core\`; CREATE DATABASE \`rms-core\`;"
	@echo "api-core를 재시작하여 마이그레이션을 실행하세요: make restart-core"

db-reset-legacy: ## rms-legacy DB 초기화 (주의!)
	docker compose exec mysql mysql -uroot -proot -e "DROP DATABASE IF EXISTS \`rms-legacy\`; CREATE DATABASE \`rms-legacy\`;"
	@echo "api-legacy를 재시작하여 Liquibase를 실행하세요: make restart-legacy"

db-clear-login-attempts: ## 로그인 시도 기록 삭제 (brute force 잠금 해제)
	docker compose exec mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts; DELETE FROM \`rms-legacy\`.login_attempts;" 2>/dev/null || true

##@ API 테스트
test-api: ## API 호환성 테스트 실행 (core-only 모드)
	@docker compose exec -T mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts;" 2>/dev/null || true
	python3 scripts/api-compatibility-test.py --core-only

test-api-golden: ## Golden file 비교 테스트
	@docker compose exec -T mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts;" 2>/dev/null || true
	python3 scripts/api-compatibility-test.py --compare-golden

test-api-save-golden: ## 현재 응답을 Golden file로 저장
	@docker compose exec -T mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts;" 2>/dev/null || true
	python3 scripts/api-compatibility-test.py --save-golden

test-api-compare: ## api-legacy vs api-core 전체 비교
	@docker compose exec -T mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts; DELETE FROM \`rms-legacy\`.login_attempts;" 2>/dev/null || true
	python3 scripts/api-compatibility-test.py

api-test: ## 단일 API 테스트 (사용법: make api-test PATH=/api/v1/users)
	python3 scripts/api-test.py $(PATH)

##@ 헬스 체크
health: ## 전체 서비스 헬스 체크
	@echo "$(CYAN)=== Health Check ===$(RESET)"
	@echo -n "MySQL:      " && docker compose exec -T mysql mysqladmin ping -uroot -proot 2>/dev/null | grep -q alive && echo "$(GREEN)OK$(RESET)" || echo "$(RED)FAIL$(RESET)"
	@echo -n "Redis:      " && docker compose exec -T redis redis-cli ping 2>/dev/null | grep -q PONG && echo "$(GREEN)OK$(RESET)" || echo "$(RED)FAIL$(RESET)"
	@echo -n "API Core:   " && curl -sf http://localhost:8080/actuator/health >/dev/null && echo "$(GREEN)OK$(RESET)" || echo "$(RED)FAIL$(RESET)"
	@echo -n "API Legacy: " && curl -sf http://localhost:8081/actuator/health >/dev/null && echo "$(GREEN)OK$(RESET)" || echo "$(RED)FAIL$(RESET)"
	@echo -n "Frontend:   " && curl -sf http://localhost:9000 >/dev/null && echo "$(GREEN)OK$(RESET)" || echo "$(RED)FAIL$(RESET)"

health-core: ## api-core 상세 헬스 체크
	@curl -s http://localhost:8080/actuator/health | python3 -m json.tool

health-legacy: ## api-legacy 상세 헬스 체크
	@curl -s http://localhost:8081/actuator/health | python3 -m json.tool

##@ 유틸리티
create-test-data: ## 테스트 데이터 생성 (api-core)
	@docker compose exec -T mysql mysql -uroot -proot -e "DELETE FROM \`rms-core\`.login_attempts;" 2>/dev/null || true
	@TOKEN=$$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
		-H 'Content-Type: application/json' \
		-d '{"username":"testadmin","password":"testadmin123"}' | \
		python3 -c 'import sys,json; print(json.load(sys.stdin).get("value",{}).get("accessToken",""))') && \
	curl -s -X POST http://localhost:8080/api/v1/dev/test-data \
		-H "Authorization: Bearer $$TOKEN" \
		-H "Content-Type: application/json" \
		-d '{"type":"all"}' | python3 -m json.tool

create-test-user: ## 테스트 사용자 생성
	@curl -s -X POST http://localhost:8080/api/v1/auth/register \
		-H 'Content-Type: application/json' \
		-d '{"userId":"testadmin","email":"testadmin@test.com","name":"Test Admin","password":"testadmin123"}' | python3 -m json.tool
	@docker compose exec -T mysql mysql -uroot -proot -e "UPDATE \`rms-core\`.user SET role=127 WHERE user_id='testadmin';" 2>/dev/null || true
	@echo "$(GREEN)testadmin 계정 생성 완료 (SUPER_ADMIN)$(RESET)"

##@ 개발 워크플로우
dev: up ## 개발 환경 시작 (up의 별칭)

dev-core: up-core logs-core ## api-core 개발 모드 (시작 + 로그)

dev-legacy: up-legacy logs-legacy ## api-legacy 개발 모드 (시작 + 로그)

dev-frontend: up logs-frontend ## frontend 개발 모드 (전체 시작 + 로그)

restart-core: ## api-core만 재시작
	docker compose restart api-core

restart-legacy: ## api-legacy만 재시작
	docker compose restart api-legacy

restart-frontend: ## frontend만 재시작
	docker compose restart frontend

##@ 정보
info: ## 프로젝트 정보 및 URL
	@echo "$(CYAN)=== Resort Management System ===$(RESET)"
	@echo ""
	@echo "$(YELLOW)서비스 URL:$(RESET)"
	@echo "  Frontend:   http://localhost:9000"
	@echo "  API Core:   http://localhost:8080"
	@echo "  API Legacy: http://localhost:8081"
	@echo "  MySQL:      localhost:3306"
	@echo "  Redis:      localhost:6379"
	@echo ""
	@echo "$(YELLOW)API 문서:$(RESET)"
	@echo "  Swagger:    http://localhost:8080/docs/swagger-ui"
	@echo "  OpenAPI:    http://localhost:8080/docs/schema"
	@echo ""
	@echo "$(YELLOW)DB 접속 정보:$(RESET)"
	@echo "  User:       rms"
	@echo "  Password:   rms123"
	@echo "  Root Pass:  root"
	@echo "  Core DB:    rms-core"
	@echo "  Legacy DB:  rms-legacy"

version: ## 각 서비스 버전 확인
	@echo "$(CYAN)=== Version Info ===$(RESET)"
	@echo -n "Go:      " && docker compose exec -T api-core go version 2>/dev/null | cut -d' ' -f3 || echo "N/A"
	@echo -n "Java:    " && docker compose exec -T api-legacy java -version 2>&1 | head -1 || echo "N/A"
	@echo -n "Node:    " && docker compose exec -T frontend node -v 2>/dev/null || echo "N/A"
	@echo -n "MySQL:   " && docker compose exec -T mysql mysql -V 2>/dev/null | cut -d' ' -f3 || echo "N/A"
	@echo -n "Redis:   " && docker compose exec -T redis redis-server -v 2>/dev/null | cut -d' ' -f3 || echo "N/A"
