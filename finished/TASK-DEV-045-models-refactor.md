# TASK-DEV-045: Models 리팩토링

## 목표
Phase 구조체 제거, Task 구조체 수정, Memo scope 변경.

## 파일
`internal/model/models.go`

## 변경 내용

### 1. Phase 구조체 삭제
```go
// 삭제
type Phase struct {
    ID          string
    ProjectID   string
    ...
}
```

### 2. Task 구조체 수정
```go
type Task struct {
    ID             string
    FeatureID      int64   // 필수 (pointer 아님)
    // PhaseID 제거
    // ParentID 유지
    Status         string
    Title          string
    Content        string
    TargetFile     string
    TargetLine     *int
    TargetFunction string
    Result         string
    Error          string
    SkeletonID     *int64
    CreatedAt      time.Time
    StartedAt      *time.Time
    CompletedAt    *time.Time
    FailedAt       *time.Time
    // Level, Skill, References 제거 (specs에 없음)
}
```

### 3. Memo 구조체 주석 수정
```go
type Memo struct {
    Scope     string // project, feature, task (phase 제거)
    ...
}
```

## 완료 조건
- [ ] Phase 구조체 삭제
- [ ] Task에서 PhaseID 제거, FeatureID 필수화
- [ ] Task에서 Level, Skill, References 제거
- [ ] Memo.Scope 주석 수정
