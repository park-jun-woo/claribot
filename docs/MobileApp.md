# Claribot Mobile App

> **Status**: Planning

---

## 1. Overview

### 1.1 Purpose

A native mobile client that provides direct access to the Claribot service without relying on the Telegram bot. While the Telegram bot serves as a convenient quick-access interface, the mobile app offers a richer, purpose-built experience for managing projects, tasks, and Claude Code interactions on the go.

### 1.2 Relationship with Existing Interfaces

```
┌──────────────────────────────────────────────────────────┐
│                     claribot (daemon)                     │
│                                                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │ Telegram │  │   CLI    │  │  Web UI  │  │   TTY    │ │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Manager  │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘ │
│       │             │             │              │        │
│       └─────────────┼─────────────┼──────────────┘        │
│                     ▼             │                        │
│              ┌──────────┐   ┌────▼─────┐                  │
│              │    DB    │   │ RESTful  │                   │
│              │ (SQLite) │   │   API    │                   │
│              └──────────┘   │ + Auth   │                   │
│                             └──────────┘                  │
└──────────────────────────────────────────────────────────┘
       ▲              ▲              ▲              ▲
       │ Message       │ HTTP         │ HTTP          │ HTTP
  [Telegram]     [clari CLI]     [Web UI]     [Mobile App]
```

| Interface | Strengths | Limitations |
|-----------|-----------|-------------|
| **Telegram** | Always accessible, quick commands, group-shareable | Limited UI, no tree views, no dashboard |
| **CLI** | Scriptable, dev-friendly, fast | Terminal only, no visual feedback |
| **Web UI** | Full visual dashboard, rich interactions | Requires browser, no push notifications |
| **Mobile App** | Native push notifications, offline access, touch-optimized, always in pocket | Requires separate installation |

### 1.3 Key Differentiators of Mobile App

- **Push Notifications**: Real-time alerts for task completion, cycle progress, and errors
- **Offline Access**: Cached project/task data available without network
- **Native Experience**: Smooth animations, gesture navigation, system-level integrations
- **Always Available**: Quick access from home screen without opening a browser

---

## 2. Technology Stack

### 2.1 Framework Comparison

| Criteria | React Native | Flutter | Go Mobile |
|----------|-------------|---------|-----------|
| **Code Reuse with Web UI** | High (React + TS) | Low (Dart) | None |
| **Ecosystem** | Mature, large | Growing, Google-backed | Minimal |
| **Performance** | Good (JSI/Fabric) | Excellent (Skia) | Native-level |
| **Development Speed** | Fast (hot reload) | Fast (hot reload) | Slow |
| **Team Familiarity** | High (existing React/TS stack) | Low | Medium (Go backend) |

### 2.2 Recommendation: React Native

**React Native** is the recommended choice for the following reasons:

1. **Maximum Code Reuse**: The existing Web UI is built with React + TypeScript. The following can be directly shared:
   - `types/index.ts` - All TypeScript type definitions
   - `api/client.ts` - API client functions (with minor URL adjustments)
   - `hooks/useClaribot.ts` - TanStack Query hooks (identical caching/invalidation logic)
   - `hooks/useAuth.ts` - Authentication hooks
   - Business logic and state management patterns

2. **Shared Dependencies**: TanStack Query, React Router (via React Navigation mapping), react-markdown

3. **Developer Experience**: Same language (TypeScript), same paradigm (React components), same tooling (npm/yarn)

### 2.3 Tech Stack

| Category | Choice | Reason |
|----------|--------|--------|
| **Framework** | React Native + TypeScript | Code reuse with Web UI |
| **Navigation** | React Navigation v7 | Standard RN navigation, maps to React Router concepts |
| **UI Library** | React Native Paper or Tamagui | Material Design / cross-platform styling |
| **State Management** | TanStack Query | Same as Web UI, server state caching |
| **Push Notifications** | Firebase Cloud Messaging (FCM) | Cross-platform push support |
| **Offline Storage** | MMKV or AsyncStorage | Fast key-value cache for offline data |
| **Markdown** | react-native-markdown-display | Spec/Plan/Report rendering |
| **QR Code** | react-native-qrcode-svg | TOTP setup QR generation |
| **Secure Storage** | react-native-keychain | JWT token secure storage |

