# TASK-EXT-020: Types 확장

## 목표
Project 탭 관련 메시지 타입 추가

## 작업 내용

### types.ts 수정
```typescript
// vscode-extension/webview-ui/src/types.ts

// 기존 타입 유지...

// MessageFromWebview 확장
export type MessageFromWebview =
  | { type: 'save'; table: string; id: number; data: any; version: number }
  | { type: 'refresh' }
  | { type: 'addEdge'; fromId: number; toId: number }
  | { type: 'removeEdge'; fromId: number; toId: number }
  | { type: 'createTask'; featureId: number; title: string; content: string }
  | { type: 'createFeature'; name: string; description: string }
  // 신규 추가
  | { type: 'saveContext'; data: Record<string, any> }
  | { type: 'saveTech'; data: Record<string, any> }
  | { type: 'saveDesign'; data: Record<string, any> };

// MessageToWebview 확장
export type MessageToWebview =
  | { type: 'sync'; data: ProjectData; timestamp: number }
  | { type: 'error'; message: string }
  | { type: 'saveResult'; success: boolean; table?: string; id?: number; error?: string }
  | { type: 'conflict'; table: string; id: number }
  | { type: 'edgeResult'; success: boolean; action?: string; error?: string }
  | { type: 'createResult'; success: boolean; table?: string; id?: number; error?: string }
  // 신규 추가
  | { type: 'settingSaveResult'; section: 'context' | 'tech' | 'design'; success: boolean; error?: string };
```

## 파일
- `vscode-extension/webview-ui/src/types.ts` (수정)

## 의존성
- 없음

## 완료 조건
- MessageFromWebview에 saveContext, saveTech, saveDesign 추가
- MessageToWebview에 settingSaveResult 추가
