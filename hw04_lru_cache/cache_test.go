package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	title string
	items []*Value
	expectedFront *Value
	expectedBack *Value
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestCacheOverloaded(t *testing.T) {
	a, b, c := &Value{"A", 123}, &Value{"B", 124}, &Value{"C", 125}
	check := func (t *testing.T, cache *lruCache, front Value, back Value) {
		require.Equal(t, front, cache.queue.Front().Value, "Check front")
		require.Equal(t, back, cache.queue.Back().Value, "Check back")
	}
	testcases := []TestCase{
		{"Last element replaces first", []*Value{a, b, c}, c, b},
		{"Last element replaces first, but first persists due to recent set action", []*Value{a, b, a, c}, c, a},
		{"Nothing replaced, because set applied multiple times for the same two items", []*Value{a, a, a, b, b, b}, b, a},
	}
	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			testCache := NewCache(2)
			for _, v := range tc.items {
				v.V = rand.Int()
				testCache.Set(v.Key, v.V)
			}
			check(t, testCache.(*lruCache), *tc.expectedFront, *tc.expectedBack)
		})
	}
}

func TestCacheClear(t *testing.T) {
	testCache := NewCache(10)
	testCache.Set("X", 1)
	testCache.Set("Y", 2)
	v, _ := testCache.Get("X")
	require.Equal(t, 1, v)
	v, _ = testCache.Get("Y")
	require.Equal(t, 2, v)
	testCache.Clear()
	v, e := testCache.Get("X")
	require.Equal(t, nil, v)
	require.Equal(t, false, e)
	v, e = testCache.Get("Y")
	require.Equal(t, nil, v)
	require.Equal(t, false, e)
}

func TestAllCacheActionsPerformed(t *testing.T) {
	testCache := NewCache(3)
	testCache.Set("a", 1)
	testCache.Set("b", 2)
	testCache.Set("c", 3)
	element, _ := testCache.Get("a")
	require.Equal(t, 1, element)
	testCache.Set("b", 4)
	element, _ = testCache.Get("b")
	require.Equal(t, 4, element)
	testCache.Get("c")
	testCache.Set("d", 999)
	element, _ = testCache.Get("d")
	require.Equal(t, 999, element)
	_, exists := testCache.Get("a")
	require.Equal(t, false, exists)
	for _, k := range []string{"b", "c", "d"} {
		_, exists = testCache.Get(Key(k))
		require.Equal(t, true, exists)
	}
}