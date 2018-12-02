package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/kpawlik/gdrive/gdrive"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	fileID     string //= flag.String("id", "", "id")
	outputFile string //= flag.String("o", "", "Output filename/path if empty original name will be used")
	outputDir  string
)

func init() {
	flag.StringVar(&fileID, "id", "", "file id")
	flag.StringVar(&outputFile, "f", "", "Output filename/path if empty original name will be used")
	flag.StringVar(&outputDir, "d", "", "Path to directory where store file")
	flag.Parse()
}

func downloadFile(d *drive.Service, filename string) error {
	var (
		err   error
		input *os.File
	)

	sf := d.Files.Get(fileID)

	if filename == "" {
		f, err := sf.Do()
		if err != nil {
			return err
		}
		filename = f.OriginalFilename
	}

	r, err := sf.Download()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}

	input, err = os.Create(path.Join(outputDir, filename))
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	f := bufio.NewWriter(input)
	io.Copy(f, r.Body)
	f.Flush()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := gdrive.GetClient(ctx, config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	downloadFile(srv, outputFile)

	/*
		r, err := srv.Files.List().MaxResults(10).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve files.", err)
		}

		fmt.Println("Files:")
		if len(r.Items) > 0 {
			for _, i := range r.Items {
				fmt.Printf("%s (%s)\n", i.Title, i.Id)
			}
		} else {
			fmt.Print("No files found.")
		}
	*/
}
