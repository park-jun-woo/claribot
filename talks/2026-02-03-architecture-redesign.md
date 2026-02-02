# Claritask 아키텍처 재설계

날짜: 2026-02-03

## 1. 제어 역전 (Control Inversion)

### 문제
- 기존: Claude Code가 Claritask를 메모장처럼 사용
- Claude Code 세션 내에서 Task 순회 → 컨텍스트 누적 → 폭발
- `/clear`는 사용자만 실행 가능 → 자동화 불가

### 해결
- **Claritask가 오케스트레이터**, Claude는 실행기
- `clari project start` 실행 시 Claritask가 `claude --print` 반복 호출
- 매 Task마다 독립 컨텍스트 → 누적 없음
- Task 100개든 10000개든 처리 가능

```
Claritask (Orchestrator)
    │
    ├─▶ Task 1 ─▶ claude --print (독립 컨텍스트)
    ├─▶ Task 2 ─▶ claude --print (독립 컨텍스트)
    └─▶ Task N ─▶ claude --print (독립 컨텍스트)
```

---

## 2. Phase → Feature 변경

### 이유
- Phase: 시간 축 (설계 → 개발 → 테스트) - 추상적
- Feature: 기능 축 (로그인, 결제, 블로그) - 구체적, 사용자 가치 단위

### Feature 정의
> "사용자가 인지하는 가치 단위"
> "릴리즈 노트에 쓸 만한 것"

### 규모
- MVP/스타트업: 5-10개
- 중형 서비스: 10-20개
- 대형 플랫폼: 20-40개

---

## 3. 트리 → 그래프 구조 (Edge)

### 문제
- 트리는 수직 관계만 표현
- 실제 코드 의존성은 수평 관계가 더 많음
- Model Task 실행 시 SQL Table Task 결과가 필요

### 해결
Task 간 Edge (의존성) 도입:

```
SQL Table Task ←───┬─── Model Task
                   │
Auth Config Task ←─┴─── API Task
```

### 장점
- 컨텍스트 정밀 주입: 해당 Task + 의존 Task result만 주입
- 실행 순서 자동 결정: Topological Sort
- 토큰 최소화: 전체 manifest 대신 필요한 것만

### 제한
- 각 Task의 Edge는 최대 4-7개
- 너무 많으면 Task 분할 필요 신호

---

## 4. Edge 추론 메커니즘

### 핵심 아이디어
Feature가 **Edge 추론의 경계**가 됨

### Flow
```
1. Feature 내 Task 목록을 LLM에 제시
2. LLM이 Task 간 의존성 분석
3. Edge 목록 반환 (JSON)
4. Claritask가 Edge 저장
```

### 왜 가능한가
- Feature 내 Task: 5-15개 → LLM 컨텍스트에 충분히 들어감
- Feature 목록: 10-30개 → 한 번에 분석 가능
- LLM이 코드 의존성 패턴 잘 이해 (SQL → Model → Service → API)

---

## 5. 구조화된 Planning 워크플로우

### Flow
```
Project Desc
    ▼ (LLM 1회)
Feature 목록 산출
    ▼ (LLM N회, 대화형)
Feature별 Spec 수립
    ▼ (LLM 1회)
Feature 간 Edge 추출
    ▼ (LLM N회)
Feature별 Task 생성
    ▼ (LLM N회)
Feature별 Task Edge 추출
    ▼
실행 준비 완료
```

### LLM 호출 횟수 (Feature 20개 기준)
| 단계 | 호출 수 |
|------|---------|
| Feature 목록 | 1회 |
| Feature Spec | 20회 |
| Feature Edge | 1회 |
| Task 생성 | 20회 |
| Task Edge | 20회 |
| **총 Planning** | **~60회** |

---

## 6. 데이터 구조 변경

### 신규 테이블
```sql
-- Feature 테이블 (phases 대체)
CREATE TABLE features (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id TEXT NOT NULL,
    name TEXT NOT NULL,
    spec TEXT DEFAULT ''
);

-- Feature 간 의존성
CREATE TABLE feature_edges (
    from_feature_id INTEGER,
    to_feature_id INTEGER,
    PRIMARY KEY (from_feature_id, to_feature_id)
);

-- Task 간 의존성
CREATE TABLE task_edges (
    from_task_id INTEGER,
    to_task_id INTEGER,
    PRIMARY KEY (from_task_id, to_task_id)
);
```

### Manifest 변경
```json
{
  "task": { ... },
  "dependencies": [
    {"id": 41, "title": "user_model", "result": "..."}
  ],
  "manifest": { ... }
}
```

`dependencies`: 의존 Task들의 `result`가 자동 포함됨

---

## 7. 핵심 Feature 정의

> **"LLM의 컨텍스트 윈도우를 무한대로 확장하는 Task Runner"**

기존 LLM 도구들은 컨텍스트가 차면 멈춘다. Claritask는 Task 단위로 쪼개서 LLM을 반복 호출하므로 Task가 몇 개든 상관없이 프로젝트를 완성한다.
