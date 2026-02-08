# Platform Development Node Graph Generation Algorithm

**작성일**: 2026-02-07
**배경**: gozip 프로젝트 382개 UI 기획서 생성 경험에서 도출

---

## 1. 문제 정의

### 1.1 배경

gozip 프로젝트에서 382개 UI 기획서를 Claude Code로 생성하는 과정에서 발견한 문제:

- 사람이 매번 "이 기획서 만들 때 이것도 참고해"라고 컨텍스트를 수동 큐레이션
- Claude Code가 중간에 Task를 누락하거나 일관성을 잃음
- 세션 끊기면 어디까지 했는지 추적 불가

### 1.2 핵심 질문

> 기획서(명세서)가 있다면, 어떻게 그것을 구체화해서 돌아가는 코드 결과물로 변환하느냐?

### 1.3 변환 파이프라인의 세 단계

| 단계 | 변환 | 성격 | 자동화 가능성 |
|------|------|------|-------------|
| 명세서 → 기획서 | "뭘 만들지" 구체화 | 창작 (도메인 전문성 필요) | 낮음 (사람이 해야 함) |
| 기획서 → 계획 | "어떤 순서로 만들지" | 분석 (의존성 파악) | 중간 → **이 알고리즘의 대상** |
| 계획 → 코드 | "명확한 지시를 코드로" | 번역 (기계적 변환) | 높음 (Claude Code) |

**명세서 → 기획서는 도메인 전문가의 영역.** gozip의 382개 기획서는 임대관리 10년 경험의 산물이며, 이 단계는 자동화 대상이 아님.

**이 문서는 기획서 → 계획 → 코드 구간의 자동화를 다룸.**

---

## 2. 기존 Claribot의 문제

### 2.1 원래 아이디어

> 연관 Edge를 기계적으로 엮어서, 컨텍스트를 10개 이내로 억제하면서도 필요한 맥락은 누락없이 주입해준다.

### 2.2 구현 과정에서의 변질

| 원래 | 실제 |
|------|------|
| Feature별 Edge 기반 컨텍스트 주입 | 범용 Task + Context Map("알아서 찾아봐") |
| 플랫폼 개발 특화 | 범용 (소설, 기획, 코딩 다 하려 함) |
| Node naming convention으로 예측 가능한 탐색 | Task title에 자연어로 기술 |

### 2.3 범용 Task의 한계

```sql
-- 현재: 평면적 Task
tasks (id, parent_id, title, spec, plan, report, status)
```

- Task 간 의존성(Edge) 정보 없음
- Layer 구분 없음 (DB 작업인지 UI 작업인지 구별 불가)
- 컨텍스트 자동 주입 불가 (Claude가 매번 탐색 반복)

---

## 3. Node Graph 알고리즘

### 3.1 핵심 개념

**Feature ID에서 Layer별 Node를 생성하고, naming convention으로 Edge를 자동 탐색하여, 2-pass로 그래프를 완성한 뒤 토폴로지 순서로 실행한다.**

### 3.2 구성 요소

#### Node

Feature와 Layer의 조합. naming convention으로 ID가 결정됨.

```
{layer}-{feature_id}

예시:
  table-OWN001     → buildings 테이블 DDL
  model-OWN001     → Building struct
  service-OWN001   → BuildingService.Create()
  api-OWN001       → POST /api/v1/owner/buildings handler
  page-OWN001      → BuildingCreatePage.tsx
```

#### Edge

Node 간 의존 관계. 두 가지 유형:

```
수직 Edge (같은 Feature, 다른 Layer):
  page-OWN001 → api-OWN001 → service-OWN001 → model-OWN001 → table-OWN001

수평 Edge (다른 Feature, cross-reference):
  page-OWN006 → api-OWN001  (건물 사진 등록이 건물 조회 API를 호출)

공통 컴포넌트 Edge:
  page-OWN001 → component-RadioGroup
  page-OWN010 → component-RadioGroup
```

#### Layer 정의

플랫폼 개발의 수직 슬라이스:

| Layer | 설명 | 출력 |
|-------|------|------|
| table | DB 스키마 | DDL, migration |
| model | 데이터 모델 | Go struct, TypeScript type |
| service | 비즈니스 로직 | Service 함수 |
| api | HTTP 핸들러 | Handler, Router |
| page | 화면 | React 컴포넌트 |
| component | 공통 UI 컴포넌트 | 재사용 컴포넌트 |

### 3.3 알고리즘

```
입력:
  - Feature 목록 (OWN-001 ~ OWN-207, TEN-001 ~ TEN-027, ...)
  - Layer 정의 (table, model, service, api, page)
  - 기획서 (docs/ui/{role}/{id}-spec.md)

Pass 1 — 그래프 생성:
  for each feature in features:
    for each layer in layers:
      create node "{layer}-{feature_id}"

    parse spec(feature):
      extract "관련 기능" → cross-feature edges
      extract "API Endpoint" → api layer node edges
      extract "컴포넌트 구조" → component node edges
      extract "화면 진입 경로" → page-to-page edges
      if target node not exists → create empty node (forward reference)

Pass 2 — 빈 노드 해결:
  for each empty node (topological order):
    resolve spec from:
      - 원본 기획서의 해당 layer 정보
      - 연결된 edge의 기존 node 참조
    fill node with plan

실행 — 토폴로지 순서:
  component-* → table-* → model-* → service-* → api-* → page-*
  (같은 layer 내에서는 병렬 실행 가능)
```

