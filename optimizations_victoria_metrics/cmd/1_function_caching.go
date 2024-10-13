package main

import "sync"

/*
	FUNCTION CACHING
1. As an example, some of our scrape targets may generate metric labels with underscores "_", and some of our targets may generate labels with hyphens "-". Relabeling allows us to make this
	consistent, making DB queries easire to write.

2. Relabeling, if defined, happens every time vmagent scrapes metrics from our targets, but as we've seen before, vmagent is likely to see the same metric label many times. That means if we once
	saw foo-bar-baz and changed it to foo_bar_baz, then it's very likely we'll have to do the same transformation on the next scrape as well. In this case, caching the results of the relabeling
	function is likely to reduce CPU usage.

3. The type Transformer contains a sync.Map for thread-safe accest to cached results, and a function transformFunc that will do the actual relabeling.

4. Transformer implements function Transform which we use during during relabeling.

	The Transform function first checks the cache usign the Load function. If a cached result is found, then it returns the results from the cache. Otherwise, it'll call transformFunc to do the
	transformation, store the result in the cache, and return it.

5. Now we can use our "hot path" to make fast ops using cache / a function defined.

6. Function result caching allows us to trade off reduced CPU time for increased memory usage in certain cases. It works best when caching CPU-heavy functions that take a limited amount of
	possible values. Examples of CPU-heavy functions include those that do string or regex matching.

7. Summary
	VM uses function result caching for its relabeling feature, but doesn't use it for caching database queries. In the cases of database queries, the range of possible values is too large and
	it's likely our cache hit rate would be low. As with strings interning, functions results caching works the best if number of cached variants is limited, so we can achieve high cache hit rate.

	
*/

type Transformer struct {
	transformFunc func(s string) string
	cache         sync.Map
}

func (t *Transformer) Transform(s string) string {
	v, ok := t.cache.Load(s)
	if ok {
		// Fast path - the transformed `s` is found in the cache.
		return v.(string)
	}

	// Slow path - transform `s` and store it in the cache.
	sTransformed := t.transformFunc(s)
	t.cache.Store(s, sTransformed)
	return sTransformed
}
