# undertaker

Get rid of your old and unused PHP code.

## How it works

`undertaker` works by preloading all your classes in a given directory and then utilizes the [krakjoe/tombs](https://github.com/krakjoe/tombs)
extension to keep track of the unused functions over time. It provides CLI and HTTP interface for querying the unused functions.

## Prerequisites

- PHP 7.2+
- Target project managed by `composer`
- PHP files specified for preloading do not have side effects when loaded. For example, defining a class and invoking
a function at the bottom of the file.

## Setup

### Enabling tombs

Make sure you have the `tombs.so` extension loaded and configured to serve the unused functions through TCP. An example
`tombs` configuration file that has TCP listener is:

```ini
zend_extension=tombs.so
tombs.slots=10000
tombs.strings=64M
tombs.socket=tcp://0.0.0.0:12345
tombs.graveyard_format=function
tombs.namespace=Undertaker
```

*Note*: You will also have to have `opcache` enabled as there is a bug in the `tombs` extension which duplicates functions
if it is not enabled.

### Setting up preloader

Put the `docker/fpm/undertaker/undertaker.php` somewhere on the disk where your FPM process resides and make sure it's
readable by the FPM process. Usually FPM master process runs as root, so in most cases this will work out of the box.

### Run undertaker

Run `undertaker` with flags or environment variables pointing to:

- the `tombs` address (`-tombs` or `TOMBS_ADDRESS`, i.e. `localhost:12345`)
- the `php-fpm` address (`-fpm` or `FPM_ADDRESS`, i.e. `localhost:9000`)
- the path to the preload file from the last step (i.e. `/var/http/web/undertaker.php`)

If you want to enable the HTTP server from `undertaker` also run it with `-port` or `HTTP_PORT` env variable.

### Wait for requests and collect unused functions

Sending a request to `undertaker` to `/collect` (or running the command with `-collect` flag) will return the functions
which have not been yet used over the lifetime of the `php-fpm` process.

## Example

This repository comes with an example. Simply run `docker-compose up -d` and you'll have bootstrapped and preloaded
`undertaker` with the PHP repository found in `docker/fpm/project`.

Collecting functions right after `undertaker` has preloaded all the files:

```bash
> curl localhost:8888/collect
Undertaker\Dummy\Calculator\PriceCalculator::allRoomsPrice
Undertaker\Dummy\Calculator\PriceCalculator::total
Undertaker\Dummy\Model\House::__construct
Undertaker\Dummy\Model\House::rooms
Undertaker\Dummy\Model\House::exteriorType
Undertaker\Dummy\Model\Building::rooms
Undertaker\Dummy\Model\InhabitableInterface::exteriorType
```

Collecting after we've invoked the `index.php` file through nginx once:

```bash
> curl localhost:8888/collect
Undertaker\Dummy\Calculator\PriceCalculator::total
Undertaker\Dummy\Model\House::exteriorType
Undertaker\Dummy\Model\Building::rooms
Undertaker\Dummy\Model\InhabitableInterface::exteriorType
```

## Future work

- Collect functions over time in a persistence engine
- Integration tests
- Detect files which might have side effects when loaded
- Running in production...