### 3.4 컨텍스트 주입

각 Node 실행 시, Claude Code에 주입하는 컨텍스트:

```
Node: page-OWN006 실행 시

자동 수집되는 Edge 기반 컨텍스트:
  1. docs/ui/owner/OWN-006-spec.md         ← 기획서 (spec → node)
  2. api-OWN006의 output                   ← 수직 edge (자기 API)
  3. api-OWN001의 output                   ← 수평 edge (건물 조회)
  4. page-OWN003의 output                  ← 진입점 edge
  5. component-FileUpload의 output         ← 컴포넌트 edge
  6. frontend/apps/owner/src/App.tsx        ← 라우팅 (고정 참조)

→ 6개. 10개 이내. 누락 없음.
```

### 3.5 컴파일러 비유

이 알고리즘은 컴파일러의 forward reference 해결과 동일한 구조:

```
컴파일러:
  Pass 1 (파싱): 심볼 수집, 미해결 참조는 심볼 테이블에 빈 항목 등록
  Pass 2 (링킹): 빈 항목 해결, 실제 주소 연결, 코드 생성

Node Graph:
  Pass 1 (계획): Node/Edge 수집, 미해결 Node는 빈 노드로 등록
  Pass 2 (실행): 빈 노드 채우기, 토폴로지 순서로 코드 생성
```

---

## 4. 데이터 모델

### 4.1 스키마

```sql
nodes (
  id TEXT PRIMARY KEY,       -- "page-OWN001"
  feature_id TEXT NOT NULL,  -- "OWN-001"
  layer TEXT NOT NULL,        -- "page"
  spec TEXT,                  -- 기획서에서 이 layer에 해당하는 부분
  plan TEXT,                  -- 실행 계획
  output TEXT,                -- 생성된 파일 경로
  status TEXT DEFAULT 'empty' -- empty | planned | done | failed
)

edges (
  from_node TEXT NOT NULL,    -- "page-OWN006"
  to_node TEXT NOT NULL,      -- "api-OWN001"
  type TEXT NOT NULL,         -- "calls" | "uses" | "extends" | "enters_from"
  PRIMARY KEY (from_node, to_node)
)
```

### 4.2 기존 Task와의 관계

| 기존 Task | Node Graph |
|-----------|-----------|
| `tasks.title` | `nodes.id` (naming convention으로 자동 생성) |
| `tasks.spec` | `nodes.spec` (기획서에서 layer별로 추출) |
| `tasks.plan` | `nodes.plan` |
| `tasks.parent_id` (트리) | `edges` (그래프) |
| 의존성 없음 | `edges`로 명시적 의존성 |

---

## 5. gozip 적용 시 규모 추정

```
Features: 382개
Layers: 6개 (table, model, service, api, page, component)
예상 Nodes: ~2,000개 (일부 feature는 모든 layer가 필요하지 않음)
예상 Edges: ~4,000개 (수직 + 수평 + 컴포넌트)

실행 순서:
  1. component-*  : ~20개 (병렬)
  2. table-*      : ~100개 (병렬)
  3. model-*      : ~100개 (병렬)
  4. service-*    : ~200개 (병렬)
  5. api-*        : ~300개 (병렬)
  6. page-*       : ~386개 (병렬)

같은 layer 내 병렬 3개 기준: 약 370 세션
```

---

## 6. Claribot 변경 범위

### 6.1 유지

- Claude Code PTY 실행 엔진 (`pkg/claude/`)
- 병렬 실행 worker pool
- Telegram/CLI 인터페이스
- Report 파일 감지 패턴

### 6.2 변경

| 항목 | 현재 | 변경 |
|------|------|------|
| 데이터 모델 | `tasks` 단일 테이블 | `nodes` + `edges` 테이블 |
| 작업 식별 | 자연어 title | naming convention (`{layer}-{feature_id}`) |
| 의존성 | 없음 (parent_id만) | Edge 기반 그래프 |
| 컨텍스트 | Context Map (전체 요약) | Edge 기반 선택적 주입 (10개 이내) |
| 실행 순서 | ID 순 또는 depth 순 | 토폴로지 정렬 |
| 1회차 | 범용 plan/split | 기획서 파싱 → 그래프 생성 |
| 2회차 | 범용 execute | 빈 노드 해결 → 토폴로지 실행 |

### 6.3 추가

- 기획서 파서 (관련 기능, API, 컴포넌트, 진입경로 추출)
- Node naming convention 엔진
- Edge 자동 탐색
- 토폴로지 정렬 실행기

---

## 7. 설계 원칙

1. **플랫폼 개발에 특화한다.** 범용 Task 시스템이 아닌, DB→Model→Service→API→Page 수직 슬라이스에 최적화.
2. **naming convention이 핵심이다.** 예측 가능한 Node ID로 Edge를 기계적으로 탐색.
3. **컨텍스트는 10개 이내로 억제한다.** Edge로 연결된 Node만 주입, 전체 Context Map 불필요.
4. **같은 layer는 병렬 실행한다.** 수직 의존성만 순서 강제, 수평은 독립.
5. **기획서는 사람이 만든다.** "무엇을 만드는가"는 도메인 전문가의 영역. 시스템은 "어떻게 만드는가"만 담당.
