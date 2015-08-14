package main

import (
	"fmt"
	// "net/url"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/relaxgo/go-cx-extractor"
)

type Options struct {
	FilePath string `short:"f" long:"file" description:"input file path"`

	Url string `short:"u" long:"url" description:"url"`
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
	fmt.Println(options)
	if f := options.FilePath; f != "" {
		data, err := File_get_contents(f)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
			return
		}
		html := string(data)
		// fmt.Println(html)
		fmt.Println("title", extractor.ExtactTitle(html))
		// fmt.Println("html", extractor.ExtractText(html))
		return
	}

	if u := options.Url; u != "" {
		response, err := http.Get(u)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			html := string(contents)
			fmt.Println("title", extractor.ExtactTitle(html))
			fmt.Println("html", extractor.ExtractText(html))
		}

		return
	}
	fmt.Println("must need file or url")
}

// Get bytes to file.
// if non-exist, create this file.
func File_get_contents(filename string) (data []byte, e error) {
	f, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if e != nil {
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		return
	}
	data = make([]byte, stat.Size())
	result, e := f.Read(data)
	if e != nil || int64(result) != stat.Size() {
		return nil, e
	}
	return
}
