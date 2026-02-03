# plan_service.go Report

**파일 경로**: `internal/service/plan_service.go`
**라인 수**: 161줄

## 요약
LLM 기반 Feature 계획 서비스. 프로젝트 컨텍스트에서 Feature 추천.

## 주요 타입
```go
type PlannedFeature struct {
    Name         string   // 기능명
    Description  string   // 설명
    Priority     int      // 1: core, 2: important, 3: nice-to-have
    Dependencies []string // 의존 Feature
}

type FeaturePlan struct {
    Features   []PlannedFeature
    TotalCount int
    Reasoning  string  // LLM 분석 논리
    Prompt     string  // 생성된 프롬프트
}
```

## 컨텍스트 구조
```go
type FeaturePlanContext struct {
    ProjectName string
    Context     map[string]interface{}
    Tech        map[string]interface{}
    Design      map[string]interface{}
    Existing    []ExistingFeature
}
```

## 주요 함수
```go
func PreparePlanFeatures(db) (*FeaturePlan, error)
```
1. 프로젝트 정보 조회
2. Context/Tech/Design 수집
3. 기존 Feature 목록 수집
4. LLM 프롬프트 생성

## LLM 프롬프트 구성
1. 프로젝트 컨텍스트
2. 기술 스택
3. 설계 결정
4. 기존 Feature (중복 방지)
5. 출력 형식 지정 (JSON)

## 우선순위 가이드
| Priority | 설명 |
|----------|------|
| 1 (Core) | MVP에 필수 |
| 2 (Important) | 좋은 UX에 필요 |
| 3 (Nice-to-have) | 추가 기능 |

## 출력 예시
```json
{
  "features": [
    {"name": "user_auth", "priority": 1, "dependencies": []},
    {"name": "user_profile", "priority": 2, "dependencies": ["user_auth"]}
  ],
  "reasoning": "..."
}
```
