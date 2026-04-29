// Package formatter — clamp transformer
//
// The clamp transformer constrains a numeric JSON field to a closed [min, max]
// interval. Values below min are raised to min; values above max are lowered
// to max. Non-numeric fields and non-JSON lines are passed through unchanged.
//
// Usage:
//
//	--clamp field,min,max
//
// Examples:
//
//	# Keep latency between 0 and 5000 ms
//	logdrain --clamp latency,0,5000
//
//	# Clamp a percentage field to [0, 100]
//	logdrain --clamp pct,0,100
//
// Multiple --clamp flags may be provided to clamp several fields at once.
package formatter
