# plan.go Report

**파일 경로**: `internal/cmd/plan.go`
**라인 수**: 110줄

## 요약
LLM 기반 Feature 계획 명령어.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `plan features` | Feature 목록 생성 |

## plan features 플래그
```
--auto-create    # 제안된 Feature 자동 생성
```

## 사용 방법

### 1. 프롬프트 생성
```bash
clari plan features
```
출력:
```json
{
  "success": true,
  "prompt": "...",
  "instructions": "Use the prompt to generate features..."
}
```

### 2. LLM 결과로 Feature 생성
```bash
clari plan features '{"features": [...]}'
```

### 3. 자동 생성
```bash
clari plan features --auto-create '{"features": [...]}'
```
출력:
```json
{
  "success": true,
  "created_features": [...],
  "total_created": 5,
  "message": "Created 5 features automatically"
}
```

## LLM 입력 형식
```json
{
  "features": [
    {
      "name": "user_auth",
      "description": "User authentication",
      "priority": 1,
      "dependencies": []
    }
  ],
  "reasoning": "Analysis explanation"
}
```
