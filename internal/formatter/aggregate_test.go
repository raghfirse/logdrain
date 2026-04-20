package formatter

import (
	"testing"
)

func TestAggregator_Count(t *testing.T) {
	a := NewAggregator("latency", AggregateCount)
	a.Add(`{"latency": 10}`)
	a.Add(`{"latency": 20}`)
	a.Add(`{"latency": 30}`)
	report := a.Report()
	if len(report) != 1 || report[0] != "latency=3" {
		t.Errorf("expected latency=3, got %v", report)
	}
}

func TestAggregator_Sum(t *testing.T) {
	a := NewAggregator("bytes", AggregateSum)
	a.Add(`{"bytes": 100}`)
	a.Add(`{"bytes": 200}`)
	report := a.Report()
	if len(report) != 1 || report[0] != "bytes=300" {
		t.Errorf("expected bytes=300, got %v", report)
	}
}

func TestAggregator_Min(t *testing.T) {
	a := NewAggregator("duration", AggregateMin)
	a.Add(`{"duration": 50}`)
	a.Add(`{"duration": 10}`)
	a.Add(`{"duration": 80}`)
	report := a.Report()
	if len(report) != 1 || report[0] != "duration=10" {
		t.Errorf("expected duration=10, got %v", report)
	}
}

func TestAggregator_Max(t *testing.T) {
	a := NewAggregator("score", AggregateMax)
	a.Add(`{"score": 5}`)
	a.Add(`{"score": 99}`)
	a.Add(`{"score": 42}`)
	report := a.Report()
	if len(report) != 1 || report[0] != "score=99" {
		t.Errorf("expected score=99, got %v", report)
	}
}

func TestAggregator_IgnoresNonJSON(t *testing.T) {
	a := NewAggregator("val", AggregateSum)
	a.Add("not json")
	a.Add(`{"val": 5}`)
	report := a.Report()
	if len(report) != 1 || report[0] != "val=5" {
		t.Errorf("unexpected report: %v", report)
	}
}

func TestAggregator_IgnoresMissingField(t *testing.T) {
	a := NewAggregator("missing", AggregateCount)
	a.Add(`{"other": 1}`)
	report := a.Report()
	if len(report) != 0 {
		t.Errorf("expected empty report, got %v", report)
	}
}

func TestAggregator_Reset(t *testing.T) {
	a := NewAggregator("x", AggregateSum)
	a.Add(`{"x": 10}`)
	a.Reset()
	report := a.Report()
	if len(report) != 0 {
		t.Errorf("expected empty after reset, got %v", report)
	}
}

func TestAggregator_EmptyReport(t *testing.T) {
	a := NewAggregator("n", AggregateCount)
	report := a.Report()
	if len(report) != 0 {
		t.Errorf("expected empty report on fresh aggregator")
	}
}
