# Test Files Summary Report

**디렉토리**: `test/`
**파일 수**: 10개

## 테스트 파일 목록

| 파일 | 테스트 대상 |
|------|------------|
| `models_test.go` | model 구조체 JSON 직렬화 |
| `db_test.go` | DB 연결, 마이그레이션 |
| `project_service_test.go` | Project/Context/Tech/Design CRUD |
| `phase_service_test.go` | Phase 생명주기 |
| `task_service_test.go` | Task CRUD, Pop, 상태 전이 |
| `memo_service_test.go` | Memo CRUD, 스코프별 조회 |
| `state_service_test.go` | State key-value 관리 |
| `feature_service_test.go` | Feature CRUD, Edge, FDL 관리 |
| `edge_service_test.go` | 의존성 Edge, 사이클 감지, 위상 정렬 |
| `fdl_service_test.go` | FDL 파싱, 검증, Task 매핑 |

## 테스트 패턴
```go
func TestXxx(t *testing.T) {
    // Setup: 임시 DB 생성
    db := setupTestDB(t)
    defer db.Close()

    // Execute: 함수 호출
    result, err := service.Xxx(db, ...)

    // Assert: 결과 검증
    if err != nil { t.Fatalf(...) }
    if result != expected { t.Errorf(...) }
}
```

## 테스트 실행
```bash
go test ./test/... -v
```