### 2.4 Server Communication

```
┌──────────────┐         ┌──────────────┐
│  Mobile App  │  HTTP   │   Claribot   │
│              │────────▶│   Server     │
│  TanStack    │  JSON   │  :9847/api/* │
│  Query       │◀────────│              │
│              │         │              │
│  FCM Token   │  Push   │  Push Agent  │
│  ──────────▶ │◀────────│  (future)    │
└──────────────┘         └──────────────┘
```

- **Primary**: HTTP REST API (`/api/*`) - identical endpoints to Web UI
- **Authentication**: JWT cookie → Bearer token adaptation (mobile-friendly)
- **Real-time Updates**: Polling via TanStack Query (same intervals as Web UI: 5s active, 15-30s idle)
- **Push Notifications**: FCM for task completion, cycle events, error alerts (Phase 2)

---

## 3. Core Features

### 3.1 Feature Matrix

| Feature | Web UI Equivalent | Mobile Behavior |
|---------|-------------------|-----------------|
| **Dashboard** | Dashboard.tsx | Project summary cards, cycle status, quick stats |
| **Projects** | Projects.tsx | List with search, sort, pin, category filter |
| **Project Edit** | ProjectEdit.tsx | Edit description, parallel, category, delete |
| **Tasks** | Tasks.tsx | Tree/list view, status filter, detail drawer |
| **Messages** | Messages.tsx | Chat UI with bubbles, send messages |
| **Specs** | Specs.tsx | Dual-panel list + detail, status lifecycle |
| **Schedules** | Schedules.tsx | Schedule list, enable/disable toggle, run history |
| **Settings** | Settings.tsx | Server connection, notification preferences |
| **Login/Setup** | Login.tsx, Setup.tsx | Password + TOTP, QR code setup |

### 3.2 Project Management

- Browse all projects with search, sort (last accessed / created / task count), and category filter
- Pin/unpin projects for quick access
- Switch active project
- View project details and edit settings (description, parallel workers, category)
- Project task statistics overview

### 3.3 Task Management

- **Tree View**: Hierarchical task display with collapsible nodes (touch-optimized indentation)
- **List View**: Flat list with status filter (todo, planned, split, done, failed)
- **Task Detail**: Full-screen view with Spec, Plan, Report rendered as markdown
- **Actions**: Create task, edit title/spec/status/priority, delete
- **Bulk Operations**: Plan All, Run All, Cycle, Stop - with confirmation dialogs
- **Status Bar**: Visual task status distribution (same as Web UI top bar)

### 3.4 Messages & Chat

- Chat bubble interface (user messages right, bot responses left)
- Send messages to Claude Code for the active project
- View message status (pending → processing → done/failed)
- Message detail with full result/report view
- Project-scoped message filtering

### 3.5 Specs Management

- List specs with status badges (draft, review, approved, deprecated)
- Create/edit specs with title and markdown content
- Status lifecycle management
- Priority ordering

### 3.6 Schedule Management

- List schedules with cron expression display
- Enable/disable toggle switches
- View execution history (runs) for each schedule
- Schedule type indicator (Claude / Bash)
- Project-scoped schedule filtering

### 3.7 Push Notifications (Phase 2)

| Event | Notification Content |
|-------|---------------------|
| Task Completed | "Task #42 'Fix login bug' completed successfully" |
| Task Failed | "Task #43 'Add auth' failed: timeout" |
| Cycle Completed | "Cycle completed: 12/12 tasks done" |
| Cycle Interrupted | "Cycle interrupted at task #44 - auth error" |
| Message Processed | "Response ready for your message in project 'api-server'" |
| Schedule Run Failed | "Schedule 'daily-backup' failed" |

### 3.8 Authentication

- **Setup Flow**: 3-step wizard (password → QR code → TOTP verify) - same as Web UI
- **Login**: Password + 6-digit TOTP code
- **Token Storage**: JWT stored in secure keychain (not AsyncStorage)
- **Session**: Auto-refresh, logout on expiry

---

## 4. Screen Design

### 4.1 Navigation Structure

