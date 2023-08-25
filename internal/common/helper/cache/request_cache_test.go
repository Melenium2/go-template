package cache

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCache_Set_Should_set_new_values_concurrent(t *testing.T) {
	cache := NewCache[uuid.UUID, []string](100)

	parallels := 3

	var wg sync.WaitGroup
	wg.Add(parallels)

	for i := 0; i < parallels; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				rand, _ := uuid.NewRandom()

				cache.Set(rand, []string{})
			}

			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, 100, len(cache.c))
}

func TestCache_Set_Should_purge_map_tail_if_exceeds_limit(t *testing.T) {
	cache := NewCache[uuid.UUID, string](3)

	var (
		rnd1, _ = uuid.NewRandom()
		rnd2, _ = uuid.NewRandom()
		rnd3, _ = uuid.NewRandom()
		rnd4, _ = uuid.NewRandom()
	)

	for _, key := range []uuid.UUID{rnd1, rnd2, rnd3, rnd4} {
		cache.Set(key, "")
	}

	_, ok := cache.Get(rnd1)
	assert.False(t, ok)

	for _, key := range []uuid.UUID{rnd2, rnd3, rnd4} {
		_, ok := cache.Get(key)
		assert.True(t, ok)
	}
}

func TestCache_Get_Should_concurrent_read_write_to_map(t *testing.T) {
	cache := NewCache[uuid.UUID, []string](100)

	parallels := 10

	var wgWrite sync.WaitGroup
	wgWrite.Add(parallels)

	for i := 0; i < parallels; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				rand, _ := uuid.NewRandom()

				cache.Set(rand, []string{})
			}

			wgWrite.Done()
		}()
	}

	var wgRead sync.WaitGroup
	wgRead.Add(parallels)

	for i := 0; i < parallels; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				rand, _ := uuid.NewRandom()

				cache.Get(rand)
			}
		}()

		wgRead.Done()
	}

	wgWrite.Wait()
	wgRead.Wait()
}
