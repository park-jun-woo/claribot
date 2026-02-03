# TASK-DEV-098: FDL Logic Layer 검증 강화

## 개요
FDL Service의 Logic Layer(서비스) 파싱에서 step 구조 검증 강화

## 배경
- **보고서**: reports/2026-02-03-total-report.md 섹션 4.3
- **현재 구현율**: 55%
- **스펙**: specs/FDL/02-B-LogicLayer.md
- **구현 파일**: cli/internal/service/fdl_service.go

## 현재 지원 상태

### 지원됨
- 서비스 기본 구조 (name, desc, input, output, throws, transaction)
- Auth 레벨 검증 (required, optional, none)
- 입력 파라미터 검증 규칙
- Step 타입 감지: validate, db, event, call, cache, log, transform, condition, loop, return

### 미지원 (구현 대상)
- db step: 작업 타입 검증 (insert, select, update, delete)
- event step: 이벤트 이름/페이로드 검증
- call step: URL, 헤더, 타임아웃 형식 검증
- cache step: 키 템플릿, TTL 형식 검증
- condition step: 조건식 파싱
- loop step: 반복 구조 파싱

## 작업 내용

### 1. db step 검증 강화
```yaml
steps:
  - db:
      action: insert|select|update|delete
      table: users
      data: {email: input.email}
      where: {id: input.id}
```
- `action` 필드가 insert/select/update/delete 중 하나인지 검증
- action에 따른 필수 필드 검증 (insert: data, select/update/delete: where)

### 2. event step 검증
```yaml
steps:
  - event:
      name: user.created
      payload: {userId: result.id}
```
- `name` 필드 필수 검증
- `payload` 형식 검증

### 3. call step 검증
```yaml
steps:
  - call:
      url: https://api.example.com/notify
      method: POST
      headers: {Authorization: "Bearer {token}"}
      timeout: 30s
```
- `url` 형식 검증 (http/https)
- `method` 유효값 검증
- `timeout` 형식 검증 (숫자+단위)

### 4. cache step 검증
```yaml
steps:
  - cache:
      action: get|set|delete
      key: "user:{userId}"
      ttl: 3600
```
- `action` 유효값 검증
- `key` 템플릿 형식 확인
- `ttl` 숫자 형식 검증

### 5. condition step 파싱
```yaml
steps:
  - condition:
      if: input.role == 'admin'
      then:
        - db: {...}
      else:
        - return: {error: "unauthorized"}
```
- `if` 조건식 존재 검증
- `then` 필수 검증
- 중첩 step 재귀 파싱

### 6. loop step 파싱
```yaml
steps:
  - loop:
      each: input.items
      as: item
      do:
        - db: {...}
```
- `each`, `as`, `do` 필수 검증
- 중첩 step 재귀 파싱

## 구현 방법

fdl_service.go에 다음 함수 추가/수정:
- `validateDbStep(step map[string]interface{}) error`
- `validateEventStep(step map[string]interface{}) error`
- `validateCallStep(step map[string]interface{}) error`
- `validateCacheStep(step map[string]interface{}) error`
- `parseConditionStep(step map[string]interface{}) (*ConditionStep, error)`
- `parseLoopStep(step map[string]interface{}) (*LoopStep, error)`

## 검증
```bash
cd cli && go test -v ./test/fdl_service_test.go -run TestLogicLayer
```

## 완료 기준
- [ ] db step action 타입 검증 구현
- [ ] event step name/payload 검증 구현
- [ ] call step url/timeout 검증 구현
- [ ] cache step action/ttl 검증 구현
- [ ] condition step 재귀 파싱 구현
- [ ] loop step 재귀 파싱 구현
- [ ] 테스트 케이스 추가
