# Notiflex 여정 기록

이 파일은 독자가 실제로 진행한 내용을 기록한다. AI가 각 챕터 완료 시 자동으로 업데이트한다.

## 진행 현황

| 챕터 | 서브챕터 | 상태 | 완료일 | 비고 |
|------|---------|------|--------|------|
| ch2 | 2.2 설치 확인 | ✅ | 2026-08-13 | Hermes Agent v0.20.0 |
| ch2 | 2.3 gcloud 설정 | ⬜ | | GKE 미사용 |
| ch2 | 2.4 GitHub 저장소 | ✅ | 2026-08-13 | yoonjaeyol/_Book_GitAIOps |
| ch2 | 2.5 클러스터 (kind→k3s) | ✅ | 2026-08-19 | kind v0.21.0 → k3s v1.34.10 (k8s v1.34), yjy-ubuntu20-test-3 |
| ch2 | 2.6 빌드/배포 | ✅ | 2026-08-13 | Docker scratch, 2 Pod Running |
| ch2 | 2.7 첫 커밋 | ✅ | 2026-08-13 | |
| ch3 | 3.2 GitOps 도구 | ✅ | 2026-08-15 | ArgoCD v3.5.1 (kind, NetworkPolicy 제거) |
| ch3 | 3.3 기능 추가 | ✅ | 2026-08-15 | /version endpoint (v0.1.1), rolling update, revert 롤백 테스트 |
| ch3 | 3.4 CI | ✅ | 2026-08-15 | GitHub Actions, Docker Hub push (firewood2002) |
| ch3 | 3.5 CI-CD 연결 | ✅ | 2026-08-15 | push → build → tag update → ArgoCD auto-sync |
| ch4 | 4.2 메트릭 모니터링 | ✅ | 2026-08-17 | kube-prometheus-stack 88.3.0, Notiflex 대시보드 4패널 (k3s에서 재검증 08-19: target 13, admin 200) |
| ch4 | 4.3 로그 수집 | ✅ | 2026-08-17 | Loki 6.55.0 SingleBinary + Fluent Bit 2.6.0, {namespace="notiflex"} 조회 확인 (k3s에서 재검증 08-19: 2 스트림) |
| ch4 | 4.4 알림 | ✅ | 2026-08-17 | PrometheusRule PodRestartTooMany, liveness 실패로 firing→Alertmanager active 확인 (k3s에서 재검증 08-19: firing 2건→active) |
| ch5 | 5.2 트래픽 관리 | ✅ | 2026-08-21 | Gateway API v1.6.1(experimental) + Traefik 3.7 + MetalLB L2 (VIP 10.10.20.100). Gateway Programmed=True, HTTPRoute Accepted/ResolvedRefs, ArgoCD notiflex-infra Synced/Healthy, VIP·NodePort(32639)로 /health·/version 실측 통과 (08-21 검증) |
| ch5 | 5.3 무중단 배포 | ⬜ | | |
| ch6 | 6.1 캐시 | ⬜ | | |
| ch6 | 6.2 시크릿 관리 | ⬜ | | |
| ch6 | 6.3 Canary 전환 | ⬜ | | |
| ch7 | 7.2 멀티 노드풀 | ⬜ | | |
| ch7 | 7.3 App of Apps | ⬜ | | |
| ch7 | 7.4 멀티테넌시 | ⬜ | | |
| ch8 | 8.1 메시징 | ⬜ | | |
| ch8 | 8.2 트레이싱 | ⬜ | | |
| ch8 | 8.3 CronJob | ⬜ | | |
| ch9 | 9.1 저장소 분석 | ⬜ | | |
| ch9 | 9.2 회고 | ⬜ | | |
| ch9 | 9.3 온보딩 문서 | ⬜ | | |
| ch9 | 9.4 GitAIOps 분석 | ⬜ | | |
| ch9 | 9.5 마무리 | ⬜ | | |

## 도구 선택 기록

| 영역 | 선택 | 검토한 대안 | 선택 이유 |
|------|------|-----------|----------|
| K8s 플랫폼 | k3s v1.34.10 | kind, GKE | 내부 인프라, 비용 제로, 빠른 실험 — 단일 바이너리 경량 (2026-08-19 kind v0.21.0 → k3s v1.34.10 전환) |
| 에이전트 | Hermes | Claude Code, Codex | 크로스세션 메모리, subagent 병렬 |
| GitOps | ArgoCD | Flux, Jenkins X, Spinnaker | Web UI로 배포 상태 시각화, selfHeal, CNCF Graduated (ch3) |
| CI | GitHub Actions | Cloud Build, GitLab CI, Jenkins | GitHub 네이티브, YAML 선언적, 별도 서버 불필요 (ch3) |
| 이미지 registry | Docker Hub | GCP Artifact Registry, ghcr.io | kind에서 anonymous pull 가능 확인 (ch3) |
| 메트릭 모니터링 | kube-prometheus-stack | Datadog, New Relic | 무료, K8s 네이티브, Grafana 통합 (ch4.2) |
| 로깅 | Loki + Fluent Bit | ELK (Elasticsearch), Datadog Logs | 경량, 라벨 기반 인덱싱, Grafana 네이티브 통합 (ch4.3) |
| 알림 | PrometheusRule + Alertmanager | Grafana Alert, Datadog Monitor | K8s 네이티브 (Prometheus Operator), GitOps로 룰 버전 관리 (ch4.4) |

