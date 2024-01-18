package patterns

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

func ShardingDemo() {
	fmt.Println("Sharding Pattern Demo...")

	shardedMap := NewShardedMap(5)

	shardedMap.Set("alpha", 1)
	shardedMap.Set("beta", 2)
	shardedMap.Set("gamma", 3)
	shardedMap.Set("delta", 4)
	shardedMap.Set("epsilon", 5)
	shardedMap.Set("zeta", 6)

	fmt.Println(shardedMap.Get("alpha"))
	fmt.Println(shardedMap.Get("beta"))
	fmt.Println(shardedMap.Get("gamma"))
	fmt.Println(shardedMap.Get("delta"))
	fmt.Println(shardedMap.Get("epsilon"))
	fmt.Println(shardedMap.Get("zeta"))

	keys := shardedMap.Keys()
	for _, k := range keys {
		fmt.Println(k)
	}
}

type Shard struct {
	sync.RWMutex                        // Compose from sync.RWMutex
	m            map[string]interface{} // m contains the shard's data
}

type ShardedMap []*Shard

// NewShardedMap create a new ShardedMap, which make use of Vertical Sharding to
// reduce lock contention by splitting the underlying data structure (usually a
// map) into several individually lockable maps. An abstracion layer provides
// access to the underlying shards as if they were a single structure.
func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards)

	for i := 0; i < nshards; i++ {
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}

	return shards
}

// getShardedIndex since a byte-sized value it's being used as the hash value,
// it can only handle up to 255 shards. If for some reason it's necessary morce
// than that, then its possible to accomplish it using a bit of byte arithmetic.
// E.g., `hash := int(sum[13]) << 8 | int(sum[17])`
func (m ShardedMap) getShardedIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17]) // Pick an arbitrary byte as the hash

	return hash % len(m)
}

func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardedIndex(key)

	return m[index]
}

func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key string, v interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = v
}

func (m ShardedMap) Keys() []string {

	var mu sync.Mutex
	keys := make([]string, 0)

	wg := sync.WaitGroup{} // Create a wait group and add a
	wg.Add(len(m))         // wait value for each slice

	for _, shard := range m { // Run a goroutine for each shard
		go func(s *Shard) {
			s.RLock() // Establish a read lock on shard

			for key := range s.m {
				mu.Lock()
				keys = append(keys, key)
				mu.Unlock()
			}

			s.RUnlock() // Release the read lock
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return keys
}
