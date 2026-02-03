# fdl.go Report

**파일 경로**: `internal/cmd/fdl.go`
**라인 수**: 575줄

## 요약
FDL(Feature Definition Language) 관리 명령어. 전체 워크플로우 지원.

## 서브명령어
| 명령어 | 설명 |
|--------|------|
| `fdl create <name>` | FDL 템플릿 생성 |
| `fdl register <file>` | FDL 파일을 Feature로 등록 |
| `fdl validate <id>` | FDL 검증 |
| `fdl show <id>` | FDL 내용 표시 |
| `fdl skeleton <id>` | 스켈레톤 파일 생성 |
| `fdl tasks <id>` | FDL에서 Task 생성 |
| `fdl verify <id>` | 구현이 FDL과 일치하는지 검증 |
| `fdl diff <id>` | FDL과 실제 코드 차이 표시 |

## fdl create
- `features/` 디렉토리에 `<name>.fdl.yaml` 생성
- 기본 템플릿 포함 (models, service, api, ui)

## fdl register
1. FDL 파일 파싱
2. FDL 검증
3. Feature 생성
4. FDL 내용 및 해시 저장

## fdl skeleton 플래그
```
--dry-run    # 생성할 파일 목록만 출력
--force      # 기존 파일 덮어쓰기
```

## fdl tasks
1. FDL 파싱
2. Task 매핑 추출
3. Phase 생성/조회
4. Task 생성 (feature_id, target_file 설정)
5. 의존성 Edge 자동 생성

## fdl verify 결과
```json
{
  "valid": false,
  "errors": [...],
  "warnings": [...],
  "functions_missing": [...],
  "files_missing": [...]
}
```

## fdl diff 결과
```json
{
  "differences": [
    {"file_path": "...", "layer": "service", "changes": [...]}
  ],
  "total_changes": 3
}
```
