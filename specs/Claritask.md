# Claritask - Task And LLM Operating System

## 개요

LLM 기반 프로젝트 자동 실행 시스템

**목표**:
- 프로젝트 수동 세팅 자동화 (30-50분 절약)
- 무제한 무인 작업 가능 (Task 수 제한 없음)
- 컨텍스트 한계 완전 극복 (매 Task마다 초기화)

**철학**:
- **Claritask가 오케스트레이터**, Claude는 실행기
- Task 단위 독립 실행으로 컨텍스트 격리
- 한 줄 명령으로 프로젝트 완성

---

## 아키텍처: 제어 역전

### 기존 구조의 한계

기존에는 Claude Code가 Claritask를 도구로 사용했다. Claude Code가 대화형 세션에서 `clari task list`, `clari memo add` 명령어를 호출하며 작업을 관리했다. 이 구조는 단일 작업이나 탐색적 개발에는 적합하지만, 대규모 자동화에는 치명적인 한계가 있다.

- **컨텍스트 누적**: Task를 처리할수록 대화 컨텍스트가 쌓인다
- **세션 의존성**: Claude Code 세션이 끊기면 작업도 중단된다
- **수동 개입 필요**: `/clear`로 컨텍스트를 비우려면 사람이 직접 입력해야 한다
- **확장 불가**: Task가 100개, 1000개로 늘어나면 단일 세션으로 처리 불가능

### 새로운 구조: Claritask가 오케스트레이터

**제어권을 역전한다.** Claritask가 드라이버가 되고, Claude는 순수 실행기가 된다.

```
┌─────────────────────────────────────────────────────────────┐
│                        Claritask                            │
│                     (Orchestrator)                          │
│                                                             │
│   ┌─────────┐    ┌─────────┐    ┌─────────┐                │
│   │ Task 1  │───▶│ Task 2  │───▶│ Task N  │───▶ 완료       │
│   └────┬────┘    └────┬────┘    └────┬────┘                │
│        │              │              │                      │
│        ▼              ▼              ▼                      │
│   claude --print claude --print claude --print              │
│   (독립 컨텍스트) (독립 컨텍스트) (독립 컨텍스트)              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

`clari project start`를 실행하면:

1. Claritask가 pending 상태의 Task 목록을 조회한다
2. 각 Task를 순서대로 `claude --print` 비대화형 모드로 전달한다
3. Claude는 해당 Task만 처리하고 종료한다 (컨텍스트 초기화)
4. Claritask가 결과를 확인하고 Task를 완료 처리한다
5. 다음 Task로 넘어간다
6. 모든 Task가 완료될 때까지 반복한다

### 왜 이 구조가 강력한가

| 측면 | 기존 (Claude 드라이버) | 신규 (Claritask 드라이버) |
|------|------------------------|---------------------------|
| 컨텍스트 | 누적되어 폭발 | 매 Task마다 초기화 |
| 세션 | 끊기면 중단 | 프로세스 기반, 복구 가능 |
| 확장성 | 수십 개 한계 | 수천 개도 가능 |
| 상태 관리 | Claude 메모리 의존 | DB에 영속화 |
| 재시작 | 처음부터 다시 | 마지막 Task부터 재개 |

### 두 가지 모드 공존

Claritask는 두 가지 사용 방식을 모두 지원한다.

**1. 자동화 모드 (Claritask 드라이버)**
```bash
clari project start
# → Task 전체 순회, claude --print 반복 호출
# → 사람 개입 없이 완료까지 실행
```

**2. 대화형 모드 (Claude/사용자 드라이버)**
```bash
# Claude Code 세션 또는 사람이 직접 실행
clari task list
clari task get 3
clari memo add --scope task --id 3 "JWT 만료 시간 수정"
```

자동화 모드는 대량 작업을, 대화형 모드는 탐색과 디버깅을 담당한다. 두 모드는 동일한 DB를 공유하므로 언제든 전환 가능하다.

---

## 기술 스택

- **Go + SQLite**: 단일 바이너리, 고성능
- **파일**: `.claritask/db` 하나로 모든 것 관리
- **성능**: 1000개 Task도 1ms

---

## 데이터 구조: 그래프 기반

### project → feature → task (with edges)

```
project: Blog Platform
├─ feature: 로그인
│  ├─ task: user_table_sql
│  ├─ task: user_model ─────────depends_on────▶ user_table_sql
│  ├─ task: auth_service ───────depends_on────▶ user_model
│  └─ task: login_api ──────────depends_on────▶ auth_service
│
├─ feature: 결제 ───────────────depends_on────▶ 로그인 (Feature Edge)
│  ├─ task: payment_table_sql
│  ├─ task: payment_model ──────depends_on────▶ payment_table_sql, user_model
│  └─ task: payment_api ────────depends_on────▶ payment_model
│
└─ feature: 블로그
   ├─ task: post_table_sql
   ├─ task: post_model ─────────depends_on────▶ post_table_sql
   └─ task: post_api ───────────depends_on────▶ post_model, auth_service
