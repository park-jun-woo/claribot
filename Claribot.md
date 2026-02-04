# Claritask v0.2 - 전역 서비스 아키텍처

> **버전**: v0.2.0 (Draft)

---

## 개요

LLM 기반 프로젝트 자동 실행 시스템

**v0.1 대비 변경점**:
- 프로젝트 로컬 → **전역 서비스** 아키텍처
- **텔레그램 봇 통합** (원격 제어)
- clari CLI → **서비스 클라이언트** 방식
- localhost HTTP 통신

---

## 아키텍처 개요

### 전체 구조

```
┌─────────────────────────────────────────────────────────┐
│              claribot (daemon)                 │
│                                                         │
│  ┌───────────┐  ┌───────────┐  ┌───────────────────┐    │
│  │  텔레그램  │  │   CLI     │  │    TTY 매니저      │    │
│  │  핸들러    │  │  핸들러    │  │   (클로드 코드)    │    │
│  │(goroutine)│  │(goroutine)│  │   (goroutine)     │    │
│  └─────┬─────┘  └─────┬─────┘  └─────────┬─────────┘    │
│        │              │                   │             │
│        └──────────────┼───────────────────┘             │
│                       ▼                                 │
│               ┌──────────────┐                          │
│               │      DB      │                          │
│               │   (mutex)    │                          │
│               └──────────────┘                          │
└─────────────────────────────────────────────────────────┘
         ▲                              ▲
         │ 메시지                        │ HTTP 요청
         │                              │
    [텔레그램]                     [clari CLI]
```

### 컴포넌트 역할

| 컴포넌트 | 역할 | 실행 |
|----------|------|------|
| claribot | 텔레그램 봇 + CLI 핸들러 + TTY 관리 | 항상 실행 (systemctl) |
| clari CLI | 서비스에 HTTP 요청 보내는 클라이언트 | 일시 실행 |
| 클로드 코드 (전역) | claritask 루트에서 실행, 라우터 역할 | TTY 호출 → 완료 후 kill |
| 클로드 코드 (프로젝트) | 프로젝트 폴더에서 실행, 실제 작업 | Task tool로 호출 → 완료 후 종료 |

---

## 폴더 구조

### Claritask 홈 (~/.claritask/)

```
~/.claritask/                     ← claritask 홈 (자동 생성)
├── db                           ← 전역 DB (프로젝트 라우팅 테이블)
├── config.yaml                  ← 서비스 설정
└── logs/                        ← 로그 (옵션)
```

### 프로젝트 폴더 (어디든 가능)

```
~/projects/blog/                 ← 프로젝트 (위치 자유)
├── .git/                        ← 독립적 git 관리
├── .claritask/
│   └── db                      ← 로컬 DB (tasks, memos, features...)
└── src/

~/work/api-server/               ← 다른 위치도 가능
├── .git/
├── .claritask/
│   └── db
└── src/
```

### DB 분리

| DB | 위치 | 역할 | git 관리 |
|----|------|------|----------|
| 전역 DB | `~/.claritask/db` | 프로젝트 목록, 경로 매핑, 텔레그램 설정 | X |
| 로컬 DB | `프로젝트/.claritask/db` | tasks, memos, features, experts 등 | O |

> **Note**: `.clt` 확장자 폐기. 전역/로컬 모두 `db` 파일명 사용.

### 전역 DB projects 테이블

