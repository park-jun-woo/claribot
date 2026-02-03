# TASK-DEV-047: Phase 명령어 삭제

## 목표
phase.go 명령어 파일 삭제 및 root.go에서 등록 제거.

## 파일
- `internal/cmd/phase.go` - 삭제
- `internal/cmd/root.go` - phase 명령어 등록 제거

## 변경 내용

### 1. phase.go 삭제
- 파일 전체 삭제

### 2. root.go 수정
```go
// 삭제
rootCmd.AddCommand(phaseCmd)
```

## 완료 조건
- [ ] phase.go 파일 삭제
- [ ] root.go에서 phaseCmd 등록 제거
