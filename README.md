# Cuckoo Filter Implementation in Go

A basic implementation of a Cuckoo Filter probabilistic data structure for learning purposes.

This readme was generated with chatgpt, after I gave it my half done code.

## What is a Cuckoo Filter?

A Cuckoo Filter is a probabilistic data structure that efficiently tests whether an element is in a set. It's similar to a Bloom filter but supports deletions and has better space efficiency in many cases.

### Key Properties

- **False positives**: Possible (but rare)
- **False negatives**: Never
- **Deletions**: Supported (unlike Bloom filters)
- **Space efficient**: Uses fingerprints instead of storing full items
- **Constant lookup time**: O(1) with at most 2 hash table lookups

## How It Works

1. **Fingerprinting**: Each item is hashed to create a small fingerprint
2. **Dual Hashing**: Each fingerprint can be stored in one of two buckets:
   - `bucket1 = hash1(item) % num_buckets`
   - `bucket2 = bucket1 ⊕ hash2(fingerprint) % num_buckets`
3. **Cuckoo Eviction**: When both buckets are full, items are evicted to make room (like cuckoo birds kicking other eggs out of nests)

## Implementation Details

### Data Structure

- 2D grid: `buckets × slots_per_bucket`
- Each slot stores a fingerprint (int32)
- Empty slots are represented as 0

### Core Functions

- `fingerprint(s)`: Generates fingerprint from string using MurmurHash3
- `bucket1(s)`: Primary bucket location
- `bucket2(fp, i1)`: Secondary bucket using XOR relationship
- `insert()`: Attempts insertion with cuckoo eviction if needed

### Current Status

- ✅ Basic insertion with linear probing
- ✅ Duplicate detection
- ✅ Fingerprint generation
- ⚠️ Eviction logic is simplified (linear vs recursive)

## Usage

```bash
go run main.go
```

The program tests insertion of duplicate strings and shows the internal state after each operation.

## Example Output

```
After Insert...
{
  "text": "Hello World",
  "fingerprint": 123,
  "index_1": 2,
  "index_2": 1,
  "element_1": null,
  "element_2": null
}
```

## Learning Notes

This implementation prioritizes clarity over performance. Key learning points:

- **XOR Magic**: The XOR relationship between buckets allows finding the "other" bucket for any fingerprint without storing additional data
- **Eviction Strategy**: Current implementation uses linear probing; production versions use recursive displacement
- **Load Factor**: Real cuckoo filters can achieve ~95% load factor with proper eviction
- **Fingerprint Size**: Larger fingerprints reduce false positives but increase memory usage

## Limitations

- No deletion support yet
- Simplified eviction (not full cuckoo eviction)
- Fixed grid size (no dynamic resizing)
- No false positive rate optimization

## Dependencies

```
go get github.com/cespare/xxhash/v2
go get github.com/spaolacci/murmur3
```

## References

- [Cuckoo Filter Paper](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf)
- [Probabilistic Data Structures Guide](https://github.com/tylertreat/BoomFilters)

---

_This is a learning project to understand probabilistic data structures and cuckoo hashing mechanics._
