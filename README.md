# Cities service

This service contains basic information about cities, and provides
auto-suggestion based on the city names.

Information about the cities is taken from the [geonames.org dump](http://download.geonames.org/export/dump/).

Altername names are taken for all supported locales (if given in the
dump), as well as all official names.

## Configuration

| Option             | Description                               |
|--------------------|-------------------------------------------|
| Port               | On which port the server is running       |
| Timeout            | Server timeout, in seconds                |
| CORSOrigins        | The list of CORS origins                  |
| Locales            | The list of locales to support            |
| CitiesFile         | Full filename of the cities dump          |
| AlternateNamesFile | Full filename to the alternate names dump |

## Running

To run the service, you need to:

  * put the compiled binary to some directory
  * create the `config.json` file in the same directory
  * download cities and alternate names dump files
  * run the binary

On the first run the service will create a database and parse the files.
