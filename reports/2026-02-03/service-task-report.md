# task_service.go Report

**파일 경로**: `internal/service/task_service.go`
**라인 수**: 634줄

## 요약
Task CRUD, 상태 관리, Pop 동작 구현. 의존성 기반 실행 순서 지원.

## 주요 타입
```go
type TaskCreateInput struct {
    PhaseID, FeatureID, SkeletonID, ParentID *int64
    Title, Content, Level, Skill string
    References []string
    TargetFile, TargetFunction string
    TargetLine *int
}
```

## 핵심 함수

### CRUD
| 함수 | 설명 |
|------|------|
| `CreateTask(db, input)` | Task 생성 |
| `GetTask(db, id)` | Task 조회 (string/int64 ID 지원) |
| `ListTasks(db, phaseID)` | Phase별 Task 목록 |
| `ListTasksByFeature(db, featureID)` | Feature별 Task 목록 |

### 상태 전이
| 함수 | 전이 | 조건 |
|------|------|------|
| `StartTask(db, id)` | pending → doing | status == pending |
| `CompleteTask(db, id, result)` | doing → done | status == doing |
| `FailTask(db, id, errMsg)` | doing → failed | status == doing |
| `ResetTaskToPending(db, id)` | * → pending | - |

### Pop 함수
```go
// 기본 Pop - 의존성 무시
func PopTask(db) (*TaskPopResult, error)

// 의존성 기반 Pop - 완료되지 않은 의존성이 있으면 스킵
func PopTaskFull(db) (*TaskPopResponse, error)

// 실행 가능한 다음 Task 조회
func GetNextExecutableTask(db) (*Task, error)
```

### Manifest 구성
PopTask는 다음을 포함하는 Manifest 반환:
- Context, Tech, Design
- State (key-value)
- High Priority Memos
- Feature 정보 (있으면)
- FDL 정보 (있으면)
- Skeleton 정보 (있으면)
- Dependencies 결과

### 상태 요약
```go
type TaskStatusResult struct {
    Total, Pending, Doing, Done, Failed int
    Progress float64 // (Done/Total)*100
}

func GetTaskStatus(db) (*TaskStatusResult, error)
```

## 의존성 기반 쿼리
```sql
SELECT * FROM tasks
WHERE status = 'pending'
AND NOT EXISTS (
    SELECT 1 FROM task_edges e
    JOIN tasks dep ON e.to_task_id = dep.id
    WHERE e.from_task_id = tasks.id
    AND dep.status != 'done'
)
```
