# VSCode Extension Project 탭 구현

> **날짜**: 2026-02-03
> **주제**: Project 탭 UI 설계 및 구현

---

## 1. 프로젝트 초기화

claritask 프로젝트에 대한 정보를 CLI로 입력:

```bash
clari init claritask --skip-analysis --skip-specs --force

clari project set '{
  "context": {
    "project_name": "Claritask",
    "description": "AI 에이전트를 위한 구조화된 작업 관리 CLI 도구...",
    "target_users": "AI 에이전트 (Claude Code 등), 개발자"
  },
  "tech": {
    "backend": "Go 1.21+",
    "frontend": "CLI (Cobra)",
    "database": "SQLite (mattn/go-sqlite3)"
  },
  "design": {
    "architecture": "CLI Application",
    "auth_method": "None (로컬 도구)",
    "api_style": "JSON 입출력"
  }
}'

clari required  # → ready: true
```

---

## 2. VscodeGUI.md 설계 업데이트

### 2.1 Project 탭 설계 (v0.0.2)

기존 `Features`, `Tasks` 2탭에서 **`Project`, `Features`, `Tasks`** 3탭 구조로 확장.

**Project 탭 레이아웃:**
```
┌─ Project Info ─────────────────────────────────────┐
│  ID / Name / Status / Created (읽기 전용)          │
└────────────────────────────────────────────────────┘
┌─ Context ──────────────────────────────────[Edit]──┐
│  project_name* / description* / target_users / ... │
└────────────────────────────────────────────────────┘
┌─ Tech Stack ───────────────────────────────[Edit]──┐
│  backend* / frontend* / database* / ...            │
└────────────────────────────────────────────────────┘
┌─ Design Decisions ─────────────────────────[Edit]──┐
│  architecture* / auth_method* / api_style* / ...   │
└────────────────────────────────────────────────────┘
```

### 2.2 토론: 대화 모드 탭 추가 여부

**논의 내용:**
- AI 에이전트 대화 인터페이스 → Claude Code 확장 활용 (중복 방지)
- 실행 로그 뷰어 → 필요함 (Project 탭에 통합)

**결정:**
- 별도 탭 추가 안 함
- Project 탭에 Execution Status 섹션 추가
- Start 버튼 제외 (CLI 또는 Claude Code로 직접 실행)

### 2.3 Execution Status 섹션 추가 (v0.0.3)

```
┌─ Execution Status ────────────────────────────────┐
│  Progress: ████████░░░░░░░░░░░░  42% (17/40)      │
│  Status:   ● Running                               │
│                                                    │
│  Recent Tasks:                                     │
│  ✓ #15 user_table_sql    "테이블 생성 완료"        │
│  ✓ #16 user_model        "모델 구현 완료"          │
│  ● #17 auth_service      "JWT 구현 중..."          │
│  ○ #18 login_endpoint    (pending)                 │
└────────────────────────────────────────────────────┘
```

**상태 표시:**
| 상태 | 색상 | 조건 |
|------|------|------|
| Running | 파란색 | doing Task 있음 |
| Idle | 회색 | pending만 있음 |
| Completed | 녹색 | 모든 Task 완료 |
| Has Failures | 빨간색 | failed Task 있음 |

---

## 3. 개발 실행

### 3.1 Task 목록 (TASK-EXT-010 ~ 021)

| Task | 설명 |
|------|------|
| EXT-010 | Project 탭 기본 구조 (App.tsx, ProjectPanel.tsx) |
| EXT-011 | SectionCard 공통 컴포넌트 |
| EXT-012 | ProjectInfo 컴포넌트 |
| EXT-013 | EditableField, AddFieldButton 공통 컴포넌트 |
| EXT-014 | ContextSection 편집 컴포넌트 |
| EXT-015 | TechSection 편집 컴포넌트 |
| EXT-016 | DesignSection 편집 컴포넌트 |
| EXT-017 | ProgressBar 컴포넌트 |
| EXT-018 | RecentTaskList 컴포넌트 |
| EXT-019 | ExecutionStatus 메인 컴포넌트 |
| EXT-020 | types.ts 메시지 타입 확장 |
| EXT-021 | CltEditorProvider 메시지 핸들러 확장 |

### 3.2 생성/수정된 파일

**신규 생성 (10개):**
- `components/ProjectPanel.tsx`
- `components/SectionCard.tsx`
- `components/ProjectInfo.tsx`
- `components/EditableField.tsx`
- `components/ContextSection.tsx`
- `components/TechSection.tsx`
- `components/DesignSection.tsx`
- `components/ProgressBar.tsx`
- `components/RecentTaskList.tsx`
- `components/ExecutionStatus.tsx`

**수정 (6개):**
- `App.tsx` - Project 탭 추가, view 타입 확장
- `types.ts` - saveContext/Tech/Design, settingSaveResult 메시지 추가
- `vscode.ts` - vscode 객체 export 추가
- `useSync.ts` - settingSaveResult 핸들러 추가
- `src/CltEditorProvider.ts` - handleSaveContext/Tech/Design 메서드
- `src/database.ts` - saveContext/saveTech/saveDesign 메서드

### 3.3 빌드 결과

```bash
# webview-ui
npm run build  # ✓ 성공 (177.90 kB)

# extension
npm run compile  # ✓ 성공 (94.3 kB)
```

---

## 4. 메시지 프로토콜

### Webview → Extension (신규)
```typescript
{ type: 'saveContext', data: Record<string, any> }
{ type: 'saveTech', data: Record<string, any> }
{ type: 'saveDesign', data: Record<string, any> }
```

### Extension → Webview (신규)
```typescript
{ type: 'settingSaveResult', section: 'context' | 'tech' | 'design', success: boolean, error?: string }
```

---

## 5. 다음 단계

- VSCode에서 실제 테스트
- 필수 필드 검증 UI 개선
- 사용자 정의 필드 추가/삭제 UX 개선
