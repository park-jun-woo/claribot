# Claritask Commands Reference

모든 Claritask 명령어의 상세 사용법

---

## 목차

0. [초기화](#초기화) (1개)
1. [Project 관리](#project-관리) (6개)
2. [Project 실행](#project-실행) (5개)
3. [Feature 관리](#feature-관리) (5개)
4. [Task 관리](#task-관리) (7개)
5. [Edge 관리](#edge-관리) (4개)
6. [Memo 관리](#memo-관리) (4개)

**총 32개 명령어**

---

## 초기화

### clari init

**설명**: 새 프로젝트 폴더와 Claritask 환경 초기화

**사용법**:
```bash
clari init <project-id> ["<project-description>"]
```

**인자**:
- `project-id` (필수): 프로젝트 ID (영문 소문자, 숫자, 하이픈, 언더스코어만 허용)
- `project-description` (선택): 프로젝트 설명

**동작**:
1. 현재 위치에 `<project-id>` 이름의 폴더 생성
2. 폴더 내 `CLAUDE.md` 파일 생성 (기본 템플릿)
3. 폴더 내 `.claritask/` 디렉토리 생성
4. `.claritask/db` SQLite 파일 생성 및 스키마 초기화
5. `projects` 테이블에 project id와 description 자동 입력

**생성되는 구조**:
```
<project-id>/
├── CLAUDE.md          # 프로젝트 설정 템플릿
└── .claritask/
    └── db             # SQLite 데이터베이스
```

**응답**:
```json
{
  "success": true,
  "project_id": "blog-api",
  "path": "/path/to/blog-api",
  "message": "Project initialized successfully"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Directory already exists: blog-api"
}
```

```json
{
  "success": false,
  "error": "Invalid project ID: Blog-API",
  "hint": "Use lowercase letters, numbers, hyphens, and underscores only"
}
```

**예시**:
```bash
# description 없이
clari init blog-api

# description 포함
clari init blog-api "Developer blogging platform with markdown support"

# 하이픈과 언더스코어 사용
clari init my_ecommerce-api "E-commerce REST API"
```

**CLAUDE.md 템플릿**:
```markdown
# <project-id>

## Description
<project-description>

## Tech Stack
- Backend:
- Frontend:
- Database:

## Commands
- `clari project set '<json>'` - 프로젝트 설정
- `clari required` - 필수 입력 확인
- `clari plan features` - Feature 목록 산출
- `clari project start` - 실행 시작
```

---

## Project 관리

### clari project set

**설명**: 프로젝트 생성 또는 전체 업데이트 (싱글톤)

**사용법**:
```bash
clari project set '<json>'
```

**JSON 포맷**:
```json
{
  "name": "Blog Platform",
  "description": "Developer blogging platform",
  "context": {
    "project_name": "Blog Platform",
    "description": "Developer blogging platform with markdown",
    "target_users": "Tech bloggers and readers",
    "deadline": "2026-03-01",
    "constraints": "Must support 10k concurrent users"
  },
  "tech": {
    "backend": "FastAPI",
    "frontend": "React 18",
    "database": "PostgreSQL",
    "cache": "Redis",
    "auth_service": "Auth0",
    "deployment": "Docker + AWS ECS"
  },
  "design": {
    "architecture": "Microservices",
    "auth_method": "JWT with 1h expiry",
    "api_style": "RESTful",
    "db_schema_users": "id, email, password_hash, role, created_at",
    "caching_strategy": "Cache-aside pattern",
    "rate_limiting": "100 req/min per user"
  }
}
```

**필수 필드**:
- `name`: 프로젝트 이름
- `context.project_name`: 프로젝트명
- `context.description`: 프로젝트 설명
- `tech.backend`: 백엔드 프레임워크
- `tech.frontend`: 프론트엔드 프레임워크
- `tech.database`: 데이터베이스
- `design.architecture`: 아키텍처 패턴
- `design.auth_method`: 인증 방식
- `design.api_style`: API 스타일

**응답**:
```json
{
  "success": true,
  "project_id": "P001",
  "message": "Project created successfully"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Missing required field: tech.backend",
  "required_fields": [
    "name",
    "context.project_name",
    "context.description",
    "tech.backend",
    "tech.frontend",
    "tech.database",
    "design.architecture",
    "design.auth_method",
    "design.api_style"
  ]
}
```

**예시**:
```bash
clari project set '{
  "name": "Blog Platform",
  "description": "A modern blogging platform",
  "context": {
    "project_name": "Blog Platform",
    "description": "Developer blogging with markdown",
    "target_users": "Tech bloggers"
  },
  "tech": {
    "backend": "FastAPI",
    "frontend": "React",
    "database": "PostgreSQL"
  },
  "design": {
    "architecture": "Monolithic",
    "auth_method": "JWT",
    "api_style": "RESTful"
  }
}'
```

---

### clari project get

**설명**: 프로젝트 정보 조회

**사용법**:
```bash
clari project get
```

**응답**:
```json
{
  "project": {
    "id": "P001",
    "name": "Blog Platform",
    "description": "Developer blogging platform",
    "status": "active",
    "created_at": "2026-02-02T10:00:00Z"
  },
  "context": {
    "project_name": "Blog Platform",
    "description": "Developer blogging platform",
    "target_users": "Tech bloggers"
  },
  "tech": {
    "backend": "FastAPI",
    "frontend": "React",
    "database": "PostgreSQL"
  },
  "design": {
    "architecture": "Microservices",
    "auth_method": "JWT",
    "api_style": "RESTful"
  }
}
```

**에러**:
```json
{
  "success": false,
  "error": "No project found"
}
```

---

### clari context set

**설명**: Context 정보만 수정

**사용법**:
```bash
clari context set '<json>'
```

**JSON 포맷**:
```json
{
  "project_name": "Blog Platform",
  "description": "Developer blogging platform with markdown",
  "target_users": "Tech bloggers and readers",
  "deadline": "2026-03-01",
  "constraints": "Must support 10k concurrent users",
  "business_goal": "Monthly revenue $10k"
}
```

**필수 필드**:
- `project_name`
- `description`

**응답**:
```json
{
  "success": true,
  "message": "Context updated successfully"
}
```

**예시**:
```bash
clari context set '{
  "project_name": "Blog Platform",
  "description": "Developer blogging platform",
  "deadline": "2026-04-01"
}'
```

---

### clari tech set

**설명**: Tech 정보만 수정

**사용법**:
```bash
clari tech set '<json>'
```

**JSON 포맷**:
```json
{
  "backend": "FastAPI",
  "frontend": "React 18",
  "database": "PostgreSQL",
  "cache": "Redis",
  "queue": "Celery",
  "auth_service": "Auth0",
  "payment_service": "Stripe",
  "email_service": "SendGrid",
  "storage": "AWS S3",
  "deployment": "Docker + AWS ECS",
  "ci_cd": "GitHub Actions",
  "monitoring": "Datadog"
}
```

**필수 필드**:
- `backend`
- `frontend`
- `database`

**응답**:
```json
{
  "success": true,
  "message": "Tech updated successfully"
}
```

**예시**:
```bash
clari tech set '{
  "backend": "FastAPI",
  "frontend": "React",
  "database": "PostgreSQL",
  "cache": "Redis"
}'
```

---

### clari design set

**설명**: Design 정보만 수정

**사용법**:
```bash
clari design set '<json>'
```

**JSON 포맷**:
```json
{
  "architecture": "Microservices",
  "auth_method": "JWT with 1h expiry",
  "api_style": "RESTful",
  "db_schema_users": "id, email, password_hash, role, created_at",
  "db_schema_posts": "id, user_id, title, content, published_at",
  "caching_strategy": "Cache-aside pattern",
  "rate_limiting": "100 req/min per user",
  "error_handling": "Centralized error handler",
  "logging": "Structured JSON logs"
}
```

**필수 필드**:
- `architecture`
- `auth_method`
- `api_style`

**응답**:
```json
{
  "success": true,
  "message": "Design updated successfully"
}
```

**예시**:
```bash
clari design set '{
  "architecture": "Monolithic",
  "auth_method": "JWT",
  "api_style": "RESTful"
}'
```

---

### clari required

**설명**: 필수 입력 중 누락된 항목 확인

**사용법**:
```bash
clari required
```

**응답** (누락 시):
```json
{
  "ready": false,
  "missing_required": [
    {
      "field": "context.project_name",
      "prompt": "What is the project name?",
      "examples": ["Blog Platform", "E-commerce API"]
    },
    {
      "field": "tech.backend",
      "prompt": "What backend framework?",
      "options": ["FastAPI", "Django", "Flask", "Express"],
      "custom_allowed": true
    },
    {
      "field": "design.architecture",
      "prompt": "What architecture pattern?",
      "options": ["Monolithic", "Microservices", "Serverless"],
      "custom_allowed": true
    }
  ],
  "total_missing": 3,
  "message": "Please configure required settings"
}
```

**응답** (완료 시):
```json
{
  "ready": true,
  "message": "All required fields configured"
}
```

---

## Project 실행

### clari plan features

**설명**: LLM을 통해 프로젝트 설명 기반 Feature 목록 산출

**사용법**:
```bash
clari plan features
```

**동작**:
1. 필수 설정 확인
2. Project description을 LLM에 전달
3. Feature 목록 자동 산출
4. 결과를 features 테이블에 저장

**응답**:
```json
{
  "success": true,
  "features": [
    {"id": 1, "name": "로그인", "description": "사용자 인증 기능"},
    {"id": 2, "name": "결제", "description": "상품 결제 기능"},
    {"id": 3, "name": "블로그", "description": "포스트 작성 및 관리"}
  ],
  "total": 3,
  "message": "Features generated successfully"
}
```

**에러** (필수 누락):
```json
{
  "success": false,
  "ready": false,
  "missing_required": [...],
  "message": "Configure required settings first"
}
```

---

### clari project start

**설명**: 전체 프로젝트 자동 실행 (Claritask 드라이버)

**사용법**:
```bash
clari project start                # pending Task 전체 자동 실행
clari project start --feature 2    # 특정 Feature만 실행
clari project start --dry-run      # 실행 없이 Task 목록만 출력
```

**옵션**:
- `--feature <id>`: 특정 Feature의 Task만 실행
- `--dry-run`: 실제 실행 없이 실행 순서만 출력

**동작**:
1. Feature Edge 기반 Feature 실행 순서 결정 (Topological Sort)
2. Feature 내 Task Edge 기반 Task 실행 순서 결정
3. 각 Task마다 `claude --print` 비대화형 호출
4. 의존 Task의 `result`를 컨텍스트에 포함
5. Task 완료/실패 처리
6. 모든 Task 완료 시 종료

**응답**:
```json
{
  "success": true,
  "mode": "execution",
  "message": "Execution started",
  "total_tasks": 42,
  "next_task": {
    "id": 1,
    "title": "user_table_sql"
  }
}
```

**응답** (--dry-run):
```json
{
  "success": true,
  "mode": "dry-run",
  "execution_order": [
    {"id": 1, "title": "user_table_sql", "feature": "로그인"},
    {"id": 2, "title": "user_model", "feature": "로그인", "depends_on": [1]},
    {"id": 3, "title": "auth_service", "feature": "로그인", "depends_on": [2]}
  ],
  "total_tasks": 3
}
```

**에러** (task 없음):
```json
{
  "success": false,
  "error": "No pending tasks",
  "message": "All tasks completed or no tasks created"
}
```

---

### clari project stop

**설명**: 실행 중인 프로젝트 중단 (현재 Task 완료 후)

**사용법**:
```bash
clari project stop
```

**응답**:
```json
{
  "success": true,
  "message": "Execution will stop after current task completes",
  "current_task": {
    "id": 5,
    "title": "auth_service"
  }
}
```

**에러**:
```json
{
  "success": false,
  "error": "No execution in progress"
}
```

---

### clari project status

**설명**: 프로젝트 실행 상태 조회

**사용법**:
```bash
clari project status
```

**응답**:
```json
{
  "execution": {
    "running": true,
    "started_at": "2026-02-02T10:00:00Z",
    "elapsed": "1h 30m"
  },
  "progress": {
    "total_features": 5,
    "completed_features": 2,
    "total_tasks": 42,
    "pending": 25,
    "doing": 1,
    "done": 15,
    "failed": 1
  },
  "current_task": {
    "id": 16,
    "title": "payment_api",
    "feature": "결제",
    "started_at": "2026-02-02T11:25:00Z"
  },
  "failed_tasks": [
    {
      "id": 10,
      "title": "external_auth_integration",
      "error": "API timeout"
    }
  ]
}
```

---

### clari project plan

**설명**: 전체 프로젝트 플래닝 시작 (대화형)

**사용법**:
```bash
clari project plan
```

**동작**:
1. 필수 설정 확인
2. Feature 생성 대화
3. 각 Feature별 Spec 수립
4. Task 생성
5. Edge 추론

**응답**:
```json
{
  "success": true,
  "mode": "planning",
  "message": "Planning mode started",
  "next_steps": [
    "Generate features: clari plan features",
    "Set feature specs: clari feature <id> spec",
    "Generate tasks: clari feature <id> tasks",
    "Infer edges: clari edge infer --project"
  ]
}
```

---

## Feature 관리

### clari feature list

**설명**: Feature 목록 조회

**사용법**:
```bash
clari feature list
```

**응답**:
```json
{
  "features": [
    {
      "id": 1,
      "name": "로그인",
      "description": "사용자 인증 기능",
      "spec": "JWT 기반 인증. Access token 1시간...",
      "status": "done",
      "tasks_total": 4,
      "tasks_done": 4
    },
    {
      "id": 2,
      "name": "결제",
      "description": "상품 결제 기능",
      "status": "active",
      "tasks_total": 3,
      "tasks_done": 1,
      "depends_on": [1]
    },
    {
      "id": 3,
      "name": "블로그",
      "description": "포스트 작성 및 관리",
      "status": "pending",
      "tasks_total": 3,
      "tasks_done": 0
    }
  ],
  "total": 3
}
```

---

### clari feature add

**설명**: Feature 추가

**사용법**:
```bash
clari feature add '<json>'
```

**JSON 포맷**:
```json
{
  "name": "로그인",
  "description": "사용자 인증 기능"
}
```

**필수 필드**:
- `name`
- `description`

**응답**:
```json
{
  "success": true,
  "feature_id": 1,
  "name": "로그인",
  "message": "Feature created successfully"
}
```

**예시**:
```bash
clari feature add '{
  "name": "결제",
  "description": "Stripe 기반 상품 결제 기능"
}'
```

---

### clari feature \<id\> spec

**설명**: Feature spec 대화 시작 (LLM과 상세 스펙 수립)

**사용법**:
```bash
clari feature 1 spec
```

**동작**:
1. Feature 정보 조회
2. LLM과 대화형으로 spec 상세화
3. 결과를 feature.spec에 저장

**응답**:
```json
{
  "success": true,
  "feature_id": 1,
  "feature_name": "로그인",
  "mode": "spec_conversation",
  "message": "Spec conversation started",
  "prompts": [
    "로그인 방식은? (JWT vs Session)",
    "소셜 로그인 필요?",
    "토큰 만료 시간은?"
  ]
}
```

**spec 저장 후 응답**:
```json
{
  "success": true,
  "feature_id": 1,
  "spec": "JWT 기반 인증. Access token 1시간, Refresh token 7일. 소셜 로그인 미지원.",
  "message": "Spec saved successfully"
}
```

---

### clari feature \<id\> tasks

**설명**: Feature 하위 Task 생성 (LLM 추론)

**사용법**:
```bash
clari feature 1 tasks
```

**동작**:
1. Feature spec 확인
2. LLM이 Task 목록 산출
3. tasks 테이블에 저장

**응답**:
```json
{
  "success": true,
  "feature_id": 1,
  "feature_name": "로그인",
  "tasks_created": [
    {"id": 1, "title": "user_table_sql"},
    {"id": 2, "title": "user_model"},
    {"id": 3, "title": "auth_service"},
    {"id": 4, "title": "login_api"}
  ],
  "total": 4,
  "message": "Tasks generated successfully"
}
```

**에러** (spec 없음):
```json
{
  "success": false,
  "error": "Feature spec not set",
  "message": "Run 'clari feature 1 spec' first"
}
```

---

### clari feature \<id\> start

**설명**: 특정 Feature 실행 시작

**사용법**:
```bash
clari feature 2 start
```

**동작**:
1. Feature의 pending Task만 실행
2. Task Edge 기반 실행 순서 결정
3. Feature 완료 시 다음 Feature로 이동 안 함 (명시적 호출 필요)

**응답**:
```json
{
  "success": true,
  "feature_id": 2,
  "feature_name": "결제",
  "mode": "execution",
  "next_task": {
    "id": 5,
    "title": "payment_table_sql"
  },
  "message": "Feature execution started"
}
```

**에러** (task 없음):
```json
{
  "success": false,
  "error": "No pending tasks in feature 2"
}
```

**에러** (의존성 미해결):
```json
{
  "success": false,
  "error": "Feature dependencies not resolved",
  "blocked_by": [
    {"id": 1, "name": "로그인", "status": "active"}
  ]
}
```

---

## Task 관리

### clari task list

**설명**: Task 목록 조회

**사용법**:
```bash
clari task list                    # 전체 Task
clari task list --feature 1        # 특정 Feature의 Task
clari task list --status pending   # 특정 상태의 Task
```

**옵션**:
- `--feature <id>`: 특정 Feature의 Task만 조회
- `--status <status>`: 특정 상태의 Task만 조회 (pending, doing, done, failed)

**응답**:
```json
{
  "tasks": [
    {
      "id": 1,
      "feature_id": 1,
      "title": "user_table_sql",
      "status": "done",
      "depends_on": []
    },
    {
      "id": 2,
      "feature_id": 1,
      "title": "user_model",
      "status": "done",
      "depends_on": [1]
    },
    {
      "id": 3,
      "feature_id": 1,
      "title": "auth_service",
      "status": "doing",
      "depends_on": [2]
    }
  ],
  "total": 3
}
```

---

### clari task add

**설명**: Task 추가

**사용법**:
```bash
clari task add '<json>'
```

**JSON 포맷**:
```json
{
  "feature_id": 1,
  "title": "user_table_sql",
  "content": "CREATE TABLE users (id, email, password_hash, created_at)"
}
```

**필수 필드**:
- `feature_id`
- `title`
- `content`

**응답**:
```json
{
  "success": true,
  "task_id": 1,
  "title": "user_table_sql",
  "message": "Task created successfully"
}
```

**에러** (validation 실패):
```json
{
  "success": false,
  "error": "Missing required field: title",
  "required_fields": ["feature_id", "title", "content"]
}
```

**예시**:
```bash
clari task add '{
  "feature_id": 1,
  "title": "auth_service",
  "content": "JWT 기반 인증 서비스 구현. login, logout, refresh 메서드 포함."
}'
```

---

### clari task pop

**설명**: 다음 실행 가능 Task 조회 (의존성 해결된 것만)

**사용법**:
```bash
clari task pop
```

**응답**:
```json
{
  "task": {
    "id": 3,
    "feature_id": 1,
    "title": "auth_service",
    "content": "JWT 기반 인증 서비스 구현",
    "status": "pending",
    "created_at": "2026-02-02T10:00:00Z"
  },
  "dependencies": [
    {
      "id": 1,
      "title": "user_table_sql",
      "result": "테이블 생성 완료. 컬럼: id, email, password_hash, created_at"
    },
    {
      "id": 2,
      "title": "user_model",
      "result": "User 모델 구현 완료. 필드: id, email, password_hash, created_at"
    }
  ],
  "manifest": {
    "context": {
      "project_name": "Blog Platform",
      "description": "Developer blogging platform",
      "target_users": "Tech bloggers"
    },
    "tech": {
      "backend": "FastAPI",
      "frontend": "React",
      "database": "PostgreSQL"
    },
    "design": {
      "architecture": "Monolithic",
      "auth_method": "JWT",
      "api_style": "RESTful"
    },
    "feature": {
      "id": 1,
      "name": "로그인",
      "spec": "JWT 기반 인증. Access token 1시간, Refresh token 7일..."
    },
    "memos": [
      {
        "scope": "project",
        "key": "jwt_security",
        "value": "Use httpOnly cookies for refresh tokens"
      },
      {
        "scope": "feature",
        "scope_id": 1,
        "key": "token_expiry",
        "value": "Access 1h, Refresh 7d"
      }
    ]
  }
}
```

**응답** (task 없음):
```json
{
  "success": false,
  "error": "No pending tasks",
  "message": "All tasks completed or blocked by dependencies"
}
```

**Manifest 포함 내용**:
1. `dependencies`: 의존 Task들의 `result` (핵심!)
2. `context`: 프로젝트 컨텍스트
3. `tech`: 기술 스택
4. `design`: 설계 결정
5. `feature`: 현재 Feature 정보 및 spec
6. `memos`: priority 1인 메모만

---

### clari task start

**설명**: Task 실행 시작 (pending → doing)

**사용법**:
```bash
clari task start <task_id>
```

**동작**:
1. status: pending → doing
2. started_at 기록
3. state 자동 업데이트

**응답**:
```json
{
  "success": true,
  "task_id": 3,
  "status": "doing",
  "started_at": "2026-02-02T10:30:00Z",
  "message": "Task started"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Task not found: 999"
}
```

```json
{
  "success": false,
  "error": "Task already started",
  "current_status": "doing"
}
```

```json
{
  "success": false,
  "error": "Task dependencies not resolved",
  "blocked_by": [
    {"id": 2, "title": "user_model", "status": "pending"}
  ]
}
```

**예시**:
```bash
clari task start 3
```

---

### clari task complete

**설명**: Task 완료 처리 (doing → done)

**사용법**:
```bash
clari task complete <task_id> '<json>'
```

**JSON 포맷**:
```json
{
  "result": "인증 서비스 구현 완료. login(), logout(), refresh() 메서드 포함. JWT 토큰 발급 및 검증 로직 구현."
}
```

**필수 필드**:
- `result`: 결과 요약 (의존하는 Task에 전달됨)

**result의 중요성**:
- Task 완료 시 `result`에 결과 요약 저장
- 의존하는 Task 실행 시 이 `result`가 `dependencies`에 포함됨
- 다음 Task가 이전 Task의 결과를 알 수 있음

**응답**:
```json
{
  "success": true,
  "task_id": 3,
  "status": "done",
  "completed_at": "2026-02-02T13:00:00Z",
  "message": "Task completed successfully"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Task not in doing status",
  "current_status": "pending"
}
```

**예시**:
```bash
clari task complete 3 '{
  "result": "auth_service 구현 완료. JWT 기반 login/logout/refresh 구현. 토큰 만료: access 1h, refresh 7d"
}'
```

---

### clari task fail

**설명**: Task 실패 처리 (doing → failed)

**사용법**:
```bash
clari task fail <task_id> '<json>'
```

**JSON 포맷**:
```json
{
  "error": "Database connection failed",
  "details": "Connection timeout after 30s. Database server unreachable."
}
```

**필수 필드**:
- `error`: 에러 설명

**응답**:
```json
{
  "success": true,
  "task_id": 3,
  "status": "failed",
  "failed_at": "2026-02-02T11:00:00Z",
  "message": "Task marked as failed"
}
```

**예시**:
```bash
clari task fail 3 '{
  "error": "External API timeout",
  "details": "Auth0 API did not respond within 30s"
}'
```

---

### clari task status

**설명**: 전체 Task 진행 상황 조회

**사용법**:
```bash
clari task status
```

**응답**:
```json
{
  "summary": {
    "total": 42,
    "pending": 25,
    "doing": 1,
    "done": 15,
    "failed": 1
  },
  "progress": 35.7,
  "by_feature": [
    {
      "id": 1,
      "name": "로그인",
      "total": 4,
      "done": 4,
      "progress": 100
    },
    {
      "id": 2,
      "name": "결제",
      "total": 3,
      "done": 1,
      "progress": 33.3
    }
  ],
  "current_task": {
    "id": 6,
    "title": "payment_model",
    "feature": "결제",
    "status": "doing",
    "started_at": "2026-02-02T10:30:00Z"
  },
  "failed_tasks": [
    {
      "id": 10,
      "title": "external_auth",
      "error": "API timeout"
    }
  ]
}
```

---

## Edge 관리

### clari edge add

**설명**: 의존성(Edge) 추가

**사용법**:
```bash
# Task 간 의존성
clari edge add --from <task_id> --to <task_id>

# Feature 간 의존성
clari edge add --feature --from <feature_id> --to <feature_id>
```

**의미**:
- `--from A --to B`: A가 B에 의존 (B가 먼저 완료되어야 함)

**응답** (Task Edge):
```json
{
  "success": true,
  "type": "task",
  "from_id": 3,
  "to_id": 2,
  "message": "Task edge created: auth_service depends on user_model"
}
```

**응답** (Feature Edge):
```json
{
  "success": true,
  "type": "feature",
  "from_id": 2,
  "to_id": 1,
  "message": "Feature edge created: 결제 depends on 로그인"
}
```

**에러** (순환 의존성):
```json
{
  "success": false,
  "error": "Circular dependency detected",
  "cycle": [3, 2, 1, 3]
}
```

**예시**:
```bash
# user_model이 user_table_sql에 의존
clari edge add --from 2 --to 1

# 결제 Feature가 로그인 Feature에 의존
clari edge add --feature --from 2 --to 1
```

---

### clari edge list

**설명**: 의존성 목록 조회

**사용법**:
```bash
clari edge list                    # 전체
clari edge list --feature          # Feature Edge만
clari edge list --task             # Task Edge만
clari edge list --feature 1        # 특정 Feature 내 Task Edge
```

**응답**:
```json
{
  "feature_edges": [
    {
      "from": {"id": 2, "name": "결제"},
      "to": {"id": 1, "name": "로그인"}
    }
  ],
  "task_edges": [
    {
      "from": {"id": 2, "title": "user_model"},
      "to": {"id": 1, "title": "user_table_sql"}
    },
    {
      "from": {"id": 3, "title": "auth_service"},
      "to": {"id": 2, "title": "user_model"}
    },
    {
      "from": {"id": 6, "title": "payment_model"},
      "to": {"id": 5, "title": "payment_table_sql"}
    },
    {
      "from": {"id": 6, "title": "payment_model"},
      "to": {"id": 2, "title": "user_model"}
    }
  ],
  "total_feature_edges": 1,
  "total_task_edges": 4
}
```

---

### clari edge infer --feature

**설명**: Feature 내 Task Edge LLM 추론

**사용법**:
```bash
clari edge infer --feature <feature_id>
```

**동작**:
1. Feature 내 Task 목록을 LLM에 전달
2. LLM이 Task 간 의존성 분석
3. Edge 목록 반환 및 저장

**응답**:
```json
{
  "success": true,
  "feature_id": 1,
  "feature_name": "로그인",
  "edges_created": [
    {"from": 2, "to": 1, "reason": "user_model needs user_table_sql schema"},
    {"from": 3, "to": 2, "reason": "auth_service uses user_model"},
    {"from": 4, "to": 3, "reason": "login_api calls auth_service"}
  ],
  "total": 3,
  "message": "Task edges inferred successfully"
}
```

**예시**:
```bash
clari edge infer --feature 1
clari edge infer --feature 2
```

---

### clari edge infer --project

**설명**: Feature 간 Edge LLM 추론

**사용법**:
```bash
clari edge infer --project
```

**동작**:
1. 전체 Feature 목록을 LLM에 전달
2. LLM이 Feature 간 의존성 분석
3. Feature Edge 목록 반환 및 저장

**응답**:
```json
{
  "success": true,
  "edges_created": [
    {"from": 2, "to": 1, "reason": "결제 requires 로그인 for user authentication"},
    {"from": 3, "to": 1, "reason": "블로그 requires 로그인 for post authorship"}
  ],
  "total": 2,
  "message": "Feature edges inferred successfully"
}
```

---

## Memo 관리

### clari memo set

**설명**: Memo 저장

**사용법**:
```bash
# Project 전역
clari memo set '<json>'

# Feature 귀속
clari memo set <feature_id>:<key> '<json>'

# Task 귀속
clari memo set <feature_id>:<task_id>:<key> '<json>'
```

**JSON 포맷**:
```json
{
  "key": "jwt_security",
  "value": "Use httpOnly cookies for refresh tokens to prevent XSS attacks",
  "priority": 1
}
```

**필수 필드**:
- `key`: 메모 키
- `value`: 메모 내용

**선택 필드**:
- `priority`: 1 (중요), 2 (보통), 3 (사소함) - 기본: 2
- `feature`: Feature ID (영역 지정용)
- `task`: Task ID (영역 지정용)

**Priority 의미**:
- `1`: 중요 - manifest에 자동 포함
- `2`: 보통
- `3`: 사소함

**응답**:
```json
{
  "success": true,
  "scope": "project",
  "key": "jwt_security",
  "priority": 1,
  "message": "Memo saved successfully"
}
```

**예시**:
```bash
# Project 전역
clari memo set '{
  "key": "jwt_best_practice",
  "value": "Always use httpOnly cookies",
  "priority": 1
}'

# Feature 귀속
clari memo set 1:api_conventions '{
  "value": "Use plural nouns for REST resources",
  "priority": 1
}'

# Task 귀속
clari memo set 1:3:implementation_notes '{
  "value": "Used bcrypt with 12 rounds for password hashing",
  "priority": 2
}'
```

---

### clari memo get

**설명**: Memo 조회

**사용법**:
```bash
# Project 전역
clari memo get <key>

# Feature 귀속
clari memo get <feature_id>:<key>

# Task 귀속
clari memo get <feature_id>:<task_id>:<key>
```

**응답**:
```json
{
  "scope": "project",
  "key": "jwt_security",
  "data": {
    "value": "Use httpOnly cookies for refresh tokens",
    "priority": 1
  },
  "created_at": "2026-02-02T10:00:00Z",
  "updated_at": "2026-02-02T10:00:00Z"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Memo not found: jwt_security"
}
```

**예시**:
```bash
clari memo get jwt_best_practice
clari memo get 1:api_conventions
clari memo get 1:3:implementation_notes
```

---

### clari memo list

**설명**: Memo 목록 조회

**사용법**:
```bash
# 전체
clari memo list

# Feature 메모만
clari memo list <feature_id>

# Task 메모만
clari memo list <feature_id>:<task_id>
```

**응답** (전체):
```json
{
  "memos": {
    "project": [
      {
        "key": "jwt_security",
        "priority": 1,
        "value": "Use httpOnly cookies...",
        "created_at": "2026-02-02T10:00:00Z"
      }
    ],
    "feature": {
      "1": [
        {
          "key": "api_conventions",
          "priority": 1,
          "value": "Use plural nouns...",
          "created_at": "2026-02-02T12:00:00Z"
        }
      ]
    },
    "task": {
      "1:3": [
        {
          "key": "implementation_notes",
          "priority": 2,
          "value": "Used bcrypt...",
          "created_at": "2026-02-02T13:00:00Z"
        }
      ]
    }
  },
  "total": 3
}
```

**예시**:
```bash
clari memo list
clari memo list 1
clari memo list 1:3
```

---

### clari memo del

**설명**: Memo 삭제

**사용법**:
```bash
# Project 전역
clari memo del <key>

# Feature 귀속
clari memo del <feature_id>:<key>

# Task 귀속
clari memo del <feature_id>:<task_id>:<key>
```

**응답**:
```json
{
  "success": true,
  "scope": "project",
  "key": "jwt_security",
  "message": "Memo deleted successfully"
}
```

**에러**:
```json
{
  "success": false,
  "error": "Memo not found: jwt_security"
}
```

**예시**:
```bash
clari memo del old_note
clari memo del 1:deprecated
clari memo del 1:3:temp_note
```

---

## 명령어 요약

### 초기화 (1개)
```bash
clari init <id> ["<desc>"]         # 프로젝트 폴더 및 환경 초기화
```

### Project (6개)
```bash
clari project set '<json>'         # 프로젝트 생성/업데이트
clari project get                  # 프로젝트 조회
clari context set '<json>'         # Context 수정
clari tech set '<json>'            # Tech 수정
clari design set '<json>'          # Design 수정
clari required                     # 필수 입력 확인
```

### Project 실행 (5개)
```bash
clari plan features                # Feature 목록 산출 (LLM)
clari project plan                 # 전체 플래닝
clari project start [--feature N]  # 전체/Feature 실행
clari project stop                 # 실행 중단
clari project status               # 실행 상태 조회
```

### Feature (5개)
```bash
clari feature list                 # Feature 목록
clari feature add '<json>'         # Feature 추가
clari feature <id> spec            # Feature spec 대화
clari feature <id> tasks           # Task 생성 (LLM)
clari feature <id> start           # Feature 실행
```

### Task (7개)
```bash
clari task list [--feature N]      # Task 목록
clari task add '<json>'            # Task 추가
clari task pop                     # 다음 실행 가능 Task
clari task start <id>              # Task 시작
clari task complete <id> '<json>'  # Task 완료
clari task fail <id> '<json>'      # Task 실패
clari task status                  # 진행 상황
```

### Edge (4개)
```bash
clari edge add --from <id> --to <id>       # Task Edge 추가
clari edge add --feature --from <id> --to <id>  # Feature Edge 추가
clari edge list                            # Edge 목록
clari edge infer --feature <id>            # Task Edge 추론 (LLM)
clari edge infer --project                 # Feature Edge 추론 (LLM)
```

### Memo (4개)
```bash
clari memo set '<json>'            # Memo 저장
clari memo get <key>               # Memo 조회
clari memo list [<scope>]          # Memo 목록
clari memo del <key>               # Memo 삭제
```

---

## 일반적인 워크플로우

### 0. 프로젝트 초기화
```bash
# 새 프로젝트 생성
clari init my-project "My awesome project"

# 생성된 폴더로 이동
cd my-project
```

### 1. 프로젝트 초기 설정
```bash
# 필수 항목 확인
clari required

# 전체 설정 (한 번에)
clari project set '{...}'

# 또는 개별 설정
clari context set '{...}'
clari tech set '{...}'
clari design set '{...}'
```

### 2. Planning (구조화된 LLM 호출)
```bash
# Feature 목록 산출 (LLM 1회)
clari plan features

# Feature별 Spec 수립 (대화형, Feature당 1회)
clari feature 1 spec
clari feature 2 spec

# Feature 간 Edge 추론 (LLM 1회)
clari edge infer --project

# Feature별 Task 생성 (Feature당 1회)
clari feature 1 tasks
clari feature 2 tasks

# Feature별 Task Edge 추론 (Feature당 1회)
clari edge infer --feature 1
clari edge infer --feature 2
```

### 3. Execution (자동화 모드)
```bash
# 전체 자동 실행
clari project start

# 또는 특정 Feature만
clari project start --feature 1

# 실행 전 순서 확인
clari project start --dry-run

# 실행 상태 확인
clari project status

# 필요시 중단
clari project stop
```

### 4. Execution (수동 모드)
```bash
# Task 조회 (의존성 해결된 것만)
clari task pop

# Task 시작
clari task start 3

# ... 작업 ...

# Task 완료 (result 중요!)
clari task complete 3 '{
  "result": "auth_service 구현 완료. JWT 기반 인증."
}'
```

### 5. Memo 활용
```bash
# 중요한 발견 저장 (priority 1)
clari memo set '{
  "key": "security_note",
  "value": "Always validate JWT signature",
  "priority": 1
}'

# 다음 task pop 시 manifest에 자동 포함됨
```

---

## 에러 처리

모든 명령어는 다음 형식으로 에러 반환:

```json
{
  "success": false,
  "error": "에러 메시지",
  "details": "상세 설명 (선택)",
  "code": "ERROR_CODE (선택)"
}
```

**일반적인 에러 코드**:
- `MISSING_REQUIRED`: 필수 필드 누락
- `INVALID_JSON`: JSON 파싱 실패
- `NOT_FOUND`: 리소스 없음
- `INVALID_STATUS`: 잘못된 상태 전이
- `VALIDATION_ERROR`: Validation 실패
- `CIRCULAR_DEPENDENCY`: 순환 의존성 감지
- `DEPENDENCY_NOT_RESOLVED`: 의존성 미해결

---

**Claritask Commands Reference v2.0**
