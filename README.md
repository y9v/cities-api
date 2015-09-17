# Cities service

[![Circle CI](https://circleci.com/gh/lebedev-yury/cities.svg?style=svg&circle-token=025787958f4452dd681fc6bcab3c52fe66a79598)](https://circleci.com/gh/lebedev-yury/cities)

This service contains basic information about cities, and provides
auto-suggestion based on the city names.

Information about the cities is taken from the [geonames.org dump](http://download.geonames.org/export/dump/).

Alternate names are taken for all supported locales (if given in the
dump), as well as all official names.

## Configuration

| Option             | Description                               | Default            |
|--------------------|-------------------------------------------|--------------------|
| Port               | On which port the server is running       | 8080               |
| Timeout            | Server timeout, in seconds                | 5                  |
| CORSOrigins        | The list of CORS origins                  | http://localhost   |
| Locales            | The list of locales to support            | en                 |
| CitiesFile         | Full filename of the cities dump          | data/cities.txt    |
| AlternateNamesFile | Full filename to the alternate names dump | data/alternate.txt |

## API

#### GET application/status

Returns the application status and basic statistics.

**Resourse URL:**

`/1.0/application/status`

**Example result:**

```json
{
  "status": "ok",
    "statistics": {
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
  "id": "3164603",
  "name": "Venice",
  "country_code": "IT",
  "population": "51298",
  "latitude": "45.43713",
  "longitude": "12.33265",
  "timezone": "Europe/Rome"
}
```

#### GET search/cities

Returns 5 matching cities, sorted by population.

If there are several matching citynames for one city, the first one is
taken based on the locale priority.

**Resourse URL:**

`/1.0/search/cities`

**Parameters:**

| Parameter | Requirement | Description |
| :--- | :--- | :--- |
| q | required | A string query.

**Example result:**

```json
{
  "cities": [
    {
      "id": "3827407",
      "name": "Venustiano Carranza",
      "country_code": "MX",
      "population": "447459",
      "latitude": "19.44286",
      "longitude": "-99.09724",
      "timezone": "America/Mexico_City"
    },
    {
      "id": "5405878",
      "name": "Ventura",
      "country_code": "US",
      "population": "96769",
      "latitude": "34.27834",
      "longitude": "-119.29317",
      "timezone": "America/Los_Angeles"
    },
    {
      "id": "2745641",
      "name": "Venlo",
      "country_code": "NL",
      "population": "92403",
      "latitude": "51.37",
      "longitude": "6.16806",
      "timezone": "Europe/Amsterdam"
    },
    {
      "id": "3833062",
      "name": "Venado Tuerto",
      "country_code": "AR",
      "population": "72340",
      "latitude": "-33.74556",
      "longitude": "-61.96885",
      "timezone": "America/Argentina/Cordoba"
    },
    {
      "id": "3164603",
      "name": "Venice",
      "country_code": "IT",
      "population": "51298",
      "latitude": "45.43713",
      "longitude": "12.33265",
      "timezone": "Europe/Rome"
    }
  ]
}
```

## Development

This project uses [gom](https://github.com/mattn/gom) dependency manager.

After installing gom you need to install the dependencies:

```
gom -test install
```

Before committing:

```
go vet ./...
gom test -cover ./...
```
