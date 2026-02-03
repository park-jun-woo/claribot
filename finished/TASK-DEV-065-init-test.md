# TASK-DEV-065: Init Service 테스트

## 목표
Init Service 단위 테스트 작성.

## 파일
`test/init_service_test.go`

## 테스트 케이스

### 1. TestInitPhase1_DBInit

- 임시 디렉토리에서 Phase 1 실행
- .claritask/db 생성 검증
- projects 테이블에 레코드 검증
- state 테이블에 init_phase=1 검증

### 2. TestInitPhase1_DBInit_AlreadyExists

- 이미 .claritask/db 존재
- Force=false → 에러
- Force=true → 덮어쓰기

### 3. TestSaveInitState_LoadInitState

- InitState 저장
- 로드 후 동일 데이터 검증

### 4. TestInitConfig_Validation

- 유효하지 않은 ProjectID
- 에러 반환 검증

### 5. TestInitPhase3_Approval_NonInteractive

- NonInteractive=true
- 자동 승인 검증

### 6. TestInitPhase5_Feedback_Approval

- 승인 시 파일 저장 검증
- specs/project-id.md 생성 검증

## 참고
- Phase 2, 4는 LLM 호출 필요하므로 mock 또는 통합 테스트
- 대화형 입력 테스트는 io.Reader mock 사용

## 완료 조건
- [ ] 6개 테스트 케이스 구현
- [ ] go test 통과