```

**특징**:
- **project**: 프로젝트 전체
- **feature**: 기능 단위 (로그인, 결제, 블로그 등) - 사용자가 인지하는 가치 단위
- **task**: 실제 실행 단위
- **edge**: Task 간 의존성 (그래프 구조)

### 왜 트리가 아니라 그래프인가

트리 구조는 **수직 관계**만 표현한다. 하지만 실제 코드 의존성은 **수평 관계**가 더 많다:

```
SQL Table Task ←───┬─── Model Task
                   │
Auth Config Task ←─┴─── API Task
```

그래프 구조의 장점:
- **컨텍스트 정밀 주입**: 해당 Task + 의존 Task 결과만 주입
- **실행 순서 자동 결정**: Topological Sort로 의존성 해결된 Task부터 실행
- **병렬 처리 가능**: 의존성 없는 Task 동시 실행
- **토큰 최소화**: 전체 manifest 대신 필요한 것만

### Edge 제한

각 Task의 Edge는 **최대 4-7개**로 제한:
- 너무 많으면 Task 분할 필요 신호
- 컨텍스트 크기 예측 가능
- LLM 토큰 한계 내 관리 가능

---

## 데이터베이스 스키마

### projects
```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'active',
    created_at TEXT NOT NULL
);
```

### features
```sql
CREATE TABLE features (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    spec TEXT DEFAULT '',           -- Feature 상세 스펙 (LLM 대화로 수립)
    status TEXT DEFAULT 'pending'
        CHECK(status IN ('pending', 'active', 'done')),
    created_at TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

### feature_edges (Feature 간 의존성)
```sql
CREATE TABLE feature_edges (
    from_feature_id INTEGER NOT NULL,
    to_feature_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (from_feature_id, to_feature_id),
    FOREIGN KEY (from_feature_id) REFERENCES features(id),
    FOREIGN KEY (to_feature_id) REFERENCES features(id)
);
```

### tasks
```sql
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK(status IN ('pending', 'doing', 'done', 'failed')),
    title TEXT NOT NULL,
    content TEXT DEFAULT '',
    result TEXT DEFAULT '',         -- Task 완료 시 결과 (의존 Task에 전달됨)
    error TEXT DEFAULT '',
    created_at TEXT NOT NULL,
    started_at TEXT,
    completed_at TEXT,
    failed_at TEXT,
    FOREIGN KEY (feature_id) REFERENCES features(id)
);
```

### task_edges (Task 간 의존성)
```sql
CREATE TABLE task_edges (
    from_task_id INTEGER NOT NULL,
    to_task_id INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (from_task_id, to_task_id),
    FOREIGN KEY (from_task_id) REFERENCES tasks(id),
    FOREIGN KEY (to_task_id) REFERENCES tasks(id)
);
```

**Edge 의미**:
- `from_task_id` → `to_task_id`: from이 to에 의존
- Task 실행 시 의존 Task들의 `result`가 컨텍스트에 포함됨

### context (싱글톤)
```sql
CREATE TABLE context (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    data TEXT NOT NULL,  -- JSON
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
```

**JSON 포맷**:
```json
{
  "project_name": "Blog Platform",
  "description": "Developer blogging platform",
  "target_users": "Tech bloggers",
  "deadline": "2026-03-01"
}
```

### tech (싱글톤)
```sql
CREATE TABLE tech (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    data TEXT NOT NULL,  -- JSON
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
```

**JSON 포맷**:
```json
{
  "backend": "FastAPI",
  "frontend": "React",
  "database": "PostgreSQL",
  "cache": "Redis",
  "deployment": "Docker + AWS"
}
```

### design (싱글톤)
```sql
CREATE TABLE design (
    id INTEGER PRIMARY KEY CHECK(id = 1),
    data TEXT NOT NULL,  -- JSON
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
```

**JSON 포맷**:
```json
{
  "architecture": "Microservices",
  "auth_method": "JWT",
  "api_style": "RESTful",
  "db_schema_users": "id, email, password_hash, created_at",
  "caching_strategy": "Cache-aside"
}
```

### state (자동 관리)
```sql
CREATE TABLE state (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

**자동 저장 항목**:
- `current_project`: 현재 프로젝트 ID
- `current_feature`: 현재 feature ID
- `current_task`: 현재 task ID

**관리**: Task 명령 실행 시 Claritask가 자동 업데이트

### memos
```sql
CREATE TABLE memos (
    scope TEXT NOT NULL,     -- 'project', 'feature', 'task'
    scope_id TEXT NOT NULL,  -- project_id, feature_id, task_id
    key TEXT NOT NULL,
    data TEXT NOT NULL,      -- JSON
    priority INTEGER DEFAULT 2
        CHECK(priority IN (1, 2, 3)),
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    PRIMARY KEY (scope, scope_id, key)
);
```

**영역**:
- `project`: 프로젝트 전역 메모
- `feature`: 특정 feature 메모
- `task`: 특정 task 메모

**Priority**:
- `1`: 중요 (manifest에 자동 포함)
- `2`: 보통
- `3`: 사소함

**JSON 포맷**:
```json
{
  "value": "실제 메모 내용",
  "summary": "간단한 요약 (선택)",
  "tags": ["tag1", "tag2"]
}
```

---

## 명령어 레퍼런스

### Project 관리
```bash
clari project '<json>' # 프로젝트 정보 입력. Claritask는 클로드 코드 내에서 작동하므로 프로젝트는 싱글톤.
```

**JSON 포맷**:
```json
{
  "name": "Blog Platform",
  "description": "Developer blogging platform",
  "context":{
    "project_name": "Blog Platform",
    "description": "Developer blogging platform with markdown",
    "target_users": "Tech bloggers and readers",
    "deadline": "2026-03-01",
    "constraints": "Must support 10k concurrent users"
  },
  "tech":{
    "backend": "FastAPI",
    "frontend": "React 18",
    "database": "PostgreSQL",
    "cache": "Redis",
    "auth_service": "Auth0",
    "deployment": "Docker + AWS ECS"
  },
  "design":{
    "architecture": "Microservices",
    "auth_method": "JWT with 1h expiry",
    "api_style": "RESTful",
    "db_schema_users": "id, email, password_hash, role, created_at",
    "caching_strategy": "Cache-aside pattern",
    "rate_limiting": "100 req/min per user"
  }
}
```

### Project 자동 실행 (Claritask 드라이버)
```bash
clari project start               # pending Task 전체 자동 실행
clari project start --feature 2   # 특정 Feature만 실행
clari project start --dry-run     # 실행 없이 Task 목록만 출력
clari project stop                # 실행 중단 (현재 Task 완료 후)
clari project status              # 실행 상태 조회
```

**동작 원리**:
1. Feature Edge 기반 Feature 실행 순서 결정
2. Feature 내 Task Edge 기반 Task 실행 순서 결정
3. 각 Task마다 `claude --print` 비대화형 호출
4. 의존 Task의 `result`를 컨텍스트에 포함
5. Task 완료/실패 처리
6. 모든 Task 완료 시 종료

**실패 처리**:
- Task 실패 시 해당 Task에서 멈추고 로그 출력
- `clari project start` 재실행 시 실패한 Task부터 재개

### Feature 관리
```bash
clari feature list             # Feature 목록 조회
clari feature add '<json>'     # Feature 등록
clari feature <id> spec        # Feature spec 대화 시작
clari feature <id> tasks       # Feature 하위 Task 생성
clari feature <id> start       # Feature 하위 Task 실행 시작
```

**JSON 포맷**:
```json
{
  "name": "로그인",
  "description": "사용자 인증 기능"
}
```

### Task 관리
```bash
clari task list                   # Task 목록 조회
clari task add '<json>'           # Task 추가
clari task pop                    # 다음 실행 가능 Task (의존성 해결된 것)
clari task start <task_id>        # pending → doing
clari task complete <task_id> '<json>'  # doing → done
clari task fail <task_id> '<json>'      # doing → failed
clari task status                 # 진행 상황
```

**add JSON 포맷**:
```json
{
  "feature_id": 1,
  "title": "user_table_sql",
  "content": "CREATE TABLE users ..."
}
```

**complete JSON 포맷**:
```json
{
  "result": "테이블 생성 완료. 컬럼: id, email, password_hash, created_at"
}
```

**fail JSON 포맷**:
```json
{
  "error": "Database connection failed",
  "details": "Connection timeout after 30s"
}
```

### Edge 관리
```bash
clari edge add --from <task_id> --to <task_id>     # Task 의존성 추가
clari edge add --feature --from <id> --to <id>     # Feature 의존성 추가
clari edge list                                     # 의존성 목록 조회
clari edge infer --feature <id>                    # Feature 내 Task Edge LLM 추론
clari edge infer --project                         # Feature 간 Edge LLM 추론
```

### Memo 관리
```bash
clari memo set '<json>'
clari memo get [feature]:[task]:<key>
clari memo list [feature]:[task]
clari memo del [feature]:[task]:<key>
```

**영역 지정**:
```bash
# Project 전역
clari memo get jwt_config

# Feature 귀속
clari memo get 1:api_decisions

# Task 귀속
clari memo get 1:42:implementation_notes
```

**JSON 포맷**:
```json
{
  "feature": 1,
  "task": 42,
  "key": "jwt_config",
  "value": "Use httpOnly cookies for refresh tokens",
  "priority": 1
}
```

**조회**:
```bash
# 전체
clari memo list

# Feature 메모만
clari memo list 1

# Task 메모만
clari memo list 1:42
```

**특징**:
- 덮어쓰기 가능 (최신 값으로 업데이트)
- 한 번만 설정하면 됨
- Task 반환 시 자동 포함 (manifest)

### 유틸리티
```bash
clari required                  # 필수 입력 중 입력하지 않은 항목 안내.
```

---

## Manifest 자동 반환

### pop 명령 응답

`clari task pop` 실행 시 Task + 의존 Task 결과 + Manifest 함께 반환

```json
{
  "task": {
    "id": 42,
    "feature_id": 2,
    "title": "auth_service",
    "content": "JWT 기반 인증 서비스 구현",
    "status": "pending"
  },
  "dependencies": [
    {
      "id": 41,
      "title": "user_model",
      "result": "User 모델 구현 완료. 필드: id, email, password_hash, created_at"
    }
  ],
  "manifest": {
    "context": {
      "project_name": "Blog Platform",
      "description": "Developer blogging platform"
    },
    "tech": {
      "backend": "Go",
      "frontend": "React",
      "database": "PostgreSQL"
    },
    "design": {
      "architecture": "Monolithic",
      "auth_method": "JWT",
      "api_style": "RESTful"
    },
    "feature": {
      "id": 2,
      "name": "로그인",
      "spec": "JWT 기반 인증. Access token 1시간, Refresh token 7일..."
    },
    "memos": [
      {
        "scope": "project",
        "key": "jwt_security",
        "value": "Use httpOnly cookies"
      },
      {
        "scope": "feature",
        "scope_id": 2,
        "key": "token_expiry",
        "value": "Access 1h, Refresh 7d"
      }
    ]
  }
}
```

**Manifest 포함 내용**:
1. `dependencies`: 의존 Task들의 `result` (핵심!)
2. `context`: 프로젝트 컨텍스트
3. `tech`: 기술 스택
4. `design`: 설계 결정
5. `feature`: 현재 Feature 정보 및 spec
6. `memos`: priority 1인 메모만

**장점**:
- 의존 Task 결과가 자동 주입 → 정보 누락 없음
- 컨텍스트 최소화 → 토큰 절약
- Feature spec으로 일관성 유지

---

## 필수 입력 시스템

### 필수 항목

**context** (필수):
- `project_name`
- `description`

**tech** (필수):
- `backend`
- `frontend`
- `database`

**design** (필수):
- `architecture`
- `auth_method`
- `api_style`

### 워크플로우

```
User: "clari plan tasks"
    ↓
Claude: clari required
    ↓
Claritask: Check required
    ↓
Missing → Return required items
    ↓
Claude: Interactive collection
    - 옵션 제시
    - 사용자 선택
    ↓
Claude: clari project '<json>'
    ↓
User: "clari plan tasks" (재시도)
    ↓
Claritask: Ready → Proceed
```

### 대화 예시

```
Claude: "프로젝트 설정이 필요합니다.

**1. 백엔드 프레임워크**
A) FastAPI - 빠르고 현대적
B) Django - 풀스택
C) Express - Node.js

