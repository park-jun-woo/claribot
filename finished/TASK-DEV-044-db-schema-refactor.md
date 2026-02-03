# TASK-DEV-044: DB 스키마 리팩토링

## 목표
specs/Claritask.md에 맞게 DB 스키마 수정. Phase 제거, Task-Feature 직접 연결.

## 파일
`internal/db/db.go`

## 변경 내용

### 1. phases 테이블 제거
```sql
-- 삭제
CREATE TABLE IF NOT EXISTS phases (...)
```

### 2. tasks 테이블 수정
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_id INTEGER NOT NULL,  -- phase_id 대신 feature_id 필수
    -- phase_id 제거
    parent_id INTEGER DEFAULT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK(status IN ('pending', 'doing', 'done', 'failed')),
    title TEXT NOT NULL,
    content TEXT DEFAULT '',
    target_file TEXT DEFAULT '',
    target_line INTEGER,
    target_function TEXT DEFAULT '',
    result TEXT DEFAULT '',
    error TEXT DEFAULT '',
    skeleton_id INTEGER DEFAULT NULL,
    created_at TEXT NOT NULL,
    started_at TEXT,
    completed_at TEXT,
    failed_at TEXT,
    FOREIGN KEY (feature_id) REFERENCES features(id),
    FOREIGN KEY (parent_id) REFERENCES tasks(id),
    FOREIGN KEY (skeleton_id) REFERENCES skeletons(id)
);
```

### 3. memos 테이블 scope 주석 업데이트
- scope: 'project', 'feature', 'task' (phase 제거)

## 주의사항
- 기존 DB 마이그레이션 필요 (별도 마이그레이션 스크립트 고려)
- tasks 테이블의 level, skill, references 필드는 specs에 없으므로 제거 검토

## 완료 조건
- [ ] phases 테이블 정의 제거
- [ ] tasks 테이블에서 phase_id 제거, feature_id NOT NULL
- [ ] 불필요한 필드(level, skill, references) 제거
