# memo.go Report

**파일 경로**: `internal/cmd/memo.go`
**라인 수**: ~150줄

## 요약
메모 관리 명령어. Project/Phase/Task 스코프 지원.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `memo set <key> '<json>'` | 메모 설정 |
| `memo get <key>` | 메모 조회 |
| `memo list` | 메모 목록 |
| `memo delete <key>` | 메모 삭제 |

## 키 포맷
- Project: `key`
- Phase: `1:key`
- Task: `1:5:key`

## memo set 플래그
```
--priority <1|2|3>  # 1: high, 2: default, 3: low
--scope <scope>     # project, phase, task
--scope-id <id>     # scope의 ID
```

## 사용 예시
```bash
clari memo set api_convention '{"value": "REST API 명명 규칙..."}'
clari memo get api_convention
clari memo list
```
