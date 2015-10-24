# Cities service

[![Circle CI](https://circleci.com/gh/lebedev-yury/cities.svg?style=svg&circle-token=025787958f4452dd681fc6bcab3c52fe66a79598)](https://circleci.com/gh/lebedev-yury/cities)

This service contains basic information about cities, and provides
auto-suggestion based on the city names.

Information about the cities is taken from the [geonames.org dump](http://download.geonames.org/export/dump/).

Alternate names are taken for all supported locales (if given in the
dump), as well as all official names.

This application uses the [Gin](https://github.com/gin-gonic/gin) framework, and [Bolt DB](https://github.com/boltdb/bolt).

## Configuration

| Option             | Description                          | Default                 |
|--------------------|--------------------------------------|-------------------------|
| Port               | On which port the server is running  | 8080                    |
| Timeout            | Server timeout, in seconds           | 5                       |
| CORSOrigins        | The list of CORS origins             | http://localhost        |
| Locales            | The list of locales to support       | en                      |
| MinPopulation      | Lower limit for the population       | 2000                    |
| CountriesFile      | Filename to the countries dump       | data/countryInfo.txt    |
| CitiesFile         | Filename to the cities dump          | data/cities.txt         |
| AlternateNamesFile | Filename to the alternate names dump | data/alternateNames.txt |

## API

#### GET application/status

Returns the application status and basic statistics.

**Resourse URL:**

`/1.0/application/status`

**Example result:**

```json
{
  "meta": {
    "status": "ok",
    "countries_count": 250,
    "cities_count": 145314,
    "city_names_count": 224695
  }
}
```

#### GET cities/:id

Returns a single city by the requested ID parameter.

**Resourse URL:**

`/1.0/cities/:id`

**Parameters:**

| Parameter | Requirement | Description |
| :--- | :--- | :--- |
| id | required | The numerical ID of the desired city.

**Example result:**

```json
{
  "data": [
    {
      "type": "cities",
      "id": 3164603,
      "attributes": {
        "name": "Venice",
        "country_code": "IT",
        "population": 51298,
        "timezone": "Europe/Rome",
        "latitude": 45.43713,
        "longitude": 12.33265
      },
      "links": {
        "self": "http://localhost:8082/1.0/cities/3164603"
      }
    }
  ]
}
```

#### GET cities

Returns 5 matching cities, sorted by population.

If there are several matching citynames for one city, the first one is
taken based on the locale priority.

If there are cities with identical names in the search results, the
country name is added to the city.

**Resourse URL:**

`/1.0/cities`

**Parameters:**

| Parameter | Requirement | Description |
| :--- | :--- | :--- |
| q | required | A string query.

**Example result:**

```json
{
  "data": [{
    "type": "cities",
    "id": 3164603,
    "attributes": {
      "name": "Venice",
      "population": 51298,
      "timezone": "Europe/Rome",
      "latitude": 45.43713,
      "longitude": 12.33265
    },
    "links": {
      "self": "http://localhost:8082/1.0/cities/3164603"
    }
  },
  {
    "type": "cities",
    "id": 4176380,
    "attributes": {
      "name": "Venice, United States",
      "population": 20748,
      "timezone": "America/New_York",
      "latitude": 27.09978,
      "longitude": -82.45426
    },
    "links": {
      "self": "http://localhost:8082/1.0/cities/4176380"
    }
  },
  {
    "type": "cities",
    "id": 4176387,
    "attributes": {
      "name": "Venice Gardens",
      "population": 7104,
      "timezone": "America/New_York",
      "latitude": 27.07311,
      "longitude": -82.4076
    },
    "links": {
      "self": "http://localhost:8082/1.0/cities/4176387"
    }
  }]
}
```

## Make commands

To install dependencies:

```
make setup
```

To get the dumpfiles:

```
make getdumpfiles
```

Before committing:

```
make test
```

To build and run the docker instance:

```
make dockerrun
```

## License

Copyright (c) 2015 Yury Lebedev

This project is released under the [MIT License](http://www.opensource.org/licenses/MIT).
