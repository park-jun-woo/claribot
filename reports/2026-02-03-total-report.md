# Claritask Spec vs Implementation 전체 비교 보고서

**작성일**: 2026-02-03
**버전**: v0.0.4

---

## 1. 요약

| 영역 | 구현율 | 주요 이슈 |
|------|--------|-----------|
| CLI 명령어 | 95% | feature get 응답 필드 누락, 문서 오류 |
| DB 스키마 | 90% | experts 테이블 구조 불일치 |
| FDL 파싱 | 62% | 고급 검증 기능 미구현 |

---

## 2. CLI 명령어 비교

### 2.1 심각도 높은 이슈

#### Issue #1: `feature get` 응답 필드 누락
**Spec (CLI/07-Feature.md)**:
```json
{
  "feature": {
    "fdl": "...",
    "fdl_hash": "abc123...",
    "skeleton_generated": true
  }
}
```

**Implementation (feature.go:169-179)**: `fdl`, `fdl_hash`, `skeleton_generated` 필드 없음

**상태**: MISSING - 수정 필요

---

#### Issue #2: Expert 명령어 문서 오류
**Spec (CLI/01-Overview.md:62)**: `Expert | 7 | 미구현`

**실제 구현**: 7개 명령어 전부 구현됨 (expert.go)
- `expert add`, `list`, `get`, `edit`, `remove`, `assign`, `unassign`

**상태**: 문서 오류 - 스펙 수정 필요

---

### 2.2 중간 심각도 이슈

| # | 명령어 | 이슈 | 상태 |
|---|--------|------|------|
| 3 | expert assign/unassign | `--feature` 플래그 미문서화 | 문서 추가 필요 |
| 4 | feature list | 응답 구조 검증 필요 | 확인 필요 |
| 5 | project status | 응답 래퍼 구조 차이 가능 | 확인 필요 |

### 2.3 낮은 심각도 (기능 확장)

| # | 명령어 | 설명 |
|---|--------|------|
| 6 | task push | 추가 필드: parent_id, target_file, target_line, target_function |
| 7 | plan features | JSON 입력 지원 (미문서화) |
| 8 | expert edit | 메시지 텍스트 차이 ("Editor closed" vs "Opening editor...") |

---

## 3. DB 스키마 비교

### 3.1 일치하는 테이블 (10개)

| 테이블 | 상태 | 비고 |
|--------|------|------|
| projects | ✅ 일치 | |
| features | ✅ 일치 | 추가: UNIQUE(project_id, name) |
| feature_edges | ✅ 일치 | |
| task_edges | ✅ 일치 | |
| context | ✅ 일치 | |
| tech | ✅ 일치 | |
| design | ✅ 일치 | |
| state | ✅ 일치 | |
| memos | ✅ 일치 | |
| skeletons | ✅ 일치 | |

### 3.2 확장된 테이블

#### tasks 테이블
**스펙에 없는 추가 컬럼**:
- `parent_id INTEGER` - 태스크 계층 구조
- `level TEXT` - leaf/parent
- `skill TEXT` - 필요 기술
- `refs TEXT` - 참조 목록

### 3.3 불일치 테이블

#### experts 테이블 - 심각

| 항목 | Spec | Implementation |
|------|------|----------------|
| Primary Key | `INTEGER AUTOINCREMENT` | `TEXT` |
| name 제약 | `UNIQUE` | 없음 |
| 파일경로 | `file_path` | `path` |
| 추가 컬럼 | - | version, domain, language, framework |

#### expert_assignments 테이블
- `expert_id` 타입: Spec은 INTEGER, 구현은 TEXT (experts.id와 일관성)

### 3.4 스펙에 없는 테이블

**project_experts** - 프로젝트 레벨 전문가 할당
```sql
CREATE TABLE project_experts (
    project_id TEXT NOT NULL,
    expert_id TEXT NOT NULL,
    assigned_at TEXT NOT NULL,
    PRIMARY KEY (project_id, expert_id)
);
```

---

## 4. FDL 파싱 비교

### 4.1 레이어별 구현율

