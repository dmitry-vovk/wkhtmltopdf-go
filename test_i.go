package main

import (
	"fmt"
	"io/ioutil"

	"github.com/dmitry-vovk/wkhtmltopdf-go/wkhtmltoimage"
)

const ext = "png"

func main() {
	c := wkhtmltoimage.NewGlobalSettings().
		// http://www.cs.au.dk/~jakobt/libwkhtmltox_0.10.0_doc/pagesettings.html#pageImageGlobal
		Set("in", "https://www.google.com").
		Set("fmt", ext).
		Set("web.enableJavascript", "true").
		Set("load.stopSlowScript", "true").
		Set("load.loadErrorHandling", "skip").
		Set("load.jsDelay", "1000").
		//gs.Set("load.proxy", "proxy here")
		NewConverter()
	c.ProgressChanged = func(c *wkhtmltoimage.Converter, b int) {
		fmt.Printf("Progress> %d\n", b)
	}
	c.Error = func(c *wkhtmltoimage.Converter, msg string) {
		fmt.Printf("Error> %s\n", msg)
	}
	c.Warning = func(c *wkhtmltoimage.Converter, msg string) {
		fmt.Printf("Warning> %s\n", msg)
	}
	c.Finished = func(c *wkhtmltoimage.Converter, s int) {
		fmt.Printf("Finished: %d\n", s)
	}
	c.Phase = func(c *wkhtmltoimage.Converter) {
		phaseNumber, phaseDescription := c.CurrentPhase()
		fmt.Printf("Phase %d: %s\n", phaseNumber, phaseDescription)
	}
	if err := c.Convert(); err != nil {
		fmt.Printf("Error converting: %s\n", err)
	}
	payload, _ := c.Payload()
	ioutil.WriteFile("test."+ext, payload, 0644)
	fmt.Printf("Got error code: %d\n", c.ErrorCode())
}
