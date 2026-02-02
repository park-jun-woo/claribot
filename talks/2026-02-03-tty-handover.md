# 2026-02-03 TTY Handover 검토 대화

## 주제
Claritask에서 디버깅 필요 시 Claude Code에 터미널 제어권을 넘기는 방식 검토

---

## 1. 제안된 아이디어

**TTY Handover (터미널 핸드오버)**:
- Claritask가 평소에는 `claude --print`로 비대화형 실행
- 테스트 실패 등 디버깅 필요 시 대화형 모드로 전환
- stdin/stdout/stderr를 Claude에 연결하여 사용자가 모니터링/개입 가능
- Claude 종료 후 Claritask가 제어권 회수

```
Claritask (비대화형) → 실패 감지 → Claude (대화형) → 종료 → Claritask 복귀
```

---

## 2. 검토 결과

### `-p` 옵션 해석 오류

제안서에서 `-p`를 "프롬프트 주입"으로 해석했으나, 실제로는 `--print`의 약자 (비대화형 모드).

```bash
# 틀림
claude -p "context"

# 맞음: 대화형 + 프롬프트
claude "테스트 실행해"
claude --system-prompt "역할" "첫 메시지"
```

### 권한 모드 필요

자동 실행하려면 `--permission-mode` 사용:
```bash
claude --permission-mode acceptEdits "버그 고쳐"
claude --permission-mode bypassPermissions "pytest 실행"
```

### 기술적 가능 여부

| 항목 | 가능 여부 |
|------|----------|
| TTY 핸드오버 | **가능** (Go exec 패키지) |
| 프롬프트 주입 | **가능** (positional arg + --system-prompt) |
| 자동 실행 | **가능** (--permission-mode) |
| 사용자 개입 | **가능** (stdin 연결) |
| 제어권 복귀 | **가능** (cmd.Run() 블로킹) |

---

## 3. 수정된 Go 구현

```go
func RunInteractiveDebugging(task Task, contextPacket string) error {
    systemPrompt := `You are in Claritask debugging mode.
Run test → analyze error → edit code → repeat until pass.
When fixed, exit with /exit.`

    initialPrompt := fmt.Sprintf(`
[Task: %s] [Target: %s] [Test: %s]
Context: %s
Start by running the test now.
`, task.ID, task.TargetFile, task.TestCmd, contextPacket)

    cmd := exec.Command("claude",
        "--system-prompt", systemPrompt,
        "--permission-mode", "acceptEdits",
        initialPrompt,
    )

    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}
```

---

## 4. 저장된 문서

`specs/TTY-Handover.md` 생성:
- 아키텍처 다이어그램
- Claude CLI 옵션 정리
- Go 구현 코드 전체
- CLI 명령어 (`clari task debug`)
- 프롬프트 전략 (Auto-Pilot Trigger)
- 고려사항 (타임아웃, 무한 루프 방지, 로깅)
- 사용 시나리오 예시

---

## 결론

**"평소에는 자동, 막히면 수동"**

TTY Handover 방식으로 비대화형/대화형 모드의 장점을 모두 활용 가능.
기술적으로 구현 가능하며, specs/TTY-Handover.md에 상세 설계 저장 완료.
