# edge.go Report

**파일 경로**: `internal/cmd/edge.go`
**라인 수**: 341줄

## 요약
의존성 Edge 관리 명령어. Feature/Task 간 의존성 추가/제거/조회.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `edge add` | 의존성 추가 |
| `edge list` | 의존성 목록 |
| `edge remove` | 의존성 제거 |
| `edge infer` | LLM으로 의존성 추론 |

## edge add 플래그
```
--from <id>    # 의존하는 쪽 ID
--to <id>      # 의존받는 쪽 ID
--feature      # Feature Edge (기본: Task Edge)
```

## edge list 플래그
```
--feature      # Feature Edge만
--task         # Task Edge만
--phase <id>   # Phase별 필터
```

## edge infer 플래그
```
--feature <id>     # Feature 내 Task Edge 추론
--project          # 프로젝트 Feature Edge 추론
--min-confidence   # 최소 신뢰도 (기본: 0.7)
```

## 사이클 감지
Edge 추가 시 사이클이 감지되면:
```json
{
  "success": false,
  "error": "Circular dependency detected",
  "cycle": [1, 2, 3]
}
```

## LLM 추론 결과
```json
{
  "success": true,
  "type": "task",
  "items": [...],
  "existing_edges": [...],
  "prompt": "...",
  "instructions": "Use the prompt to infer edges..."
}
```
