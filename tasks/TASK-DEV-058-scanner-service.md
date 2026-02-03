# TASK-DEV-058: Scanner Service 구현

## 목표
프로젝트 디렉토리를 스캔하여 메타 파일 정보를 수집하는 서비스 구현.

## 파일
`internal/service/scanner_service.go`

## 구현 내용

### 1. 타입 정의

```go
type ScannedFile struct {
    Path     string `json:"path"`
    Type     string `json:"type"`     // "package_json", "go_mod", "readme", etc.
    Content  string `json:"content"`  // 파일 내용 (일부)
}

type ScanResult struct {
    Files       []ScannedFile `json:"files"`
    Directories []string      `json:"directories"`  // src/, internal/, lib/ 등
}
```

### 2. ScanProjectFiles 함수

```go
func ScanProjectFiles(dir string) (*ScanResult, error)
```

- 스캔 대상 파일:
  - `package.json` → type: "package_json"
  - `go.mod` → type: "go_mod"
  - `requirements.txt` → type: "requirements_txt"
  - `pyproject.toml` → type: "pyproject_toml"
  - `Cargo.toml` → type: "cargo_toml"
  - `pom.xml` → type: "pom_xml"
  - `build.gradle` → type: "build_gradle"
  - `README.md` → type: "readme"
  - `docker-compose.yml` → type: "docker_compose"
  - `Makefile` → type: "makefile"
  - `.env.example` → type: "env_example"

- 디렉토리 구조 스캔:
  - `src/`, `internal/`, `lib/`, `pkg/`, `cmd/` 존재 여부

- 각 파일의 처음 500줄만 읽기
- 바이너리 파일 제외

### 3. FormatScanResultForLLM 함수

```go
func FormatScanResultForLLM(result *ScanResult) string
```

- LLM 프롬프트에 포함할 형식으로 스캔 결과 포맷팅
- 파일별 내용 요약

## 완료 조건
- [ ] ScannedFile, ScanResult 타입 정의
- [ ] ScanProjectFiles 함수 구현
- [ ] FormatScanResultForLLM 함수 구현
- [ ] 테스트 작성