## 현재 버전

| 컴포넌트 | 버전 | 변경 이력 |
|---------|------|----------|
| Go | 1.25 | golang:1.25-alpine |
| Notiflex 이미지 | firewood2002/notiflex-api:sha-cc3f896 | ch2.6 v0.1.0 → ch3.3 v0.1.1 → ch3.4 CI SHA 태그 |
| k3s | v1.34.10+k3s1 | 2026-08-19 kind(v0.21.0) → k3s 전환 (내장 containerd 2.2.5) |
| Kubernetes | v1.34.10 | k3s v1.34 라인 (flannel 10.42.0.0/16) |
| MetalLB | 0.16.1 (frr-k8s 0.0.25) | k8s/infra (ch5 대비), pool 10.10.20.100-140, L2 br0 |
| Traefik | v3.7.10 | GatewayController (k8s/infra), web=80/websecure=443 (NET_BIND_SERVICE) |
| Gateway API | v1.6.1 (experimental) | CRD (k8s/infra) — v1.2.1(standard)에서 v1.6.1로 승격: Traefik 3.7이 TLSRoute/BackendTLSPolicy를 v1으로 watch |
| ArgoCD | v3.5.1 | ch3.2 설치 (stable manifest, NetworkPolicy 제거) |
| kube-prometheus-stack | 88.3.0 (Prometheus 3.13.2, Grafana 13.1.3) | ch4.2 설치 (monitoring ns, 100m/256Mi 경량) |
| Loki | 3.6.7 (chart 6.55.0, SingleBinary) | ch4.3 설치 (filesystem, RF=1, auth off) |
| Fluent Bit | 2.1.0 (chart 2.6.0) | ch4.3 설치 (DaemonSet, Loki push, namespace/pod/container 라벨) |
| Alertmanager | (kube-prometheus-stack) | ch4.4 — PrometheusRule PodRestartTooMany (k8s/monitoring/, release: kube-prometheus) |

## 현재 리소스

| 노드 | 호스트 | 머신 | 노드 수 | 주요 워크로드 |
|------|--------|------|---------|-------------|
| yjy-ubuntu20-test-3 | 10.10.20.60 (k3s server, br0) | 4C 15G | 1 | 모든 워크로드 (k3s) |

## 트러블슈팅 이력

