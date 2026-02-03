# task.go Report

**파일 경로**: `internal/cmd/task.go`
**라인 수**: 388줄

## 요약
Task 관리 명령어 그룹 구현. Push/Pop 기반 작업 큐.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `task push '<json>'` | 새 Task 추가 |
| `task pop` | 다음 pending Task 가져오기 + 시작 |
| `task start <id>` | Task 시작 (pending → doing) |
| `task complete <id> '<json>'` | Task 완료 (doing → done) |
| `task fail <id> '<json>'` | Task 실패 (doing → failed) |
| `task status` | Task 진행 현황 |
| `task get <id>` | 특정 Task 조회 |
| `task list <phase_id>` | Phase별 Task 목록 |

## task push 입력
```go
type taskPushInput struct {
    PhaseID    int64    `json:"phase_id"`
    ParentID   *int64   `json:"parent_id"`
    Title      string   `json:"title"`
    Content    string   `json:"content"`
    Level      string   `json:"level"`
    Skill      string   `json:"skill"`
    References []string `json:"references"`
}
```

## task pop 응답
```json
{
  "success": true,
  "task": { "id": "1", "title": "...", "status": "doing" },
  "manifest": {
    "context": {...},
    "tech": {...},
    "design": {...},
    "state": {...},
    "memos": [...]
  }
}
```

## task complete 입력
```go
type taskCompleteInput struct {
    Result string `json:"result"`
    Notes  string `json:"notes"`
}
```

## task fail 입력
```go
type taskFailInput struct {
    Error   string `json:"error"`
    Details string `json:"details"`
}
```

## 상태 전이
- `pending` → `doing` (start/pop)
- `doing` → `done` (complete)
- `doing` → `failed` (fail)
