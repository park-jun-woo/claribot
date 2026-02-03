# TASK-DEV-055: State Service 리팩토링

## 목표
State에서 current_phase 관련 제거.

## 파일
`internal/service/state_service.go`

## 변경 내용

### 1. 상태 키 정리
유지:
- `current_project`
- `current_feature`
- `current_task`
- `next_task`

제거:
- `current_phase` (있다면)

### 2. 관련 함수들 수정
- GetCurrentPhase → 제거
- SetCurrentPhase → 제거

## 완료 조건
- [ ] current_phase 관련 코드 제거
- [ ] current_feature가 올바르게 동작하는지 확인
