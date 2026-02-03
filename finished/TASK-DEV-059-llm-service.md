# TASK-DEV-059: LLM Service 구현

## 목표
`claude --print` 비대화형 모드 호출 및 응답 파싱 서비스 구현.

## 파일
`internal/service/llm_service.go`

## 구현 내용

### 1. 타입 정의

```go
type LLMRequest struct {
    Prompt  string
    Timeout time.Duration  // 기본 60초
    Retries int            // 기본 3회
}

type LLMResponse struct {
    Output    string
    Success   bool
    Error     string
    Duration  time.Duration
}

type ContextAnalysisResult struct {
    Tech    map[string]interface{} `json:"tech"`
    Design  map[string]interface{} `json:"design"`
    Context map[string]interface{} `json:"context"`
}
```

### 2. CallClaude 함수

```go
func CallClaude(request LLMRequest) (*LLMResponse, error)
```

- `exec.Command("claude", "--print", prompt)` 실행
- 타임아웃 처리
- 재시도 로직 (실패 시 3회)
- stdout/stderr 캡처

### 3. ParseContextAnalysis 함수

```go
func ParseContextAnalysis(output string) (*ContextAnalysisResult, error)
```

- LLM 출력에서 JSON 블록 추출
- \`\`\`json ... \`\`\` 패턴 파싱
- JSON 파싱 및 검증

### 4. ParseSpecsDocument 함수

```go
func ParseSpecsDocument(output string) (string, error)
```

- LLM 출력에서 Markdown 문서 추출
- \`\`\`markdown ... \`\`\` 패턴 또는 전체 출력

### 5. 프롬프트 템플릿

```go
const ContextAnalysisPromptTemplate = `...`
const SpecsGenerationPromptTemplate = `...`
const SpecsRevisionPromptTemplate = `...`
```

### 6. BuildContextAnalysisPrompt 함수

```go
func BuildContextAnalysisPrompt(scanResult *ScanResult, description string) string
```

### 7. BuildSpecsGenerationPrompt 함수

```go
func BuildSpecsGenerationPrompt(projectID, name, description string,
    tech, design, context map[string]interface{}) string
```

### 8. BuildSpecsRevisionPrompt 함수

```go
func BuildSpecsRevisionPrompt(currentSpecs, feedback string) string
```

## 완료 조건
- [ ] LLMRequest, LLMResponse 타입 정의
- [ ] CallClaude 함수 구현 (exec.Command)
- [ ] ParseContextAnalysis 함수 구현
- [ ] ParseSpecsDocument 함수 구현
- [ ] 프롬프트 템플릿 정의
- [ ] Build*Prompt 함수들 구현
- [ ] 테스트 작성
