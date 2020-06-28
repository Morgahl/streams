package main

import (
	"io"
	"log"

	"github.com/curlymon/streams/convert/hex"
	"github.com/curlymon/streams/sink"
	"github.com/curlymon/streams/source"
	"github.com/curlymon/streams/tap"
)

func main() {
	var stream io.Reader
	var err error
	{
		if stream, err = source.File("./data_file"); err != nil {
			log.Fatal(err)
		}
		if stream, err = tap.File(stream, "./data_file_copy"); err != nil {
			log.Fatal(err)
		}
	}

	{
		stream = hex.NewHexEncodingReader(stream)
		if stream, err = tap.File(stream, "./data_file_to_hex"); err != nil {
			log.Fatal(err)
		}
	}

	{
		stream = hex.NewHexDecodingReader(stream)
		if stream, err = tap.File(stream, "./data_file_from_hex"); err != nil {
			log.Fatal(err)
		}
	}

	{
		_, err = sink.File(stream, "./raw_data_file")
	}

	log.Println(err)
}
