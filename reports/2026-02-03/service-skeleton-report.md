# skeleton_service.go Report

**파일 경로**: `internal/service/skeleton_service.go`
**라인 수**: 422줄

## 요약
스켈레톤 파일 관리, Python 제너레이터 실행, TODO 위치 추출.

## Skeleton CRUD
| 함수 | 설명 |
|------|------|
| `CreateSkeleton(db, featureID, path, layer)` | 스켈레톤 레코드 생성 |
| `GetSkeleton(db, id)` | 스켈레톤 조회 |
| `ListSkeletonsByFeature(db, featureID)` | Feature별 스켈레톤 목록 |
| `DeleteSkeleton(db, id)` | 스켈레톤 삭제 |
| `DeleteSkeletonsByFeature(db, featureID)` | Feature 스켈레톤 전체 삭제 |

## 파일 체크섬
```go
func CalculateFileChecksum(filePath) (string, error)  // SHA256
func UpdateSkeletonChecksum(db, id, checksum) error
func HasSkeletonChanged(db, id) (bool, error)
```

## 파일 내용 읽기
```go
func ReadSkeletonContent(filePath) (string, error)
func GetSkeletonAtLine(filePath, line, contextLines) (string, error)
```

## TODO 위치 추출
```go
type TODOLocation struct {
    Line     int
    Function string
    Content  string
}

func ExtractTODOLocations(filePath) ([]TODOLocation, error)
```
지원 패턴: `# TODO`, `// TODO`, `/* TODO`, `* TODO`

## 함수명 추출
```go
func GetFunctionAtLine(filePath, targetLine) (string, error)
```
지원 언어: Python (`def`), Go (`func`), JavaScript (`function`, `const/let`)

## Python 제너레이터 실행
```go
type SkeletonGeneratorResult struct {
    GeneratedFiles []GeneratedFile
    Errors         []string
}

func RunSkeletonGenerator(fdlPath, outputDir, backend, frontend, force) (*SkeletonGeneratorResult, error)
func RunSkeletonGeneratorDryRun(fdlPath, outputDir, backend, frontend) (*SkeletonGeneratorResult, error)
func RunSkeletonGeneratorForFeature(db, featureID, outputDir, force) (*SkeletonGeneratorResult, error)
```

**명령어**:
```bash
python3 scripts/skeleton_generator.py \
  --fdl <fdl_path> \
  --output-dir <dir> \
  --backend <backend> \
  --frontend <frontend> \
  [--force] [--dry-run]
```

## Task Pop용 스켈레톤 정보
```go
func GetSkeletonInfo(db, skeletonID, targetLine) (*SkeletonInfo, error)
```
- 전체 파일 또는 특정 라인 주변 컨텍스트 반환
- 2000자 초과 시 truncate
