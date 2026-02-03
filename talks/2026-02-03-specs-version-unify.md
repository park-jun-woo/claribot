# 2026-02-03 Specs 버전 통일 및 프로젝트 구조 정리

## 세션 개요

specs/ 문서 검수 및 버전/형식 통일, 프로젝트 구조 재정리 세션.

---

## 1. FDL 파일명 변경

**변경 내용:**
- `02a-DataLayer.md` → `02-A-DataLayer.md`
- `02b-LogicLayer.md` → `02-B-LogicLayer.md`
- `02c-InterfaceLayer.md` → `02-C-InterfaceLayer.md`
- `02d-PresentationLayer.md` → `02-D-PresentationLayer.md`

**참조 업데이트:** FDL/01-Overview.md, FDL/02-Schema.md

---

## 2. Specs 문서 검수 결과

### 발견된 문제

| 심각도 | 문제 | 영향 파일 |
|--------|------|----------|
| CRITICAL | CLI/04-14 버전 v0.0.3 유지 | 11개 |
| CRITICAL | FDL/02-A~D, TTY/02-07 버전 v0.0.1 유지 | 13개 |
| HIGH | 버전 헤더 형식 3가지 혼재 | 전체 |
| HIGH | 초단문 문서 (50줄 미만) | 8개 |
| MEDIUM | 푸터 형식 불일치 (날짜 포함 여부) | 다수 |

### 정상 확인 항목
- 모든 링크 참조 정상
- Markdown 포맷 일관성
- 문서 간 중복은 의도적 (요약 vs 상세)

---

## 3. 버전 및 형식 통일 작업

### 수정된 파일 수

| 폴더 | 파일 수 |
|------|--------|
| CLI/ | 14개 |
| FDL/ | 8개 |
| TTY/ | 7개 |
| VSCode/ | 15개 |
| **총합** | **44개** |

### 통일된 형식

**헤더:**
```markdown
> **현재 버전**: v0.0.4 ([변경이력](../HISTORY.md))

---
```

**푸터:**
```markdown
*Claritask Commands Reference v0.0.4*
```
(날짜 제거)

---

## 4. 프로젝트 구조 재정리

### CLI 소스코드 이동

**Before:**
```
claritask/
├── cmd/
├── internal/
├── test/
├── go.mod
├── go.sum
├── Makefile
└── scripts/
```

**After:**
```
claritask/
├── cli/
│   ├── cmd/
│   ├── internal/
│   ├── test/
│   ├── scripts/
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
├── vscode-extension/
└── specs/
```

### CLAUDE.md 업데이트
- Project Structure 경로 업데이트
- Report Process 경로 업데이트

### Makefile 정리
- VSCode Extension 빌드 타겟 제거 (`ext-build`, `ext-package`, `ext-install`)

---

## 5. HISTORY.md 업데이트

v0.0.4 변경 내역:
- CLI/* (14개): 버전 v0.0.4 통일, 헤더/푸터 형식 통일
- FDL/* (8개): 버전 v0.0.4 통일, 파일명 변경 (02a→02-A)
- TTY/* (7개): 버전 v0.0.4 통일, 헤더/푸터 형식 통일
- VSCode/* (15개): 버전 v0.0.4 통일, 헤더/푸터 형식 통일

---

## 보류된 작업

- 3순위: 초단문 문서 내용 보강 (8개 파일)

---

*2026-02-03 대화 저장*
