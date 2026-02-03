# fdl_service.go Report

**파일 경로**: `internal/service/fdl_service.go`
**라인 수**: 850줄

## 요약
FDL(Feature Definition Language) 파싱, 검증, Task 매핑, 구현 검증 기능.

## FDL 구조
```go
type FDLSpec struct {
    Feature     string
    Description string
    Models      []FDLModel
    Service     []FDLService
    API         []FDLAPI
    UI          []FDLUI
}
```

### FDLModel
```go
type FDLModel struct {
    Name   string
    Table  string
    Fields []FDLField  // Name, Type, Constraints
}
```

### FDLService
```go
type FDLService struct {
    Name   string
    Desc   string
    Input  map[string]interface{}
    Output string
    Steps  []string
}
```

### FDLAPI
```go
type FDLAPI struct {
    Path, Method, Summary string
    Use      string  // service.FunctionName
    Request  map[string]interface{}
    Response map[string]interface{}
}
```

### FDLUI
```go
type FDLUI struct {
    Component, Type, Parent string
    Props     map[string]interface{}
    State     []string
    Init      []string
    View      []map[string]interface{}
}
```

## 주요 함수

### 파싱
| 함수 | 설명 |
|------|------|
| `ParseFDL(yaml)` | YAML 문자열 파싱 |
| `ParseFDLFile(path)` | 파일에서 파싱 |

### 검증
| 함수 | 검증 내용 |
|------|----------|
| `ValidateFDL(spec)` | 전체 검증 |
| `validateFeatureName` | 이름 규칙 |
| `validateModels` | 중복, 필드 체크 |
| `validateServices` | 중복, Steps 체크 |
| `validateAPIs` | 메서드, service.use 참조 |
| `validateUIs` | 컴포넌트 타입 체크 |

### Task 매핑
```go
type FDLTaskMapping struct {
    Title, Content string
    TargetFile, TargetFunction string
    Layer        string  // model, service, api, ui
    Dependencies []string
}

func ExtractTaskMappings(spec, tech) ([]FDLTaskMapping, error)
```

### 구현 검증
```go
type VerifyResult struct {
    Valid             bool
    Errors, Warnings  []string
    FunctionsMissing, FunctionsExtra []string
    FilesMissing      []string
    SignatureMismatch []SignatureDiff
    ModelsMissing, APIsMissing []string
}

func VerifyFDLImplementation(db, featureID) (*VerifyResult, error)
```

### Diff 기능
```go
type DiffResult struct {
    FeatureID   int64
    FeatureName string
    Differences []FileDiff
    TotalChanges int
}

func DiffFDLImplementation(db, featureID) (*DiffResult, error)
```

## 파일 경로 규칙 (백엔드별)
| 백엔드 | Model | Service | API |
|--------|-------|---------|-----|
| Go | internal/model/*.go | internal/service/*_service.go | internal/api/*_handler.go |
| Python | models/*.py | services/*_service.py | api/*_router.py |
