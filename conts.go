package alfred

import "fmt"

const (
	NoAnyResult          = "No any result"
	OpenInTheBrowser     = "Open in the browser"
	OpenInTheBrowserTips = "Press ⌥ " + OpenInTheBrowser
	OpenInTheTerminal    = "Open in the Terminal"
)

func BuildCopyTips(v string) string {
	return fmt.Sprintf("Copy '%s' to the clipboard", v)
}
