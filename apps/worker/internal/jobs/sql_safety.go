package jobs

import (
	"strings"
	"time"
)

var (
	deniedKeywords = []string{
		"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE",
		"TRUNCATE", "GRANT", "REVOKE", "COPY", "CALL", "EXEC",
		"MERGE", "VACUUM", "ANALYZE",
	}
)

type SQLSafetyChecker struct {
	maxTimeout time.Duration
	maxRows    int
}

func NewSQLSafetyChecker(maxTimeout time.Duration, maxRows int) *SQLSafetyChecker {
	return &SQLSafetyChecker{
		maxTimeout: maxTimeout,
		maxRows:    maxRows,
	}
}
func (c *SQLSafetyChecker) ValidateSQL(sql string, allowExplainAnalyze bool) error {
	sqlUpper := strings.ToUpper(strings.TrimSpace(sql))
	// Check for denied keywords
	for _, keyword := range deniedKeywords {
		if strings.Contains(sqlUpper, keyword) {
			// Exception: EXPLAIN ANALYZE is allowed if explicitly permitted
			if keyword == "ANALYZE" && allowExplainAnalyze && strings.HasPrefix(sqlUpper, "EXPLAIN") {
				continue
			}
			return &SQLSafetyError{
				Reason:  "denied_keyword",
				Keyword: keyword,
				Message: "SQL contains denied keyword: " + keyword,
			}
		}
	}
	// Must start with SELECT or EXPLAIN
	if !strings.HasPrefix(sqlUpper, "SELECT") && !strings.HasPrefix(sqlUpper, "EXPLAIN") {
		return &SQLSafetyError{
			Reason:  "invalid_statement_type",
			Message: "SQL must be a SELECT or EXPLAIN statement",
		}
	}
	// Check for multiple statements (semicolon check)
	statements := strings.Split(sql, ";")
	nonEmptyStatements := 0
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			nonEmptyStatements++
		}
	}
	if nonEmptyStatements > 1 {
		return &SQLSafetyError{
			Reason:  "multiple_statements",
			Message: "Only single-statement queries are allowed",
		}
	}
	return nil
}
func (c *SQLSafetyChecker) EnforceSampleLimit(sql string, sampleMode bool) (string, error) {
	if !sampleMode {
		return sql, nil
	}
	sqlUpper := strings.ToUpper(sql)
	if !strings.Contains(sqlUpper, "LIMIT") {
		return "", &SQLSafetyError{
			Reason:  "missing_limit",
			Message: "Sample mode requires LIMIT clause",
		}
	}
	return sql, nil
}

type SQLSafetyError struct {
	Reason  string
	Keyword string
	Message string
}

func (e *SQLSafetyError) Error() string {
	return e.Message
}
