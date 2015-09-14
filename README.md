# Panorma

Simple app to archive the images/videos into proper folder structure.

## Features

* Similarity detection based on pHash and Hamming distance
* Keep a record of the archived images in the database
* EXIF data extraction

## Requirements

* golang >= 1.4
* Postgresql database
* pHash binary in path

## Installation

`go get github.com/jaischeema/panorma`

## Usage

**Import command**

`panorma import -s /source_dir -e /destination_dir -d "dbname=panorma_dev sslmode=disable"`

**Server Command**

*Not yet implemented*

## TODO

* Add tests - *test 'em all*
* Simple server component that could be used with a frontend app
* Upload archived images to flickr or Google Photos
* Generate thumbnails

## License

MIT
