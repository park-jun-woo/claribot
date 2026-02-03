# db.go Report

**파일 경로**: `internal/db/db.go`
**라인 수**: 193줄

## 요약
SQLite 데이터베이스 래퍼와 스키마 마이그레이션 제공.

## 주요 타입
```go
type DB struct {
    *sql.DB
    path string
}
```

## 주요 함수
| 함수 | 설명 |
|------|------|
| `Open(path)` | DB 연결 열기, 디렉토리 자동 생성 |
| `Close()` | DB 연결 닫기 |
| `Migrate()` | 테이블 스키마 생성 |
| `TimeNow()` | 현재 시간 RFC3339 포맷 반환 |
| `ParseTime(s)` | RFC3339 문자열 파싱 |

## 테이블 스키마 (12개)
1. **projects** - 프로젝트 (TEXT PK)
2. **phases** - 작업 단계 (AUTOINCREMENT)
3. **features** - Feature (AUTOINCREMENT, UNIQUE name per project)
4. **feature_edges** - Feature 간 의존성
5. **skeletons** - 스켈레톤 파일 추적
6. **tasks** - 작업 (AUTOINCREMENT)
7. **task_edges** - Task 간 의존성
8. **context** - 프로젝트 컨텍스트 (싱글톤)
9. **tech** - 기술 스택 (싱글톤)
10. **design** - 설계 결정 (싱글톤)
11. **state** - 상태 key-value
12. **memos** - 메모 (scope 기반 복합키)

## 특징
- Foreign key 활성화 (`PRAGMA foreign_keys = ON`)
- Status CHECK 제약조건 적용
- SQLite 파일 기반 DB 사용
