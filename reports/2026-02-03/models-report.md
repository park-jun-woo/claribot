# models.go Report

**파일 경로**: `internal/model/models.go`
**라인 수**: 185줄

## 요약
모든 데이터 모델과 응답 구조체 정의.

## 핵심 모델

### Project
```go
type Project struct {
    ID, Name, Description string
    Status    string // active, archived
    CreatedAt time.Time
}
```

### Phase
```go
type Phase struct {
    ID, ProjectID, Name, Description string
    OrderNum  int
    Status    string // pending, active, done
    CreatedAt time.Time
}
```

### Task
```go
type Task struct {
    ID, PhaseID, Status, Title, Level, Skill string
    ParentID       *string
    FeatureID      *int64
    SkeletonID     *int64
    TargetFile, TargetFunction string
    TargetLine     *int
    References     []string
    Content, Result, Error string
    CreatedAt      time.Time
    StartedAt, CompletedAt, FailedAt *time.Time
}
```

### Feature
```go
type Feature struct {
    ID                int64
    ProjectID, Name, Description, Spec string
    FDL, FDLHash      string
    SkeletonGenerated bool
    Status            string // pending, active, done
    CreatedAt         time.Time
}
```

### Edge 모델
- `FeatureEdge` - Feature 의존성 (FromFeatureID → ToFeatureID)
- `TaskEdge` - Task 의존성 (FromTaskID → ToTaskID)

### 싱글톤 모델
- `Context` - 프로젝트 컨텍스트 (ID=1)
- `Tech` - 기술 스택 (ID=1)
- `Design` - 설계 결정 (ID=1)
- `State` - Key-Value 상태

## 응답 타입
| 타입 | 용도 |
|------|------|
| `Response` | 기본 JSON 응답 |
| `MemoData` | 메모 JSON 출력용 |
| `FDLInfo` | FDL 정보 (task pop 시) |
| `SkeletonInfo` | 스켈레톤 정보 |
| `Dependency` | 의존 Task 정보 |
| `TaskPopResponse` | task pop 전체 응답 |
| `Manifest` | 컨텍스트 매니페스트 |

## Task 상태 흐름
`pending` → `doing` → `done` | `failed`