추천: FastAPI"

User: "A"

Claude: [모든 필수 항목 수집 후]

clari context set '{
  "project_name": "Blog Platform",
  "description": "Developer blogging platform"
}'

clari tech set '{
  "backend": "FastAPI",
  "frontend": "React",
  "database": "PostgreSQL"
}'

clari design set '{
  "architecture": "Monolithic",
  "auth_method": "JWT",
  "api_style": "RESTful"
}'

"✅ 설정 완료! 이제 'clari plan tasks' 가능합니다."
```

---

## 워크플로우

### Planning Phase: 구조화된 LLM 호출

Planning 단계에서 LLM 호출을 구조화하여 호출 횟수를 최소화한다.

```
Project Description
        │
        ▼ (LLM 1회)
Feature 목록 산출
        │
        ▼ (LLM N회, 대화형)
Feature별 Spec 수립
        │
        ▼ (LLM 1회)
Feature 간 Edge 추출
        │
        ▼ (LLM N회)
Feature별 Task 생성
        │
        ▼ (LLM N회)
Feature별 Task Edge 추출
        │
        ▼
실행 준비 완료
```

**LLM 호출 횟수 (Feature 20개 기준)**:
| 단계 | 호출 수 |
|------|---------|
| Feature 목록 산출 | 1회 |
| Feature Spec 수립 | 20회 (대화) |
| Feature Edge 추출 | 1회 |
| Task 생성 | 20회 |
| Task Edge 추출 | 20회 |
| **총 Planning** | **~60회** |

### 1. 프로젝트 초기화

```bash
clari init blog-platform "개발자 블로그 플랫폼"

