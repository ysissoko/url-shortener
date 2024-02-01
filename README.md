# url-shortener
A url shortener developed with Golang to generate short url from long URLs.

## Prerequisites

You need to install docker and docker-compose to launch the project.

## Run the project

To run the project, go to the root folder and launch the following command:
`docker-compose up -d --build`
This command will fetch the docker images, build the containers (the webserver will be built from the Dockerfile in the root folder) and the redis database container from the redis official image.

## Usage

The http request payload should contain the following fields:

- url: The long URL to be shortened.
- ttl: The "time to live" attribute correspond to the maximum time the shortened URL will be available (set to 0 for no limit).

An example usage sending a post request to generate a uuid for the provided URL:

`curl -X POST http://localhost:8080/generate -H "Content-Type: application/json" -d '{"url": "https://google.com", "ttl": 60}'`

The response will be a text plain uuid corresponding to the shortened URL.

To use this uuid, you need to browse the following URL:

http://localhost:8080/<uuid>

You will be redirected to the associated URL.
