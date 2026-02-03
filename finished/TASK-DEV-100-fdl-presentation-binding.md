# TASK-DEV-100: FDL Presentation Layer 바인딩/루프 파싱

## 개요
FDL Service의 Presentation Layer(UI) 파싱에서 바인딩 구문 및 루프 구조 파싱 구현

## 배경
- **보고서**: reports/2026-02-03-total-report.md 섹션 4.5
- **현재 구현율**: 50%
- **스펙**: specs/FDL/02-D-PresentationLayer.md
- **구현 파일**: cli/internal/service/fdl_service.go

## 현재 지원 상태

### 지원됨
- 컴포넌트 타입: Page, Template, Organism, Molecule, Atom
- Props, State, Computed 파싱
- Init, Methods, View 파싱
- 조건부 렌더링 (if/then/else)
- 액션 타입: call, set, navigate, show, validate, confirm, emit, parallel, redirect

### 미지원 (구현 대상)
- 바인딩 구문 파싱 ({variable})
- for 루프 파싱
- 이벤트 핸들러 검증 (@click, @change)
- 컴포넌트 계층 순환 참조 검증
- Props 타입 검증 (ReactNode, 제네릭)

## 작업 내용

### 1. 바인딩 구문 파싱
```yaml
view:
  - Text:
      content: "Hello, {user.name}!"
      class: "text-{theme.color}"
```
- `{expression}` 패턴 추출
- 바인딩된 변수가 state/props/computed에 존재하는지 검증
- 중첩 바인딩 지원 ({user.profile.name})

구현:
```go
type Binding struct {
    Expression string
    Path       []string  // ["user", "name"]
}

func extractBindings(value string) []Binding {
    re := regexp.MustCompile(`\{([^}]+)\}`)
    matches := re.FindAllStringSubmatch(value, -1)
    // ...
}
```

### 2. for 루프 파싱
```yaml
view:
  - for:
      each: items
      as: item
      index: idx
      key: item.id
      render:
        - ListItem:
            data: "{item}"
```
- `each`, `as`, `render` 필수 검증
- `index`, `key` 선택적 파싱
- 루프 변수(item, idx)가 render 내에서만 유효함 확인

### 3. 이벤트 핸들러 검증
```yaml
view:
  - Button:
      text: "Submit"
      @click: handleSubmit
      @hover: showTooltip
```
- `@` 접두사 이벤트 핸들러 파싱
- 핸들러가 methods에 정의되어 있는지 검증
- 지원 이벤트: click, change, submit, hover, focus, blur, keydown, keyup

### 4. 컴포넌트 순환 참조 검증
```yaml
# ComponentA uses ComponentB
# ComponentB uses ComponentA -> 순환 참조!
```
- 컴포넌트 의존성 그래프 구축
- 순환 감지 알고리즘 (DFS)
- 순환 발견 시 에러 반환

### 5. Props 타입 검증 확장
```yaml
props:
  - name: children
    type: ReactNode
  - name: items
    type: List<User>
  - name: onClick
    type: () => void
```
- 제네릭 타입 파싱 (List<T>, Map<K,V>)
- 함수 타입 파싱 (() => void, (arg: T) => R)
- ReactNode, ReactElement 등 특수 타입 인식

## 구현 방법

fdl_service.go에 다음 함수 추가:
- `extractBindings(value string) []Binding`
- `validateBindings(bindings []Binding, scope *ComponentScope) error`
- `parseForLoop(node map[string]interface{}) (*ForLoop, error)`
- `parseEventHandler(key string, value interface{}) (*EventHandler, error)`
- `detectCircularDependency(components map[string]*Component) error`
- `parseGenericType(typeStr string) (*GenericType, error)`

### 새로운 구조체
```go
type Binding struct {
    Raw        string
    Expression string
    Path       []string
}

type ForLoop struct {
    Each   string
    As     string
    Index  string
    Key    string
    Render []interface{}
}

type EventHandler struct {
    Event   string  // click, change, etc.
    Handler string  // method name
}
```

## 검증
```bash
cd cli && go test -v ./test/fdl_service_test.go -run TestPresentationLayer
```

## 완료 기준
- [ ] 바인딩 구문 추출 및 검증 구현
- [ ] for 루프 파싱 구현
- [ ] 이벤트 핸들러 (@event) 파싱 구현
- [ ] 순환 참조 검증 구현
- [ ] 제네릭 타입 파싱 구현
- [ ] 테스트 케이스 추가
