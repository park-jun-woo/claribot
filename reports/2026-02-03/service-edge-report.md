# edge_service.go Report

**파일 경로**: `internal/service/edge_service.go`
**라인 수**: 794줄

## 요약
Task/Feature 간 의존성 그래프 관리. 사이클 감지, 위상 정렬, LLM 의존성 추론 지원.

## Task Edge 관리
| 함수 | 설명 |
|------|------|
| `AddTaskEdge(db, from, to)` | 의존성 추가 (사이클 체크) |
| `RemoveTaskEdge(db, from, to)` | 의존성 제거 |
| `GetTaskEdges(db)` | 전체 Edge 조회 |
| `GetTaskEdgesByPhase(db, phaseID)` | Phase별 Edge |
| `GetTaskEdgesByFeature(db, featureID)` | Feature별 Edge |
| `GetTaskDependencies(db, taskID)` | Task가 의존하는 목록 |
| `GetTaskDependents(db, taskID)` | Task에 의존하는 목록 |

## 사이클 감지 (DFS)
```go
func CheckTaskCycle(db, fromID, toID) (bool, []string, error)
func DetectAllCycles(db) ([][]string, error)
```
- DFS 알고리즘 사용
- 사이클 발견 시 경로 반환

## 위상 정렬 (Kahn's Algorithm)
```go
func TopologicalSortTasks(db, phaseID) ([]Task, error)
func TopologicalSortFeatures(db, projectID) ([]Feature, error)
```
- In-degree 기반 정렬
- 의존성 순서대로 정렬된 목록 반환

## 실행 가능 Task 조회
```go
func GetExecutableTasks(db) ([]Task, error)
func IsTaskExecutable(db, taskID) (bool, []Task, error)
```
- 모든 의존성이 `done` 상태인 pending Task

## LLM 의존성 추론
```go
type InferenceContext struct {
    Type     string          // "task" or "feature"
    Items    []InferenceItem
    Existing []InferenceExisting
    Prompt   string
}

func PrepareTaskEdgeInference(db, featureID) (*InferenceContext, error)
func PrepareFeatureEdgeInference(db) (*InferenceContext, error)
```

### 추론 규칙 (프롬프트에 포함)
- Model → Service → API → UI 순서
- 동일 엔티티 참조 시 의존성
- 인증 Feature는 보호 Feature보다 먼저

### 레이어 추론
```go
func inferLayerFromPath(path string) string
// model, service, api, ui, unknown
```

## Edge 목록 조회
```go
type EdgeListResult struct {
    FeatureEdges []FeatureEdgeItem
    TaskEdges    []TaskEdgeItem
    TotalFeature, TotalTask int
}

func ListAllEdges(db) (*EdgeListResult, error)
```
