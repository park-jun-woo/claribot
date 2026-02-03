# phase_service.go Report

**파일 경로**: `internal/service/phase_service.go`
**라인 수**: 129줄

## 요약
Phase 생명주기 관리 서비스. CRUD 및 상태 전이 기능.

## 주요 타입
```go
type PhaseCreateInput struct {
    ProjectID   string
    Name        string
    Description string
    OrderNum    int
}

type PhaseListItem struct {
    ID, Name, Description string
    OrderNum, TasksTotal, TasksDone int
    Status string
}
```

## CRUD 함수
| 함수 | 설명 |
|------|------|
| `CreatePhase(db, input)` | Phase 생성 (status=pending) |
| `GetPhase(db, id)` | Phase 조회 |
| `ListPhases(db, projectID)` | 프로젝트 Phase 목록 (Task 통계 포함) |
| `UpdatePhaseStatus(db, id, status)` | 상태 업데이트 |

## 상태 전이 함수
| 함수 | 전이 | 조건 |
|------|------|------|
| `StartPhase(db, id)` | pending → active | status == pending |
| `CompletePhase(db, id)` | active → done | status == active |

## Phase 상태 흐름
```
pending → active → done
```

## 목록 쿼리 (Task 통계 포함)
```sql
SELECT p.*,
    (SELECT COUNT(*) FROM tasks WHERE phase_id = p.id) as tasks_total,
    (SELECT COUNT(*) FROM tasks WHERE phase_id = p.id AND status = 'done') as tasks_done
FROM phases p
WHERE project_id = ?
ORDER BY order_num
```
