# 2026-02-04 Message Report Markdown 렌더링 구현

## 주요 작업 내용

### 1. TTY 완료 파일 버그 수정
**문제:** `clari message send` 실행 후 "signal: killed" 에러 발생, 메시지가 failed 상태로 저장됨

**원인:**
- `watchCompleteFileWithSignal`이 `.claritask/complete` 파일 감지 시 즉시 삭제
- `RunMessageAnalysisWithTTY`가 파일 내용을 읽기 전에 파일이 사라짐
- `cmd.Wait()`가 프로세스 kill로 인한 에러 반환

**해결:**
1. `tty_service.go`: `watchCompleteFileWithSignal`에서 파일 삭제 제거 (caller가 읽은 후 삭제)
2. `RunWithTTYHandoverEx`: 완료 채널 시그널 확인 후 성공으로 처리

```go
// 완료 채널로 kill이 의도적인지 판단
completedChan := make(chan bool, 1)
go watchCompleteFileWithSignal(completeFile, cmd, completedChan)

waitErr := cmd.Wait()

select {
case <-completedChan:
    return nil  // 완료 기반 종료 = 성공
default:
    return waitErr  // 다른 이유
}
```

### 2. 완료 파일명 변경
**변경 전:** `.claritask/complete`
**변경 후:**
- Message: `.claritask/complete-message-<id>.md`
- Feature: `.claritask/complete-feature-<id>.md`

**수정 파일:**
- `cli/internal/service/message_service.go`
- `cli/internal/service/tty_service.go`

### 3. Message 프롬프트 수정
- "modification request" → "request" (범용적 용도로 변경)
- `MessageAnalysisSystemPrompt()` 업데이트

### 4. VSCode Extension Report Markdown 렌더링

**문제:** Report 섹션에 원문 텍스트 그대로 출력

**해결:**
1. `react-markdown` → `marked` 라이브러리로 교체 (더 가볍고 안정적)
2. `ReportSection` 컴포넌트 추가

```tsx
function ReportSection({ content }: { content: string }) {
  const html = useMemo(() => {
    return marked.parse(content) as string;
  }, [content]);

  return (
    <div>
      <h3 className="text-sm font-semibold mb-1 text-green-400">Report</h3>
      <div
        className="markdown-content"
        dangerouslySetInnerHTML={{ __html: html }}
      />
    </div>
  );
}
```

3. `index.css`에 `.markdown-content` 스타일 추가
   - 헤더 (h1-h4): 녹색, 굵게
   - 리스트, 코드 블록, 인용문, 테이블 스타일링

**번들 크기:** 311KB → 234KB (react-markdown 제거로 감소)

## 수정된 파일 목록

### CLI
- `cli/internal/service/tty_service.go`
  - `watchCompleteFileWithSignal`: 파일 삭제 제거
  - `RunWithTTYHandoverEx`: 완료 채널 처리 추가
  - `FDLGenerationSystemPrompt`: 완료 파일 프롬프트 업데이트
  - `RunFDLGenerationWithTTY`: 새 파일명 사용
  - `BuildFDLPrompt`: completeFile 파라미터 추가

- `cli/internal/service/message_service.go`
  - `MessageAnalysisSystemPrompt`: "modification" 제거
  - `RunMessageAnalysisWithTTY`: 새 파일명 사용
  - `BuildMessageAnalysisPrompt`: completeFile 파라미터 추가

### VSCode Extension
- `webview-ui/src/components/MessagesPanel.tsx`
  - `marked` 라이브러리 사용
  - `ReportSection` 컴포넌트 추가

- `webview-ui/src/index.css`
  - `.markdown-content` 스타일 추가

- `webview-ui/package.json`
  - `marked` 의존성 추가

## 버전
- VSCode Extension: v0.0.15
