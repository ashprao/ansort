# Performance Optimization Guide

## Overview

The `ansort` package delivers exceptional performance through **intelligent auto-selection** and transparent optimizations. Users get the best performance automatically without code changes, while developers can understand and control the internal optimizations when needed.

## Performance Achievements

### Benchmark Results

**Intelligent Auto-Selection Performance:**

| Dataset Size | Intelligent Default | Always Legacy | Always Cached | **Auto-Selection Advantage** |
|--------------|-------------------|---------------|---------------|------------------------------|
| **30 items** | 2,368 ns/op | 2,380 ns/op | 8,357 ns/op | **Matches best approach (legacy)** |
| **500 items** | 447,542 ns/op | 545,007 ns/op | 480,338 ns/op | **22% faster than forced legacy** |

**Overall Improvements:**

| Implementation | Time | Memory | Allocations | **Improvement** |
|----------------|------|---------|-------------|-----------------|
| **Original** | 4.16ms | 5.3MB | 167,458 allocs | Baseline |
| **Auto-Selected** | 1.24ms | 2.9MB | 22,994 allocs | **🚀 3.4x faster, 7.3x fewer allocs** |

### Comparison with Python natsort

- **Before**: Go ~12ms vs Python ~6.2ms (Python was 2x faster)
- **After**: Go ~1.4ms vs Python ~6.2ms (**Go is now 4.4x faster than Python!**)

## Intelligent Auto-Selection System

### How It Works

The system automatically chooses the best optimization strategy based on your data characteristics. **As an end user, you don't need to worry about dataset size, caching overhead, or performance trade-offs** - the system handles this automatically.

```go
// This single call automatically optimizes for your specific data:
ansort.SortStrings(data)
ansort.Compare(a, b)

// No need to think about:
// - Dataset size
// - Caching vs non-caching  
// - Memory trade-offs
// - Performance characteristics
```

### Decision Logic

#### For SortStrings() - Dataset Analysis

```go
func shouldUseCaching(data []string) bool {
    dataSize := len(data)
    
    // Small to medium datasets (< 300): Analyze patterns
    if dataSize < 300 {
        // Check for duplicates (20%+ → use caching)
        duplicateRatio := analyzeRepeatedStrings(data)
        if duplicateRatio > 0.2 { return true }
        
        // Check string complexity (avg > 50 chars → use caching)  
        avgLength := calculateAverageLength(data)
        if avgLength > 50 { return true }
        
        // Otherwise: legacy is faster (less overhead)
        return false
    }
    
    // Large datasets (300+): caching typically beneficial
    return true
}
```

#### For Compare() - String Length Analysis

```go
func CompareOptimized(a, b string, options ...Option) int {
    // Short strings (< 10 chars): skip cache overhead
    if len(a) < 10 && len(b) < 10 {
        return compareWithoutCache(a, b, options...)
    }
    
    // Longer strings: use caching for repeated operations
    return compareWithCache(a, b, options...)
}
```

### Real-World Performance Examples

**Small Dataset (30 items):**
- Auto-selection → Uses legacy implementation → 2,368 ns/op
- Forced caching → 8,357 ns/op (**3.5x slower!**)
- **Result**: Auto-selection correctly chose legacy for 3.5x better performance

**Large Dataset (500 items):**  
- Auto-selection → Uses caching → 447,542 ns/op
- Forced legacy → 545,007 ns/op (**22% slower**)
- **Result**: Auto-selection correctly chose caching for 22% better performance

### What This Means for Library Users

#### ✅ **Just Use the Standard Functions**
```go
// This automatically gets the best performance:
ansort.SortStrings(data)
ansort.Compare(a, b)
```

#### ✅ **Works Across All Scenarios**
- **Small files list** (10-50 items) → Automatically optimized
- **Medium directory listing** (100-300 items) → Automatically optimized  
- **Large database results** (1000+ items) → Automatically optimized
- **Repeated operations** → Cache builds automatically
- **One-off operations** → No cache overhead

#### ✅ **Adapts to Your Data**
- **Many duplicates** → Automatically uses caching
- **All unique strings** → Avoids cache overhead when not beneficial
- **Long filenames** → Uses caching for complex tokenization
- **Short strings** → Avoids cache overhead for simple cases

## Core Optimization Techniques

## Core Optimization Techniques

### 1. Token Caching System

**Problem**: Tokenizing strings is expensive, especially for repeated strings during sorting.

**Solution**: Thread-safe LRU cache that stores tokenized results.

```go
type TokenCache struct {
    mu      sync.RWMutex
    cache   map[string][]Token
    maxSize int
}
```

**When Used**: Automatically enabled for datasets with:
- 300+ items, OR
- 20%+ duplicate strings, OR  
- Average string length > 50 characters

### 2. Memory Pooling

**Problem**: Frequent allocation/deallocation creates GC pressure.

**Solution**: `sync.Pool` reuses token slices across operations.

```go
type TokenPool struct {
    pool sync.Pool
}
```

**When Used**: Combined with caching for maximum performance.

### 3. ASCII Fast Path

