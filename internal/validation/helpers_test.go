package validation

import "testing"

func findCheck(t *testing.T, checks []CheckResult, id string) CheckResult {
	t.Helper()
	for _, check := range checks {
		if check.ID == id {
			return check
		}
	}
	t.Fatalf("expected check %q to exist", id)
	return CheckResult{}
}

func hasCheckID(checks []CheckResult, id string) bool {
	for _, check := range checks {
		if check.ID == id {
			return true
		}
	}
	return false
}
