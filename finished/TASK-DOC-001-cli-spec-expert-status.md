# TASK-DOC-001: CLI 스펙 Expert 명령어 상태 업데이트

## 개요
Expert 명령어가 구현되었으므로 CLI 스펙 문서 업데이트

## 배경
- **보고서**: reports/2026-02-03-total-report.md Issue #2, #3
- **대상 파일**:
  - specs/CLI/01-Overview.md
  - specs/CLI/11-Expert.md

## 작업 내용

### 1. CLI/01-Overview.md 수정

현재 (예상):
```markdown
| Expert | 7 | 미구현 |
```

수정:
```markdown
| Expert | 7 | 구현됨 |
```

### 2. CLI/11-Expert.md에 --feature 플래그 문서화

#### expert assign
```bash
clari expert assign <expert_id> <task_id>
clari expert assign <expert_id> --feature <feature_id>
```

**옵션:**
- `--feature <id>`: Feature 레벨에 전문가 할당

#### expert unassign
```bash
clari expert unassign <expert_id> <task_id>
clari expert unassign <expert_id> --feature <feature_id>
```

**옵션:**
- `--feature <id>`: Feature 레벨에서 전문가 해제

## 검증
- 문서 일관성 확인
- 마크다운 형식 검증

## 완료 기준
- [ ] Overview에서 Expert 상태를 "구현됨"으로 수정
- [ ] Expert.md에 --feature 플래그 문서화
