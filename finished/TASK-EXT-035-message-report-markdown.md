# TASK-EXT-035: Message Report Markdown Rendering

## 목표
Messages 탭의 상세 화면에서 Content 바로 아래에 Report 섹션을 추가하고,
`.claritask/complete` 파일에서 받은 보고서를 Markdown 렌더링하여 표시

## 변경 파일
- `vscode-extension/webview-ui/src/components/MessagesPanel.tsx`
- `vscode-extension/webview-ui/package.json` (react-markdown 의존성 추가)

## 구현 내용

### 1. react-markdown 설치
```bash
cd vscode-extension/webview-ui
npm install react-markdown
```

### 2. MessagesPanel.tsx 수정

#### MessageDetail 컴포넌트에서:
- "Content" 섹션 바로 아래에 "Report" 섹션 추가
- 조건: `message.status === 'completed' && message.response`
- 기존 "Response" 섹션을 "Report" 섹션으로 변경
- `react-markdown`으로 Markdown 렌더링
- 스타일: 녹색 배경, prose 스타일 적용

### 3. Report 섹션 UI
```tsx
{message.status === 'completed' && message.response && (
  <div>
    <h3 className="text-sm font-semibold mb-1 text-green-400">Report</h3>
    <div className="prose prose-invert prose-sm max-w-none p-3 bg-green-900/30 border border-green-800 rounded">
      <ReactMarkdown>{message.response}</ReactMarkdown>
    </div>
  </div>
)}
```

## 테스트
1. VSCode Extension에서 Message send
2. Claude Code가 처리 완료 후 .claritask/complete에 보고서 작성
3. Messages 탭에서 해당 메시지 선택
4. Report 섹션에 Markdown이 렌더링되어 표시되는지 확인
