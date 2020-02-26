# go-locust-client

![Build Status](https://github.com/amila-ku/go-locust-client/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/amila-ku/go-locust-client)](https://goreportcard.com/report/github.com/amila-ku/go-locust-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

go-locust-client is a library to control a locust deployment and get statistics. This is client allows to start/stop a locust load test and ramp up load test.

Currently, go-locust-client requires Go version 1.13 or greater and Locust 0.14 or higher. I will try my best to test with older versions of Go and Locust, but due time constraints, I haven't tested with older versions.

This does not process client stats and presents information as it is.

## Usage 

```

	client, err := lc.New(HostUrl)
	if err != nil {
		return err
	}
	_, err = client.StartLoad(Users, HatchRate)
	if err != nil {
		return err
	}

```

* HostUrl : locust endpoiint to connect to, ex : http://locust.loadenv.io:8089

* Users : Number of users to simulate

* HatchRate : How many users to be added per second


complete example:

```
package main 

import (
    lc "github.com/amila-ku/locust-client"
)

func main(){

    // create new client and check for errors
	client, err := lc.New(HostUrl)
	if err != nil {
		return err
	}

    // start load generation
	_, err = client.startLoad(Users, HatchRate)
	if err != nil {
		return err
	}

    // stop load generation
	_, err = client.stopLoad()
	if err != nil {
		return err
	}

    // get loadtest status
	_, err = client.getStatus()
	if err != nil {
		return err
	}
}

```

## Authintication

Library go-locust-client will not directly handle authentication. When creating http client configure http.Client to handle authentication for you. 
Check https://github.com/golang/oauth2 for implementation


## Run Locust

install locust 

```

pip3 install locust

```

and configure loadtest file as described here https://docs.locust.io/en/stable/quickstart.html


```

locust -f locust_files/my_locust_file.py

```

## To Do

* add example