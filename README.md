# Cities service

[![Circle CI](https://circleci.com/gh/lebedev-yury/cities.svg?style=svg&circle-token=025787958f4452dd681fc6bcab3c52fe66a79598)](https://circleci.com/gh/lebedev-yury/cities)

This service contains basic information about cities, and provides
auto-suggestion based on the city names.

Information about the cities is taken from the [geonames.org dump](http://download.geonames.org/export/dump/).

Alternate names are taken for all supported locales (if given in the
dump), as well as all official names.

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
  "status": "ok",
    "statistics": {
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
  "id": 3164603,
  "name": "Venice",
  "population": 51298,
  "latitude": 45.43713,
  "longitude": 12.33265,
  "timezone": "Europe/Rome",
  "country": {
    "code": "IT",
    "name": "Italy",
    "translations": {
      "be": "Італія",
      "de": "Italien",
      "en": "Italy",
      "ru": "Италия",
      "uk": "Італія"
    }
  }
}
```

#### GET search/cities

Returns 5 matching cities, sorted by population.

If there are several matching citynames for one city, the first one is
taken based on the locale priority.

If there are cities with identical names in the search results, the
country name is added to the city.

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
      "id": 3164603,
      "name": "Venice",
      "population": 51298,
      "latitude": 45.43713,
      "longitude": 12.33265,
      "timezone": "Europe/Rome"
    },
    {
      "id": 4176380,
      "name": "Venice, United States",
      "population": 20748,
      "latitude": 27.09978,
      "longitude": -82.45426,
      "timezone": "America/New_York"
    },
    {
      "id": 4176387,
      "name": "Venice Gardens",
      "population": 7104,
      "latitude": 27.07311,
      "longitude": -82.4076,
      "timezone": "America/New_York"
    }
  ]
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
