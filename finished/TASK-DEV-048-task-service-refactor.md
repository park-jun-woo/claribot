# TASK-DEV-048: Task Service 리팩토링

## 목표
phase_id 관련 로직 제거, feature_id 필수화.

## 파일
`internal/service/task_service.go`

## 변경 내용

### 1. TaskCreateInput 수정
```go
type TaskCreateInput struct {
    FeatureID      int64   // 필수 (pointer 아님)
    // PhaseID 제거
    ParentID       *int64
    SkeletonID     *int64
    Title          string
    Content        string
    TargetFile     string
    TargetFunction string
    TargetLine     *int
    // Level, Skill, References 제거
}
```

### 2. CreateTask 함수 수정
- phase_id 컬럼 제거
- feature_id를 NOT NULL로 처리
- level, skill, references 컬럼 제거

### 3. ListTasks 함수 수정
- `ListTasks(db, phaseID)` → `ListTasksByFeature(db, featureID)`로 통합
- 또는 phaseID 파라미터 제거

### 4. PopTask / PopTaskFull 수정
- phase 관련 조회 제거

### 5. 기타 함수들
- GetTaskStatus: phase 관련 제거
- scanTask: 필드 변경 반영

## 완료 조건
- [ ] TaskCreateInput에서 PhaseID 제거
- [ ] FeatureID를 필수로 변경
- [ ] level, skill, references 관련 코드 제거
- [ ] 모든 SQL 쿼리에서 phase_id 제거
