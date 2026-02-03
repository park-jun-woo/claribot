# project_service.go Report

**파일 경로**: `internal/service/project_service.go`
**라인 수**: 336줄

## 요약
프로젝트, 컨텍스트, 기술스택, 설계 결정 관리 서비스.

## 주요 함수

### Project CRUD
| 함수 | 설명 |
|------|------|
| `CreateProject(db, id, name, desc)` | 프로젝트 생성 |
| `GetProject(db)` | 프로젝트 조회 (단일) |
| `UpdateProject(db, p)` | 프로젝트 업데이트 |

### Context/Tech/Design (싱글톤)
| 함수 | 설명 |
|------|------|
| `SetContext(db, data)` | 컨텍스트 설정 (UPSERT) |
| `GetContext(db)` | 컨텍스트 조회 |
| `SetTech(db, data)` | 기술스택 설정 (UPSERT) |
| `GetTech(db)` | 기술스택 조회 |
| `SetDesign(db, data)` | 설계결정 설정 (UPSERT) |
| `GetDesign(db)` | 설계결정 조회 |

### 필수 필드 체크
```go
type MissingField struct {
    Field   string   // 필드명
    Prompt  string   // 입력 프롬프트
    Options []string // 선택 옵션
}

func CheckRequired(db) (*RequiredResult, error)
```

**필수 필드**:
- `context.project_name`, `context.description`
- `tech.backend`, `tech.frontend`, `tech.database`
- `design.architecture`, `design.auth_method`, `design.api_style`

### 통합 설정
```go
type ProjectSetInput struct {
    Name, Description string
    Context, Tech, Design map[string]interface{}
}

func SetProjectFull(db, input) error
```

## 옵션 값
| 필드 | 옵션 |
|------|------|
| backend | go, node, python, java |
| frontend | react, vue, angular, none |
| database | postgresql, mysql, sqlite, mongodb |
| architecture | monolith, microservice, serverless |
| auth_method | jwt, session, oauth, none |
| api_style | rest, graphql, grpc |
