# TASK-DEV-061: Init Service 구현

## 목표
`clari init` 전체 프로세스를 관리하는 서비스 구현.

## 파일
`internal/service/init_service.go`

## 구현 내용

### 1. 타입 정의

```go
type InitConfig struct {
    ProjectID       string
    Name            string
    Description     string
    SkipAnalysis    bool
    SkipSpecs       bool
    NonInteractive  bool
    Force           bool
}

type InitState struct {
    Phase           int       `json:"phase"`        // 1-5
    StartedAt       time.Time `json:"started_at"`
    Tech            map[string]interface{} `json:"tech,omitempty"`
    Design          map[string]interface{} `json:"design,omitempty"`
    Context         map[string]interface{} `json:"context,omitempty"`
    CurrentSpecs    string    `json:"current_specs,omitempty"`
    SpecsRevision   int       `json:"specs_revision"`
}

type InitResult struct {
    Success     bool
    ProjectID   string
    DBPath      string
    SpecsPath   string
    Error       string
}
```

### 2. Phase 상수

```go
const (
    InitPhaseDBInit       = 1
    InitPhaseAnalysis     = 2
    InitPhaseApproval     = 3
    InitPhaseSpecsGen     = 4
    InitPhaseFeedback     = 5
    InitPhaseComplete     = 6
)
```

### 3. RunInit 함수 (메인 진입점)

```go
func RunInit(config InitConfig) (*InitResult, error)
```

- Phase 1-5 순차 실행
- 상태 저장 및 복구
- 에러 처리 및 롤백

### 4. Phase 1: InitPhase1_DBInit

```go
func InitPhase1_DBInit(config InitConfig) (*db.DB, error)
```

- .claritask/ 디렉토리 생성
- DB 생성 및 마이그레이션
- 프로젝트 레코드 생성
- 상태 초기화

### 5. Phase 2: InitPhase2_Analysis

```go
func InitPhase2_Analysis(database *db.DB, dir, description string) (*ContextAnalysisResult, error)
```

- ScanProjectFiles 호출
- BuildContextAnalysisPrompt로 프롬프트 생성
- CallClaude로 LLM 호출
- ParseContextAnalysis로 결과 파싱

### 6. Phase 3: InitPhase3_Approval

```go
func InitPhase3_Approval(database *db.DB, result *ContextAnalysisResult) error
```

- PrintAnalysisResult로 결과 출력
- PromptApproval로 사용자 선택
- 승인 시 DB에 저장
- 수정 시 PromptEdit → 다시 출력
- 재분석 시 Phase 2로 돌아감

### 7. Phase 4: InitPhase4_SpecsGen

```go
func InitPhase4_SpecsGen(database *db.DB, config InitConfig) (string, error)
```

- DB에서 tech, design, context 조회
- BuildSpecsGenerationPrompt로 프롬프트 생성
- CallClaude로 LLM 호출
- ParseSpecsDocument로 결과 파싱

### 8. Phase 5: InitPhase5_Feedback

```go
func InitPhase5_Feedback(database *db.DB, specs string) (string, error)
```

- PrintSpecs로 스펙 출력
- PromptApproval로 사용자 선택
- 승인 시 파일 저장 및 완료
- 피드백 시 PromptMultilineInput → 수정 → 반복

### 9. SaveInitState / LoadInitState

```go
func SaveInitState(database *db.DB, state *InitState) error
func LoadInitState(database *db.DB) (*InitState, error)
```

- state 테이블에 init 상태 저장/로드
- JSON 직렬화

### 10. ResumeInit

```go
func ResumeInit(database *db.DB) (*InitResult, error)
```

- 저장된 상태에서 재개

## 완료 조건
- [ ] InitConfig, InitState, InitResult 타입 정의
- [ ] Phase 상수 정의
- [ ] RunInit 메인 함수 구현
- [ ] InitPhase1_DBInit 구현
- [ ] InitPhase2_Analysis 구현
- [ ] InitPhase3_Approval 구현
- [ ] InitPhase4_SpecsGen 구현
- [ ] InitPhase5_Feedback 구현
- [ ] SaveInitState, LoadInitState 구현
- [ ] ResumeInit 구현