| 챕터 | 문제 | 해결 |
|------|------|------|
| ch2.5 | kind 바이너리 다운로드 403/404 | GitHub rate limit 우회 (Accept header + v0.21.0), 로컬→scp 전송 |
| ch2.5 | kubectl --short flag unknown | v1.29.1은 --short 미지원, 일반 버전 사용 |
| ch4.2 | retentionSize 2Gi 오류 | unit 없는 "2Gi" 불가, "2GiB"로 |
| ch4.3 | loki chart 7.x SimpleScalable 기본값 | 6.55.0(마지막 SingleBinary) 고정, read/write/backend replicas 0, test.enabled false |
| ch4.3 | mkdir /var/loki: read-only file system | persistence off 시 쓰기 불가 → persistence 1Gi (kind standard/local-path) |
| ch4.3 | 401 Unauthorized (3.x 기본 auth) | loki.auth_enabled: false |
| ch4.3 | too many unhealthy instances in the ring | commonConfig.replication_factor 3→1 (단일 ingester) |
| ch4.3 | notiflex 로그가 Loki에 없음 | tail은 파일 끝부터 — 앱 기동(19h 전) 이전 로그는 수집 안 됨. Pod 재시동으로 신규 로그 생성 후 수집 확인 |
| ch4.4 | kubectl delete pod은 알람 안 발생 | restarts_total은 동일 Pod의 컨테이너 재시작만 센다 — Pod 삭제는 새 Pod(0)이라 증가 0. liveness probe를 실패(__fail__)로 바꿔 crash loop로 실제 재시작 5회 유도 |
| ch4.4 | ArgoCD selfHeal이 테스트를 방해 | ArgoCD v3.5.1 CRD에 spec.suspend가 없음(unknown field). syncPolicy.automated.selfHeal: false로 임시 비활성화 후 테스트, 끝나고 true로 복원 |
| ch4.4 | 재현 시 selfHeal:false만으론 부족 | automated sync(주기적)이 kubectl 패치(liveness `__fail__`)를 계속 revert. **신뢰적 재현은 application controller(StatefulSet `argocd-application-controller`)를 replicas=0로 scale-down해 sync를 완전히 정지**한 뒤 crashloop 유도, 종료 후 replicas=1 복원. Prometheus `/api/v1/alerts` firing + Alertmanager `/api/v2/alerts` active 2건 확인 (2026-08-19) |
| ch2.5 (k3s) | k3s 설치 스크립트가 K3S_VERSION 무시 (v1.36.3 stable) | 특정 릴리스 바이너리 직접 다운로드 + sha256sum 대조 (GitHub release v1.34.10+k3s1) |
| ch2.5 (k3s) | k3s 기동 "6443 bind: address already in use" | ArgoCD UI용 port-forward 프로세스(8/14 잔존)가 6443 점유 — 종료 |
| ch2.5 (k3s) | k3s 내장 containerd CRI 플러그인 "too many open files" (fsnotify watcher) | inotify 한계는 **uid당** (`max_user_instances=128`) — kind control-plane(=docker bridge)의 k8s 스택(모두 uid 0)이 소진. kind 삭제 후 inotify 137→21로 해소, CRI 정상 |
| ch2.5 (k3s) | kind→k3s 인프라 매니페스트 불일치 | L2Advertisement interfaces `eth0`(kind docker bridge)→`br0`, MetalLB pool 172.18.0.x→10.10.20.100-140 (LAN 공중 대역), loki storageClass standard→local-path |
| ch4 (k3s) | Kube*Down 알림이 항상 pending/firing | k3s는 control-plane이 호스트 단일 프로세스라 Prometheus 타겟 미노출 — 차트 기본 룰의 특이점, 유해하지 않음 (kind에서도 존재) |
| ch5.2 (k3s) | Gateway 11시간째 Programmed=Unknown, "Waiting for controller" | **Gateway API 버전 불일치**: v1.2.1(standard)은 TLSRoute/BackendTLSPolicy를 v1alpha2/v1alpha3로 서빙하는데 Traefik 3.7은 v1으로 watch → reflector cache sync 미완료 → gateway provider reconcile 불발. v1.6.1(experimental, v1 서빙)으로 승격 후 해결 |
| ch5.2 (k3s) | traefik LB EXTERNAL-IP `<pending>` (MetalLB pool autoAssign:false) | pool이 autoAssign:false라 svc에 IP 명시 필요 — traefik svc에 `loadBalancerIP: 10.10.20.100` 추가 (재현성 유지) |
| ch5.2 (k3s) | k3s 내장 svclb가 MetalLB와 LB IP를 경쟁(flapping) | k3s server에 `--disable service-loadbalancer` 추가 + 재기동, 잔존 svclb pod/DaemonSet 정리 |
| ch5.2 (k3s) | Gateway listener "no matching entryPoint for port 80" | traefik web entrypoint가 :8000(매니페스트)인데 Gateway가 port 80 — web=80/websecure=443로 통일, non-root가 80/443 바인딩에 `NET_BIND_SERVICE` 추가 |
| ch5.2 (k3s) | VIP(10.10.20.100) 외부 flap/타임아웃, /32가 br0에 안 걸림 | **MetalLB L2 speaker가 NET_RAW만 있고 NET_ADMIN 없음** — VIP를 bridge에 secondary로 올리는 netlink 작업 거부. speaker에 NET_ADMIN 추가. 추가: bridge(br0) 환경 특성상 /32가 간헐적으로만 걸림, 외부 안정 접근은 NodePort(32639)로. 노드 `rp_filter` 2→0(MetalLB L2 필요, /etc/sysctl.d/90-metallb-l2.conf) |
| ch5.2 (k3s) | syncOptions에 `ServerSideApply=true`를 넣었는데도 SyncError가 그대로 남음 | ArgoCD는 **실패한 sync를 동일 revision으로 자동 재시도하지 않는다** — spec을 바꿔 generation이 올라가도(471→479) 재실행되지 않음. 수동 sync는 `spec.operation`이 아니라 **top-level `operation`** 필드를 patch해야 걸린다 (`spec.operation`은 v3.5.1 CRD에서 unknown field) |
| ch5.2 (k3s) | SSA 적용 후 httproutes CRD `forbidden: ValidatingAdmissionPolicy 'safe-upgrades' ... experimental on top of standard` | SSA로 annotation 262144바이트 한도는 해소됐지만(SSA는 `last-applied-configuration`을 쓰지 않음) 다음 벽이 드러남 — live httproutes만 **v1.2.1 standard** 잔존(나머지 12개 CRD는 v1.6.1 experimental 승격 완료)이라 v1.6 번들이 함께 설치하는 `safe-upgrades` VAP가 standard→experimental 전환을 차단. **일회성 마이그레이션**: VAPBinding 임시 삭제 → httproutes를 SSA(field-manager `argocd-controller`)로 experimental v1.6.1 적용 → 잔존 193KB last-applied annotation 제거 → ArgoCD sync가 VAPBinding을 Git에서 복원. 이후엔 `oldObject.channel == 'experimental'`이라 VAP를 정상 통과 |