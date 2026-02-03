# 2026-02-03: VSCode Extension Expert 기능 및 TTY Handover 구현

## 개요
- VSCode Extension에 Expert 탭 기능 추가
- CLI TTY Handover 기능 구현

---

## 1. VSCode Extension Expert 기능 (TASK-EXT-023 ~ 029)

### 배경
- VSCode Extension에 Experts 탭이 미구현 상태
- specs/VSCode/07-ExpertsTab.md, 08-ExpertSync.md 스펙 존재

### 완료된 작업

#### TASK-EXT-023: Expert 타입 정의 및 Store 확장
- `vscode-extension/src/types.ts` - Expert 인터페이스, 메시지 타입 추가
- `vscode-extension/webview-ui/src/types.ts` - Expert 인터페이스, ProjectData 확장
- `vscode-extension/webview-ui/src/store.ts` - experts, selectedExpertId 상태 추가

#### TASK-EXT-024: Expert DB 쿼리 함수
- `vscode-extension/src/database.ts`에 Expert CRUD 함수 추가
  - `getExperts()`, `getProjectExperts()`
  - `assignExpert()`, `unassignExpert()`
  - `createExpert()`, `updateExpertContent()`, `updateExpertMetadata()`

#### TASK-EXT-025: Expert 메시지 프로토콜
- `CltEditorProvider.ts` - Expert 메시지 핸들러 추가
- `vscode.ts` - `assignExpert()`, `unassignExpert()`, `openExpertFile()`, `createExpert()` 함수

#### TASK-EXT-026: ExpertsPanel 컴포넌트
- `components/ExpertsPanel.tsx` 생성
- `App.tsx`에 Experts 탭 추가 (4번째 탭)

#### TASK-EXT-027: ExpertCard 컴포넌트
- `components/ExpertCard.tsx` 생성
- Assign/Unassign, Open File 버튼, 메타 정보 표시

#### TASK-EXT-028: CreateExpertDialog 컴포넌트
- `components/CreateExpertDialog.tsx` 생성
- Expert ID 유효성 검증 (영문 소문자, 숫자, 하이픈만)

#### TASK-EXT-029: Expert 파일 동기화 (FileSystemWatcher)
- `CltEditorProvider.ts`에 FileSystemWatcher 추가
- EXPERT.md 파일 변경 감지 → DB 자동 동기화
- 파일 삭제 시 DB 백업에서 복구

---

## 2. TTY Handover 기능 (TASK-TTY-001 ~ 005)

### 배경
- specs/TTY/ 스펙 존재하나 구현 미비
- `debug_service.go`에 부분 구현만 존재
- Phase 1 (대화형 요구사항 수립), Phase 2 (자동 실행) 모두 미완성

### 기존 상태
```
| 기능                    | 상태       |
|------------------------|-----------|
| TTY Handover 기본 구조   | 구현됨 (debug_service.go) |
| Headless 모드           | 구현됨 (claude --print) |
| Phase 1 대화형 init     | 미구현     |
| Phase 2 Fallback 연결   | TODO 상태  |
| task run/retry 명령어   | 미구현     |
```

### 완료된 작업

#### TASK-TTY-001: tty_service.go 생성
**파일**: `cli/internal/service/tty_service.go`

핵심 함수:
- `RunWithTTYHandover(systemPrompt, initialPrompt, permissionMode)` - Claude에 TTY 제어권 넘김
- `Phase1SystemPrompt()` - 요구사항 수립용 시스템 프롬프트
- `Phase2SystemPrompt()` - Task 실행용 시스템 프롬프트
- `RunTaskWithTTY(database, task)` - Task 단위 TTY 실행
- `BuildTaskPromptForTTY(database, task)` - 초기 프롬프트 빌더
- `InferTestCommand(targetFile)` - 파일 확장자 기반 테스트 명령 추론
- `VerifyTask(task)` - 사후 검증
- `RunInteractiveInit(database, projectID, name, description)` - Phase 1 실행

#### TASK-TTY-002: Phase 1 TTY (init -i)
**파일**: `cli/internal/cmd/init.go`

- `--interactive` (`-i`) 플래그 추가
- 프로젝트 초기화 후 `RunInteractiveInit()` 호출

사용법:
```bash
clari init my-project "프로젝트 설명" -i
```

#### TASK-TTY-003: task run/retry 명령어
**파일**: `cli/internal/cmd/task.go`

- `taskRunCmd` - 개별 Task TTY 실행
- `taskRetryCmd` - 실패한 Task 재시도 (pending으로 리셋 후 실행)

사용법:
```bash
clari task run <task_id>
clari task retry <task_id>
```

#### TASK-TTY-004: Phase 2 Orchestrator TTY 연결
**파일**: `cli/internal/service/orchestrator_service.go`

`ExecuteAllTasks()` 수정:
- Headless 실패 시 `RunTaskWithTTY()` 호출
- 사후 검증 `VerifyTask()` 연동
- 진행 상황 콘솔 출력

**파일**: `cli/internal/cmd/project.go`

- `project start`를 foreground 실행으로 변경 (기존: goroutine)
- 최종 결과 JSON 출력

사용법:
```bash
clari project start --fallback-interactive
```

#### TASK-TTY-005: debug_service.go 정리
- 중복 코드 제거
- deprecated 함수로 래핑 (하위 호환성)
- `RunInteractiveDebugging()` → `RunTaskWithTTY()` 위임
- `VerifyAfterDebugging()` → `VerifyTask()` 위임

---

## 3. 새로운 CLI 명령어 요약

```bash
# Phase 1: 대화형 요구사항 수립
clari init <project-id> "<description>" -i

# Phase 2: 자동 실행 (headless 실패 시 TTY 전환)
clari project start --fallback-interactive

# 개별 Task TTY 실행
clari task run <task_id>

# 실패한 Task 재시도
clari task retry <task_id>
```

---

## 4. TTY Handover 플로우

### Phase 1
```
$ clari init my-project "중고거래 플랫폼" -i
    ↓
[프로젝트 초기화]
    ↓
[TTY Handover → Claude Code]
    ↓
사용자: "당근마켓처럼"
Claude: "기능 제안..." → clari feature add '...'
사용자: "좋아, 개발해"
Claude: clari project start → /exit
    ↓
[Claude Code 종료]
```

### Phase 2
```
$ clari project start --fallback-interactive
    ↓
[Task 목록 조회]
    ↓
[Task N 실행 (headless: claude --print)]
    ↓ (실패 시)
[TTY Handover → Claude Code]
    ↓
[사후 검증]
    ↓
[다음 Task로]
```

---

## 5. 빌드/테스트 결과
- CLI 빌드: 성공
- VSCode Extension 빌드: 성공
- Webview UI 빌드: 성공
- 테스트: 모두 통과 (`ok parkjunwoo.com/claritask/test 0.909s`)
