# Claribot v0.2

> **버전**: v0.2.1

---

## 핵심 개념

**Claribot = 장기 주행 컨텍스트 오케스트레이터**

- Claude Code 단독: 컨텍스트 창 안에서 해결 가능한 작업
- Claribot: 컨텍스트 창을 넘어서는 거대 작업의 분할 정복

### Task 기반 분할 정복

1. 사용자가 메시지 입력 (→ 최상위 Task로 등록)
2. claribot이 Claude Code에게 Task 분할 지시
3. Claude Code가 하위 Task들 등록
4. 생성된 Task들을 순서대로 실행
5. Task 수행 중 추가 하위 Task 등록 가능
6. 모든 Task 완료 → 결과 반환

### 핵심 기술: 그래프 기반 컨텍스트 선별 주입

- 전체 히스토리 덤프가 아닌 Edge로 연결된 Task만 주입
- O(n) → O(k) 컨텍스트 사용량 최적화

---

## 아키텍처

```
┌─────────────────────────────────────────────────────────┐
│                    claribot (daemon)                    │
│                                                         │
│  ┌───────────┐  ┌───────────┐  ┌───────────────────┐   │
│  │ Telegram  │  │    CLI    │  │    TTY Manager    │   │
│  │  Handler  │  │  Handler  │  │  (Claude Code)    │   │
│  └─────┬─────┘  └─────┬─────┘  └─────────┬─────────┘   │
│        │              │                   │             │
│        └──────────────┼───────────────────┘             │
│                       ▼                                 │
│               ┌──────────────┐                          │
│               │      DB      │                          │
│               └──────────────┘                          │
└─────────────────────────────────────────────────────────┘
         ▲                              ▲
         │ 메시지                        │ HTTP
    [Telegram]                     [clari CLI]
```

### 컴포넌트

| 컴포넌트 | 역할 | 실행 |
|----------|------|------|
| claribot | Telegram + CLI 핸들러 + TTY 관리 | systemctl 서비스 |
| clari CLI | 서비스에 HTTP 요청 | 일시 실행 |
| Claude Code (전역) | ~/.claribot/에서 실행, 라우터 | TTY → kill |
| Claude Code (프로젝트) | 프로젝트 폴더에서 실행, 실제 작업 | TTY → kill |

---

## 프로젝트 구조

```
claribot/
├── Makefile
├── claribot.service.template
├── bot/
│   ├── cmd/claribot/main.go
│   ├── internal/db/db.go
│   └── pkg/
│       ├── claude/claude.go
│       └── telegram/telegram.go
└── cli/
    └── cmd/clari/main.go
```

---

## DB 스키마

### 전역 DB (`~/.claribot/db.clt`)

```sql
projects (
    id TEXT PRIMARY KEY,      -- 'blog', 'api-server'
    name TEXT,
    path TEXT UNIQUE,         -- 프로젝트 경로
    description TEXT,
    status TEXT,              -- 'active', 'archived'
    created_at, updated_at
)
```

### 로컬 DB (`프로젝트/.claribot/db.clt`)

```sql
tasks (
    id INTEGER PRIMARY KEY,
    parent_id INTEGER,        -- NULL이면 최상위 Task (=사용자 메시지)
    source TEXT,              -- 'telegram', 'cli', 'agent', ''
    title TEXT,
    content TEXT,
    status TEXT,              -- 'pending', 'running', 'done', 'failed'
    result TEXT,
    error TEXT,
    created_at, started_at, completed_at
)

task_edges (
    from_task_id INTEGER,     -- 선행 작업
    to_task_id INTEGER,       -- 후행 작업
    created_at
)
```

**핵심 단순화**: Message를 별도 테이블로 두지 않고 Task로 통합. `parent_id=NULL`이면 최상위 Task.

---

## 설정

### 파일 위치

```
~/.claribot/
├── config.yaml              -- 서비스 설정
└── db.clt                   -- 전역 DB
```

### config.yaml

```yaml
service:
  port: 9847
  host: 127.0.0.1

telegram:
  token: "BOT_TOKEN"
  allowed_users:
    - 123456789

claude:
  timeout: 1200              # idle timeout (초)
  max: 3                     # 동시 실행 최대 개수
```

---

## 설치

```bash
# 빌드 및 설치
make install

# 서비스 관리
make status    # 상태 확인
make restart   # 재시작
make logs      # 로그 확인

# 제거
make uninstall
```

---

## 텔레그램 문법

| 기호 | 용도 | 예시 |
|------|------|------|
| `$` | 프로젝트 지정 | `$blog 태스크 뭐있어?` |
| `!` | clari 명령어 | `$blog !task list` |

```
$blog 버그 수정해줘          # 프로젝트 + 자연어
$api !task list              # 프로젝트 + 명령어
!project list                # 명령어만
```

---

## Claude Code 실행 모델

### 2-Depth 제한

```
claribot
    └──TTY──▶ Claude Code [전역] (~/.claribot/)
                    └──Task──▶ Claude Code [프로젝트] (프로젝트 경로)
                                    └──▶ 작업 수행 (더 이상 네스트 안 함)
```

---

## 구현 현황

- [x] Makefile, systemd 서비스 설정
- [x] DB 스키마 (전역/로컬)
- [x] Telegram 패키지
- [x] claribot 메인 (Telegram Echo 연동)
- [ ] 메시지 → Task 생성 로직
- [ ] Claude Code TTY 연동
- [ ] Task 실행 및 결과 반환
- [ ] clari CLI HTTP 클라이언트

---

*Claribot v0.2.1*
