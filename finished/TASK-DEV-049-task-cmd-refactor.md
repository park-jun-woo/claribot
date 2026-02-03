# TASK-DEV-049: Task 명령어 리팩토링

## 목표
task 명령어에서 phase_id 관련 제거, feature_id 필수화.

## 파일
`internal/cmd/task.go`

## 변경 내용

### 1. task push 입력 수정
```go
type taskPushInput struct {
    FeatureID  int64  `json:"feature_id"`  // 필수
    // PhaseID 제거
    ParentID   *int64 `json:"parent_id"`
    Title      string `json:"title"`
    Content    string `json:"content"`
    // Level, Skill, References 제거
}
```

### 2. 필수 필드 검증 수정
```go
if input.FeatureID == 0 {
    return errors.New("feature_id is required")
}
```

### 3. task list 수정
- `task list <phase_id>` → `task list [feature_id]`
- 또는 feature_id 옵션으로 변경

### 4. JSON 응답에서 phase 관련 필드 제거

## 완료 조건
- [ ] taskPushInput에서 PhaseID 제거
- [ ] FeatureID 필수 검증 추가
- [ ] level, skill, references 제거
- [ ] task list의 phase_id 파라미터 → feature_id로 변경
