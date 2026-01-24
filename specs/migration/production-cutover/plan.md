---
id: production-cutover-plan
status: approved
updated: 2026-01-12
---

# 블루-그린 마이그레이션 실행 계획

> api-legacy → api-core 프로덕션 전환 체크리스트

---

## 검증 상태 (2026-01-12 완료)

### API 호환성
- [x] 40/40 엔드포인트 테스트 통과
- [x] Golden file 기반 응답 검증
- [x] JWT 토큰 호환성 확인

### 프론트엔드 연동
- [x] 로그인/로그아웃 정상
- [x] 대시보드 (예약 캘린더, 테이블) 정상
- [x] 객실 관리 (58개 객실 표시) 정상
- [x] 예약 관리 (18개 예약 표시) 정상
- [x] 통계 (매출/예약/인원 차트) 정상
- [x] 결제 수단 (5개) 정상
- [x] 계정 관리 (5개 계정) 정상
- [x] 콘솔 에러 없음

### DB 스키마
- [x] 테이블 구조 100% 동일 확인
- [x] 동일 DB 직접 사용 가능 확인

---

## Phase 1: 사전 준비 (D-1)

### 1.1 프로덕션 DB 백업
```bash
# 전체 백업
mysqldump -h <prod-db-host> -u <user> -p \
  --single-transaction --quick \
  rms-legacy > backup_$(date +%Y%m%d_%H%M%S).sql

# 백업 검증
mysql -h <backup-db> -u <user> -p rms-test < backup_*.sql
```

### 1.2 api-core 이미지 준비
```bash
# 프로덕션 이미지 빌드
docker build -t registry.example.com/api-core:v1.0.0 \
  -f apps/api-core/Dockerfile .

# 이미지 푸시
docker push registry.example.com/api-core:v1.0.0

# 이미지 검증
docker run --rm registry.example.com/api-core:v1.0.0 --version
```

### 1.3 환경 설정 준비
```yaml
# api-core ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-core-config
data:
  PROFILE: "production"
  DATABASE_MYSQL_DATABASE: "rms-legacy"  # 기존 DB 사용

# api-core Secret (기존과 동일)
apiVersion: v1
kind: Secret
metadata:
  name: api-core-secret
type: Opaque
data:
  DATABASE_MYSQL_HOST: <base64>
  DATABASE_MYSQL_USER: <base64>
  DATABASE_MYSQL_PASSWORD: <base64>
  JWT_SECRET: <base64>  # 기존과 동일해야 함
  REDIS_HOST: <base64>
```

---

## Phase 2: 그린 환경 배포 (D-Day)

### 2.1 api-core 배포 (트래픽 없음)
```bash
# api-core Deployment 적용
kubectl apply -f deploy/production/api-core/

# Pod 상태 확인
kubectl get pods -l app=api-core -w

# 헬스체크 확인
kubectl exec -it <api-core-pod> -- curl localhost:8080/actuator/health
```

### 2.2 내부 테스트
```bash
# 포트포워딩으로 내부 테스트
kubectl port-forward svc/api-core 8080:8080

# API 테스트 실행
python3 scripts/api-compatibility-test.py --core-only

# 프론트엔드 테스트 (로컬)
VITE_API_URL=http://localhost:8080 yarn dev
```

---

## Phase 3: 트래픽 전환 (Canary)

### 3.1 10% 트래픽 전환
```yaml
# Ingress canary 설정
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-canary
  annotations:
    nginx.ingress.kubernetes.io/canary: "true"
    nginx.ingress.kubernetes.io/canary-weight: "10"
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /api/v1
            pathType: Prefix
            backend:
              service:
                name: api-core
                port:
                  number: 8080
```

### 3.2 모니터링 (10분 관찰)
```bash
# 에러율 확인
kubectl logs -l app=api-core --tail=100 | grep -i error

# 응답 시간 확인 (Prometheus)
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
```

### 3.3 점진적 전환
| 시간 | 트래픽 비율 | 확인 사항 |
|------|-------------|-----------|
| T+0 | 10% | 에러율, 응답시간 |
| T+10min | 30% | 동일 |
| T+20min | 50% | 동일 |
| T+30min | 100% | 완전 전환 |

---

## Phase 4: 완전 전환

### 4.1 api-legacy 트래픽 제거
```yaml
# Ingress 수정: api-core만 사용
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api
spec:
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /api/v1
            pathType: Prefix
            backend:
              service:
                name: api-core
                port:
                  number: 8080
```

### 4.2 api-legacy Pod 유지 (롤백 대비)
```bash
# api-legacy를 0 replica로 줄이지 않음
kubectl scale deployment api-legacy --replicas=1
```

---

## Phase 5: 안정화 (D+7)

### 5.1 모니터링 계속
- 에러율 < 0.1%
- 응답시간 p99 < 500ms
- 메모리 사용량 안정

### 5.2 api-legacy 제거 (안정화 확인 후)
```bash
# 1주일 후 api-legacy 제거
kubectl delete deployment api-legacy
kubectl delete service api-legacy
```

---

## 롤백 계획

### 즉시 롤백 (5분 이내)
```bash
# Ingress를 api-legacy로 복구
kubectl apply -f deploy/production/ingress-legacy.yaml

# 또는 canary 비율 0%로
kubectl annotate ingress api-canary \
  nginx.ingress.kubernetes.io/canary-weight="0" --overwrite
```

### 롤백 기준
- 에러율 > 1%
- 응답시간 p99 > 2s
- 프론트엔드 크리티컬 기능 실패

---

## 커뮤니케이션

### 전환 전
- [ ] 팀 공지: 전환 일정 및 롤백 계획
- [ ] 온콜 담당자 지정

### 전환 중
- [ ] 슬랙 채널에 실시간 상태 공유
- [ ] 이상 징후 시 즉시 보고

### 전환 후
- [ ] 전환 완료 공지
- [ ] 모니터링 대시보드 공유

---

## 연락처

| 역할 | 담당자 | 연락처 |
|------|--------|--------|
| 전환 담당 | TBD | TBD |
| DB 담당 | TBD | TBD |
| 프론트엔드 | TBD | TBD |
| 온콜 | TBD | TBD |
