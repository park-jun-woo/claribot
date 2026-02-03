# state_service.go Report

**파일 경로**: `internal/service/state_service.go`
**라인 수**: 122줄

## 요약
Key-Value 상태 관리 서비스. 현재 실행 위치 추적.

## 상태 키 상수
```go
const (
    StateCurrentProject = "current_project"
    StateCurrentPhase   = "current_phase"
    StateCurrentTask    = "current_task"
    StateNextTask       = "next_task"
)
```

## CRUD 함수
| 함수 | 설명 |
|------|------|
| `SetState(db, key, value)` | 상태 설정 (UPSERT) |
| `GetState(db, key)` | 상태 조회 |
| `GetAllStates(db)` | 전체 상태 맵 반환 |
| `DeleteState(db, key)` | 상태 삭제 |

## 상태 업데이트
```go
// Task 실행 시 현재 상태 업데이트
func UpdateCurrentState(db, projectID, phaseID string, taskID, nextTaskID int64) error

// 새 프로젝트 초기화
func InitState(db, projectID string) error
```

## 용도
- Task pop/start 시 현재 위치 추적
- 실행 재개 시 컨텍스트 복원
- Manifest에 포함되어 Claude에게 전달
