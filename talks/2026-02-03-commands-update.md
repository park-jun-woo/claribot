# Commands.md 업데이트

날짜: 2026-02-03

## 1. 대화 검토

### 이전 대화 (architecture-redesign.md) 확인
- 7가지 아키텍처 재설계 내용 검토
- 제어 역전, Phase→Feature, 트리→그래프, Edge 추론 등

### Claritask.md 반영 여부 확인
- **결과**: 100% 반영됨
- 스펙 문서가 대화 내용을 충실히 담고 있음

---

## 2. Commands.md 변경 사항

### 버전
- v1.0 → v2.0

### 명령어 수 변경
- 23개 → 32개

### 주요 변경 내용

| 변경 항목 | 상세 |
|----------|------|
| Phase → Feature | 전체 섹션 전환 |
| Edge 섹션 | 신규 추가 (4개 명령어) |
| Task | push→add, list 추가, dependencies 응답 |
| Project 실행 | start 옵션, stop, status 추가 |
| Memo | scope: phase → feature |

### 신규 명령어

**Project 실행**:
- `clari plan features` - Feature 목록 산출 (LLM)
- `clari project start --feature N` - 특정 Feature만 실행
- `clari project start --dry-run` - 실행 없이 순서 출력
- `clari project stop` - 실행 중단
- `clari project status` - 실행 상태 조회

**Feature**:
- `clari feature list`
- `clari feature add '<json>'`
- `clari feature <id> spec` - LLM과 spec 대화
- `clari feature <id> tasks` - Task 생성 (LLM)
- `clari feature <id> start`

**Edge**:
- `clari edge add --from <id> --to <id>` - Task Edge
- `clari edge add --feature --from <id> --to <id>` - Feature Edge
- `clari edge list`
- `clari edge infer --feature <id>` - Task Edge 추론
- `clari edge infer --project` - Feature Edge 추론

**Task**:
- `clari task list [--feature N]` - 신규
- `clari task add` (push에서 변경)

### task pop 응답 변경

```json
{
  "task": { ... },
  "dependencies": [
    {"id": 1, "title": "...", "result": "..."}
  ],
  "manifest": {
    "feature": { ... },
    ...
  }
}
```

- `dependencies`: 의존 Task들의 result 포함 (핵심!)
- `manifest.feature`: phase 대신 feature 정보

### 에러 코드 추가
- `CIRCULAR_DEPENDENCY`: 순환 의존성 감지
- `DEPENDENCY_NOT_RESOLVED`: 의존성 미해결

---

## 3. 다음 단계

스펙 문서 업데이트 완료. 코드 구현 필요:
1. Phase 관련 코드 → Feature로 변경
2. Edge 테이블 및 서비스 구현
3. task pop에 dependencies 추가
4. Topological Sort 구현
