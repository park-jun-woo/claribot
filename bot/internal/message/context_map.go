package message

import (
	"database/sql"
	"fmt"
	"strings"
	"unicode/utf8"

	"parkjunwoo.com/claribot/internal/db"
	"parkjunwoo.com/claribot/internal/task"
)

// BuildContextMap builds a context map text combining recent messages and task tree.
// Returns empty string on failure to avoid disrupting message processing.
func BuildContextMap(globalDB *db.DB, projectPath string, projectID *string, contextMax int) string {
	var sb strings.Builder

	// Section 1: Recent messages (reports only) from project or global
	msgSection := buildMessageSection(globalDB, projectID, contextMax)
	if msgSection != "" {
		sb.WriteString("## 최근 대화 이력\n\n")
		sb.WriteString(msgSection)
		sb.WriteString("\n")
	}

	// Section 2: Task tree from local DB (only if project path exists)
	if projectPath != "" {
		taskSection := buildTaskSection(projectPath)
		if taskSection != "" {
			sb.WriteString("## Task 현황\n\n")
			sb.WriteString("```\n")
			sb.WriteString(taskSection)
			sb.WriteString("```\n\n")
		}
	}

	if sb.Len() == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString("# Context Map\n\n")
	result.WriteString(sb.String())
	return result.String()
}

// buildMessageSection queries recent done messages (with reports) from project or global.
func buildMessageSection(globalDB *db.DB, projectID *string, contextMax int) string {
	var rows *sql.Rows
	var err error

	if projectID != nil && *projectID != "" {
		// Filter by project
		rows, err = globalDB.Query(`
			SELECT id, content, COALESCE(result, '')
			FROM messages
			WHERE project_id = ? AND status = 'done' AND result != ''
			ORDER BY id DESC
			LIMIT ?
		`, *projectID, contextMax)
	} else {
		// Global (all projects)
		rows, err = globalDB.Query(`
			SELECT id, content, COALESCE(result, '')
			FROM messages
			WHERE status = 'done' AND result != ''
			ORDER BY id DESC
			LIMIT ?
		`, contextMax)
	}

	if err != nil {
		return ""
	}
	defer rows.Close()

	var sb strings.Builder
	count := 0
	for rows.Next() {
		var id int
		var content, result string
		if err := rows.Scan(&id, &content, &result); err != nil {
			continue
		}

		contentFirst := firstLine(content, 60)
		// Show more of the report (first 200 chars)
		reportSummary := truncateText(result, 200)

		sb.WriteString(fmt.Sprintf("### #%d: %s\n", id, contentFirst))
		sb.WriteString(reportSummary)
		sb.WriteString("\n\n")
		count++
	}

	if count == 0 {
		return ""
	}
	return sb.String()
}

// buildTaskSection opens the local DB and builds a task tree summary with stats.
func buildTaskSection(projectPath string) string {
	localDB, err := db.OpenLocal(projectPath)
	if err != nil {
		return ""
	}
	defer localDB.Close()

	contextMap, err := task.BuildContextMap(localDB)
	if err != nil {
		return ""
	}

	// Append task stats
	stats := buildTaskStats(localDB)
	if stats != "" {
		return contextMap + stats
	}
	return contextMap
}

// buildTaskStats returns a one-line task status summary.
func buildTaskStats(localDB *db.DB) string {
	rows, err := localDB.Query(`
		SELECT status, COUNT(*) as cnt
		FROM tasks
		GROUP BY status
	`)
	if err != nil {
		return ""
	}
	defer rows.Close()

	statusCounts := make(map[string]int)
	total := 0
	for rows.Next() {
		var status string
		var cnt int
		if err := rows.Scan(&status, &cnt); err != nil {
			continue
		}
		statusCounts[status] = cnt
		total += cnt
	}

	if total == 0 {
		return ""
	}

	var parts []string
	for _, s := range []string{"todo", "planned", "split", "done", "failed"} {
		if c, ok := statusCounts[s]; ok && c > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", s, c))
		}
	}

	return fmt.Sprintf("통계: 총 %d개 (%s)\n", total, strings.Join(parts, ", "))
}

// firstLine returns the first line of text, truncated to maxRunes.
func firstLine(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// Take first line only
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)

	if utf8.RuneCountInString(s) > maxRunes {
		s = string([]rune(s)[:maxRunes]) + "..."
	}
	return s
}

// truncateText truncates text to maxRunes, preserving word boundaries where possible.
func truncateText(s string, maxRunes int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}

	return string(runes[:maxRunes]) + "..."
}
