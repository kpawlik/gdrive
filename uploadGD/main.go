package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kpawlik/gdrive/gdrive" //"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	inputPath  = flag.String("i", "", "input file path")
	outputFile = flag.String("o", "", "output filename")
	folderName = flag.String("f", "", "folder name")
)

func Comma(v int64) string {
	sign := ""
	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

func FileSizeFormat(bytes int64, forceBytes bool) string {
	if forceBytes {
		return fmt.Sprintf("%v B", bytes)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

	var i int
	value := float64(bytes)

	for value > 1000 {
		value /= 1000
		i++
	}
	return fmt.Sprintf("%.1f %s", value, units[i])
}

func MeasureTransferRate() func(int64) string {
	start := time.Now()

	return func(bytes int64) string {
		seconds := int64(time.Now().Sub(start).Seconds())
		if seconds < 1 {
			return fmt.Sprintf("%s/s", FileSizeFormat(bytes, false))
		}
		bps := bytes / seconds
		return fmt.Sprintf("%s/s", FileSizeFormat(bps, false))
	}
}

func getOrCreateFolder(d *drive.Service, folderName string) string {
	folderId := ""
	if folderName == "" {
		return ""
	}
	q := fmt.Sprintf("title=\"%s\" and mimeType=\"application/vnd.google-apps.folder\"", folderName)

	r, err := d.Files.List().Q(q).MaxResults(1).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve foldername.", err)
	}

	if len(r.Items) > 0 {
		folderId = r.Items[0].Id
	} else {
		// no folder found create new
		fmt.Printf("Folder not found. Create new folder : %s\n", folderName)
		f := &drive.File{Title: folderName, Description: "Auto Create by gdrive-upload", MimeType: "application/vnd.google-apps.folder"}
		r, err := d.Files.Insert(f).Do()
		if err != nil {
			fmt.Printf("An error occurred when create folder: %v\n", err)
		}
		folderId = r.Id
	}
	return folderId
}

func uploadFile(d *drive.Service, title string, description string,
	parentName string, mimeType string, filename string) (*drive.File, error) {
	input, err := os.Open(filename)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}
	// Grab file info
	inputInfo, err := input.Stat()
	if err != nil {
		return nil, err
	}

	parentId := getOrCreateFolder(d, parentName)
	fmt.Println("Start upload")
	f := &drive.File{Title: title, Description: description, MimeType: mimeType}
	if parentId != "" {
		p := &drive.ParentReference{Id: parentId}
		f.Parents = []*drive.ParentReference{p}
	}
	getRate := MeasureTransferRate()

	// progress call back
	showProgress := func(current, total int64) {
		fmt.Printf("Uploaded at %s, %s/%s\r", getRate(current), Comma(current), Comma(total))
	}

	r, err := d.Files.Insert(f).ResumableMedia(context.Background(), input, inputInfo.Size(), mimeType).ProgressUpdater(showProgress).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return nil, err
	}

	// Total bytes transferred
	bytes := r.FileSize
	// Print information about uploaded file
	fmt.Printf("Uploaded '%s' at %s, total %s\n", r.Title, getRate(bytes), FileSizeFormat(bytes, false))
	fmt.Printf("Upload Done. ID : %s\n", r.Id)
	return r, nil
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

	fmt.Printf("Read file: %s\n", *inputPath)
	outputTitle := *outputFile
	if outputTitle == "" {
		outputTitle = filepath.Base(*inputPath)
	}
	fmt.Printf("Output name: %s\n", outputTitle)

	ext := filepath.Ext(*inputPath)
	mimeType := "application/octet-stream"
	if ext != "" {
		mimeType = mime.TypeByExtension(ext)
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	fmt.Printf("Mime : %s\n", mimeType)

	uploadFile(srv, outputTitle, "", *folderName, mimeType, *inputPath)
}
