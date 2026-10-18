package browser_tests

import (
	"LukeWinikates/january-twenty-five/lib/server"
	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/browser/testing/gosttest"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWebBrowser(t *testing.T) {
	t.Parallel()
	b := browser.New(
		browser.WithHandler(server.New()),
		browser.WithLogger(gosttest.NewTestingLogger(t)),
		browser.WithContext(t.Context()),
	)
	window, err := b.Open("http://example.com/") // Host is ignored
	assert.NoError(t, err)
	window.Document().GetElementById("test-button").(html.HTMLElement).Click()

	resultField := win.Document().GetElementById("output-element")
	assert.Equal(t, "The button was clicked", resultField.TextContent())
}
