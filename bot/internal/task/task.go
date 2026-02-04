package task

// MaxDepth is the maximum depth for task subdivision
const MaxDepth = 5

// Task represents a task
type Task struct {
	ID        int
	ParentID  *int
	Title     string
	Spec      string // 요구사항 명세서 (불변)
	Plan      string // 계획서 (1회차 순회에서 생성, leaf만)
	Report    string // 완료 보고서 (2회차 순회 후 생성)
	Status    string // spec_ready → subdivided/plan_ready → done
	Error     string
	IsLeaf    bool // true: 실행 대상, false: 분할됨
	Depth     int  // 트리 깊이 (root=0)
	CreatedAt string
	UpdatedAt string
}
