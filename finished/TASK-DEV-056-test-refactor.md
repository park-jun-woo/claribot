# TASK-DEV-056: 테스트 리팩토링

## 목표
Phase 제거에 따른 테스트 파일 수정.

## 파일
- `test/phase_service_test.go` - 삭제
- `test/task_service_test.go` - 수정
- `test/memo_service_test.go` - 수정
- `test/edge_service_test.go` - 수정
- `test/db_test.go` - 수정 (스키마 테스트)
- `test/models_test.go` - 수정

## 변경 내용

### 1. phase_service_test.go
- 파일 삭제

### 2. task_service_test.go
- Phase 생성 → Feature 생성으로 변경
- phase_id → feature_id로 변경
- level, skill, references 관련 테스트 제거

### 3. memo_service_test.go
- phase scope → feature scope로 변경
- 키 포맷 테스트 수정

### 4. edge_service_test.go
- GetTaskEdgesByPhase → GetTaskEdgesByFeature

### 5. db_test.go
- phases 테이블 관련 테스트 제거
- tasks 테이블 스키마 테스트 수정

### 6. models_test.go
- Phase 모델 테스트 제거
- Task 모델 테스트 수정

## 완료 조건
- [ ] phase_service_test.go 삭제
- [ ] 모든 테스트에서 phase → feature 변경
- [ ] `go test ./test/...` 통과
