gcm [![Build Status](https://travis-ci.org/Aorioli/gcm.svg?branch=master)](https://travis-ci.org/Aorioli/gcm)[![Coverage Status](https://coveralls.io/repos/Aorioli/gcm/badge.svg?branch=master&service=github)](https://coveralls.io/github/Aorioli/gcm?branch=master)
===

## NOTE: I no longer maintain this library. Feel free to fork! :)

The Android SDK provides a nice convenience library ([com.google.android.gcm.server](https://github.com/google/gcm/tree/master/client-libraries/java/rest-client/src/com/google/android/gcm/server)) that greatly simplifies the interaction between Java-based application servers and Google's GCM servers. However, Google has not provided much support for application servers implemented in languages other than Java, specifically those written in the Go programming language. The `gcm` package helps to fill in this gap, providing a simple interface for sending GCM messages and automatically retrying requests in case of service unavailability.

Documentation: http://godoc.org/github.com/Aorioli/gcm

Getting Started
---------------

To install gcm, use `go get`:

```bash
go get github.com/Aorioli/gcm
```

Import gcm with the following:

```go
import "github.com/Aorioli/gcm"
```

Sample Usage
------------

Here is a quick sample illustrating how to send a message to the GCM server:

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aorioli/gcm"
)

func main() {
	// Create the message to be sent.
	data := map[string]interface{}{"score": "5x1", "time": "15:10"}
	regIDs := []string{"4", "8", "15", "16", "23", "42"}
	m := new(gcm.Message)
	m.RegistrationIDs = regIDs
	m.Data = data

	// Create a Sender to send the message.
	sender := gcm.NewSender("sample_api_key", 2, time.Minute)

	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(m)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	/* ... */
}
```

## Contributing

Open up an issue describing the bug, enhancement or whatever so it's open for discussion.

Fork the repo, make your changes on a separate branch from master and create a pull request.
