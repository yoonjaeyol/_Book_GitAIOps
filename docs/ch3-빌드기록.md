# ch3. 배포 자동화 구축 기록 (GitOps + CI/CD)

> 2026-08-15 완료. kind 클러스터 기반.
> 원본 가드레일은 GKE 기준이므로, 본 문서는 **kind 환경에서 실제 실행한 내용**을 기록한다.

## 1. 전체 아키텍처

```
app/main.go 변경 (git push)
        │
        ▼
┌──────────────────────┐
│  GitHub Actions (CI) │
│  1. checkout         │
│  2. docker build     │
│  3. docker push ─────┼──→ Docker Hub (firewood2002)
│  4. deployment.yaml  │
│     image tag 갱신   │
│  5. git push ────────┼──→ main 브랜치
└──────────────────────┘
        │
        ▼
┌──────────────────────┐
│   ArgoCD (CD)        │
│  1. Git poll (3분)   │
│  2. diff 감지        │
│  3. k8s/smb apply    │
│  4. rolling update   │
└──────────────────────┘
        │
        ▼
   notiflex-api Pod 교체 (서비스 중단 없음)
```

## 2. ch3.2 — ArgoCD 설치 (GitOps)

### 설치
```bash
kubectl --context kind-notiflex create namespace argocd
kubectl --context kind-notiflex apply -n argocd \
  -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml \
  --server-side=true --force-conflicts=true
```
- `--server-side=true --force-conflicts=true` 필수 (CRD annotations가 너무 길어 client-side apply가 실패함)
- 결과: ArgoCD v3.5.1, 7개 Pod 모두 Running

### kind 환경 필수 조치: NetworkPolicy 삭제
ArgoCD v3 stable manifest에 포함된 NetworkPolicy가 repo-server의 GitHub 접근을 차단한다.
```bash
kubectl --context kind-notiflex delete networkpolicy -n argocd --all
kubectl --context kind-notiflex rollout restart deployment -n argocd
kubectl --context kind-notiflex rollout restart statefulset -n argocd
```

### ArgoCD CLI 설치 (원격 노드)
```bash
# v3.5.1 (서버와 동일 버전)
scp argocd root@10.10.20.60:/usr/local/bin/argocd
chmod +x /usr/local/bin/argocd
```
> 주의: kind 환경에서 argocd CLI는 클러스터 내부 DNS(argocd-server.argocd)로 접속 불가.
> UI는 port-forward로 접속. 운영은 kubectl로 수행.

### GitHub 저장소 연결
```yaml
# argocd/repo-secret.yaml — git에 커밋하지 않음 (7장 App of Apps에서 placeholder로 덮어쓰김)
apiVersion: v1
kind: Secret
metadata:
  name: repo-notiflex-platform
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: repository
stringData:
  url: https://github.com/yoonjaeyol/_Book_GitAIOps.git
  username: yoonjaeyol
  password: <GH_PAT>
  type: git
  forceHttpBasicAuth: "true"   # ArgoCD v3에서 private repo 필수
```

```yaml
# argocd/notiflex-smb.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: notiflex-smb
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/yoonjaeyol/_Book_GitAIOps.git
    targetRevision: main
    path: k8s/smb
  destination:
    server: https://kubernetes.default.svc
    namespace: notiflex
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

> 주의: repo secret 적용 후에도 Sync Unknown이면 repo-server가 secret을 캐싱한 것.
> 전체 `rollout restart`(deployment + statefulset)로 해결.

### ArgoCD 정보
| 항목 | 값 |
|---|---|
| 버전 | v3.5.1 (quay.io/argoproj/argocd) |
| 네임스페이스 | argocd |
| Admin ID/PW | admin / GBpYLJ2Sij2A1bEJ |
| UI 접속 | `kubectl port-forward svc/argocd-server -n argocd 8443:443` → https://localhost:8443 |
| 확인 명령 | `kubectl get application notiflex-smb -n argocd -o custom-columns='NAME:.metadata.name,SYNC:.status.sync.status,HEALTH:.status.health.status,REV:.status.sync.revision'` |

## 3. ch3.3 — 기능 추가 + 롤링 업데이트 + 롤백

### 새 기능: /version endpoint
```go
// app/main.go
var version = "v0.1.1"

http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"version":"%s","pod":"%s"}`, version, podName)
})
```
- `k8s/smb/deployment.yaml`의 image tag를 `notiflex-api:v0.1.1`로 변경
- `git push` → ArgoCD가 3분 내 감지 → Pod 하나씩 교체 (rolling update, 서비스 중단 없음)

### 호출 방법 (kind)
kind는 NodePort가 호스트 IP로 노출되지 않으므로 **port-forward**가 유일 방법:
```bash
# 원격 노드에서
nohup kubectl --context kind-notiflex port-forward svc/notiflex-api -n notiflex 8080:80 &
curl -s http://localhost:8080/version
# {"version":"v0.1.1","pod":"notiflex-api-5564c8b744-vqr6g"}
```
- kind extraPortMappings나 Ingress Controller가 없는 한 NodePort는 클러스터 내부에서만 동작
- Pod 이미지(scratch)에는 wget/curl이 없음 → kubectl exec로 HTTP 테스트 불가

### 롤백 테스트 (GitOps의 핵심)
```bash
# 1. 롤백: 마지막 커밋을 되돌린다
cd /Users/yjy/yjy/_Book_GitAIOps
git revert HEAD --no-edit
git push

# 2. ArgoCD가 revert 커밋(db23e9e)을 감지 → v0.1.0으로 자동 롤백
#    확인: REV=db23e9e, Pod image=notiflex-api:v0.1.0

# 3. 복원: revert를 다시 revert
git revert HEAD --no-edit
git push
#    확인: REV=5e7170b, Pod image=notiflex-api:v0.1.1
```