| 레이어 | 구현율 | 설명 |
|--------|--------|------|
| Data (Model) | 70% | 기본 타입/제약 지원, 고급 검증 미흡 |
| Logic (Service) | 55% | Step 타입 감지만, 구조 검증 미흡 |
| Interface (API) | 65% | 기본 구조 지원, 경로 파라미터 미추출 |
| Presentation (UI) | 50% | 기본 파싱 지원, 바인딩/루프 미지원 |

### 4.2 Data Layer (70%)

**지원됨**:
- 14가지 필드 타입 (uuid, string, text, int, bigint, float, decimal, boolean, datetime, date, time, json, blob, enum)
- 타입 옵션 파싱: string(50), decimal(10,2), enum(...)
- 제약조건: pk, fk, required, unique, auto, index, nullable, default, check, onDelete, onUpdate
- 인덱스 정의, 관계 정의 (hasOne, hasMany, belongsTo)

**미지원**:
- belongsToMany 관계
- through 테이블 검증
- 부분 인덱스 where 절 검증
- 복합 제약조건 호환성 검증

### 4.3 Logic Layer (55%)

**지원됨**:
- 서비스 기본 구조 (name, desc, input, output, throws, transaction)
- Auth 레벨 검증 (required, optional, none)
- 입력 파라미터 검증 규칙
- Step 타입 감지: validate, db, event, call, cache, log, transform, condition, loop, return

**미지원**:
- db step: 작업 타입 검증 (insert, select, update, delete)
- event step: 이벤트 이름/페이로드 검증
- call step: URL, 헤더, 타임아웃 형식 검증
- cache step: 키 템플릿, TTL 형식 검증
- condition step: 조건식 파싱
- loop step: 반복 구조 파싱

### 4.4 Interface Layer (65%)

**지원됨**:
- HTTP 메서드: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- Auth 타입: required, optional, none, apikey, bearer
- Request/Response 파싱
- Rate limiting, Transform 설정

**미지원**:
- 경로 파라미터 추출 ({userId} 등)
- 파일 업로드 제약 검증
- apiGroups 구조
- 상태 코드 검증

### 4.5 Presentation Layer (50%)

**지원됨**:
- 컴포넌트 타입: Page, Template, Organism, Molecule, Atom
- Props, State, Computed 파싱
- Init, Methods, View 파싱
- 조건부 렌더링 (if/then/else)
- 액션 타입: call, set, navigate, show, validate, confirm, emit, parallel, redirect

**미지원**:
- 컴포넌트 계층 순환 참조 검증
- Props 타입 검증 (ReactNode, 제네릭)
- 바인딩 구문 파싱 ({variable})
- for 루프 파싱
- 이벤트 핸들러 검증 (@click, @change)

---

## 5. 권장 수정 사항

### 5.1 높은 우선순위

1. **feature.go**: `feature get` 응답에 fdl, fdl_hash, skeleton_generated 필드 추가
2. **CLI/01-Overview.md**: Expert 명령어 상태를 "구현됨"으로 수정
3. **db.go**: experts 테이블 구조 스펙과 일치시키기 (또는 스펙 업데이트)

### 5.2 중간 우선순위

4. **CLI/11-Expert.md**: `--feature` 플래그 문서화
5. **fdl_service.go**: 서비스 step 구조 검증 강화
6. **fdl_service.go**: API 경로 파라미터 추출 구현

### 5.3 낮은 우선순위

7. task 테이블 추가 컬럼들 스펙에 문서화
8. project_experts 테이블 스펙에 추가
9. UI 바인딩 구문 파싱 구현

---

## 6. 파일 참조

### 스펙 파일
- `/specs/CLI/*.md` - CLI 명령어 스펙
- `/specs/DB/*.md` - DB 스키마 스펙
- `/specs/FDL/*.md` - FDL 문법 스펙

### 구현 파일
- `/cli/internal/cmd/*.go` - CLI 명령어 구현
- `/cli/internal/db/db.go` - DB 스키마
- `/cli/internal/service/fdl_service.go` - FDL 파싱
- `/cli/internal/model/models.go` - 데이터 모델
