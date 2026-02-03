# TASK-DEV-066: Init 통합 테스트

## 목표
`clari init` 전체 프로세스 통합 테스트.

## 파일
`test/init_integration_test.go`

## 테스트 케이스

### 1. TestClariInit_SkipAll

```bash
clari init test-project --skip-analysis --skip-specs
```

- Phase 1만 실행
- DB 생성 및 프로젝트 레코드 확인

### 2. TestClariInit_NonInteractive

```bash
clari init test-project --non-interactive
```

- 모든 Phase 자동 진행
- (실제 claude 없으면 skip)

### 3. TestClariInit_Force

```bash
clari init test-project
clari init test-project --force
```

- 두 번째 실행 시 덮어쓰기

### 4. TestClariInit_Resume

- Phase 2에서 중단
- `clari init --resume`로 재개
- 상태 복원 검증

## 참고
- 실제 claude CLI 필요한 테스트는 `//go:build integration` 태그 사용
- CI에서는 skip

## 완료 조건
- [ ] 4개 통합 테스트 구현
- [ ] go test -tags=integration 통과 (claude 있을 때)
- [ ] go test 통과 (skip)
