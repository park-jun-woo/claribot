# TASK-DEV-097: Feature 응답 필드 추가

## 개요
`feature get`과 `feature list` 명령어 응답에 스펙에 정의된 누락 필드 추가

## 배경
- **보고서**: reports/2026-02-03-total-report.md Issue #1
- **스펙**: specs/CLI/07-Feature.md
- **현재 구현**: cli/internal/cmd/feature.go:169-179

## 스펙 요구사항

### feature get 응답
```json
{
  "feature": {
    "id": 1,
    "name": "user_auth",
    "spec": "사용자 인증 시스템",
    "fdl": "feature: user_auth\n...",
    "fdl_hash": "abc123...",
    "skeleton_generated": true
  }
}
```

### feature list 응답
```json
{
  "features": [
    {
      "id": 1,
      "name": "user_auth",
      "spec": "사용자 인증 시스템",
      "status": "active",
      "fdl_hash": "abc123...",
      "skeleton_generated": true
    }
  ]
}
```

## 작업 내용

### 1. Feature 모델 확인 (models.go)
- `fdl`, `fdl_hash`, `skeleton_generated` 필드가 모델에 있는지 확인
- 없으면 Feature 구조체에 필드 추가

### 2. feature get 수정 (feature.go)
**파일**: `cli/internal/cmd/feature.go`
**위치**: 169-179 라인

현재:
```go
outputJSON(map[string]interface{}{
    "success": true,
    "feature": map[string]interface{}{
        "id":          feature.ID,
        "name":        feature.Name,
        "description": feature.Description,
        "spec":        feature.Spec,
        "status":      feature.Status,
        "created_at":  feature.CreatedAt,
    },
})
```

수정:
```go
outputJSON(map[string]interface{}{
    "success": true,
    "feature": map[string]interface{}{
        "id":                 feature.ID,
        "name":               feature.Name,
        "description":        feature.Description,
        "spec":               feature.Spec,
        "status":             feature.Status,
        "created_at":         feature.CreatedAt,
        "fdl":                feature.FDL,
        "fdl_hash":           feature.FDLHash,
        "skeleton_generated": feature.SkeletonGenerated,
    },
})
```

### 3. feature list 수정 (feature.go)
- feature list 응답에도 `fdl_hash`, `skeleton_generated` 필드 추가

### 4. feature_service.go 확인
- GetFeature 함수가 fdl, fdl_hash, skeleton_generated를 조회하는지 확인
- 필요시 쿼리 수정

## 검증
```bash
cd cli && go build -o clari ./cmd/claritask && ./clari feature get 1
```

## 완료 기준
- [ ] feature get 응답에 fdl, fdl_hash, skeleton_generated 필드 포함
- [ ] feature list 응답에 fdl_hash, skeleton_generated 필드 포함
- [ ] 테스트 통과
