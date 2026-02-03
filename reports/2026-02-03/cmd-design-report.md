# design.go Report

**파일 경로**: `internal/cmd/design.go`
**라인 수**: ~90줄

## 요약
설계 결정 관리 명령어.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `design set '<json>'` | 설계 결정 설정 |
| `design get` | 설계 결정 조회 |

## 사용 예시
```bash
clari design set '{"architecture": "monolith", "auth_method": "jwt", "api_style": "rest"}'
```
