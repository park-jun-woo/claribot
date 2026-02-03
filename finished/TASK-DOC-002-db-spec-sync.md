# TASK-DOC-002: DB 스펙과 구현 동기화

## 개요
DB 스펙 문서를 현재 구현 상태에 맞춰 업데이트

## 배경
- **보고서**: reports/2026-02-03-total-report.md 섹션 3
- **대상 파일**: specs/DB/02-A-Core.md, specs/DB/02-C-Content.md

## 불일치 항목

### 1. experts 테이블 구조

**현재 스펙 (예상)**:
```sql
CREATE TABLE experts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    file_path TEXT NOT NULL,
    ...
);
```

**실제 구현**:
```sql
CREATE TABLE experts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    version TEXT,
    domain TEXT,
    language TEXT,
    framework TEXT,
    description TEXT,
    created_at TEXT NOT NULL
);
```

### 2. tasks 테이블 추가 컬럼

스펙에 없는 구현된 컬럼:
- `parent_id INTEGER` - 태스크 계층 구조
- `level TEXT` - leaf/parent
- `skill TEXT` - 필요 기술
- `refs TEXT` - 참조 목록

### 3. project_experts 테이블

스펙에 없는 새 테이블:
```sql
CREATE TABLE project_experts (
    project_id TEXT NOT NULL,
    expert_id TEXT NOT NULL,
    assigned_at TEXT NOT NULL,
    PRIMARY KEY (project_id, expert_id)
);
```

### 4. expert_assignments 테이블

`expert_id` 타입 불일치:
- 스펙: INTEGER (experts.id 참조)
- 구현: TEXT (experts.id가 TEXT이므로)

## 작업 내용

### specs/DB/02-A-Core.md
1. experts 테이블 스키마 수정
2. expert_assignments 테이블 expert_id 타입 수정

### specs/DB/02-C-Content.md
1. tasks 테이블에 추가 컬럼 문서화
2. project_experts 테이블 추가

## 완료 기준
- [ ] experts 테이블 스키마가 구현과 일치
- [ ] tasks 테이블 추가 컬럼 문서화
- [ ] project_experts 테이블 문서화
- [ ] expert_assignments 타입 수정
