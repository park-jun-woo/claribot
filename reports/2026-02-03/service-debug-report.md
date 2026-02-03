# debug_service.go Report

**파일 경로**: `internal/service/debug_service.go`
**라인 수**: 199줄

## 요약
TTY Handover 기반 인터랙티브 디버깅. 실패한 Task를 Claude와 대화형으로 해결.

## 핵심 함수
```go
// 인터랙티브 디버깅 세션 시작
func RunInteractiveDebugging(db, task) error

// 헤드리스 실패 시 인터랙티브로 폴백
func ExecuteWithFallback(db, task, manifest) error

// 디버깅 후 검증
func VerifyAfterDebugging(task) (bool, error)
```

## TTY Handover
```go
cmd := exec.Command("claude",
    "--system-prompt", systemPrompt,
    "--permission-mode", "acceptEdits",
    initialPrompt,
)

// TTY 연결
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr

cmd.Run()  // 블로킹
```

## 시스템 프롬프트
```
ROLE: Debug and fix failing tests autonomously.

WORKFLOW:
1. Run the test command
2. Analyze the error output
3. Read the relevant code
4. Edit the code to fix
5. Run the test again
6. Repeat until pass

CONSTRAINTS:
- Do NOT modify function signatures
- Only implement TODO sections
- Follow FDL specification

COMPLETION:
- Test passes → summarize and /exit
- 3 attempts fail → explain blocker and exit
```

## 테스트 명령어 추론
```go
func inferTestCommand(task) string
```
| 확장자 | 테스트 명령어 |
|--------|--------------|
| `.py` | `pytest <test_file> -v` |
| `.go` | `go test <dir> -v` |
| `.ts/.tsx` | `npm test` |
| `.js/.jsx` | `npm test` |

## 초기 프롬프트 구성
1. Task 정보 (ID, Target File, Function)
2. 테스트 명령어
3. FDL Specification (있으면)
4. 현재 코드 (스켈레톤)
5. 실행 지시사항

## 상수
```go
const MaxDebugAttempts = 3
```
