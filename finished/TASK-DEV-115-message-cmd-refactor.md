# TASK-DEV-115: Message 명령어 리팩토링

## 목표
Message 명령어를 Feature와 동일한 패턴으로 수정

## 변경 내용

### cli/internal/cmd/message.go

#### runMessageSend 수정
- `GetState(database, StateCurrentProject)` 대신 `GetProject(database)` 사용
- `--content`와 `--feature` 플래그 지원 (feature add의 --name, --description 패턴)

```go
func runMessageSend(cmd *cobra.Command, args []string) error {
    // ...

    // Get project from DB directly (like feature.go)
    project, err := service.GetProject(database)
    if err != nil {
        outputError(fmt.Errorf("get project: %w", err))
        return nil
    }

    // Try flags first, then positional argument
    content, _ := cmd.Flags().GetString("content")
    if content == "" && len(args) > 0 {
        content = args[0]
    }

    // ...
}
```

#### runMessageList 수정
- 동일하게 `GetProject(database)` 사용

## 완료 조건
- [ ] GetProject() 사용으로 변경
- [ ] --content 플래그 추가
- [ ] project set 없이 동작 확인