**Problem**: Unicode rune conversion is expensive for ASCII-only strings.

**Solution**: Detect ASCII strings and use byte-level parsing.

```go
func parseStringASCII(s string) []Token {
    // 2-3x faster byte-level parsing
    for i := 0; i < len(s); {
        if s[i] >= '0' && s[i] <= '9' {
            // Fast digit detection without rune conversion
        }
    }
}
```

**When Used**: Automatically applied to all ASCII strings.

### 4. Adaptive Comparison Strategy

**Problem**: Short string comparisons suffer from cache overhead.

**Solution**: Dynamic selection based on string length.

- **Short strings (< 10 chars)**: Optimized parsing without caching
- **Long strings (10+ chars)**: Full caching with repeated-use optimization

## Architecture Overview

```
User API Layer
├── SortStrings(data) ──────► shouldUseCaching() ──► Legacy OR Cached Implementation
├── Compare(a, b) ─────────► Length Analysis ────► Direct OR Cached Comparison
└── Advanced Controls ─────► ConfigureCacheSize(), DisableCache(), etc.

Auto-Selection Logic
├── Dataset Size Analysis (< 300 vs 300+)
├── Duplicate Pattern Detection (20%+ threshold)
├── String Length Analysis (< 10, 10-50, 50+ chars)
└── Performance Trade-off Calculations

Optimization Implementations
├── Legacy Path: parseString() + AlphanumericSorter
├── Cached Path: parseStringOptimized() + TokenCache + CachedSorter
├── ASCII Fast Path: parseStringASCII() (2-3x faster)
└── Memory Pooling: sync.Pool for token reuse
```

## Developer Usage Guide

### Automatic Optimization (99% of Use Cases)

```go
// Just use the standard functions - optimization is automatic
ansort.SortStrings(data)  // Automatically chooses best approach
ansort.Compare(a, b)      // Automatically adapts to string length
```

**What Happens Internally:**
- Small datasets (< 300) with simple strings → Uses legacy (no cache overhead)
- Large datasets OR complex patterns → Uses caching + memory pooling  
- ASCII strings → Uses fast byte-level parsing (2-3x faster)
- Unicode strings → Falls back to rune-based parsing

### Advanced Cache Control

```go
// Configure global cache size (default: 2000 entries)
ansort.ConfigureCacheSize(10000)  // Larger cache for big applications
ansort.ConfigureCacheSize(500)    // Smaller cache for memory-constrained

// Disable/enable caching globally
ansort.DisableCache()  // Forces legacy implementation
ansort.EnableCache()   // Re-enables intelligent selection

// Monitor cache effectiveness
hits, misses, ratio := ansort.CacheEfficiencyStats()
fmt.Printf("Cache hit ratio: %.2f%%\n", ratio*100)

// Legacy implementations (for specific needs)
ansort.SortStringsLegacy(data)  // Original implementation
ansort.CompareLegacy(a, b)      // No caching overhead
```

### Performance Monitoring

```go
// Global cache statistics
size, maxSize := ansort.GlobalCacheStats()
fmt.Printf("Cache: %d/%d entries\n", size, maxSize)

// Reset cache statistics
ansort.ResetCacheStats()

// Clear cache (if memory cleanup needed)
ansort.ClearGlobalCache()
```

## Performance Decision Matrix

### When Auto-Selection Uses Legacy Implementation

| Condition | Threshold | Reason |
|-----------|-----------|---------|
| Small dataset | < 300 items | Cache overhead > benefit |
| Few duplicates | < 20% repeated strings | Tokenization not repeated enough |
| Short strings | < 10 avg chars | Simple parsing, cache overhead high |
| Short comparisons | Both strings < 10 chars | Direct parsing faster |

### When Auto-Selection Uses Caching

| Condition | Threshold | Reason |
|-----------|-----------|---------|
| Large dataset | 300+ items | More opportunities for cache hits |
| Many duplicates | 20%+ repeated strings | Significant tokenization savings |
| Long strings | 50+ avg chars | Complex tokenization benefits from caching |
| Long comparisons | Either string 10+ chars | Cache amortizes across operations |

## Performance Tuning

### For Memory-Constrained Environments

```go
// Option 1: Smaller cache
ansort.ConfigureCacheSize(100)

// Option 2: Disable caching entirely
ansort.DisableCache()

// Option 3: Use legacy explicitly
ansort.SortStringsLegacy(data)
```

### For High-Performance Applications

```go
// Option 1: Larger cache for better hit rates
ansort.ConfigureCacheSize(10000)

// Option 2: Pre-warm cache with representative data
representativeData := loadTypicalDataset()
ansort.SortStrings(representativeData) // Populates cache

// Option 3: Monitor and tune cache effectiveness
hits, misses, ratio := ansort.CacheEfficiencyStats()
if ratio < 0.3 { // Low hit rate
    ansort.ConfigureCacheSize(ansort.GlobalCacheStats().maxSize * 2)
}
```

## Debugging Performance Issues

### Analyzing Auto-Selection Decisions