# 필수 설정
clari context set '{"project_name": "Blog Platform", ...}'
clari tech set '{"backend": "Go", "frontend": "React", ...}'
clari design set '{"architecture": "Monolithic", ...}'
```

### 2. Feature 목록 산출 (LLM 1회)

```bash
clari plan features
# → LLM이 Project description 기반으로 Feature 목록 산출
# → 결과: 로그인, 블로그, 댓글, 알림 등
```

### 3. Feature Spec 수립 (대화형, Feature별)

```bash
clari feature 1 spec
# → LLM과 대화하며 Feature spec 상세화
# → "로그인 방식은? JWT vs Session"
# → "소셜 로그인 필요?"
# → spec 저장
```

### 4. Edge 추출 (LLM 자동 추론)

```bash
# Feature 간 의존성 추출 (1회)
clari edge infer --project
# → LLM: "결제 Feature는 로그인 Feature에 의존"

# Feature 내 Task Edge 추출 (Feature별 1회)
clari edge infer --feature 1
# → LLM: "user_model은 user_table_sql에 의존"
```

**Edge 추론이 쉬운 이유**:
- Feature 내 Task는 5-15개 → LLM 컨텍스트에 충분히 들어감
- Feature 목록은 10-30개 → 한 번에 분석 가능
- LLM이 코드 의존성 패턴을 잘 이해함 (SQL → Model → Service → API)

### 5. 자동 실행 (Claritask 드라이버)

```bash
clari project start