```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY,           -- 프로젝트 ID (영문, 숫자, -, _)
    name TEXT NOT NULL,            -- 표시명
    path TEXT NOT NULL,            -- 상대경로 (루트 내) 또는 절대경로 (외부)
    telegram_chat_id TEXT,         -- 텔레그램 채팅방 매핑 (옵션)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**예시:**

| id | name | path | telegram_chat_id |
|----|------|------|------------------|
| blog | 블로그 | ~/projects/blog | -123456789 |
| api | API 서버 | ~/work/api-server | NULL |
| legacy | 레거시 | ~/company/old-project | NULL |

---

## 통신 구조

### CLI ↔ Service

```
clari CLI  ──HTTP──▶  claribot (127.0.0.1:9847)
```

- **프로토콜**: localhost HTTP
- **바인딩**: `127.0.0.1` (외부 접근 불가)
- **포트**: 설정 가능 (기본 9847)

**보안:**
```go
// 반드시 127.0.0.1로 바인딩
http.ListenAndServe("127.0.0.1:9847", handler)
```

### 텔레그램 ↔ Service

```
텔레그램 API  ◀──Long Polling──▶  claribot
```

---

## 텔레그램 문법

### 특수 기호

| 기호 | 용도 | 예시 |
|------|------|------|
| `$` | 프로젝트 지정 | `$blog`, `$api` |
| `!` | clari 명령어 실행 | `!task list`, `!memo add` |

### 메시지 패턴

```
# 프로젝트 + 자연어
$blog 태스크 뭐있어?
$api 버그 수정해줘

# 프로젝트 + 명령어
$blog !task list
$api !memo add "결제 버그 발견"

# 명령어만 (프로젝트 컨텍스트 불필요)
!project list

# 자연어만 (기본 프로젝트 또는 클로드가 판단)
진행상황 알려줘
```

### 기본 프로젝트 설정

```
/default blog       ← 기본 프로젝트를 blog로 설정
태스크 뭐있어?      ← $없으면 기본 프로젝트 (blog)
```

---

## 클로드 코드 실행 모델

### 2-Depth 제한

```
claribot
    │
    └──TTY 호출──▶ 클로드 코드 [전역]  (~/.claritask/)
                        │
                        └──Task tool──▶ 클로드 코드 [프로젝트]  (~/projects/blog/)
                                              │
                                              └──▶ 작업 수행 (더 이상 네스트 안 함)
```

- **전역 클로드**: 메시지 분석, 프로젝트 라우팅 (WorkingDirectory: ~/.claritask/)
- **프로젝트 클로드**: 실제 작업 수행 (WorkingDirectory: 프로젝트 경로)
- **최대 2depth**: 전역 → 프로젝트, 그 이상 네스트 금지

### 생명주기

```
요청 수신
    │
    ▼
클로드 코드 TTY 호출
    │
    ▼
작업 수행
    │
    ▼
결과 반환
    │
    ▼
클로드 코드 kill  ← 반드시 종료
```

**항시 실행 아님.** 요청마다 호출하고, 완료되면 반드시 종료.

---

## 메시지 처리 흐름

### 텔레그램 메시지 예시

```
사용자: $blog 태스크 정리해줘
```

### 처리 흐름

```
1. [텔레그램] 메시지 수신
       │
       ▼
2. [claribot] 파싱
   - 프로젝트: blog
   - 메시지: "태스크 정리해줘"
       │
       ▼
3. [전역 DB] projects 조회
   - blog → ./blog
       │
       ▼
4. [TTY 매니저] 클로드 코드 호출
   - WorkingDirectory: ~/.claritask/
   - 프롬프트: "$blog 태스크 정리해줘"
       │
       ▼
5. [클로드 코드 - 전역]
   - blog 프로젝트 확인
   - Task tool로 프로젝트 클로드 호출
       │
       ▼
6. [클로드 코드 - blog]
   - clari task list 실행
   - 정리 작업 수행
   - 결과 반환
       │
       ▼
7. [claribot]
   - 결과를 DB에 저장
   - 클로드 코드 kill
       │
       ▼
