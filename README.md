# gyulabot
https://off---white.com/

(work in progress)

### Instalation:
`go get -u github.com/hah/bot`

### Example:

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

### TODO:
- [x] ATC
- [ ] LOGIN (not sure)
- [ ] CHECKOUT
- [ ] Write docs
- [ ] Write tests (ouch)
- [ ] Refactor
