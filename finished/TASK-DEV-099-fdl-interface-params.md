# TASK-DEV-099: FDL Interface Layer 경로 파라미터 추출

## 개요
FDL Service의 Interface Layer(API) 파싱에서 경로 파라미터 추출 및 추가 검증 구현

## 배경
- **보고서**: reports/2026-02-03-total-report.md 섹션 4.4
- **현재 구현율**: 65%
- **스펙**: specs/FDL/02-C-InterfaceLayer.md
- **구현 파일**: cli/internal/service/fdl_service.go

## 현재 지원 상태

### 지원됨
- HTTP 메서드: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- Auth 타입: required, optional, none, apikey, bearer
- Request/Response 파싱
- Rate limiting, Transform 설정

### 미지원 (구현 대상)
- 경로 파라미터 추출 ({userId} 등)
- 파일 업로드 제약 검증
- apiGroups 구조
- 상태 코드 검증

## 작업 내용

### 1. 경로 파라미터 추출
```yaml
endpoints:
  - path: /users/{userId}/posts/{postId}
    method: GET
```
- 경로에서 `{paramName}` 패턴 추출
- 추출된 파라미터 목록 저장
- request.params와 교차 검증

구현:
```go
func extractPathParams(path string) []string {
    re := regexp.MustCompile(`\{([^}]+)\}`)
    matches := re.FindAllStringSubmatch(path, -1)
    params := make([]string, 0, len(matches))
    for _, m := range matches {
        params = append(params, m[1])
    }
    return params
}
```

### 2. 파일 업로드 제약 검증
```yaml
request:
  body:
    type: multipart
    fields:
      - name: avatar
        type: file
        maxSize: 5MB
        allowedTypes: [image/jpeg, image/png]
```
- `maxSize` 형식 검증 (숫자+단위: KB, MB, GB)
- `allowedTypes` MIME 타입 형식 검증

### 3. apiGroups 구조 지원
```yaml
apiGroups:
  - prefix: /api/v1
    auth: required
    rateLimit: {requests: 100, window: 60s}
    endpoints:
      - path: /users
        method: GET
```
- 그룹 레벨 설정 상속
- 엔드포인트별 오버라이드 지원

### 4. 상태 코드 검증
```yaml
response:
  success: 200
  errors:
    - code: 404
      message: "User not found"
    - code: 400
      message: "Invalid input"
```
- 상태 코드가 유효한 HTTP 상태 코드인지 검증 (100-599)
- 에러 코드 중복 검증

## 구현 방법

fdl_service.go에 다음 함수 추가/수정:
- `extractPathParams(path string) []string`
- `validatePathParams(path string, declaredParams []string) error`
- `validateFileUpload(field map[string]interface{}) error`
- `parseApiGroups(groups []interface{}) ([]APIGroup, error)`
- `validateHttpStatusCode(code int) error`

### ParsedAPI 구조체 확장
```go
type ParsedEndpoint struct {
    Path        string
    PathParams  []string  // 추가
    Method      string
    Auth        string
    RateLimit   *RateLimit
    Request     *Request
    Response    *Response
}
```

## 검증
```bash
cd cli && go test -v ./test/fdl_service_test.go -run TestInterfaceLayer
```

## 완료 기준
- [ ] 경로 파라미터 자동 추출 구현
- [ ] 파일 업로드 maxSize 형식 검증 구현
- [ ] apiGroups 파싱 및 설정 상속 구현
- [ ] HTTP 상태 코드 유효성 검증 구현
- [ ] 테스트 케이스 추가
