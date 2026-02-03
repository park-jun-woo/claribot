# TASK-EXT-048: Message CLI 호출 단순화

## 목표
handleSendMessageCLI를 handleCreateFeature와 동일한 패턴으로 단순화

## 변경 내용

### CltEditorProvider.ts

#### handleSendMessageCLI 수정
- `clari project set` 제거 (불필요)
- `clari message send --content '...' [--feature N]` 형식으로 변경

```typescript
private async handleSendMessageCLI(
  message: { content: string; featureId?: number },
  database: Database,
  webview: vscode.Webview
): Promise<void> {
  try {
    // Escape content for shell
    const escapedContent = message.content.replace(/'/g, "'\\''");

    // Build command (like handleCreateFeature)
    let command = `~/bin/clari message send --content '${escapedContent}'`;
    if (message.featureId) {
      command += ` --feature ${message.featureId}`;
    }

    // Build full command with cd for WSL
    const isWindows = process.platform === 'win32';
    const workspacePath = vscode.workspace.workspaceFolders?.[0].uri.fsPath;
    let fullCommand = command;

    if (isWindows && workspacePath) {
      const wslPath = windowsToWslPath(workspacePath);
      fullCommand = `cd '${wslPath}' && ${command}`;
    }

    // Use TTY Session Manager (like handleCreateFeature)
    if (ttySessionManager) {
      const sessionId = `message-send-${Date.now()}`;
      await ttySessionManager.startSession(sessionId, fullCommand);
    } else {
      const terminal = vscode.window.createTerminal({
        name: 'Claritask - Send Message',
        shellPath: isWindows ? 'wsl.exe' : undefined,
      });
      terminal.show();
      terminal.sendText(fullCommand);
    }

    webview.postMessage({
      type: 'cliStarted',
      command: 'message.send',
      message: 'Claude Code will analyze your request...',
    });
  } catch (err) {
    // ...
  }
}
```

## 완료 조건
- [ ] project set 제거
- [ ] --content 플래그 사용
- [ ] handleCreateFeature와 동일한 패턴