| 상태 | Revision | Pod image | /version |
|---|---|---|---|
| v0.1.1 배포 | e5b8c04 | v0.1.1 | ✅ |
| 롤백 | db23e9e | v0.1.0 | ❌ (없음) |
| 복원 | 5e7170b | v0.1.1 | ✅ |

## 4. ch3.4 — GitHub Actions CI

### 구조
```
app/** push → Actions runner(ubuntu-latest)
  1. checkout
  2. docker/login-action → Docker Hub (firewood2002)
  3. docker/build-push-action → firewood2002/notiflex-api:sha-<7자리SHA>
  4. sed로 deployment.yaml image 교체 → git push (GH_PAT 사용)
```

### 워크플로 파일: `.github/workflows/ci.yaml`
- 트리거: main 브랜치 push + `app/**` 경로 변경
- 이미지 태그: `sha-${GITHUB_SHA:0:7}` (git SHA 기반, unique)
- git push: `https://x-access-token:${GH_PAT}@github.com/...` (actions/checkout의 GH_TOKEN이 repo scope 부족 문제 대비)
- `permissions: contents: write` 필수 (workflow 파일 push용)

### kind에서 registry 선택
| Registry | kind anonymous pull | 채택 |
|---|---|---|
| Docker Hub | ✅ (hello-world pull 확인) | ✅ |
| ghcr.io | ❌ (403 Forbidden) | ❌ |
| GCP Artifact Registry | ❌ (GCP 인증 필요) | ❌ |

### GitHub Secrets (yoonjaeyol/_Book_GitAIOps)
| Secret | 값 |
|---|---|
| `DOCKERHUB_USERNAME` | firewood2002 |
| `DOCKERHUB_TOKEN` | dckr_pat_r6... (Docker Hub PAT) |
| `GH_PAT` | ghp_... (repo scope) |

### 설정 확인 사항
- `Settings → Actions → General → Workflow permissions`: **Read and write permissions**
- Secrets 값이 Secret 이름과 동일하게 등록되면 `Username and password required` 에러 (실제 값 등록 필수)

### 실행 결과
- `feat: v0.1.3 CI 재시도` (cc3f896) push → build → push 성공
- deployment.yaml 자동 갱신: `image: firewood2002/notiflex-api:sha-cc3f896`
- ArgoCD sync → Pod 교체 완료

## 5. ch3.5 — CI-CD 연결 (완료)

CI가 deployment.yaml을 push하면 ArgoCD가 감지하여 자동 sync하므로,
**app 코드 push 하나면 빌드 → 배포까지全自动**.

```
git push (app/**)
  → Actions: build + Docker Hub push
  → Actions: deployment.yaml tag 갱신 + push
  → ArgoCD: 3분 내 감지, apply, rolling update
  → Pod 교체 (신호: kubectl get pods -n notiflex)
```

수동 필요 없는 지점:
- Docker build: ❌ (Actions가 함)
- kind load docker-image: ❌ (Docker Hub pull로 대체)
- kubectl apply: ❌ (ArgoCD가 함)
- rollback: `git revert HEAD --no-edit && git push`만

## 6. 환경 정보

| 항목 | 값 |
|---|---|
| 클러스터 | kind 'notiflex', yjy-ubuntu20.60 (yjy-ubuntu20-test-3) |
| kubectl context | kind-notiflex |
| K8s | v1.29.1 |
| kind | v0.21.0 |
| ArgoCD | v3.5.1 |
| Notiflex | firewood2002/notiflex-api:sha-cc3f896 (v0.1.3, /health /id /version) |
| Pod | 2개 (replicas: 2) |
| Service | notiflex-api (ClusterIP 80→8080) |
| GitHub | yoonjaeyol/_Book_GitAIOps (SSH alias: yoonjaeyol-github) |
| Docker Hub | firewood2002 (PAT 인증) |

## 7. 트러블슈팅 기록 (실전)

| 증상 | 원인 | 해결 |
|---|---|---|
| `kubectl apply` CRD 실패 | annotations 길이 초과 | `--server-side=true --force-conflicts=true` |
| Application Sync Unknown (NetworkPolicy) | ArgoCD v3 manifest의 NetworkPolicy가 repo-server egress 차단 | `kubectl delete networkpolicy -n argocd --all` + 전체 restart |
| repo secret 등록 후에도 Unknown | repo-server secret 캐시 | 전체 rollout restart (deployment + statefulset) |
| argocd CLI 접속 불가 (kind) | 클러스터 내부 DNS 외부 불가 | kubectl 또는 port-forward로 UI 접속 |
| NodePort 외부 접근 실패 (kind) | kind가 6443만 host에 매핑 | port-forward 사용 |
| ghcr.io pull 403 | anonymous token 거부 | Docker Hub 사용 (anonymous pull 가능 확인) |
| CI `Username and password required` | Docker Hub PAT 인증 실패/Secret 값 오류 | PAT 재확인, Secret에 실제 값 등록 |
| Actions git push 403 | Workflow permissions read-only | `Read and write permissions` 설정 + `contents: write` |
| ArgoCD가 새 커밋을 지연 감지 | polling interval 3분 + manifest cache | repo secret 재생성 + 전체 restart (강제 refresh) |