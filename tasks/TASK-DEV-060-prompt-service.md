# TASK-DEV-060: Prompt Service 구현

## 목표
터미널 대화형 입출력 서비스 구현.

## 파일
`internal/service/prompt_service.go`

## 구현 내용

### 1. 타입 정의

```go
type PromptOption struct {
    Key         string
    Label       string
    Description string
}

type ApprovalResult struct {
    Action   string  // "approve", "edit", "reanalyze", "cancel"
    EditKey  string  // 수정할 항목 (예: "tech.database")
    EditValue string // 새 값
}
```

### 2. PrintBox 함수

```go
func PrintBox(title string, content string)
```

- 박스 형태로 내용 출력
- 스펙의 UI 형식 구현

### 3. PrintAnalysisResult 함수

```go
func PrintAnalysisResult(tech, design, context map[string]interface{})
```

- 분석 결과를 보기 좋게 출력
- Tech Stack, Design, Context 섹션별 출력

### 4. PromptApproval 함수

```go
func PromptApproval(options []PromptOption) (*ApprovalResult, error)
```

- `[A] 승인  [E] 수정  [R] 재분석  [Q] 취소` 형태
- 사용자 입력 읽기 (bufio.Scanner)
- 유효한 옵션만 허용

### 5. PromptEdit 함수

```go
func PromptEdit(currentData map[string]interface{}) (string, string, error)
```

- 수정할 항목 키 입력
- 현재 값 표시
- 새 값 입력

### 6. PromptMultilineInput 함수

```go
func PromptMultilineInput(prompt string) (string, error)
```

- 여러 줄 입력 받기
- 빈 줄 두 번으로 종료

### 7. PromptConfirm 함수

```go
func PromptConfirm(message string, defaultYes bool) (bool, error)
```

- Y/n 또는 y/N 형태 확인

### 8. PrintSpecs 함수

```go
func PrintSpecs(specs string)
```

- Markdown 스펙 문서 출력
- 페이지네이션 (내용이 길 경우)

### 9. PrintFinalResult 함수

```go
func PrintFinalResult(projectID, dbPath, specsPath string)
```

- 최종 완료 메시지 박스 출력

## 완료 조건
- [ ] PrintBox, PrintAnalysisResult 함수 구현
- [ ] PromptApproval 함수 구현
- [ ] PromptEdit 함수 구현
- [ ] PromptMultilineInput 함수 구현
- [ ] PromptConfirm 함수 구현
- [ ] PrintSpecs, PrintFinalResult 함수 구현
