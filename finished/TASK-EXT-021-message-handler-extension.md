# TASK-EXT-021: 메시지 핸들러 확장

## 목표
Extension Host에서 Context/Tech/Design 저장 메시지 처리

## 작업 내용

### CltEditorProvider.ts 수정
```typescript
// vscode-extension/src/CltEditorProvider.ts

// onDidReceiveMessage 핸들러에 추가
case 'saveContext':
  this.saveContext(message.data, panel);
  break;
case 'saveTech':
  this.saveTech(message.data, panel);
  break;
case 'saveDesign':
  this.saveDesign(message.data, panel);
  break;

// 메서드 추가
private saveContext(data: Record<string, any>, panel: vscode.WebviewPanel) {
  try {
    const db = this.getDatabase();
    // context 테이블에 저장 (upsert)
    const stmt = db.prepare(`
      INSERT OR REPLACE INTO context (project_id, key, value)
      VALUES (?, ?, ?)
    `);

    const projectId = this.getCurrentProjectId();
    for (const [key, value] of Object.entries(data)) {
      stmt.run(projectId, key, JSON.stringify(value));
    }

    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'context',
      success: true
    });

    // 전체 데이터 재동기화
    this.syncData(panel);
  } catch (error) {
    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'context',
      success: false,
      error: error.message
    });
  }
}

private saveTech(data: Record<string, any>, panel: vscode.WebviewPanel) {
  try {
    const db = this.getDatabase();
    const stmt = db.prepare(`
      INSERT OR REPLACE INTO tech (project_id, key, value)
      VALUES (?, ?, ?)
    `);

    const projectId = this.getCurrentProjectId();
    // 기존 데이터 삭제 후 새로 삽입
    db.exec(`DELETE FROM tech WHERE project_id = '${projectId}'`);
    for (const [key, value] of Object.entries(data)) {
      stmt.run(projectId, key, value);
    }

    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'tech',
      success: true
    });

    this.syncData(panel);
  } catch (error) {
    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'tech',
      success: false,
      error: error.message
    });
  }
}

private saveDesign(data: Record<string, any>, panel: vscode.WebviewPanel) {
  try {
    const db = this.getDatabase();
    const stmt = db.prepare(`
      INSERT OR REPLACE INTO design (project_id, key, value)
      VALUES (?, ?, ?)
    `);

    const projectId = this.getCurrentProjectId();
    db.exec(`DELETE FROM design WHERE project_id = '${projectId}'`);
    for (const [key, value] of Object.entries(data)) {
      stmt.run(projectId, key, value);
    }

    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'design',
      success: true
    });

    this.syncData(panel);
  } catch (error) {
    panel.webview.postMessage({
      type: 'settingSaveResult',
      section: 'design',
      success: false,
      error: error.message
    });
  }
}
```

### useSync.ts 수정
```typescript
// vscode-extension/webview-ui/src/hooks/useSync.ts

// settingSaveResult 핸들링 추가
case 'settingSaveResult':
  if (!message.success) {
    setError(`Failed to save ${message.section}: ${message.error}`);
  }
  break;
```

## 파일
- `vscode-extension/src/CltEditorProvider.ts` (수정)
- `vscode-extension/webview-ui/src/hooks/useSync.ts` (수정)

## 의존성
- TASK-EXT-020 (Types 확장)

## 완료 조건
- saveContext 메시지 처리 및 DB 저장
- saveTech 메시지 처리 및 DB 저장
- saveDesign 메시지 처리 및 DB 저장
- 저장 결과 피드백 (성공/실패)
- 저장 후 자동 데이터 재동기화
