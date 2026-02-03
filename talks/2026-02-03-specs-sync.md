# Claritask.md와 Commands.md 동기화

**날짜**: 2026-02-03

## 개요

Claritask.md와 Commands.md 간 불일치 항목을 식별하고 수정했습니다.

---

## 발견된 불일치 및 해결

### 1. Task 추가 명령어
- **불일치**: Claritask.md는 `task add`, Commands.md는 `task push`
- **결정**: `clari task push '<json>'`로 통일

### 2. Memo set 문법
- **불일치**: Claritask.md는 key가 JSON 내부, Commands.md는 별도 인자
- **결정**: `clari memo set <key> '<json>'`로 통일

### 3. Memo scope
- **불일치**: Claritask.md는 project/feature/task, Commands.md는 project/phase/task
- **결정**: `project / feature / task`로 통일

### 4. Feature spec 문법
- **불일치**: Claritask.md는 `clari feature <id> spec`, Commands.md는 `clari feature spec <id>`
- **결정**: `clari feature spec <id> '<spec_text>'`로 통일

### 5. Edge add 플래그 순서
- **불일치**: --feature 플래그 위치 상이
- **결정**: `clari edge add --feature --from <id> --to <id>`로 통일

### 6. Feature tasks 명령어
- **불일치**: Claritask.md는 `clari feature <id> tasks`, Commands.md는 `clari fdl tasks`
- **결정**: `clari fdl tasks <feature_id>`로 통일 (FDL 기반)

### 7. Feature get 명령어
- **불일치**: Claritask.md에 없음
- **결정**: `clari feature get <id>` 추가

---

## Phase 제거 결정

### 논의
- Phase와 Feature 역할이 중복됨
- Feature + Edge로 순서/의존성 표현 가능
- FDL이 Feature 단위로 설계되어 있음

### 결정
Commands.md에서 Phase 관련 명령어 전체 제거:
- `clari phase create`
- `clari phase list`
- `clari phase plan`
- `clari phase start`

### 변경 사항
- 명령어 수: 48개 → 44개
- `phase_id` → `feature_id`로 변경
- `current_phase` → `current_feature`로 변경
- memo scope: `phase` → `feature`

---

## Project set JSON 포맷 정리

### 문제
```json
{
  "name": "Blog Platform",        // 중복
  "description": "...",           // 중복
  "context": {
    "project_name": "Blog Platform",  // 중복
    "description": "..."              // 중복
  }
}
```

### 결정: 최상위 필드 제거
```json
{
  "context": {
    "project_name": "Blog Platform",
    "description": "..."
  },
  "tech": {...},
  "design": {...}
}
```

- 프로젝트 id/description은 `clari init`에서 설정
- `project set`은 context/tech/design 설정 전용

---

## 최종 명령어 구조

```
clari
├── init
├── project (set/get/plan/start/stop/status)
├── task (push/pop/start/complete/fail/status/get/list)
├── feature (list/add/get/spec/start)
├── edge (add/list/remove/infer)
├── fdl (create/register/validate/show/skeleton/tasks/verify/diff)
├── plan (features)
├── memo (set/get/list/del)
├── context (set/get)
├── tech (set/get)
├── design (set/get)
└── required
```

---

## 수정된 파일

- `specs/Claritask.md`
- `specs/Commands.md` (v3.0 → v3.1)
