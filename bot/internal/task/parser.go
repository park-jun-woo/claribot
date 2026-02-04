package task

import (
	"regexp"
	"strings"
)

// PlanResult represents the result of 1회차 순회
type PlanResult struct {
	Type    string   // "subdivided" or "planned"
	Plan    string   // Plan content (if planned)
	Children []Child // Child tasks (if subdivided)
}

// Child represents a child task created during subdivision
type Child struct {
	ID    int
	Title string
}

// ParsePlanOutput parses Claude output from 1회차 순회
func ParsePlanOutput(output string) PlanResult {
	output = strings.TrimSpace(output)

	// Check for [SUBDIVIDED]
	if strings.HasPrefix(output, "[SUBDIVIDED]") {
		return PlanResult{
			Type:     "subdivided",
			Children: parseSubdividedChildren(output),
		}
	}

	// Check for [PLANNED]
	if strings.HasPrefix(output, "[PLANNED]") {
		plan := strings.TrimPrefix(output, "[PLANNED]")
		plan = strings.TrimSpace(plan)
		return PlanResult{
			Type: "planned",
			Plan: plan,
		}
	}

	// Default: treat entire output as plan (backward compatibility)
	return PlanResult{
		Type: "planned",
		Plan: output,
	}
}

// parseSubdividedChildren extracts child task info from [SUBDIVIDED] output
// Format: - Task #<id>: <title>
func parseSubdividedChildren(output string) []Child {
	var children []Child

	// Pattern: - Task #123: Some title
	re := regexp.MustCompile(`-\s*Task\s*#(\d+):\s*(.+)`)
	matches := re.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			var id int
			if _, err := regexp.MatchString(`^\d+$`, match[1]); err == nil {
				// Parse ID
				for _, c := range match[1] {
					id = id*10 + int(c-'0')
				}
			}
			children = append(children, Child{
				ID:    id,
				Title: strings.TrimSpace(match[2]),
			})
		}
	}

	return children
}
