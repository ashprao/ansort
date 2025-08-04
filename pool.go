package ansort

import (
	"sort"
	"sync"
)

// TokenPool manages a pool of Token slices to reduce allocations
type TokenPool struct {
	pool sync.Pool
}

// NewTokenPool creates a new token pool
func NewTokenPool() *TokenPool {
	return &TokenPool{
		pool: sync.Pool{
			New: func() interface{} {
				// Pre-allocate with capacity 8, which covers most common cases
				return make([]Token, 0, 8)
			},
		},
	}
}

// Get retrieves a token slice from the pool
func (tp *TokenPool) Get() []Token {
	tokens := tp.pool.Get().([]Token)
	return tokens[:0] // Reset length but keep capacity
}

// Put returns a token slice to the pool
func (tp *TokenPool) Put(tokens []Token) {
	// Only return to pool if capacity is reasonable to avoid memory waste
	if cap(tokens) <= 32 {
		tp.pool.Put(tokens)
	}
}

// Global token pool for reuse across functions
var globalTokenPool = NewTokenPool()

// parseStringPooled uses the global token pool to reduce allocations
func parseStringPooled(s string) []Token {
	if len(s) == 0 {
		return []Token{}
	}

	// Get a token slice from the pool
	tokens := globalTokenPool.Get()
	defer globalTokenPool.Put(tokens)

	// Parse using the pooled slice
	if isASCII(s) {
		tokens = parseStringASCII(s, tokens)
	} else {
		tokens = parseStringUnicode(s, tokens)
	}

	// Return a copy since we're returning the pooled slice
	result := make([]Token, len(tokens))
	copy(result, tokens)
	return result
}

// PooledSorter uses memory pools for optimal performance
type PooledSorter struct {
	data      []string
	config    Config
	cache     *TokenCache
	tokenPool *TokenPool
}

// NewPooledSorter creates a new pooled sorter
func NewPooledSorter(data []string, options ...Option) *PooledSorter {
	config := buildConfig(options...)
	return &PooledSorter{
		data:      data,
		config:    config,
		cache:     NewTokenCache(1000),
		tokenPool: NewTokenPool(),
	}
}

// Len implements sort.Interface
func (ps *PooledSorter) Len() int {
	return len(ps.data)
}

// Less implements sort.Interface with pooling and caching
func (ps *PooledSorter) Less(i, j int) bool {
	return ps.comparePooled(ps.data[i], ps.data[j]) < 0
}

// Swap implements sort.Interface
func (ps *PooledSorter) Swap(i, j int) {
	ps.data[i], ps.data[j] = ps.data[j], ps.data[i]
}

// comparePooled performs comparison using both cache and memory pools
func (ps *PooledSorter) comparePooled(a, b string) int {
	// Handle identical strings quickly
	if a == b {
		return 0
	}

	// Try to get tokens from cache
	tokensA := ps.cache.Get(a)
	if tokensA == nil {
		tokensA = ps.parseWithPool(a)
		ps.cache.Put(a, tokensA)
	}

	tokensB := ps.cache.Get(b)
	if tokensB == nil {
		tokensB = ps.parseWithPool(b)
		ps.cache.Put(b, tokensB)
	}

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], ps.config)
		if result != 0 {
			return result
		}
	}

	// If all compared tokens are equal, the shorter string comes first
	if len(tokensA) < len(tokensB) {
		return -1
	} else if len(tokensA) > len(tokensB) {
		return 1
	}

	return 0
}

// parseWithPool uses the sorter's token pool for parsing
func (ps *PooledSorter) parseWithPool(s string) []Token {
	if len(s) == 0 {
		return nil
	}

	// Get a token slice from the pool
	tokens := ps.tokenPool.Get()
	defer ps.tokenPool.Put(tokens)

	// Parse using the pooled slice
	if isASCII(s) {
		tokens = parseStringASCII(s, tokens)
	} else {
		tokens = parseStringUnicode(s, tokens)
	}

	// Return a copy since we're returning the pooled slice
	result := make([]Token, len(tokens))
	copy(result, tokens)
	return result
}

// SortStringsPooled provides the most optimized sorting with both caching and pooling
func SortStringsPooled(data []string, options ...Option) {
	if len(data) <= 1 {
		return
	}

	// Use pooled sorter for maximum performance
	sorter := NewPooledSorter(data, options...)

	// Use Go's optimized sort algorithm
	sort.Sort(sorter)
}

// HighPerformanceSorter combines all optimizations for maximum speed
type HighPerformanceSorter struct {
	*PooledSorter
}

// NewHighPerformanceSorter creates the fastest possible sorter
func NewHighPerformanceSorter(data []string, options ...Option) *HighPerformanceSorter {
	pooledSorter := NewPooledSorter(data, options...)

	// Use a larger cache for high-performance scenarios
	pooledSorter.cache = NewTokenCache(5000)

	return &HighPerformanceSorter{
		PooledSorter: pooledSorter,
	}
}

// SortStringsHighPerformance provides the absolute fastest sorting
func SortStringsHighPerformance(data []string, options ...Option) {
	if len(data) <= 1 {
		return
	}

	sorter := NewHighPerformanceSorter(data, options...)
	sort.Sort(sorter)
}
