# ✨ gyulabot
https://off---white.com/

(work in progress)

Note: Since I'm washed, expect breaking changes.

## 📦 Install:
`go get -u github.com/hah/bot`

## ⌨️ Usage:

```go
package main

import (
	"fmt"

	"github.com/hah/bot/modules/offwhite"
)

func main() {
	var product offwhite.Product

	item := offwhite.Item{
		Name:  "SPRAY STRIPES SLIDERS",
		Color: "BLACK WHITE",
		Size:  "44",
		URL:   "",
	}

	if item.URL == "" {
		fmt.Println("no URL provided, searching on the site.")
		product = item.Search()
	} else {
		product = item.Fetch()
	}

	product.ATC()
}

```

## 🔨 Todo:
- [x] ~~ATC~~
- [ ] CLOUDFLARE
- [ ] LOGIN (not sure)
- [ ] CHECKOUT
- [ ] Write docs
- [ ] Write tests (ouch)
- [ ] Refactor

## 🤝 Contributors:
- [Hasan](https://www.github.com/except)