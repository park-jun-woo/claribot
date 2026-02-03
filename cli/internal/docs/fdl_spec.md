# FDL (Feature Definition Language) Specification

## Structure

```yaml
feature: <feature_name>        # snake_case
description: <description>

models:                        # DATA LAYER
  - name: <ModelName>          # PascalCase
    table: <table_name>        # snake_case
    fields:
      - <name>: <type> [constraints...]

service:                       # LOGIC LAYER
  - name: <functionName>       # camelCase
    input: { key: type, ... }
    steps:
      - <step_description>

api:                           # INTERFACE LAYER
  - path: /path/{param}
    method: GET|POST|PUT|DELETE
    use: service.<function>    # Wiring
    request: { ... }
    response: { ... }

ui:                            # PRESENTATION LAYER
  - component: <ComponentName> # PascalCase
    type: Page|Organism|Molecule|Atom
    state: [...]
    view:
      - <element>: <label>
        action: API.<METHOD> /path
```

## Field Types

| Type | Description | Example |
|------|-------------|---------|
| `uuid` | UUID | `id: uuid (pk)` |
| `string` | String (255) | `name: string` |
| `string(N)` | String (N chars) | `code: string(50)` |
| `text` | Long text | `content: text` |
| `int` | Integer | `count: int` |
| `bigint` | Big integer | `views: bigint` |
| `float` | Float | `rating: float` |
| `decimal(M,N)` | Decimal | `price: decimal(10,2)` |
| `boolean` | Boolean | `is_active: boolean` |
| `datetime` | DateTime | `created_at: datetime` |
| `date` | Date only | `birth_date: date` |
| `json` | JSON object | `metadata: json` |
| `enum(a,b,c)` | Enum | `status: enum(draft,published)` |

## Constraints

| Constraint | Description | Example |
|------------|-------------|---------|
| `pk` | Primary Key | `id: uuid (pk)` |
| `fk: table.field` | Foreign Key | `user_id: uuid (fk: users.id)` |
| `required` | NOT NULL | `email: string (required)` |
| `unique` | Unique | `email: string (unique)` |
| `default: value` | Default | `status: string (default: "active")` |
| `nullable` | Allow NULL | `deleted_at: datetime (nullable)` |
| `index` | Index | `email: string (index)` |
| `onDelete: cascade` | FK cascade | `fk: users.id (onDelete: cascade)` |

## Service Steps

```yaml
service:
  - name: createComment
    input: { userId: uuid, postId: uuid, content: string }
    steps:
      - validate: "content 1-1000자 검증"
      - db: "INSERT INTO comments"
      - event: "작성자에게 알림 발송"
      - return: "생성된 Comment"
```

For simple CRUD: `steps: [CRUD Standard]`

## API Wiring

```yaml
api:
  - path: /posts/{postId}/comments
    method: POST
    summary: 댓글 작성
    use: service.createComment    # Must specify service function
    request:
      body: { content: string }
    response:
      201: { id: uuid, content: string }
```

## UI Wiring

```yaml
ui:
  - component: CommentForm
    type: Molecule
    state:
      - newComment: string
    view:
      - Textarea: "댓글 입력"
        bind: newComment
      - Button: "등록"
        action: API.POST /posts/{postId}/comments
        onSuccess: "목록에 추가"
```

## Example: Comment System

```yaml
feature: comment_system
description: 게시글 댓글 작성 및 조회

models:
  - name: Comment
    table: comments
    fields:
      - id: uuid (pk)
      - content: text (required)
      - post_id: uuid (fk: posts.id, onDelete: cascade)
      - user_id: uuid (fk: users.id)
      - created_at: datetime (default: now)

service:
  - name: createComment
    input: { userId: uuid, postId: uuid, content: string }
    steps:
      - validate: "content 1-1000자"
      - db: "INSERT INTO comments"
      - return: "Comment"

  - name: listComments
    input: { postId: uuid }
    steps:
      - db: "SELECT * FROM comments WHERE post_id = ?"
      - return: "Comment[]"

api:
  - path: /posts/{postId}/comments
    method: POST
    use: service.createComment
    request:
      body: { content: string }
    response:
      201: { id: uuid, content: string, created_at: datetime }

  - path: /posts/{postId}/comments
    method: GET
    use: service.listComments
    response:
      200: [{ id: uuid, content: string, user: { name: string } }]

ui:
  - component: CommentSection
    type: Organism
    props: { postId: uuid }
    state:
      - comments: Array
      - newComment: string
    init:
      - call: API.GET /posts/{postId}/comments -> comments

  - component: CommentForm
    type: Molecule
    parent: CommentSection
    view:
      - Textarea: "댓글 입력"
        bind: newComment
      - Button: "등록"
        action: API.POST /posts/{postId}/comments
        onSuccess: "comments에 추가"
```

## Naming Convention

| Item | Convention | Example |
|------|------------|---------|
| feature | snake_case | `comment_system` |
| model | PascalCase | `Comment` |
| table | snake_case, plural | `comments` |
| field | snake_case | `created_at` |
| service | camelCase | `createComment` |
| component | PascalCase | `CommentForm` |
| API path | kebab-case | `/user-profiles` |

## Checklist

- [ ] All `api` have `use:` field
- [ ] All UI `action` have API path
- [ ] `service.input` matches `api.request`
- [ ] FK references exist
