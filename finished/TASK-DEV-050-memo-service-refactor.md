# TASK-DEV-050: Memo Service 리팩토링

## 목표
Memo scope에서 phase → feature로 변경.

## 파일
`internal/service/memo_service.go`

## 변경 내용

### 1. Scope 상수/검증 수정
```go
// 유효한 scope
validScopes := []string{"project", "feature", "task"}
// "phase" 제거
```

### 2. Key 파싱 로직 수정
현재:
- `key` → project
- `1:key` → phase (id=1)
- `1:5:key` → task (phase_id=1, task_id=5)

변경:
- `key` → project
- `1:key` → feature (id=1)
- `1:5:key` → task (feature_id=1, task_id=5)

### 3. 관련 함수들 수정
- SetMemo
- GetMemo
- ListMemos
- DeleteMemo
- GetMemosByScope

## 완료 조건
- [ ] scope 검증에서 phase → feature
- [ ] Key 파싱 로직 수정 (phase → feature)
- [ ] 관련 문서/주석 수정
