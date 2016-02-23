package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	vision "google.golang.org/api/vision/v1"

	"github.com/kaneshin/pigeon"
)

var (
	faceDetection       = flag.Bool("face", false, "This flag specifies the face detection of the feature")
	landmarkDetection   = flag.Bool("landmark", false, "This flag specifies the landmark detection of the feature")
	logoDetection       = flag.Bool("logo", false, "This flag specifies the logo detection of the feature")
	labelDetection      = flag.Bool("label", false, "This flag specifies the label detection of the feature")
	textDetection       = flag.Bool("text", false, "This flag specifies the text detection (OCR) of the feature")
	safeSearchDetection = flag.Bool("safe-search", false, "This flag specifies the safe-search of the feature")
	imageProperties     = flag.Bool("image-properties", false, "This flag specifies the image safe-search properties of the feature")
)

var defaultDetection = pigeon.LabelDetection

func features() []*vision.Feature {
	list := []int{}
	if *faceDetection {
		list = append(list, pigeon.FaceDetection)
	}
	if *landmarkDetection {
		list = append(list, pigeon.LandmarkDetection)
	}
	if *logoDetection {
		list = append(list, pigeon.LogoDetection)
	}
	if *labelDetection {
		list = append(list, pigeon.LabelDetection)
	}
	if *textDetection {
		list = append(list, pigeon.TextDetection)
	}
	if *safeSearchDetection {
		list = append(list, pigeon.SafeSearchDetection)
	}
	if *imageProperties {
		list = append(list, pigeon.ImageProperties)
	}

	// Default
	if len(list) == 0 {
		list = append(list, defaultDetection)
	}

	features := make([]*vision.Feature, len(list))
	for i := 0; i < len(list); i++ {
		features[i] = pigeon.NewFeature(list[i])
	}
	return features
}

func main() {
	// Parse arguments to run this function.
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	if args := flag.Args(); len(args) == 0 {
		fmt.Fprintf(os.Stderr, "pigeon [options] <source>\n")
		os.Exit(1)
	}

	// Initialize vision service by a credentials json.
	client, err := pigeon.New()
	if err != nil {
		log.Fatalf("Unable to retrieve vision service: %v\n", err)
	}

	// To call multiple image annotation requests.
	batch, err := pigeon.NewBatchAnnotateImageRequest(flag.Args(), features()...)
	if err != nil {
		log.Fatalf("Unable to retrieve image request: %v\n", err)
	}

	// Execute the "vision.images.annotate".
	res, err := client.ImagesService().Annotate(batch).Do()
	if err != nil {
		log.Fatalf("Unable to execute images annotate requests: %v\n", err)
	}

	// Marshal annotations from responses
	body, err := json.MarshalIndent(res.Responses, "", "  ")
	if err != nil {
		log.Fatalf("Unable to marshal the response: %v\n", err)
	}
	fmt.Println(string(body))
}
