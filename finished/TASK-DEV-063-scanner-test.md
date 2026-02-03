# TASK-DEV-063: Scanner Service 테스트

## 목표
Scanner Service 단위 테스트 작성.

## 파일
`test/scanner_service_test.go`

## 테스트 케이스

### 1. TestScanProjectFiles_GoProject

- 임시 디렉토리에 go.mod, README.md 생성
- ScanProjectFiles 호출
- 파일 타입 및 내용 검증

### 2. TestScanProjectFiles_NodeProject

- package.json, README.md 생성
- 결과 검증

### 3. TestScanProjectFiles_PythonProject

- requirements.txt, pyproject.toml 생성
- 결과 검증

### 4. TestScanProjectFiles_EmptyDirectory

- 빈 디렉토리
- 빈 결과 반환 검증

### 5. TestScanProjectFiles_DirectoryStructure

- src/, internal/, lib/ 디렉토리 생성
- Directories 필드 검증

### 6. TestScanProjectFiles_LargeFile

- 1000줄 이상 파일 생성
- 처음 500줄만 읽는지 검증

### 7. TestFormatScanResultForLLM

- ScanResult를 문자열로 포맷
- 포맷 검증

## 완료 조건
- [ ] 7개 테스트 케이스 구현
- [ ] go test 통과
