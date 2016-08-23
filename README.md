# httpmonitor
Very simple monitoring of HTTP services. Sends an e-mail when the shit hits the
fan. Handles both GET and POST endpoints. 

Made over a couple of beers so use with caution. 

# Example

```
package main

import (
	"fmt"
	"time"

	"github.com/fjukstad/httpmonitor"
)

func main() {
	s := httpmonitor.Service{"google.com", "GET", "", "", 1 *
		time.Minute, time.Time{}}
	m := httpmonitor.Monitor{"to@example.com", "from@example.com",
		"from-password", []httpmonitor.Service{s}}

	err := m.Run()
	fmt.Println(err)
}
```
