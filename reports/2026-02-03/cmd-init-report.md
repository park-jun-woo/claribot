# init.go Report

**파일 경로**: `internal/cmd/init.go`
**라인 수**: 124줄

## 요약
새 프로젝트 초기화 명령어 구현.

## 명령어
```bash
clari init <project-id> [description]
```

## 주요 기능
1. 프로젝트 ID 검증 (영문소문자, 숫자, 하이픈, 언더스코어)
2. 프로젝트 디렉토리 생성 (`{cwd}/{project-id}/`)
3. `.claritask/db` 폴더 및 DB 파일 생성
4. DB 마이그레이션 실행
5. 프로젝트 레코드 생성
6. 상태 초기화
7. `CLAUDE.md` 파일 생성

## 프로젝트 ID 검증
```go
func validateProjectID(id string) error {
    matched, _ := regexp.MatchString(`^[a-z0-9_-]+$`, id)
    // ...
}
```

## CLAUDE.md 템플릿
```markdown
# {project_name}
## Description
## Tech Stack
## Commands
```

## 출력 예시
```json
{
  "success": true,
  "project_id": "my-project",
  "path": "/path/to/my-project",
  "message": "Project initialized successfully"
}
```
