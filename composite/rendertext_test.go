package composite

import (
	"image/png"
	"os"
	"testing"
)

func TestRenderTextToPNG(t *testing.T) {
	const filename = "go_test_output.png"

	// First, test to see if we can generate the file
	if err := RenderTextToPNG("Hello, World!", filename, "font.ttf"); err != nil {
		t.Fatal("Error rendering text to png:", err)
	}

	// So far so good. Now, check to see if it exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("Error rendering text to png:", filename, "does not exist after render complete")
	}
	defer func() { // since we know the file is there
		os.Remove(filename) // we should remove it later
	}()

	// Just to be thorough, let's see if we can parse it
	imageFile, err := os.Open(filename)
	if err != nil {
		t.Fatal("Error opening rendered PNG File:", err)
	}
	defer imageFile.Close()

	if _, err := png.Decode(imageFile); err != nil {
		t.Fatal("Error decoding rendered png:", err)
	}
}
