package imageupload

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

// createGCSClient instantiates a authenicated GCS client.
func createGCSClient() (context.Context, *storage.Client) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	return ctx, client
}

// writeBytes writes a file in memory to a specified GCS bucket.
func writeBytes(data []byte, filename, contentType, bucketName string) error {
	ctx, client := createGCSClient()

	w := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	w.ObjectAttrs.ContentType = contentType

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		log.Fatalln(err)
		return err
	}
	if err := w.Close(); err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

// UploadImage extracts image from request and uploads it to bucket.
func UploadImage(w http.ResponseWriter, r *http.Request) {
	// Extract the request body for further task details.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Recieved Image")

	// Decode string from base64
	data, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Writing image to bucket...")

	writeBytes(data, "odie.png", "image/png", "cloud-functions-talk")
	if err != nil {
		log.Fatalln(err)
	}
}
