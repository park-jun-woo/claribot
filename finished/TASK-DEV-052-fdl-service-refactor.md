# TASK-DEV-052: FDL Service 리팩토링

## 목표
FDL에서 Task 생성 시 phase 대신 feature 직접 연결.

## 파일
`internal/service/fdl_service.go`

## 변경 내용

### 1. CreateTasksFromFDL 함수 수정
- Phase 생성/조회 로직 제거
- Feature에 직접 Task 연결

현재:
```go
// Phase 생성 또는 조회
phase, err := GetOrCreatePhase(db, projectID, featureName)
// Task 생성 with phase_id
```

변경:
```go
// Feature에 직접 Task 생성
input := TaskCreateInput{
    FeatureID: featureID,  // 직접 연결
    Title: ...,
    Content: ...,
}
```

### 2. ExtractTaskMappings 함수
- 반환값에서 phase 관련 제거 (있다면)

## 완료 조건
- [ ] Task 생성 시 phase_id 대신 feature_id 사용
- [ ] Phase 생성/조회 로직 제거
