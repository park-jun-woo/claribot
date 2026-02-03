# 2026-02-04 VSCode Extension 및 FDL 개선

## 1. VSCode Extension 수정사항

### TASK-EXT-035: Feature CLI 경로 수정
- **문제**: VSCode에서 feature 추가 시 WSL 터미널에서 Windows 경로 인식 불가
- **해결**: `windowsToWslPath()` 함수 추가
  - `C:\Users\mail\git\claritask` → `/mnt/c/Users/mail/git/claritask`
  - 터미널에서 `cd` 후 CLI 실행

### TASK-EXT-036: Feature Delete 버튼 수정
- **문제 1**: `window.confirm()`이 VSCode webview에서 미지원
- **해결**: `ConfirmModal` React 컴포넌트 추가
- **문제 2**: task_edges 미삭제로 FK 제약조건 위반
- **해결**: `deleteFeature()`에서 task_edges 먼저 삭제

### 버전 업데이트
- v0.0.4 → v0.0.7

## 2. CLI 개선사항

### --dangerously-skip-permissions 옵션 추가
- **위치**: `tty_service.go` - `RunWithTTYHandover()`
- **영향**: 모든 TTY handover 모드에서 사용자 확인 없이 자동 실행

### TASK-DEV-107: FDL 문서 Embed 및 Complete 파일 감지

#### FDL 핵심 문서 통합
- **파일**: `cli/internal/docs/fdl_spec.md` (약 200줄)
- 구조, 필드 타입, 제약조건, 예시, 네이밍 컨벤션 포함

#### Go Embed 설정
```go
// cli/internal/docs/embed.go
package docs
import _ "embed"
//go:embed fdl_spec.md
var FDLSpec string
```

#### 시스템 프롬프트 수정
- `FDLGenerationSystemPrompt()`에 `docs.FDLSpec` 포함
- "clari fdl register" 지시 제거
- ".claritask/complete" 파일 생성 지시 추가

#### Complete 파일 감지 로직
- `RunWithTTYHandoverEx()` 함수 추가
- `watchCompleteFile()` goroutine으로 파일 감시 (500ms 간격)
- 파일 감지 시 Claude 프로세스 종료 및 파일 삭제

## 3. 토론 내용

### FDL 문서 전달 방식
- **결론**: Go embed로 CLI에 내장하는 것이 적합
- **이유**: 단일 바이너리 배포, 84KB 추가로 무시 가능한 크기

### FDL 가이드북 준수 확인
- `features/회원기능.fdl.yaml` 분석 (919줄)
- **결과**: 거의 완벽히 준수
  - 네이밍 컨벤션 (snake_case, camelCase, PascalCase)
  - 4계층 구조 (models, service, api, ui)
  - 모든 API에 `use:` wiring 명시
  - 모든 UI action에 API 경로 명시

## 4. 파일 변경 목록

### 신규 생성
- `cli/internal/docs/embed.go`
- `cli/internal/docs/fdl_spec.md`

### 수정
- `vscode-extension/src/CltEditorProvider.ts` - WSL 경로 변환, windowsToWslPath()
- `vscode-extension/src/database.ts` - deleteFeature() task_edges 삭제 추가
- `vscode-extension/webview-ui/src/components/FeatureList.tsx` - ConfirmModal 추가
- `cli/internal/service/tty_service.go` - FDL 문서 전달, complete 파일 감지

### 완료된 TASK
- TASK-EXT-035-feature-cli-path-fix.md
- TASK-EXT-036-feature-delete-fix.md
- TASK-DEV-107-fdl-embed-complete.md