8. [텔레그램] 결과 전송
```

---

## clari CLI 변경점

### v0.1 (직접 실행)

```go
// DB 직접 접근
db, _ := db.Open(".claritask/db")
tasks, _ := db.ListTasks()
```

### v0.2 (서비스 클라이언트)

```go
// 서비스에 HTTP 요청
resp, _ := http.Get("http://127.0.0.1:9847/task/list?project=blog")
```

### 주요 명령어 변경

| 명령어 | v0.1 | v0.2 |
|--------|------|------|
| `clari init` | 로컬 DB 생성 | 서비스에 프로젝트 등록 요청 |
| `clari task list` | 로컬 DB 조회 | 서비스에 조회 요청 |
| `clari project start` | 직접 TTY 호출 | 서비스에 실행 요청 |
| `clari serve` | (없음) | **신규**: 서비스 데몬 실행 |

---

## VSCode Extension

### 동작 방식 변경

`.clt` 확장자 폐기로 진입점 변경 필요:

```
기존: db.clt 더블클릭 → Custom Editor 열기
변경: Activity Bar 아이콘 또는 Command Palette
```

**활성화 조건:**
```
VSCode에서 ~/projects/blog/ 열기
    │
    ▼
blog/.claritask/db 파일 감지 (workspaceContains:.claritask/db)
    │
    ▼
Claritask Extension 활성화
    │
    ▼
Activity Bar 아이콘 클릭 또는 Ctrl+Shift+P → "Claritask: Open"
```

- git으로 프로젝트 상태 공유 가능 (로컬 DB)

---

## 서비스 설정

### 설정 파일

```yaml
# ~/.claritask/config.yaml
service:
  port: 9847
  host: 127.0.0.1

telegram:
  token: "BOT_TOKEN"
  allowed_users:
    - 123456789

claude:
  timeout: 1200
  max: 3
```

### 설정 항목 설명

| 섹션 | 항목 | 기본값 | 설명 |
|------|------|--------|------|
| service | port | 9847 | CLI 통신용 HTTP 포트 |
| service | host | 127.0.0.1 | 바인딩 주소 (외부 접근 차단) |
| telegram | token | - | 텔레그램 봇 토큰 |
| telegram | allowed_users | - | 허용된 텔레그램 사용자 ID 목록 |
| claude | timeout | 1200 | idle timeout (초). 출력 없이 이 시간 지나면 kill |
| claude | max | 3 | 동시 실행 가능한 클로드 코드 최대 개수 |

**claude.timeout (idle timeout):**
- 클로드 코드가 출력을 생성하는 동안은 무제한 실행
- 출력 없이 20분(1200초) 지나면 문제 발생으로 간주하고 kill
- Planning이 1시간 걸려도 출력만 계속 나오면 유지

**claude.max:**
- 동시에 실행 가능한 클로드 코드 인스턴스 최대 개수
- 전역 클로드 + 프로젝트 클로드 합산
- 초과 요청은 큐에서 대기

> **Note**: 프로젝트 경로는 전역 DB의 projects 테이블에서 관리.

### systemctl 서비스

```ini
# /etc/systemd/system/claritask.service
[Unit]
Description=claribot
After=network.target

[Service]
Type=simple
User=username
WorkingDirectory=/home/username/.claritask
ExecStart=/usr/local/bin/clari serve
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

---

## 마이그레이션 (v0.1 → v0.2)

### 단계

1. **claritask 홈 초기화**
   ```bash
   clari init --global   # ~/.claritask/ 자동 생성
   ```

2. **기존 프로젝트 등록**
   ```bash
   # 프로젝트 경로 등록
   clari project add blog --path ~/projects/blog
   clari project add api --path ~/work/api-server
   ```

3. **기존 로컬 DB 파일명 변경**
   ```bash
   # 각 프로젝트에서
   mv .claritask/db.clt .claritask/db
   ```

4. **서비스 시작**
   ```bash
   clari serve  # 또는 systemctl start claritask
   ```

---

## 관련 문서

| 문서 | 내용 |
|------|------|
| Service/ | 서비스 아키텍처 상세 |
| CLI/ | 명령어 레퍼런스 (v0.2) |
| DB/ | 데이터베이스 스키마 (전역 + 로컬) |
| Telegram/ | 텔레그램 봇 명세 |
| TTY/ | TTY Handover 명세 |
| VSCode/ | VSCode Extension (변경 없음) |

---

*Claritask Specification v0.2.0 Draft*
