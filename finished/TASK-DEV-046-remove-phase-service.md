# TASK-DEV-046: Phase Service 삭제

## 목표
phase_service.go 파일 삭제.

## 파일
`internal/service/phase_service.go`

## 변경 내용
- 파일 전체 삭제

## 영향 범위
- phase.go (cmd)에서 import 제거
- 다른 서비스에서 phase_service 참조 확인 및 제거

## 완료 조건
- [ ] phase_service.go 파일 삭제
- [ ] 관련 import 정리
