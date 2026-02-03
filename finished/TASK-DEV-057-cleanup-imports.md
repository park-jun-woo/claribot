# TASK-DEV-057: Import 및 빌드 정리

## 목표
모든 파일에서 삭제된 파일 import 제거, 빌드 확인.

## 파일
- 모든 Go 파일

## 변경 내용

### 1. 삭제된 패키지 import 제거
- phase_service 관련 import
- phase 관련 상수/타입 참조

### 2. 사용하지 않는 import 정리
```bash
goimports -w .
```

### 3. 빌드 확인
```bash
go build ./...
```

### 4. 테스트 실행
```bash
go test ./...
```

## 완료 조건
- [ ] `go build ./...` 성공
- [ ] `go test ./...` 성공
- [ ] 사용하지 않는 import 없음
