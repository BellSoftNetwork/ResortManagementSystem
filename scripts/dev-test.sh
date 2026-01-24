#!/bin/bash

# 리조트 관리 시스템 - 개발 테스트 실행기

set -e

# 출력 색상
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 색상 없음

# 색상 메시지 출력
print_msg() {
    echo -e "${2}${1}${NC}"
}

# 기본값
SERVICE=""
TEST_TYPE="unit"

# 명령줄 인자 파싱
while [[ $# -gt 0 ]]; do
    case $1 in
        --service)
            SERVICE="$2"
            shift 2
            ;;
        --type)
            TEST_TYPE="$2"
            shift 2
            ;;
        --help)
            print_msg "사용법: ./dev-test.sh [옵션]" "$BLUE"
            echo "옵션:"
            echo "  --service [api-core|api-legacy|frontend|all]  테스트할 서비스 (기본값: all)"
            echo "  --type [unit|integration|all]                 테스트 종류 (기본값: unit)"
            echo ""
            echo "예시:"
            echo "  ./dev-test.sh                                 # 모든 유닛 테스트 실행"
            echo "  ./dev-test.sh --service api-core              # api-core 유닛 테스트 실행"
            echo "  ./dev-test.sh --service api-legacy --type all # api-legacy 모든 테스트 실행"
            exit 0
            ;;
        *)
            print_msg "알 수 없는 옵션: $1" "$RED"
            exit 1
            ;;
    esac
done

print_msg "🧪 리조트 관리 시스템 - 테스트 실행기" "$GREEN"
echo ""

# api-core 테스트 실행 함수
test_api_core() {
    print_msg "🔧 API Core (Go) 테스트 중..." "$YELLOW"

    if [ "$TEST_TYPE" == "unit" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "유닛 테스트 실행 중..." "$BLUE"
        docker compose exec -T api-core go test -v ./internal/... -short
    fi

    if [ "$TEST_TYPE" == "integration" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "통합 테스트 실행 중..." "$BLUE"
        docker compose exec -T api-core go test -v ./internal/... -run Integration
    fi

    print_msg "✅ API Core 테스트 완료" "$GREEN"
    echo ""
}

# api-legacy 테스트 실행 함수
test_api_legacy() {
    print_msg "☕ API Legacy (Spring Boot) 테스트 중..." "$YELLOW"

    if [ "$TEST_TYPE" == "unit" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "유닛 테스트 실행 중..." "$BLUE"
        docker compose exec -T api-legacy ./gradlew test
    fi

    if [ "$TEST_TYPE" == "integration" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "통합 테스트 실행 중..." "$BLUE"
        docker compose exec -T api-legacy ./gradlew integrationTest
    fi

    print_msg "✅ API Legacy 테스트 완료" "$GREEN"
    echo ""
}

# frontend 테스트 실행 함수
test_frontend() {
    print_msg "🎨 Frontend (Vue.js) 테스트 중..." "$YELLOW"

    if [ "$TEST_TYPE" == "unit" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "유닛 테스트 실행 중..." "$BLUE"
        docker compose exec -T frontend yarn test:unit
    fi

    if [ "$TEST_TYPE" == "integration" ] || [ "$TEST_TYPE" == "all" ]; then
        print_msg "E2E 테스트 실행 중..." "$BLUE"
        docker compose exec -T frontend yarn test:e2e
    fi

    print_msg "✅ Frontend 테스트 완료" "$GREEN"
    echo ""
}

# 서비스 실행 상태 확인 함수
check_services() {
    if ! docker compose ps | grep -q "Up"; then
        print_msg "❌ 서비스가 실행되고 있지 않습니다. 먼저 ./dev-setup.sh를 실행해주세요." "$RED"
        exit 1
    fi
}

# 메인 테스트 실행
check_services

case $SERVICE in
    api-core)
        test_api_core
        ;;
    api-legacy)
        test_api_legacy
        ;;
    frontend)
        test_frontend
        ;;
    all|"")
        test_api_core
        test_api_legacy
        test_frontend
        ;;
    *)
        print_msg "❌ 알 수 없는 서비스: $SERVICE" "$RED"
        exit 1
        ;;
esac

print_msg "🎉 모든 테스트 완료!" "$GREEN"

# 커버리지 리포트 생성 (가능한 경우)
if [ "$SERVICE" == "api-core" ] || [ "$SERVICE" == "all" ] || [ "$SERVICE" == "" ]; then
    print_msg "📊 Go 커버리지 리포트 생성 중..." "$YELLOW"
    docker compose exec -T api-core go test -coverprofile=coverage.out ./internal/... || true
    docker compose exec -T api-core go tool cover -html=coverage.out -o coverage.html || true
    print_msg "커버리지 리포트가 apps/api-core/coverage.html에 저장되었습니다" "$GREEN"
fi

if [ "$SERVICE" == "api-legacy" ] || [ "$SERVICE" == "all" ] || [ "$SERVICE" == "" ]; then
    print_msg "📊 JaCoCo 커버리지 리포트 생성 중..." "$YELLOW"
    docker compose exec -T api-legacy ./gradlew jacocoTestReport || true
    print_msg "커버리지 리포트가 apps/api-legacy/build/reports/jacoco/test/html/index.html에 저장되었습니다" "$GREEN"
fi
