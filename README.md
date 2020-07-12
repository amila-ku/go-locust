# go-locust-client

![Build Status](https://github.com/amila-ku/go-locust-client/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/amila-ku/go-locust-client)](https://goreportcard.com/report/github.com/amila-ku/go-locust-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Description

go-locust-client is a library to control a locust load generation and get statistics written in golang. This is client allows to start,stop a locust load test and ramp up load. This uses locust endpoints to cummunicate with locust.

Currently, go-locust-client requires Go version 1.13 or greater and Locust 0.14 or higher. I will try my best to test with older versions of Go and Locust, but due time constraints, I haven't tested with older versions.

This does not process client stats and presents information as it is.

## Usage 

check example folder

```

package main

import (
	"log"
	lc "github.com/amila-ku/go-locust-client"
)

const (
	 hostURL = "http://localhost:8089"
	 users   = 5
	 hatchRate = 1
)

func main(){
	client, err := lc.New(hostURL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.GenerateLoad(users, hatchRate)
	if err != nil {
		log.Fatal(err)
	}

}

```

* hostUrl : locust endpoiint to connect to, ex : http://locust.loadenv.io:8089

* users : Number of users to simulate

* hatchRate : How many users to be added per second


complete example:

```
package main 

import (
    lc "github.com/amila-ku/locust-client"
)

const (
	 hostURL = "http://localhost:8089"
	 users   = 5
	 hatchRate = 1
)

func main(){
	client, err := lc.New(hostURL)
	if err != nil {
		log.Fatal(err)
	}

	// start load generation
	_, err = client.GenerateLoad(users, hatchRate)
	if err != nil {
		log.Fatal(err)
	}

    // stop load generation
	_, err = client.stopLoad()
	if err != nil {
		log.Fatal(err)
	}

    // get loadtest status
	_, err = client.getStatus()
	if err != nil {
		log.Fatal(err)
	}
}

```

## Authentication

Library go-locust-client will not directly handle authentication since locust does not require authentication. If you have implemented authetication mechanism handle it accordingly. As an example if Oauth2 is used when creating http client configure http.Client to handle authentication for you. 
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

## Contributing

If you are using the client and willing to add new functionality to it, you are welcome.

Also check out [locust-operator](https://github.com/amila-ku/locust-operator) which makes running locust in a distributed setup makes easy.

## License

Open source licensed under the MIT license (see _LICENSE_ file for details).