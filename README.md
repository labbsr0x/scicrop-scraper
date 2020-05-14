# Scicrop Scraper

Scraper to perform data extraction from Scicrop connected stations

## How to Run locally

How to run scraper locally:

* `docker-compose up --build`

## Service Configuration

* `SCISCRAPER_CONDUCTOR_PROTOCOL` - Configure the protocol to access conductor web server
* `SCISCRAPER_CONDUCTOR_HOST` - Configure the host to access conductor web server
* `SCISCRAPER_CONDUCTOR_PORT` - Configure the port to access conductor web server
* `SCISCRAPER_API_URI` - Sets the URI to connect on Scicrop AgroDataAPI
* `SCISCRAPER_AGRO_URI` - Sets the URI to connect on AgroDataAPI.
* `SCISCRAPER_SCHEMA_URI` - Sets the URI to connect on SCHEMA API.
* `SCISCRAPER_HTTP_TIMEOUT` - Sets the http timeout IN SECONDS for Scicrop rest integration
* `SCISCRAPER_THREADS_PER_TASK` - Sets the number of threads per task provided on task-list parameter
* `SCISCRAPER_POOLING_INTERVAL` - Sets the interval IN SECONDS between attempts to pool tasks from conductor

## Usage

Regardless of how you run (locally or using docker) you should trigger scrapper task on Conductor using the following request using your **Username**,**Password** and required **Date**:

- Scraps data from scicrop agroDataAPI extracting day reads
```
curl --request POST \
  --url http://localhost:8080/api/workflow \
  --header 'content-type: application/json' \
  --data '{
    "name": "scicrop_date",
    "version": 1,
    "input": {
	"username": "user1",
	"password": "pass1",
	"date": "2019-07-01"
    }
}'
```
