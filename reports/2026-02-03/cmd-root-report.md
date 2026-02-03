# root.go Report

**파일 경로**: `internal/cmd/root.go`
**라인 수**: 68줄

## 요약
Cobra 루트 명령어와 공통 유틸리티 함수 정의.

## 루트 명령어
```go
var rootCmd = &cobra.Command{
    Use:   "clari",
    Short: "Clear Task Management for Claude Code",
    Long:  "Claude Code를 위한 장시간 자동 실행 시스템",
}
```

## 공통 유틸리티 함수
| 함수 | 설명 |
|------|------|
| `Execute()` | 루트 명령어 실행 (main에서 호출) |
| `getDB()` | 현재 디렉토리의 `.claritask/db` 열기 |
| `outputJSON(v)` | JSON 출력 (pretty print) |
| `outputError(err)` | 에러 JSON 출력 |
| `parseJSON(str, v)` | JSON 문자열 파싱 |

## 등록된 서브명령어 (14개)
```go
init, project, context, tech, design, required,
phase, task, memo, feature, edge, fdl, plan
```

## DB 경로 규칙
```
{현재_디렉토리}/.claritask/db
```
