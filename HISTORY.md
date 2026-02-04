# Claribot Version History

## v0.2.1 (Draft) - 2026-02-04

현재 개발 중인 버전.

### 변경사항
- CLAUDE.md 정리 및 프로젝트 개요 추가
- DB 파일명 `db.clt` 유지 (v0.1과 동일)

---

## v0.2.0 (Draft) - 전역 서비스 아키텍처

v0.1 대비 전면 재설계.

### 주요 변경점
- **아키텍처**: 프로젝트 로컬 → 전역 서비스 (daemon)
- **텔레그램 봇**: 원격 제어 통합
- **clari CLI**: DB 직접 접근 → 서비스 클라이언트 (HTTP)
- **통신**: localhost HTTP (127.0.0.1:9847)

### DB 분리
- 전역 DB (`~/.claribot/db`): 프로젝트 목록, 경로 매핑
- 로컬 DB (`프로젝트/.claribot/db`): tasks, memos, features 등

### 신규 명령어
- `clari serve`: 서비스 데몬 실행

---

## v0.1.x - 로컬 CLI

초기 버전. 프로젝트별 독립 실행 방식.

### 특징
- 프로젝트 폴더 내 `.claribot/db.clt`에 로컬 DB
- CLI가 DB에 직접 접근
- VSCode Extension: `.clt` 파일 더블클릭으로 Custom Editor 실행

### 한계
- 프로젝트마다 개별 실행 필요
- 원격 제어 불가
- 중앙 관리 기능 없음

---

*최초 작성: 2026-02-04*
