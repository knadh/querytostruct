# querytostruct

An extremely tiny utility for unmarshalling and scanning
querystrings into structs. It supports scanning into int*, float*, string,
and []byte types, along with []int*, []float*, and []\*string slices

## Install

`go get -u github.com/knadh/querytostruct`

## Usage

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/knadh/querytostruct"
)

func main() {
	qs := "name=John+Doe&yes=true&count=42&tag=x&tag=y"
	q, _ := url.ParseQuery(qs)

	type test struct {
		Name  string   `q:"name"`
		Yes   bool     `q:"yes"`
		Count int      `q:"count"`
		Tags  []string `q:"tag"`
	}

	var t test
	fields, err := querytostruct.Unmarshal(q, &t, "q")
	fmt.Println("fields=", fields, "error=", err)
	fmt.Println(t)
}

```

Licensed under the MIT License.
