# main.go Report

**파일 경로**: `cmd/claritask/main.go`
**라인 수**: 13줄

## 요약
CLI 애플리케이션 진입점. Cobra 명령어 시스템을 실행하는 최소 구현.

## 구조
```go
func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## 특징
- internal/cmd 패키지의 Execute() 함수 호출
- 에러 발생 시 exit code 1 반환
- 단순하고 명확한 진입점 구현

## 의존성
- `parkjunwoo.com/claritask/internal/cmd`
