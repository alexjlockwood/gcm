go-gcm
======

The Android SDK provides a nice convenience library ([com.google.android.gcm.server](http://developer.android.com/reference/com/google/android/gcm/server/package-summary.html)) that greatly simplifies the interaction between Java-based application servers and Google's GCM servers. However, Google has not provided much support for application servers implemented in languages other than Java, specifically those written in the Go programming language. go-gcm helps to fill in this gap. This library provides a simple interface for sending GCM messages and automatically retries requests in case of service unavailability using exponential backoff.

Documentation: http://godoc.org/github.com/alexjlockwood/go-gcm

Getting Started
---------------

To install go-gcm, use `go get`:

    go get github.com/alexjlockwood/go-gcm

Import go-gcm with the following:

    import "github.com/alexjlockwood/go-gcm/gcm"

Sample Usage
------------

Here is a quick sample illustrating how to send a message to the GCM server:

    package sample
    
    import (
        "fmt"
        "net/http"
        "github.com/alexjlockwood/go-gcm/gcm"
    )
    
    func main() {
        // Create the message to be sent
        regIds := []string{"4","8","15","16","23","32"}
        data := map[string]string{"score": "5x1", "time": "15:10"}
        msg := gcm.NewMessage(data, regIds...)

        // Create a Sender to send the message
        sender := &gcm.NewSender("sample_api_key")
        
        // Send the message and receive the response. If the results indicate
        // a service unavailibility error (i.e. if one or more of the result's 
        // Error field is "Unavailable"), the message will be resent at most
        // two times. Note that the message will be retried using exponential 
        // backoff, and thus may block for several seconds.
        response, err := sender.Send(msg, 2)
        if err != nil {
            fmt.Println("Failed to send message: " + err.Error())
            return       
        }
    }

Note for Google AppEngine users
-------------------------------

If your application server runs on Google AppEngine, you must import the `appengine/urlfetch` package and create the `Sender` as follows:

    import (
        "appengine"
        "appengine/urlfetch"
    )

    func handler(w http.ResponseWriter, r *http.Request) {
        /* ... */

        c := appengine.NewContext(r)
        client := urlfetch.Client(c)
        sender := gcm.NewSender("sample_api_key", client)

        /* ... */
    }        
