# feature.go Report

**파일 경로**: `internal/cmd/feature.go`
**라인 수**: 349줄

## 요약
Feature 관리 명령어 그룹.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `feature list` | Feature 목록 (통계 포함) |
| `feature add '<json>'` | Feature 추가 |
| `feature get <id>` | Feature 조회 |
| `feature spec <id> '<spec>'` | Spec 설정 |
| `feature start <id>` | Feature 실행 시작 |
| `feature tasks <id>` | Feature의 Task 목록/생성 |

## feature add 입력
```go
type featureAddInput struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

## feature tasks
- FDL이 있으면: Task 목록 반환
- FDL이 없으면: LLM 프롬프트 생성

### --generate 플래그
FDL 없이 LLM으로 Task 생성 준비

## 출력 예시
```json
{
  "success": true,
  "features": [
    {
      "id": 1,
      "name": "user_auth",
      "status": "pending",
      "tasks_total": 5,
      "tasks_done": 2,
      "depends_on": []
    }
  ]
}
```
