package ansort

import (
	"sort"
	"sync"
)

// TokenCache provides thread-safe caching of tokenized strings
type TokenCache struct {
	mu      sync.RWMutex
	cache   map[string][]Token
	maxSize int
}

// NewTokenCache creates a new token cache with the specified maximum size
func NewTokenCache(maxSize int) *TokenCache {
	if maxSize <= 0 {
		maxSize = 1000 // Default cache size
	}
	return &TokenCache{
		cache:   make(map[string][]Token, maxSize),
		maxSize: maxSize,
	}
}

// Get retrieves tokens from cache, returns nil if not found
func (tc *TokenCache) Get(s string) []Token {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if tokens, exists := tc.cache[s]; exists {
		// Return a copy to prevent modification of cached data
		result := make([]Token, len(tokens))
		copy(result, tokens)
		return result
	}
	return nil
}

// Put stores tokens in cache, evicting oldest entries if necessary
func (tc *TokenCache) Put(s string, tokens []Token) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// Simple eviction: if cache is full, clear it
	// In production, this could be LRU or other strategies
	if len(tc.cache) >= tc.maxSize {
		tc.cache = make(map[string][]Token, tc.maxSize)
	}

	// Store a copy to prevent external modification
	cached := make([]Token, len(tokens))
	copy(cached, tokens)
	tc.cache[s] = cached
}

// Size returns the current cache size
func (tc *TokenCache) Size() int {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return len(tc.cache)
}

// Clear empties the cache
func (tc *TokenCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.cache = make(map[string][]Token, tc.maxSize)
}

// CachedSorter is an optimized sorter that uses token caching
type CachedSorter struct {
	data   []string
	config Config
	cache  *TokenCache
}

// NewCachedSorter creates a new cached sorter with the specified options
func NewCachedSorter(data []string, options ...Option) *CachedSorter {
	config := buildConfig(options...)
	return &CachedSorter{
		data:   data,
		config: config,
		cache:  NewTokenCache(1000), // Default cache size
	}
}

// NewCachedSorterWithCache creates a cached sorter with a custom cache
func NewCachedSorterWithCache(data []string, cache *TokenCache, options ...Option) *CachedSorter {
	config := buildConfig(options...)
	return &CachedSorter{
		data:   data,
		config: config,
		cache:  cache,
	}
}

// Len implements sort.Interface
func (cs *CachedSorter) Len() int {
	return len(cs.data)
}

// Less implements sort.Interface with caching
func (cs *CachedSorter) Less(i, j int) bool {
	return cs.compareWithCache(cs.data[i], cs.data[j]) < 0
}

// Swap implements sort.Interface
func (cs *CachedSorter) Swap(i, j int) {
	cs.data[i], cs.data[j] = cs.data[j], cs.data[i]
}

// compareWithCache performs comparison using cached tokens when possible
func (cs *CachedSorter) compareWithCache(a, b string) int {
	// Handle identical strings quickly
	if a == b {
		return 0
	}

	// Try to get tokens from cache
	tokensA := cs.cache.Get(a)
	if tokensA == nil {
		tokensA = parseStringOptimized(a)
		cs.cache.Put(a, tokensA)
	}

	tokensB := cs.cache.Get(b)
	if tokensB == nil {
		tokensB = parseStringOptimized(b)
		cs.cache.Put(b, tokensB)
	}

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], cs.config)
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

// parseStringOptimized is an optimized version of parseString with reduced allocations
func parseStringOptimized(s string) []Token {
	if len(s) == 0 {
		return []Token{} // Return empty slice to match original parseString behavior
	}

	// Pre-allocate with estimated capacity to reduce reallocations
	// Most strings have 2-4 tokens, so start with capacity of 4
	tokens := make([]Token, 0, 4)

	// Work with byte slice when possible for ASCII strings to avoid rune conversion
	// Fall back to rune slice only when necessary
	if isASCII(s) {
		return parseStringASCII(s, tokens)
	}

	return parseStringUnicode(s, tokens)
}

// isASCII checks if string contains only ASCII characters
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 128 {
			return false
		}
	}
	return true
}

// parseStringASCII optimized ASCII-only parsing
func parseStringASCII(s string, tokens []Token) []Token {
	for i := 0; i < len(s); {
		start := i

		// Check if current character is a digit
		if s[i] >= '0' && s[i] <= '9' {
			// Collect all consecutive digits
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				i++
			}
			tokens = append(tokens, Token{
				Type:  NumericToken,
				Value: s[start:i],
			})
		} else {
			// Collect all consecutive non-digits
			for i < len(s) && (s[i] < '0' || s[i] > '9') {
				i++
			}
			tokens = append(tokens, Token{
				Type:  AlphaToken,
				Value: s[start:i],
			})
		}
	}

	return tokens
}

// parseStringUnicode handles Unicode strings (fallback)
func parseStringUnicode(s string, tokens []Token) []Token {
	runes := []rune(s)

	for i := 0; i < len(runes); {
		start := i

		// Check if current character is a digit
		if runes[i] >= '0' && runes[i] <= '9' {
			// Collect all consecutive digits
			for i < len(runes) && runes[i] >= '0' && runes[i] <= '9' {
				i++
			}
			tokens = append(tokens, Token{
				Type:  NumericToken,
				Value: string(runes[start:i]),
			})
		} else {
			// Collect all consecutive non-digits
			for i < len(runes) && (runes[i] < '0' || runes[i] > '9') {
				i++
			}
			tokens = append(tokens, Token{
				Type:  AlphaToken,
				Value: string(runes[start:i]),
			})
		}
	}

	return tokens
}

