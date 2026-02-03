# required.go Report

**파일 경로**: `internal/cmd/required.go`
**라인 수**: ~60줄

## 요약
필수 입력 확인 명령어.

## 명령어
```bash
clari required
```

## 출력
```json
{
  "success": true,
  "ready": false,
  "missing_required": [
    {"field": "context.project_name", "prompt": "프로젝트 이름을 입력하세요"},
    {"field": "tech.backend", "prompt": "백엔드 기술을 선택하세요", "options": ["go", "node"]}
  ]
}
```

## 확인 항목
- context: project_name, description
- tech: backend, frontend, database
- design: architecture, auth_method, api_style