# Claritask 내부 동작:
# 1. Feature Edge 기반 실행 순서 결정 (Topological Sort)
# 2. Feature 내 Task Edge 기반 실행 순서 결정

for feature in sorted_features:
    for task in sorted_tasks(feature):
        # 의존 Task의 result 수집
        deps = get_dependency_results(task)

        # Prompt 생성 (Task + 의존 결과 + Manifest)
        prompt = build_prompt(task, deps, manifest)

        # LLM 호출 (독립 컨텍스트)
        result = exec("claude --print", prompt)

        if result.success:
            save_result(task, result)  # result 저장 (다음 Task에 전달됨)
        else:
            mark_failed(task)
            break
```

### 6. 수동 실행 (탐색/디버깅용)

```bash
# 특정 Task만 실행
clari task pop
clari task start 42
# ... 작업 ...
clari task complete 42 '{"result": "구현 완료"}'
```

---

## Task Status

```
pending → doing → done/failed
```

**전이**:
- `clari task start`: pending → doing
- `clari task complete`: doing → done (result 저장)
- `clari task fail`: doing → failed

**크래시 복구**:
- 크래시 시 status='doing'으로 남음
- 재시작 후 감지 → 재개 가능

**result의 중요성**:
- Task 완료 시 `result`에 결과 요약 저장
- 의존하는 Task 실행 시 이 `result`가 컨텍스트에 포함됨
- 예: SQL Task의 result → Model Task에 전달

---

## 제약사항

### Task
- `title`, `content` 필수
- `feature_id` 필수
- Edge는 최대 4-7개 권장

### Feature
- `name`, `description` 필수
- `spec`은 LLM 대화로 수립

### Edge
- Task Edge: 같은 Feature 내 또는 Feature 경계 넘어 가능
- Feature Edge: Feature 간 의존성
- 순환 의존성 불가 (DAG)

### 필수 설정
- context: project_name, description
- tech: backend, frontend, database
- design: architecture, auth_method, api_style

### Memo
- 영역: project, feature, task
- Priority: 1 (중요), 2 (보통), 3 (사소함)
- JSON 포맷 필수

---

## 성능

| Tasks | JSON | SQLite |
|-------|------|--------|
| 100   | 10ms | 1ms |
| 1,000 | 150ms | 1ms |
| 10,000| 2.5s | 2ms |

---

## 설치 및 초기화

### 바이너리 설치
```bash
# Go로 빌드된 바이너리 설치
go install parkjunwoo.com/claritask/cmd/claritask@latest
```

### 프로젝트 초기화

```bash
clari init <project-id> ["<project-description>"]
```

**동작**:
1. 현재 위치에 `<project-id>` 폴더 생성
2. 폴더 내 `CLAUDE.md` 파일 생성 (기본 템플릿)
3. 폴더 내 `.claritask/db` SQLite 파일 생성
4. projects 테이블에 project id와 description 자동 입력

**예시**:
```bash
# description 없이
clari init blog-api

# description 포함
clari init blog-api "Developer blogging platform with markdown support"
```

**생성되는 구조**:
```
blog-api/
├── CLAUDE.md          # 프로젝트 설정 템플릿
└── .claritask/
    └── db             # SQLite 데이터베이스
```

---

## 핵심 가치

1. **제어 역전**: Claritask가 오케스트레이터, Claude는 실행기
2. **그래프 기반**: Task 간 의존성을 Edge로 명시, 정밀한 컨텍스트 주입
3. **컨텍스트 최소화**: 전체 manifest 대신 의존 Task result만 주입
4. **구조화된 Planning**: Feature 단위로 Edge 추론, LLM 호출 최소화
5. **무제한 확장**: Task 수천 개도 자동 처리
6. **복구 가능**: 실패 시 해당 Task부터 재개

**Claritask = LLM 컨텍스트 한계를 우회하는 프로젝트 실행 엔진**

사람은 Feature를 정의하고 시작 버튼만 누른다. Claritask가 의존성을 분석하고, 필요한 컨텍스트만 주입하며, Claude를 수천 번이고 호출해 작업을 완료한다. 컨텍스트 폭발도, 정보 누락도, 수동 개입도 없다.