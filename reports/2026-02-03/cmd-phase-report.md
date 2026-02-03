# phase.go Report

**파일 경로**: `internal/cmd/phase.go`
**라인 수**: 244줄

## 요약
Phase(단계) 관리 명령어 그룹 구현.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `phase create '<json>'` | 새 Phase 생성 |
| `phase list` | Phase 목록 조회 |
| `phase plan <id>` | Phase 플래닝 시작 |
| `phase start <id>` | Phase 실행 시작 |

## phase create 입력
```go
type phaseCreateInput struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    OrderNum    int    `json:"order_num"`
}
```

## phase start 동작
1. Phase 조회
2. 해당 Phase의 Task 조회
3. Task가 없으면 실패
4. pending Task가 없으면 완료 상태
5. Phase가 pending이면 active로 변경
6. pending Task 수 반환

## Phase 상태
- `pending` - 대기 중
- `active` - 실행 중
- `done` - 완료

## 출력 예시
```json
{
  "success": true,
  "phase_id": 1,
  "name": "MVP Development",
  "mode": "execution",
  "pending_tasks": 5,
  "total_tasks": 10
}
```
