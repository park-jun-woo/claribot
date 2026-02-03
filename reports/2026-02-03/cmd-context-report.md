# context.go Report

**파일 경로**: `internal/cmd/context.go`
**라인 수**: 97줄

## 요약
프로젝트 컨텍스트 관리 명령어. JSON 기반 설정.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `context set '<json>'` | 컨텍스트 설정 |
| `context get` | 컨텍스트 조회 |

## 필수 필드 검증
```go
func validateContext(data map[string]interface{}) error {
    required := []string{"project_name", "description"}
    // ...
}
```

## 사용 예시
```bash
clari context set '{"project_name": "MyApp", "description": "A web app"}'
clari context get
```
