package main

import (
	"encoding/json"
	"fmt"

	"github.com/cespare/xxhash/v2"
	"github.com/spaolacci/murmur3"
)

func first_or_default[T any](slice []T, filter func(*T) bool) (element *T) {
	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

func where[T any](slice []T, filter func(*T) bool) []*T {
	var ret []*T = make([]*T, 0)
	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			ret = append(ret, &slice[i])
		}
	}
	return ret
}

func insert(grid [][]int32, fingerprint int32, index1, index2 uint32, maxKicks int, number_of_slots int) bool {
	for i := 0; i < number_of_slots; i++ {
		if grid[index1][i] == 0 {
			grid[index1][i] = fingerprint
			return true
		}
		if grid[index2][i] == 0 {
			grid[index2][i] = fingerprint
			return true
		}
	}

	idx := index1
	keepKicking := true
	kickCount := 0
	for keepKicking {
		kickCount++

		slot := kickCount % number_of_slots
		evicted := grid[idx][slot]
		grid[idx][slot] = fingerprint

		fmt.Println("Evicted...")
		// Cuckoo kick-out loop (optional, capped by maxKicks)
		fmt.Println(evicted)

		if kickCount >= maxKicks {
			keepKicking = false
		}
	}

	return false // filter is full
}

func pretty_print(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(b))
}

func fingerprint(s string) int32 {
	h := murmur3.Sum32WithSeed([]byte(s), 42)
	return int32(h & 0xff)
}

func bucket1(s string, numBuckets uint32) uint32 {
	return murmur3.Sum32WithSeed([]byte(s), 1) % numBuckets
}

func bucket2(fp int32, i1 uint32, numBuckets uint32) uint32 {
	var b [1]byte
	b[0] = byte(fp)
	h := xxhash.Sum64(b[:])
	return (i1 ^ uint32(h)) % numBuckets
}

// Cuckoo Filter Basic Implementation
func main() {
	incoming_ids := []string{
		"Hello World",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
		"Hello World1",
	}

	// i can think of this as hashmap of of key an arrays....
	number_of_buckets := uint32(4) // rows
	number_of_slots := uint32(4)   // cols
	grid := make([][]int32, number_of_buckets)
	for i := range grid {
		grid[i] = make([]int32, number_of_slots)
	}

	for _, text := range incoming_ids {
		fingerprint := fingerprint(text)
		index_1 := bucket1(text, number_of_buckets)
		index_2 := bucket2(fingerprint, index_1, number_of_buckets)

		bucket_1 := grid[index_1]
		bucket_2 := grid[index_2]

		element_1 := first_or_default(bucket_1, func(v *int32) bool { return *v == fingerprint })
		element_2 := first_or_default(bucket_2, func(v *int32) bool { return *v == fingerprint })

		data := map[string]interface{}{
			"text":        text,
			"fingerprint": fingerprint,
			"index_1":     index_1,
			"index_2":     index_2,
			"element_1":   element_1,
			"element_2":   element_2,
		}

		success := insert(grid, fingerprint, index_1, index_2, 500, int(number_of_slots))
		if !success {
			fmt.Println("Failed to insert", text)
		} else {
			fmt.Println("After Insert...")
			pretty_print(data)
		}
	}
}