// shouldUseCaching determines whether to use caching based on dataset characteristics
func shouldUseCaching(data []string) bool {
	dataSize := len(data)

	// For small to medium datasets (< 300), cache overhead often outweighs benefits
	// unless there are specific patterns that benefit from caching
	if dataSize < 300 {
		// Check for string repetition patterns that would benefit from caching
		uniqueStrings := make(map[string]bool)
		duplicateCount := 0

		// Sample up to 50 strings to avoid expensive analysis
		sampleSize := dataSize
		if sampleSize > 50 {
			sampleSize = 50
		}

		for i := 0; i < sampleSize; i++ {
			if uniqueStrings[data[i]] {
				duplicateCount++
			} else {
				uniqueStrings[data[i]] = true
			}
		}

		// If we see significant duplicate strings in the sample, caching will help
		duplicateRatio := float64(duplicateCount) / float64(sampleSize)
		if duplicateRatio > 0.2 { // More than 20% duplicates
			return true
		}

		// Check average string length - much longer strings benefit more from tokenization caching
		totalLength := 0
		for i := 0; i < sampleSize; i++ {
			totalLength += len(data[i])
		}
		avgLength := totalLength / sampleSize

		// Very long strings (> 50 chars) with complex patterns benefit from caching
		if avgLength > 50 {
			return true
		}

		// For smaller datasets with typical strings, legacy is often faster due to less overhead
		return false
	}

	// For very large datasets (300+), caching benefits depend on repeated operations
	// Since single-shot sorting may not benefit much from caching, be conservative
	return true // Still default to caching for large datasets, but the threshold is higher
}

// SortStringsOptimized provides an optimized sorting function with intelligent caching
func SortStringsOptimized(data []string, options ...Option) {
	if len(data) <= 1 {
		return
	}

	// Check if caching is disabled - fallback to legacy implementation
	if globalCacheDisabled {
		SortStringsLegacy(data, options...)
		return
	}

	// Intelligent auto-selection based on dataset size and characteristics
	// For small datasets, cache overhead may outweigh benefits
	if shouldUseCaching(data) {
		// Use cached sorter for better performance on larger datasets
		sorter := NewCachedSorter(data, options...)
		sort.Sort(sorter)
	} else {
		// Use legacy implementation for small datasets where cache overhead dominates
		SortStringsLegacy(data, options...)
	}
}

// CompareOptimized provides optimized comparison with adaptive caching
var globalCache = NewTokenCache(2000)
var cacheHits, cacheMisses int64
var adaptiveCachingEnabled = true

func CompareOptimized(a, b string, options ...Option) int {
	// Handle identical strings quickly
	if a == b {
		return 0
	}

	// Check if caching is disabled - fallback to legacy implementation
	if globalCacheDisabled {
		return CompareLegacy(a, b, options...)
	}

	// For short strings, cache overhead may not be worth it
	if adaptiveCachingEnabled && len(a) < 10 && len(b) < 10 {
		// Use optimized parsing but skip caching for very short strings
		return compareWithoutCache(a, b, options...)
	}

	config := buildConfig(options...)

	// Try to get tokens from global cache
	tokensA := globalCache.Get(a)
	if tokensA == nil {
		tokensA = parseStringOptimized(a)
		globalCache.Put(a, tokensA)
		cacheMisses++
	} else {
		cacheHits++
	}

	tokensB := globalCache.Get(b)
	if tokensB == nil {
		tokensB = parseStringOptimized(b)
		globalCache.Put(b, tokensB)
		cacheMisses++
	} else {
		cacheHits++
	}

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], config)
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

// compareWithoutCache performs optimized comparison without caching
func compareWithoutCache(a, b string, options ...Option) int {
	config := buildConfig(options...)

	// Use optimized parsing but don't cache results
	tokensA := parseStringOptimized(a)
	tokensB := parseStringOptimized(b)

	// Compare token by token
	minLen := len(tokensA)
	if len(tokensB) < minLen {
		minLen = len(tokensB)
	}

	for i := 0; i < minLen; i++ {
		result := compareTokensWithConfig(tokensA[i], tokensB[i], config)
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

// ClearGlobalCache clears the global comparison cache
func ClearGlobalCache() {
	globalCache.Clear()
}

// GlobalCacheStats returns statistics about the global cache
func GlobalCacheStats() (size int, maxSize int) {
	return globalCache.Size(), globalCache.maxSize
}

// CacheEfficiencyStats returns cache hit/miss statistics
func CacheEfficiencyStats() (hits int64, misses int64, hitRatio float64) {
	total := cacheHits + cacheMisses
	if total == 0 {
		return 0, 0, 0.0
	}
	return cacheHits, cacheMisses, float64(cacheHits) / float64(total)
}

// ResetCacheStats resets the cache hit/miss counters
func ResetCacheStats() {
	cacheHits = 0
	cacheMisses = 0
}