```
┌─────────────────────────────────────┐
│  Bottom Tab Navigation              │
├─────────┬─────────┬─────────┬───────┤
│Dashboard│  Tasks  │Messages │  More │
│   🏠    │   ✅    │   💬    │  ···  │
└─────────┴─────────┴─────────┴───────┘
                                  │
                          ┌───────┴───────┐
                          │  More Menu    │
                          ├───────────────┤
                          │  Projects     │
                          │  Specs        │
                          │  Schedules    │
                          │  Settings     │
                          └───────────────┘
```

- **Bottom Tabs**: 4 primary tabs (Dashboard, Tasks, Messages, More)
- **Header**: Project selector dropdown + Claude status badge
- **Stack Navigation**: Each tab has its own navigation stack for drill-down

### 4.2 Dashboard Screen

```
┌─────────────────────────────────┐
│  [Project ▼]        🟢 Claude  │  ← Header
├─────────────────────────────────┤
│                                 │
│  ┌─────────────────────────────┐│
│  │  Cycle Status               ││
│  │  ████████░░ 8/10 Running    ││
│  │  Phase: Running  ⏱ 12m     ││
│  └─────────────────────────────┘│
│                                 │
│  ┌─────────┐ ┌─────────┐       │
│  │  todo   │ │ planned │       │
│  │   12    │ │    5    │       │
│  └─────────┘ └─────────┘       │
│  ┌─────────┐ ┌─────────┐       │
│  │  done   │ │ failed  │       │
│  │   42    │ │    1    │       │
│  └─────────┘ └─────────┘       │
│                                 │
│  Recent Messages                │
│  ┌─────────────────────────────┐│
│  │ 💬 Fix auth bug    2m ago  ││
│  │ 💬 Add logging     15m ago ││
│  └─────────────────────────────┘│
│                                 │
│  Project Stats                  │
│  ┌─────────────────────────────┐│
│  │ api-server   ██░░ 8/10     ││
│  │ blog         ████ 20/20    ││
│  │ docs         █░░░ 3/12     ││
│  └─────────────────────────────┘│
│                                 │
├─────────┬─────────┬─────────┬───┤
│Dashboard│  Tasks  │Messages │···│
└─────────┴─────────┴─────────┴───┘
```

### 4.3 Tasks Screen

```
┌─────────────────────────────────┐
│  [Project ▼]        🟢 Claude  │
├─────────────────────────────────┤
│  Status: ⬜12 🔵5 🟢42 🔴1    │  ← Status bar
├─────────────────────────────────┤
│  [🌳 Tree] [📋 List]  [Filter ▼]│  ← View toggle + filter
├─────────────────────────────────┤
│                                 │
│  ▼ #1 Build auth system  [done] │
│    ▶ #2 Design schema   [done] │
│    ▶ #3 Implement API   [done] │
│    ▼ #4 Add tests       [todo] │
│      #5 Unit tests      [todo] │
│      #6 E2E tests       [todo] │
│  ▶ #7 Add logging      [planned]│
│  #8 Fix CSS bug         [todo] │
│                                 │
│                    [+ New Task] │  ← FAB button
├─────────────────────────────────┤
│  [▶ Plan All] [▶▶ Run All]     │
│  [🔄 Cycle]   [⏹ Stop]         │
├─────────┬─────────┬─────────┬───┤
│Dashboard│  Tasks  │Messages │···│
└─────────┴─────────┴─────────┴───┘

        ↓ Tap on task #4

┌─────────────────────────────────┐
│  ← Back         Task #4        │
├─────────────────────────────────┤
│  Add tests                      │
│  Status: todo    Priority: 3    │
│  Parent: #1 Build auth system   │
│  Created: 2025-01-15 14:30      │
├─────────────────────────────────┤
│  [Spec] [Plan] [Report]        │  ← Tab bar
├─────────────────────────────────┤
│                                 │
│  ## Spec                        │
│                                 │
│  Write unit tests and E2E tests │
│  for the authentication system: │
│  - Login endpoint               │
│  - Token refresh                │
│  - Permission middleware        │
│                                 │
│                                 │
│  [✏️ Edit] [▶ Plan] [▶▶ Run]   │
└─────────────────────────────────┘
```

