# memo_service.go Report

**파일 경로**: `internal/service/memo_service.go`
**라인 수**: 232줄

## 요약
계층적 메모 관리 서비스. Project/Phase/Task 스코프 지원.

## 주요 타입
```go
type MemoSetInput struct {
    Scope    string   // "project", "phase", "task"
    ScopeID  string
    Key      string
    Value    string
    Priority int      // 1(high), 2(default), 3(low)
    Summary  string
    Tags     []string
}

type MemoListResult struct {
    Project map[string][]MemoSummary
    Phase   map[string][]MemoSummary
    Task    map[string][]MemoSummary
    Total   int
}
```

## CRUD 함수
| 함수 | 설명 |
|------|------|
| `SetMemo(db, input)` | 메모 설정 (UPSERT) |
| `GetMemo(db, scope, scopeID, key)` | 메모 조회 |
| `DeleteMemo(db, scope, scopeID, key)` | 메모 삭제 |
| `ListMemos(db)` | 전체 메모 목록 |
| `ListMemosByScope(db, scope, scopeID)` | 스코프별 메모 |

## 특수 함수
```go
// 우선순위 1인 메모만 조회 (task pop 시 사용)
func GetHighPriorityMemos(db) ([]Memo, error)

// 키 파싱: "PH001:T042:key" 형식 지원
func ParseMemoKey(input) (scope, scopeID, key, error)
```

## 키 포맷
| 레벨 | 포맷 | 예시 |
|------|------|------|
| Project | `key` | `api_convention` |
| Phase | `phase:key` | `1:progress_notes` |
| Task | `phase:task:key` | `1:5:implementation_detail` |

## 데이터 저장 구조
```json
{
  "value": "actual content",
  "summary": "optional summary",
  "tags": ["tag1", "tag2"]
}
```
