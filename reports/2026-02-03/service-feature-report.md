# feature_service.go Report

**파일 경로**: `internal/service/feature_service.go`
**라인 수**: 377줄

## 요약
Feature CRUD, FDL 관리, Feature Edge(의존성) 관리.

## Feature CRUD
| 함수 | 설명 |
|------|------|
| `CreateFeature(db, projectID, name, desc)` | Feature 생성 |
| `GetFeature(db, id)` | Feature 조회 |
| `ListFeatures(db, projectID)` | Feature 목록 |
| `UpdateFeature(db, feature)` | Feature 업데이트 |
| `DeleteFeature(db, id)` | Feature 삭제 (Edge 포함) |

## 상태 전이
| 함수 | 전이 |
|------|------|
| `StartFeature(db, id)` | pending → active |
| `CompleteFeature(db, id)` | active → done |

## Spec/FDL 관리
| 함수 | 설명 |
|------|------|
| `SetFeatureSpec(db, id, spec)` | Spec 설정 |
| `GetFeatureSpec(db, id)` | Spec 조회 |
| `SetFeatureFDL(db, id, fdl)` | FDL 설정 (해시 자동 계산) |
| `GetFeatureFDL(db, id)` | FDL 조회 |
| `CalculateFDLHash(fdl)` | SHA256 해시 계산 |

## Feature Edge 관리
| 함수 | 설명 |
|------|------|
| `AddFeatureEdge(db, from, to)` | 의존성 추가 (사이클 체크) |
| `RemoveFeatureEdge(db, from, to)` | 의존성 제거 |
| `GetFeatureEdges(db, featureID)` | Edge 조회 |
| `GetFeatureDependencies(db, id)` | 의존하는 Feature 목록 |
| `GetFeatureDependents(db, id)` | 의존받는 Feature 목록 |

## 사이클 감지
```go
// DFS로 사이클 감지
func CheckFeatureCycle(db, fromID, toID) (bool, []int64, error)
```

## 목록 조회 (통계 포함)
```go
type FeatureListItem struct {
    ID, Name, Description, Spec, Status string
    TasksTotal, TasksDone int
    DependsOn []int64
}

func ListFeaturesWithStats(db, projectID) ([]FeatureListItem, error)
```
