// Package formatter provides log line formatting, filtering, and transformation.
//
// # Aggregation
//
// The Aggregator type accumulates numeric field values across multiple JSON log
// lines and produces a summary report. It supports four modes:
//
//   - count: number of lines that contain the field
//   - sum:   total of all field values
//   - min:   smallest observed field value
//   - max:   largest observed field value
//
// Use ParseAggregateFlag to parse a "field:mode" CLI flag value, and
// AggregateFlagValue as a flag.Value implementation for use with the standard
// flag or pflag packages.
//
// # Multiple Aggregations
//
// Multiple aggregations can be applied in a single pass by creating several
// Aggregator instances and calling Add on each. The Report method returns
// results sorted by field name for deterministic output.
//
// Example:
//
//	agg := formatter.NewAggregator("latency_ms", formatter.AggregateSum)
//	for _, line := range lines {
//		agg.Add(line)
//	}
//	for _, entry := range agg.Report() {
//		fmt.Println(entry)
//	}
//
// # CLI Flag Usage
//
// AggregateFlagValue can be registered with the standard flag package:
//
//	var agg formatter.AggregateFlagValue
//	flag.Var(&agg, "aggregate", "field:mode to aggregate (e.g. latency_ms:sum)")
package formatter
