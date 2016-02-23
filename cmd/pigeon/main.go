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
	FaceDetection       = flag.Bool("face", false, "")
	LandmarkDetection   = flag.Bool("landmark", false, "")
	LogoDetection       = flag.Bool("logo", false, "")
	LabelDetection      = flag.Bool("label", false, "")
	TextDetection       = flag.Bool("text", false, "")
	SafeSearchDetection = flag.Bool("safe-search", false, "")
	ImageProperties     = flag.Bool("image-properties", false, "")
)

func features() []*vision.Feature {
	list := []int{}
	if *FaceDetection {
		list = append(list, pigeon.FaceDetection)
	}
	if *LandmarkDetection {
		list = append(list, pigeon.LandmarkDetection)
	}
	if *LogoDetection {
		list = append(list, pigeon.LogoDetection)
	}
	if *LabelDetection {
		list = append(list, pigeon.LabelDetection)
	}
	if *TextDetection {
		list = append(list, pigeon.TextDetection)
	}
	if *SafeSearchDetection {
		list = append(list, pigeon.SafeSearchDetection)
	}
	if *ImageProperties {
		list = append(list, pigeon.ImageProperties)
	}

	// Default
	if len(list) == 0 {
		list = append(list, pigeon.FaceDetection)
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