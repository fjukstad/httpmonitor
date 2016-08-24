# httpmonitor
Very simple monitoring of HTTP services. Sends an e-mail when the shit hits the
fan. Handles both GET and POST endpoints. 

Made over a couple of beers so use with caution. 

# example
The following example checks that [http://google.com](http://google.com) is up
and running once every minute. If it goes down it will send an e-mail to
`to@example.com` from `from@example.com` with an appropriate message describing
what went wrong. 

```
package main

import (
	"fmt"
	"time"

	"github.com/fjukstad/httpmonitor"
)

func main() {
	s := httpmonitor.NewGetService("google.com", 1*time.Minute)
	m := httpmonitor.Monitor{"to@example.com",
		"from@example.com",
		"from-password",
		[]httpmonitor.Service{s}}

	err := m.Run()
	fmt.Println(err)
}

```
