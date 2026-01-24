package cmd

import "testing"

func TestNormalizeFinanceReportType(t *testing.T) {
	reportType, err := normalizeFinanceReportType("financial")
	if err != nil {
		t.Fatalf("expected report type to parse, got %v", err)
	}
	if string(reportType) != "FINANCIAL" {
		t.Fatalf("expected FINANCIAL, got %q", reportType)
	}

	_, err = normalizeFinanceReportType("invalid")
	if err == nil {
		t.Fatal("expected error for invalid report type")
	}
}

func TestNormalizeFinanceReportDate(t *testing.T) {
	date, err := normalizeFinanceReportDate("2025-12")
	if err != nil {
		t.Fatalf("expected date to parse, got %v", err)
	}
	if date != "2025-12" {
		t.Fatalf("expected date to be 2025-12, got %q", date)
	}

	_, err = normalizeFinanceReportDate("2025-13")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestNormalizeFinanceReportRegion(t *testing.T) {
	region, err := normalizeFinanceReportRegion("FINANCIAL", "us")
	if err != nil {
		t.Fatalf("expected region to parse, got %v", err)
	}
	if region != "US" {
		t.Fatalf("expected region to be US, got %q", region)
	}

	region, err = normalizeFinanceReportRegion("FINANCE_DETAIL", "z1")
	if err != nil {
		t.Fatalf("expected Z1 to parse, got %v", err)
	}
	if region != "Z1" {
		t.Fatalf("expected region to be Z1, got %q", region)
	}

	_, err = normalizeFinanceReportRegion("FINANCE_DETAIL", "US")
	if err == nil {
		t.Fatal("expected error for non-Z1 FINANCE_DETAIL region")
	}
}
