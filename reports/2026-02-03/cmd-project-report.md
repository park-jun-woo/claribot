# project.go Report

**파일 경로**: `internal/cmd/project.go`
**라인 수**: 345줄

## 요약
프로젝트 관리 명령어 그룹 구현.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `project set '<json>'` | 프로젝트 설정 업데이트 |
| `project get` | 프로젝트 정보 조회 |
| `project plan` | 플래닝 모드 시작 |
| `project start` | 실행 시작 |
| `project stop` | 실행 중지 |
| `project status` | 실행 상태 조회 |

## project start 플래그
```go
--feature <id>    // 특정 Feature만 실행
--dry-run         // 실행 계획만 출력
--fallback-interactive // 실패 시 인터랙티브 모드
```

## 주요 흐름

### project set
1. JSON 파싱 (`ProjectSetInput`)
2. `service.SetProjectFull()` 호출
3. name, description, context, tech, design 한번에 설정

### project start
1. 필수 필드 확인 (`CheckRequired`)
2. dry-run이면 실행 계획만 출력
3. Task 상태 확인
4. 백그라운드에서 `ExecuteAllTasks` 실행

### project stop
1. 실행 상태 확인
2. `RequestStop()` 호출
3. 현재 Task 완료 후 중지

## 출력 예시
```json
{
  "success": true,
  "mode": "execution",
  "status": { "total": 10, "pending": 5, "done": 3 },
  "progress": 30.0
}
```
