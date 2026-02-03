# TASK-DEV-051: Memo 명령어 리팩토링

## 목표
Memo 명령어에서 phase scope를 feature로 변경.

## 파일
`internal/cmd/memo.go`

## 변경 내용

### 1. 도움말/Usage 수정
```
Key formats:
  key           - project scope
  1:key         - feature scope (feature_id=1)
  1:5:key       - task scope (feature_id=1, task_id=5)
```

### 2. --scope 플래그 옵션 수정
```go
// "phase" 대신 "feature"
cmd.Flags().StringVarP(&scope, "scope", "s", "", "Scope: project, feature, task")
```

### 3. 출력 메시지 수정
- phase 관련 문자열 → feature로 변경

## 완료 조건
- [ ] 도움말에서 phase → feature
- [ ] scope 플래그 옵션 수정
- [ ] 에러 메시지 등 수정
