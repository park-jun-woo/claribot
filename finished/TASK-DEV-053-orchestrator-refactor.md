# TASK-DEV-053: Orchestrator Service 리팩토링

## 목표
실행 오케스트레이션에서 phase 관련 로직 제거.

## 파일
`internal/service/orchestrator_service.go`

## 변경 내용

### 1. GenerateExecutionPlan 수정
- Phase별 그룹핑 제거
- Feature별 Task 실행 순서 결정

### 2. ExecuteAllTasks 수정
- Phase 순회 → Feature 순회
- 또는 의존성 기반 순서대로 실행

### 3. 상태 추적
- current_phase → 제거 또는 current_feature로 대체

## 완료 조건
- [ ] Phase 관련 로직 제거
- [ ] Feature 기반 실행으로 변경
