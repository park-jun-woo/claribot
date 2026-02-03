# TASK-DEV-054: Edge Service 리팩토링

## 목표
Edge 서비스에서 phase 관련 함수 제거/수정.

## 파일
`internal/service/edge_service.go`

## 변경 내용

### 1. GetTaskEdgesByPhase 함수
- 삭제 또는 GetTaskEdgesByFeature로 대체

### 2. TopologicalSortTasks 함수
- phaseID 파라미터 → featureID 또는 제거

### 3. PrepareTaskEdgeInference 함수
- phase 관련 로직 제거

## 완료 조건
- [ ] Phase 관련 함수 제거/수정
- [ ] Feature 기반으로 변경
