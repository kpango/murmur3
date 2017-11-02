# murmur3
MurmurHash3 go implementation

# Example
```go
	seed := 1234567
	data := []byte("data")
	h1, h2 := murmur3.New(seed).Sum128(data)
```
