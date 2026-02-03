# TASK-DEV-067: Commands 문서 업데이트

## 목표
specs/Commands.md에 새로운 init 명령어 문서 업데이트.

## 파일
`specs/Commands.md`

## 변경 내용

### 1. init 명령어 섹션 업데이트

```markdown
## clari init

프로젝트 초기화. LLM과 협업하여 프로젝트 설정 완성.

### 사용법

```bash
clari init <project-id> [options]
```

### 인자

| 인자 | 설명 |
|------|------|
| project-id | 프로젝트 ID (필수, 영문소문자/숫자/하이픈/언더스코어) |

### 옵션

| 옵션 | 단축 | 설명 |
|------|------|------|
| --name | -n | 프로젝트 이름 (기본값: project-id) |
| --description | -d | 프로젝트 설명 |
| --skip-analysis | | 컨텍스트 분석 건너뛰기 |
| --skip-specs | | Specs 생성 건너뛰기 |
| --non-interactive | | 비대화형 모드 (자동 승인) |
| --force | | 기존 DB 덮어쓰기 |
| --resume | | 중단된 초기화 재개 |

### 프로세스

1. **Phase 1**: DB 초기화 (.claritask/db 생성)
2. **Phase 2**: 프로젝트 파일 분석 (claude --print)
3. **Phase 3**: tech/design 승인 (대화형)
4. **Phase 4**: Specs 초안 생성 (claude --print)
5. **Phase 5**: 피드백 루프 (승인까지 반복)

### 예시

```bash
# 기본 사용
clari init my-api

# 옵션 지정
clari init my-api --name "My REST API" --description "사용자 관리 API"

# 빠른 초기화
clari init my-api --skip-analysis --skip-specs

# 기존 프로젝트 재초기화
clari init my-api --force

# 중단된 초기화 재개
clari init --resume
```

### 출력

각 Phase 완료 시 JSON 출력.

```json
{
  "success": true,
  "phase": "complete",
  "project_id": "my-api",
  "db_path": ".claritask/db",
  "specs_path": "specs/my-api.md"
}
```
```

## 완료 조건
- [ ] Commands.md의 init 섹션 업데이트
- [ ] 모든 옵션 문서화
- [ ] 예시 추가
