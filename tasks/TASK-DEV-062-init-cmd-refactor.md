# TASK-DEV-062: Init Command 리팩토링

## 목표
새로운 Init Service를 사용하도록 init.go 명령어 수정.

## 파일
`internal/cmd/init.go`

## 구현 내용

### 1. 플래그 추가

```go
func init() {
    initCmd.Flags().StringP("name", "n", "", "프로젝트 이름")
    initCmd.Flags().StringP("description", "d", "", "프로젝트 설명")
    initCmd.Flags().Bool("skip-analysis", false, "컨텍스트 분석 건너뛰기")
    initCmd.Flags().Bool("skip-specs", false, "Specs 생성 건너뛰기")
    initCmd.Flags().Bool("non-interactive", false, "비대화형 모드")
    initCmd.Flags().Bool("force", false, "기존 DB 덮어쓰기")
    initCmd.Flags().Bool("resume", false, "중단된 초기화 재개")
}
```

### 2. runInit 수정

```go
func runInit(cmd *cobra.Command, args []string) error {
    // 플래그 파싱
    name, _ := cmd.Flags().GetString("name")
    description, _ := cmd.Flags().GetString("description")
    skipAnalysis, _ := cmd.Flags().GetBool("skip-analysis")
    skipSpecs, _ := cmd.Flags().GetBool("skip-specs")
    nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
    force, _ := cmd.Flags().GetBool("force")
    resume, _ := cmd.Flags().GetBool("resume")

    // --resume 처리
    if resume {
        return runInitResume(cmd, args)
    }

    // InitConfig 구성
    config := service.InitConfig{
        ProjectID:      args[0],
        Name:           name,
        Description:    description,
        SkipAnalysis:   skipAnalysis,
        SkipSpecs:      skipSpecs,
        NonInteractive: nonInteractive,
        Force:          force,
    }

    // Name 기본값 설정
    if config.Name == "" {
        config.Name = config.ProjectID
    }

    // RunInit 호출
    result, err := service.RunInit(config)
    if err != nil {
        outputError(err)
        return nil
    }

    // 결과 출력 (JSON 모드가 아닐 때는 이미 출력됨)
    if !result.Success {
        outputError(fmt.Errorf(result.Error))
        return nil
    }

    return nil
}
```

### 3. runInitResume 함수

```go
func runInitResume(cmd *cobra.Command, args []string) error {
    database, err := getDB()
    if err != nil {
        outputError(fmt.Errorf("open database: %w", err))
        return nil
    }
    defer database.Close()

    result, err := service.ResumeInit(database)
    if err != nil {
        outputError(err)
        return nil
    }

    if !result.Success {
        outputError(fmt.Errorf(result.Error))
    }

    return nil
}
```

### 4. 기존 코드 제거

- 기존 단순 초기화 로직 제거
- claudeTemplate 제거 (또는 별도 파일로 이동)

### 5. --force 동작

- .claritask/db가 이미 존재할 때:
  - --force 없으면 에러
  - --force 있으면 삭제 후 재생성

## 완료 조건
- [ ] 새 플래그 추가 (--name, --description, --skip-*, --non-interactive, --force, --resume)
- [ ] runInit에서 InitConfig 구성
- [ ] service.RunInit 호출
- [ ] runInitResume 함수 구현
- [ ] 기존 단순 로직 제거
