# TASK-DEV-064: LLM Service 테스트

## 목표
LLM Service 단위 테스트 작성.

## 파일
`test/llm_service_test.go`

## 테스트 케이스

### 1. TestParseContextAnalysis_ValidJSON

```go
input := `Some text before
` + "```json" + `
{
  "tech": {"language": "Go"},
  "design": {"architecture": "MVC"},
  "context": {"project_type": "CLI"}
}
` + "```" + `
Some text after`
```

- JSON 블록 추출 및 파싱 검증

### 2. TestParseContextAnalysis_NoCodeBlock

- JSON이 코드 블록 없이 바로 출력된 경우
- 전체 출력에서 JSON 추출 시도

### 3. TestParseContextAnalysis_InvalidJSON

- 잘못된 JSON
- 에러 반환 검증

### 4. TestParseSpecsDocument_MarkdownBlock

```go
input := `Here is the specs:
` + "```markdown" + `
# Project Specs
...
` + "```"
```

- Markdown 블록 추출 검증

### 5. TestParseSpecsDocument_NoBlock

- 코드 블록 없이 Markdown 출력
- 전체 출력 반환

### 6. TestBuildContextAnalysisPrompt

- 프롬프트 구성 검증
- 필수 섹션 포함 확인

### 7. TestBuildSpecsGenerationPrompt

- 프롬프트 구성 검증
- tech, design, context 포함 확인

### 8. TestBuildSpecsRevisionPrompt

- 기존 스펙 + 피드백 포함 검증

## 참고
- CallClaude는 실제 claude CLI 필요하므로 통합 테스트에서 검증
- 또는 mock 사용

## 완료 조건
- [ ] 8개 테스트 케이스 구현
- [ ] go test 통과