### 4.4 Messages / Chat Screen

```
┌─────────────────────────────────┐
│  [Project ▼]        🟢 Claude  │
├─────────────────────────────────┤
│                                 │
│           Jan 15, 2025          │
│                                 │
│        ┌───────────────────┐    │
│        │ Fix the login bug │    │
│        │ in auth.go        │    │
│        │         14:30  ✓  │    │
│        └───────────────────┘    │
│                                 │
│  ┌───────────────────┐          │
│  │ Done. Fixed null   │          │
│  │ pointer in         │          │
│  │ ValidateToken()    │          │
│  │                    │          │
│  │ [View Report ▶]    │          │
│  │ 14:32  ✅ done     │          │
│  └───────────────────┘          │
│                                 │
│        ┌───────────────────┐    │
│        │ Add rate limiting │    │
│        │ to the API        │    │
│        │         15:01  ⏳ │    │
│        └───────────────────┘    │
│                                 │
├─────────────────────────────────┤
│  [Type a message...       ] [➤] │
├─────────┬─────────┬─────────┬───┤
│Dashboard│  Tasks  │Messages │···│
└─────────┴─────────┴─────────┴───┘
```

### 4.5 Settings Screen

```
┌─────────────────────────────────┐
│  Settings                       │
├─────────────────────────────────┤
│                                 │
│  Server Connection              │
│  ┌─────────────────────────────┐│
│  │ Server URL                  ││
│  │ [http://192.168.1.10:9847 ] ││
│  │                             ││
│  │ Status: 🟢 Connected       ││
│  │ Version: v0.2.21           ││
│  │ Uptime: 3d 12h             ││
│  └─────────────────────────────┘│
│                                 │
│  Notifications                  │
│  ┌─────────────────────────────┐│
│  │ Push Notifications    [🔵] ││
│  │ Task Completion       [🔵] ││
│  │ Cycle Events          [🔵] ││
│  │ Error Alerts          [🔵] ││
│  └─────────────────────────────┘│
│                                 │
│  Cache                          │
│  ┌─────────────────────────────┐│
│  │ Cached Data: 2.4 MB        ││
│  │ [Clear Cache]              ││
│  └─────────────────────────────┘│
│                                 │
│  Account                        │
│  ┌─────────────────────────────┐│
│  │ [Logout]                   ││
│  └─────────────────────────────┘│
│                                 │
├─────────┬─────────┬─────────┬───┤
│Dashboard│  Tasks  │Messages │···│
└─────────┴─────────┴─────────┴───┘
```

---

## 5. Architecture

