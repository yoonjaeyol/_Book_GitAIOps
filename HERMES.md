# HERMES.md

> 이 책은 **Claude Code**를 기준으로 쓰였습니다. **Hermes Agent**에서도 동작 가능합니다.
> 본 파일은 Hermes 한정 차이점만 다루며, 책 본문 가이드는 **반드시 같은 디렉터리의
> `CLAUDE.md`를 먼저 읽어주세요**.

## 0. 🚨 CRITICAL: 필수 규칙

**이 파일과 CLAUDE.md의 모든 규칙을 함께 따릅니다. 예외 없음.**

- **한국어**: 이 책은 한국어 책입니다. **반드시 한국어로 계속 진행**합니다. 영어로 전환하지 않습니다.
- **kubectl 안전**: **모든 kubectl 명령에 `--context kind-notiflex`를 반드시 지정**합니다. (GKE 대신 kind를 사용하므로 컨텍스트 이름이 다름)
- **자율 실행**: 모든 도구를 직접 호출하여 실행합니다. 사용자에게 명령 안내만 하지 않습니다.
- **인프라**: GCP/GKE 대신 **kind (Kubernetes IN Docker)**를 사용합니다. 원격 노드(`yjy-ubuntu20-test-3` 10.10.20.60)에서 동작합니다.

## 1. 실행 환경

Hermes Agent는 기본 구성으로 외부 네트워크 접근이 가능하므로 Codex의 `danger-full-access`나
Gemini의 `--yolo` 같은 별도 옵션이 **필요하지 않습니다**.

| Claude Code | Hermes Agent |
|---|---|
| `claude --dangerously-skip-permissions` | Hermes CLI 또는 Desktop 그대로 실행 |

작업 저장소(`notiflex-platform`)가 현재 저장소의 형제 디렉터리에 있습니다:
```
parent/
├── _Book_GitAIOps/    ← 이 저장소 (실습 가이드)
└── notiflex-platform/ ← 독자의 작업 저장소
```

## 2. 도구 매핑

Hermes Agent는 Claude Code의 도구 명칭과 다릅니다. 책의 가드레일 파일은 도구 명칭과 무관하게 작성되었으므로
동일하게 동작합니다. 주요 매핑:

| Claude Code 도구 | Hermes Agent 도구 |
|---|---|
| Shell | `terminal` |
| Edit | `patch` (기존 파일 수정) |
| Write | `write_file` (신규 파일 작성) |
| Read | `read_file` |
| Ask (필요시) | `clarify` |

## 3. Hermes 강점 (vs Claude/Codex/Gemini)

### 크로스세션 메모리
- Hermes는 `memory` tool로 영속 메모리 유지 (Codex의 자동 메모리 없음 문제 해소)
- `JOURNEY.md`도 계속 업데이트하되, Hermes 메모리와 중복 정보를 기록하지 않아도 됩니다.

### Subagent 병렬 처리
- Hermes는 `delegate_task`로 격리된 서브에이전트 호출 가능 (Claude의 Agent 도구와 유사)
- 독립적인 작업(예: 여러 Namespace 병렬 설정)은 subagent로 분리하여 병렬 실행 가능합니다.

### Computer Use (macOS 배경 제어)
- `computer_use` tool로 macOS 데스크톱 배경 제어 가능
- Grafana/ArgoCD GUI 확인 등 시각적 검증이 필요한 단계에서 유용합니다.

## 4. 알려진 동작 차이

### /update-docs 대응
Claude Code의 `/skill` /update-docs는 Hermes에 없습니다. 직접 실행:
- `prompt-guardrails/ch2/update-docs-skill.md` 가드레일 읽기
- `~/.hermes/skills/`에 Skill 생성 또는 `skill_manage` tool로 직접 실행

### 경로 처리
- 파일 생성 시 **절대 경로**를 명시합니다.
- `notiflex-platform/` 접근 시 형제 디렉터리임을 인지하고 절대 경로 사용:
  `/Users/yjy/notiflex-platform/...`

### 파일 편집 방식
- 기존 파일 수정: `patch(mode='replace')` (Claude Edit 대응)
- 신규 파일: `write_file` (Claude Write 대응)
- `patch` 적용 실패 시 재읽기 후 재시도, 동일 위치 2회 실패 시 `write_file`으로 전체 재작성

### GBrain 연동 (선택)
- Hermes는 GBrain knowledge base 연동 가능
- `gbrain-knowledge-base` skill 로드 후 기술 쿼리 시 자동 레트리가 활성화됩니다.
-本书의 가드레일과 충돌하지 않습니다.

## 5. 책 본문 진행

CLAUDE.md의 모든 규칙과 동일하게 따릅니다:
- `decision-guides/` — 도구 선택 근거
- `prompt-guardrails/` — 단계별 실행 지침
- `result-templates/` — 검증 체크리스트
- 3-프롬프트 패턴 (탐색 → 비교 → 실행)
- JOURNEY.md 관리 규칙

## 6. 에이전트 파일 간 참조

- Claude Code: `CLAUDE.md` (책 본문 가이드)
- Codex CLI: `AGENTS.md` (형제 파일)
- Gemini CLI: `GEMINI.md` (형제 파일)
- Hermes Agent: `HERMES.md` (본 파일)