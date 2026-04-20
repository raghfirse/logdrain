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
// Example:
//
//	agg := formatter.NewAggregator("latency_ms", formatter.AggregateSum)
//	for _, line := range lines {
//		agg.Add(line)
//	}
//	for _, entry := range agg.Report() {
//		fmt.Println(entry)
//	}
package formatter
