package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var bucket = flag.String("bucket", "bucket-name", "the S3 bucket")
var accessKeyID = flag.String("access-key-id", "xxxxxx", "the S3 access key id")
var secretAccessKey = flag.String("secret-access-key", "xxxxxx", "the S3 secret access key")
var region = flag.String("region", "us-west-1", "the S3 secret access key")
var concurrency = flag.Int("concurrency", 10, "the number of downloads at the same time")

func main() {
	flag.Parse()
	err := downloadAll(*bucket, *region, *accessKeyID, *secretAccessKey, *concurrency)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadAll(bucketName, regionName, accessKeyID, secretAccessKey string, concurrency int) error {
	auth, err := aws.GetAuth(accessKeyID, secretAccessKey)
	if err != nil {
		return err
	}

	region, ok := aws.Regions[regionName]
	if !ok {
		return errors.New("Invalid region")
	}

	client := s3.New(auth, region)
	bucket := client.Bucket(bucketName)

	log.Println("Getting bucket contents")
	content, err := bucket.GetBucketContents()
	if err != nil {
		return err
	}

	tempDir, err := ioutil.TempDir(".", "s3_backup_")

	if err != nil {
		return err
	}

	keys := make(chan string, concurrency)
	wg := &sync.WaitGroup{}

	log.Printf("Starting download with concurrency of: %d", concurrency)
	for i := 0; i < concurrency; i++ {
		go downloadFileFromS3(*bucket, keys, tempDir, wg)
	}

	for key, _ := range *content {
		wg.Add(1)
		log.Printf("Enqueuing: %s", key)
		keys <- key
	}

	wg.Wait()

	return nil
}

func downloadFileFromS3(bucket s3.Bucket, keys chan string, tempDir string, wg *sync.WaitGroup) {
	for {
		key := <-keys
		newDir := fmt.Sprintf("%s/%s", tempDir, filepath.Dir(key))
		os.MkdirAll(newDir, 0755)

		log.Printf("Downloading: %s", key)

		reader, err := bucket.GetReader(key)
		if err != nil {
			log.Printf("Failed to download: %s, %s", key, err)
			wg.Done()
			continue
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Printf("Failed to download: %s, %s", key, err)
			wg.Done()
			continue
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", tempDir, key), data, 0644)
		if err != nil {
			log.Printf("Failed to download: %s, %s", key, err)
			wg.Done()
			continue
		}
		wg.Done()
	}
}
