# murmur3 [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0) [![release](https://img.shields.io/github/release/kpango/murmur3.svg?style=flat-square)](https://github.com/kpango/murmur3/releases/latest) [![CircleCI](https://circleci.com/gh/kpango/murmur3.svg)](https://circleci.com/gh/kpango/murmur3) [![codecov](https://codecov.io/gh/kpango/murmur3/branch/master/graph/badge.svg?style=flat-square)](https://codecov.io/gh/kpango/murmur3) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/df7420e45d5d4656a95f49c85e94cff6)](https://www.codacy.com/app/i.can.feel.gravity/murmur3?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kpango/murmur3&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/kpango/murmur3)](https://goreportcard.com/report/github.com/kpango/murmur3) [![GolangCI](https://golangci.com/badges/github.com/kpango/murmur3.svg?style=flat-square)](https://golangci.com/r/github.com/kpango/murmur3) [![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/kpango/murmur3) [![GoDoc](http://godoc.org/github.com/kpango/murmur3?status.svg)](http://godoc.org/github.com/kpango/murmur3)


MurmurHash3 go implementation

# Example
```go
	seed := 1234567
	data := []byte("data")
	h1, h2 := murmur3.New(seed).Sum128(data)
```
