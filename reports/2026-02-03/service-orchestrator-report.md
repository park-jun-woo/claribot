# orchestrator_service.go Report

**파일 경로**: `internal/service/orchestrator_service.go`
**라인 수**: 429줄

## 요약
Task 실행 오케스트레이션. Claude 호출, 실행 계획 생성, 진행 상태 추적.

## 실행 상태
```go
type ExecutionState struct {
    Running        bool
    StartedAt      *time.Time
    CurrentTask    *int64
    TotalTasks     int
    CompletedTasks int
    FailedTasks    int
}
```

## 실행 제어 함수
| 함수 | 설명 |
|------|------|
| `GetExecutionState(db)` | 현재 실행 상태 조회 |
| `StartExecution(db)` | 실행 시작 |
| `StopExecution(db)` | 실행 종료 |
| `RequestStop(db)` | 중지 요청 |
| `IsStopRequested(db)` | 중지 요청 확인 |
| `UpdateExecutionCurrentTask(db, id)` | 현재 Task 업데이트 |

## 실행 계획 (Dry-Run)
```go
type ExecutionPlan struct {
    Tasks []PlannedTask
    Total int
}

func GenerateExecutionPlan(db, featureID) (*ExecutionPlan, error)
```
- 위상 정렬된 Task 목록 반환
- Feature 필터링 지원
- pending Task만 포함

## Claude 실행
```go
type ClaudeResult struct {
    Success bool
    Output  string
    Error   string
}

func ExecuteTaskWithClaude(task, manifest) (*ClaudeResult, error)
```
- `claude --print <prompt>` 명령 실행
- Manifest 컨텍스트를 프롬프트에 포함

## 전체 실행
```go
type ExecutionOptions struct {
    FeatureID           *int64 // 특정 Feature만
    DryRun              bool   // 계획만 출력
    FallbackInteractive bool   // 실패 시 인터랙티브
}

func ExecuteAllTasks(db, options) error
func ExecuteFeature(db, featureID) error
```

### 실행 루프
1. 중지 요청 확인
2. `PopTaskFull()` 호출
3. Feature 필터 적용
4. Claude 실행
5. 성공 → `CompleteTask`
6. 실패 → `FailTask` (fallback 시 계속)

## 진행 상태
```go
type ExecutionProgress struct {
    TotalFeatures, CompletedFeatures int
    TotalTasks, Pending, Doing, Done, Failed int
    Progress     float64
    CurrentTask  *TaskInfo
    FailedTasks  []TaskInfo
}

func GetExecutionProgress(db) (*ExecutionProgress, error)
```
