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
	item := offwhite.Item{
		Name:  "SPRAY STRIPES SLIDERS",
		Color: "BLACK WHITE",
		Size:  "44",
		URL:   "",
	}
	if item.URL == "" {
		fmt.Println("no URL provided, searching on the site.")
		item.Search()
	} else {
		item.Fetch()
	}
}
```

### TODO:
- [ ] ATC
- [ ] LOGIN (not sure)
- [ ] CHECKOUT
- [ ] Write docs
- [ ] Write tests (ouch)
