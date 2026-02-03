# TASK-DEV-107: FDL 문서 Embed 및 Complete 파일 감지

## 목표
1. FDL 문법 핵심을 하나의 파일로 통합
2. Go embed로 CLI에 내장
3. FDL 생성 프롬프트에 문서 전달
4. .claritask/complete 파일 감지하여 TTY 종료

## 세부 작업

### 1. FDL 핵심 문서 통합
- 위치: `cli/internal/docs/fdl_spec.md`
- 내용:
  - 02-Schema.md (전체 구조)
  - 03-Examples.md (예시)
  - 02-A-DataLayer.md (핵심 부분만)
  - 04-Guidelines.md (핵심 부분만)

### 2. Go Embed 설정
```go
// cli/internal/docs/embed.go
package docs

import _ "embed"

//go:embed fdl_spec.md
var FDLSpec string
```

### 3. tty_service.go 수정
- `FDLGenerationSystemPrompt()` 수정
  - FDL 문서 포함
  - "clari fdl register" 지시 제거
  - ".claritask/complete 파일 생성" 지시 추가

### 4. Complete 파일 감지 로직
- `RunFDLGenerationWithTTY()` 수정
- goroutine으로 .claritask/complete 파일 감시
- 파일 감지 시 프로세스 종료

## 테스트
1. `clari feature add --name 'test' --description 'test'` 실행
2. Claude가 FDL 문서를 참조하여 작성하는지 확인
3. FDL 작성 완료 후 .claritask/complete 생성하는지 확인
4. complete 파일 감지 후 TTY 종료되는지 확인