### 5.1 Application Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Mobile App                            │
│                                                         │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │  Screens │  │  Navigation  │  │   Notifications  │  │
│  │ (React   │  │  (React      │  │   (FCM + Local)  │  │
│  │  Native) │  │  Navigation) │  │                  │  │
│  └────┬─────┘  └──────────────┘  └──────────────────┘  │
│       │                                                  │
│  ┌────▼──────────────────────────────────────────────┐  │
│  │              TanStack Query Layer                  │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────┐    │  │
│  │  │ Queries  │  │Mutations │  │ Invalidation │    │  │
│  │  │ (cached) │  │(optimistic│  │ (on mutation) │    │  │
│  │  └──────────┘  └──────────┘  └──────────────┘    │  │
│  └────┬──────────────────────────────────────────────┘  │
│       │                                                  │
│  ┌────▼──────────────────────────────────────────────┐  │
│  │              API Client Layer                      │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────┐    │  │
│  │  │ HTTP     │  │  Auth    │  │   Offline    │    │  │
│  │  │ Client   │  │  (JWT)   │  │   Queue      │    │  │
│  │  └──────────┘  └──────────┘  └──────────────┘    │  │
│  └────┬──────────────────────────────────────────────┘  │
│       │                                                  │
│  ┌────▼──────────────────────────────────────────────┐  │
│  │              Storage Layer                         │  │
│  │  ┌──────────┐  ┌──────────────┐                   │  │
│  │  │ MMKV     │  │  Keychain    │                   │  │
│  │  │ (cache)  │  │  (JWT token) │                   │  │
│  │  └──────────┘  └──────────────┘                   │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                        │
                   HTTP / JSON
                        │
                        ▼
┌─────────────────────────────────────────────────────────┐
│                  Claribot Server                         │
│                  (127.0.0.1:9847)                        │
│                                                          │
│   /api/auth/*    /api/projects/*    /api/tasks/*         │
│   /api/messages/*  /api/schedules/*  /api/specs/*        │
│   /api/configs/*   /api/status       /api/usage          │
└─────────────────────────────────────────────────────────┘
```

### 5.2 API Reuse Strategy

The mobile app reuses the existing RESTful API endpoints without any backend modifications.

**Shared Code from Web UI** (via shared package or copy):

| Web UI Module | Mobile Equivalent | Reuse Level |
|---------------|-------------------|-------------|
| `types/index.ts` | `types/index.ts` | 100% - identical TypeScript types |
| `api/client.ts` | `api/client.ts` | 90% - replace `fetch` with RN fetch, URL config |
| `hooks/useClaribot.ts` | `hooks/useClaribot.ts` | 95% - identical query/mutation logic |
| `hooks/useAuth.ts` | `hooks/useAuth.ts` | 80% - adapt token storage to keychain |
| Component logic | Screen logic | 60% - adapt JSX to RN components |

**API Endpoint Mapping** (complete coverage):

```
Auth:       POST /api/auth/setup, /login, /logout, GET /status, /totp-setup
Status:     GET  /api/status, /health, /usage
Projects:   GET/POST /api/projects, GET/PATCH/DELETE /api/projects/{id}
Tasks:      GET/POST /api/tasks, GET/PATCH/DELETE /api/tasks/{id}
            POST /api/tasks/{id}/plan, /run, /plan-all, /run-all, /cycle, /stop
Messages:   GET/POST /api/messages, GET /api/messages/{id}, /status, /processing
Schedules:  GET/POST /api/schedules, GET/PATCH/DELETE /api/schedules/{id}
            POST /api/schedules/{id}/enable, /disable, GET /runs
Specs:      GET/POST /api/specs, GET/PATCH/DELETE /api/specs/{id}
Configs:    GET/PUT/DELETE /api/configs/{key}, GET/PUT /api/config-yaml
```

### 5.3 Authentication Adaptation

The Web UI uses HTTP-only JWT cookies. For the mobile app:

- **Option A (Recommended)**: Add `Authorization: Bearer <token>` header support to the existing auth middleware. The server returns the JWT in the response body on login, and the mobile app stores it in the system keychain.
- **Option B**: Continue using cookies with React Native's cookie jar.

Option A is preferred because:
1. Secure keychain storage is more appropriate for mobile
2. No cookie management complexity
3. Works better with background fetch/notifications
4. Minimal server change (one middleware check: cookie OR Authorization header)

### 5.4 Network & Offline Strategy

| Scenario | Behavior |
|----------|----------|
| **Local Network** | Direct HTTP to `192.168.x.x:9847` (same LAN) |
| **Remote Access** | VPN or reverse proxy with HTTPS required |
| **Offline - Read** | TanStack Query cache serves stale data with visual indicator |
| **Offline - Write** | Queue mutations locally, sync when reconnected |
| **Reconnection** | Auto-invalidate all queries, replay mutation queue |

---

## 6. Implementation Roadmap

### Phase 1: Core MVP

**Goal**: Basic read + send functionality, equivalent to essential Telegram bot features.

| Feature | Details |
|---------|---------|
| Authentication | Login (password + TOTP), setup wizard |
| Dashboard | Project stats, Claude status, cycle status |
| Projects | List, switch, search, filter |
| Tasks | Tree/list view, status filter, task detail (read-only) |
| Messages | Chat UI, send messages, view responses |
| Settings | Server URL configuration, logout |

**Estimated Scope**: 6 screens, API client, auth flow, basic navigation

### Phase 2: Full CRUD + Notifications

**Goal**: Feature parity with Web UI, plus mobile-exclusive push notifications.

| Feature | Details |
|---------|---------|
| Task CRUD | Create, edit, delete tasks |
| Task Actions | Plan, Run, Plan All, Run All, Cycle, Stop |
| Specs CRUD | Create, edit, delete specs with status management |
| Schedules | Full CRUD, enable/disable, run history |
| Project Edit | Edit description, parallel, category, delete |
| Push Notifications | FCM integration for task/cycle/message events |
| Server Push Agent | Backend component to send FCM messages on events |

### Phase 3: Advanced Features

**Goal**: Mobile-optimized enhancements beyond Web UI feature parity.

| Feature | Details |
|---------|---------|
| Offline Cache | MMKV-based persistent cache, offline read access |
| Mutation Queue | Queue writes when offline, sync on reconnect |
| Home Screen Widget | Quick task stats / cycle status widget (iOS/Android) |
| Quick Actions | 3D Touch / Long Press shortcuts (New Task, Send Message) |
| Biometric Auth | Face ID / Fingerprint unlock (stored JWT + biometric gate) |
| Dark Mode | System theme integration |
| Haptic Feedback | Tactile response for task completion, cycle events |

---

## 7. Project Structure

```
claribot-mobile/
├── src/
│   ├── api/
│   │   └── client.ts           # API client (adapted from gui/src/api/client.ts)
│   ├── screens/
│   │   ├── DashboardScreen.tsx
│   │   ├── TasksScreen.tsx
│   │   ├── TaskDetailScreen.tsx
│   │   ├── MessagesScreen.tsx
│   │   ├── ProjectsScreen.tsx
│   │   ├── ProjectEditScreen.tsx
│   │   ├── SpecsScreen.tsx
│   │   ├── SchedulesScreen.tsx
│   │   ├── SettingsScreen.tsx
│   │   ├── LoginScreen.tsx
│   │   └── SetupScreen.tsx
│   ├── components/
│   │   ├── ChatBubble.tsx
│   │   ├── MarkdownRenderer.tsx
│   │   ├── ProjectSelector.tsx
│   │   ├── TaskTreeItem.tsx
│   │   ├── StatusBar.tsx
│   │   └── CycleProgress.tsx
│   ├── hooks/
│   │   ├── useClaribot.ts      # TanStack Query hooks (shared with Web UI)
│   │   └── useAuth.ts          # Auth hooks (adapted for keychain)
│   ├── navigation/
│   │   ├── TabNavigator.tsx
│   │   ├── StackNavigator.tsx
│   │   └── AuthNavigator.tsx
│   ├── types/
│   │   └── index.ts            # Shared types (identical to gui/src/types/index.ts)
│   ├── storage/
│   │   ├── cache.ts            # MMKV cache wrapper
│   │   └── keychain.ts         # Secure JWT storage
│   └── App.tsx
├── android/
├── ios/
├── package.json
└── tsconfig.json
```

---

## 8. Server-Side Changes Required

Minimal backend modifications needed:

| Change | File | Description |
|--------|------|-------------|
| Bearer token auth | `bot/internal/handler/auth.go` | Accept `Authorization: Bearer <token>` header in addition to JWT cookie |
| Login response body | `bot/internal/handler/auth.go` | Return JWT token in login response JSON (not just Set-Cookie) |
| Push notification agent | `bot/internal/push/` (new) | FCM integration to send push on task/cycle events (Phase 2) |
| FCM token registration | `POST /api/auth/fcm-token` (new) | Store device FCM tokens for push delivery (Phase 2) |

Phase 1 requires only the Bearer token authentication change (2 files modified).

---

## 9. Design Guidelines

### 9.1 Touch Targets

- Minimum touch target: 44x44 points (Apple HIG) / 48x48 dp (Material Design)
- This aligns with the Web UI's existing mobile touch target work (tasks #33-#36)

### 9.2 Typography

- Use system fonts (San Francisco on iOS, Roboto on Android)
- Monospace font for code blocks, cron expressions, and task IDs

### 9.3 Color System

- Match Web UI's color palette (shadcn/ui defaults)
- Support system dark/light mode toggle
- Status colors: todo (gray), planned (blue), split (purple), done (green), failed (red)

### 9.4 Platform Conventions

- iOS: Large title navigation bars, swipe-back gesture, bottom tab bar
- Android: Material Design 3, navigation drawer option, FAB for primary actions
- Both: Pull-to-refresh on all list screens