```go
// Test what approach would be selected for your data
data := yourDataset()

// This will show whether caching would be used
fmt.Printf("Dataset size: %d\n", len(data))
// For datasets < 300, check patterns manually:
if len(data) < 300 {
    uniqueCount := countUniqueStrings(data)
    duplicateRatio := 1.0 - float64(uniqueCount)/float64(len(data))
    avgLength := calculateAverageLength(data)
    
    fmt.Printf("Duplicate ratio: %.2f%%\n", duplicateRatio*100)
    fmt.Printf("Average length: %.1f chars\n", avgLength)
    fmt.Printf("Would use caching: %t\n", 
        duplicateRatio > 0.2 || avgLength > 50)
}
```

### Cache Effectiveness Monitoring

```go
// Reset stats before your operation
ansort.ResetCacheStats()

// Perform your sorting operations
for i := 0; i < 1000; i++ {
    ansort.SortStrings(yourData[i])
}

// Check cache effectiveness
hits, misses, ratio := ansort.CacheEfficiencyStats()
fmt.Printf("Cache hits: %d, misses: %d, ratio: %.2f%%\n", 
    hits, misses, ratio*100)

// If ratio < 30%, consider disabling cache or adjusting size
if ratio < 0.3 {
    fmt.Println("Consider disabling cache for this workload")
}
```

### Performance Profiling

```go
import _ "net/http/pprof"

func BenchmarkYourWorkload(b *testing.B) {
    data := yourActualData() // Use your real data patterns
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        testData := make([]string, len(data))
        copy(testData, data)
        ansort.SortStrings(testData)
    }
}

// Compare different approaches:
// go test -bench=BenchmarkYourWorkload
// go test -bench=. -cpuprofile=cpu.prof
// go tool pprof cpu.prof
```

## Internal API Reference

### Auto-Selection Functions

```go
// These are used internally by SortStrings() and Compare()
func shouldUseCaching(data []string) bool
func compareWithoutCache(a, b string, options ...Option) int
func parseStringOptimized(s string) []Token
func parseStringASCII(s string, tokens []Token) []Token
```

### Cache Management

```go
// Global cache control
func ConfigureCacheSize(size int)
func DisableCache()
func EnableCache()
func ClearGlobalCache()

// Statistics and monitoring
func GlobalCacheStats() (size int, maxSize int)
func CacheEfficiencyStats() (hits int64, misses int64, hitRatio float64)
func ResetCacheStats()
```

### Legacy Implementations

```go
// Original implementations (always available)
func SortStringsLegacy(data []string, options ...Option)
func CompareLegacy(a, b string, options ...Option) int

// Explicit optimization variants
func SortStringsOptimized(data []string, options ...Option)
func CompareOptimized(a, b string, options ...Option) int
```

## Best Practices for Developers

### 1. Trust the Auto-Selection

```go
// ✅ Recommended - Let the system optimize automatically
ansort.SortStrings(data)

// ❌ Usually unnecessary - Manual optimization selection
if len(data) > 1000 {
    ansort.SortStringsOptimized(data)
} else {
    ansort.SortStringsLegacy(data)
}
```

### 2. Monitor in Production

```go
// Periodically check cache effectiveness in long-running apps
go func() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        hits, misses, ratio := ansort.CacheEfficiencyStats()
        if ratio < 0.2 { // Very low hit rate
            log.Printf("ansort cache hit ratio low: %.2f%%", ratio*100)
            // Consider adjusting cache size or disabling
        }
        ansort.ResetCacheStats() // Reset for next measurement
    }
}()
```

### 3. Memory Management

```go
// For memory-constrained applications
if memoryConstrained {
    ansort.ConfigureCacheSize(100) // Smaller cache
    // OR
    ansort.DisableCache() // No caching overhead
}

// For high-memory applications
if highMemoryAvailable {
    ansort.ConfigureCacheSize(10000) // Larger cache for better hit rates
}
```

### 4. Testing with Real Data

```go
func TestWithRealData(t *testing.T) {
    // Use your actual data patterns for testing
    realData := loadProductionDataSample()
    
    // Test that auto-selection works correctly
    ansort.SortStrings(realData)
    
    // Verify it's sorted correctly
    for i := 1; i < len(realData); i++ {
        if ansort.Compare(realData[i-1], realData[i]) > 0 {
            t.Errorf("Incorrect sort order at %d", i)
        }
    }
}
```

## Summary

**The Goal Achieved**: *"Users will always expect the performance to be at its best unless there are actual tradeoffs in terms of required memory etc, in which case they want control over when to use the high performance routines."*

The intelligent auto-selection system provides:

- **✅ Best performance automatically** - System chooses optimal approach for your data
- **✅ Zero-configuration optimization** - No API changes required  
- **✅ Smart about trade-offs** - Avoids cache overhead when it doesn't help
- **✅ Adaptive behavior** - Responds to actual data characteristics  
- **✅ Transparent operation** - Same simple API works everywhere
- **✅ Advanced control available** - Power users can override when needed

**Your code gets faster automatically, without any changes required!** 🚀

Developers can trust the system to make good decisions while retaining full control when specialized behavior is needed.
