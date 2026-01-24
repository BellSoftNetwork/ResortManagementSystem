# GitOps 배포용 예제 구성

이 디렉토리는 ArgoCD를 통한 GitOps 배포를 위한 예제 Kubernetes 리소스 정의 파일들을 포함합니다.

## 디렉토리 구조

```
deploy/example/
├── README.md                    # 이 파일
├── base/                        # 기본 리소스 정의
│   ├── kustomization.yaml      # Kustomize 기본 설정
│   ├── namespace.yaml          # 네임스페이스 정의
│   ├── api-core/              # Go API (api-core) 리소스
│   │   ├── deployment.yaml    # Deployment 정의
│   │   ├── service.yaml       # Service 정의
│   │   └── configmap.yaml     # ConfigMap 정의
│   ├── api-legacy/            # Kotlin API (api-legacy) 리소스
│   │   ├── deployment.yaml    # Deployment 정의
│   │   └── service.yaml       # Service 정의
│   ├── frontend-web/          # Vue.js Frontend 리소스
│   │   ├── deployment.yaml    # Deployment 정의
│   │   └── service.yaml       # Service 정의
│   └── ingress.yaml           # Ingress 통합 설정
└── overlays/                   # 환경별 오버레이
    ├── production/             # 프로덕션 환경
    │   ├── kustomization.yaml  # 프로덕션 특화 설정
    │   └── values.yaml         # 환경별 값 오버라이드
    └── staging/                # 스테이징 환경
        ├── kustomization.yaml  # 스테이징 특화 설정
        └── values.yaml         # 환경별 값 오버라이드
```

## GitOps 레포지토리 구성 방법

1. GitOps 레포지토리(`https://gitlab.bellsoft.net/devops/bsn-gitops`)에 다음 경로로 복사:
   ```
   kubernetes/clusters/bsn-main/public/bsn-resort-management-system-production/
   ```

2. 프로덕션 환경용 Kustomization 파일 예시:
   ```yaml
   # kubernetes/clusters/bsn-main/public/bsn-resort-management-system-production/kustomization.yaml
   apiVersion: kustomize.config.k8s.io/v1beta1
   kind: Kustomization
   
   namespace: bsn-resort-management-system-production
   
   resources:
     - ../../../../../../base/resort-management-system/
   
   patchesStrategicMerge:
     - patches/production-values.yaml
   
   images:
     - name: api-core
       newName: registry.gitlab.bellsoft.net/bell/resort-management-system/api-core
       newTag: main-abcd1234  # CI/CD에서 업데이트
     - name: api-legacy
       newName: registry.gitlab.bellsoft.net/bell/resort-management-system/api-legacy
       newTag: main-abcd1234  # CI/CD에서 업데이트
     - name: frontend-web
       newName: registry.gitlab.bellsoft.net/bell/resort-management-system/frontend-web
       newTag: main-abcd1234  # CI/CD에서 업데이트
   ```

## CI/CD 파이프라인 연동

GitLab CI/CD에서 이미지 빌드 완료 후, GitOps 레포지토리의 `kustomization.yaml` 파일에서 `images` 섹션의 `newTag`만 업데이트하여 배포를 트리거합니다.

## ArgoCD Application 설정

ArgoCD에서 다음과 같이 Application을 생성합니다:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: bsn-resort-management-system-production
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://gitlab.bellsoft.net/devops/bsn-gitops
    targetRevision: HEAD
    path: kubernetes/clusters/bsn-main/public/bsn-resort-management-system-production
  destination:
    server: https://kubernetes.default.svc
    namespace: bsn-resort-management-system-production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```