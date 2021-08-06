package main

import (
	"fmt"
	"os/exec"

	"github.com/egregors/splsh/icon"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

// TODO: extract it to OS specific modules
// getScreenResolution returns the screen resolution of the current monitor (works on MacOS)
func getScreenResolution() string {
	cmd := "system_profiler SPDisplaysDataType | awk '/Resolution/{print $2, $3, $4}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Failed to execute command: %s", cmd)
	}
	return string(out)
}

func main() {
	println("Hello, World!")
	systray.Run(onReady, onExit)
}

func onReady() {
	// TODO: find icon in modern MacOS style
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Splsh")
	mImgSourceURL := systray.AddMenuItem("Unsplash", "All images from https://unsplash.com/")
	mScreenResolution := systray.AddMenuItem(getScreenResolution(), "")
	mScreenResolution.Disable()
	systray.AddSeparator()
	mNext := systray.AddMenuItem("Next", "Set next image")
	systray.AddSeparator()
	mGrayscaleMode := systray.AddMenuItemCheckbox("Grayscale", "", false)
	mBlurMode := systray.AddMenuItemCheckbox("Blur", "", false)
	systray.AddSeparator()
	mClearCache := systray.AddMenuItem("Clear Cache", "")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	for {
		select {
		case <-mImgSourceURL.ClickedCh:
			_ = open.Run("https://unsplash.com/")
		case <-mNext.ClickedCh:
			println("Next")
			// TODO:
			// 		- [ ] get cache folder (get or create)
			// 		- [ ] download image
			// 		- [ ] set image from file
		// ---------------------------
		case <-mGrayscaleMode.ClickedCh:
			println("Gray Mood")
			if !mGrayscaleMode.Checked() {
				println("gray mode on")
				mGrayscaleMode.Check()
			} else {
				println("gray mode off")
				mGrayscaleMode.Uncheck()
			}
		case <-mBlurMode.ClickedCh:
			println("blur Mood")
			if !mBlurMode.Checked() {
				println("blur mode on")
				mBlurMode.Check()
			} else {
				println("blur mode off")
				mBlurMode.Uncheck()
			}
		// ---------------------------
		case <-mClearCache.ClickedCh:
			println("Clear Cache")
		// ---------------------------
		case <-mQuit.ClickedCh:
			systray.Quit()
			fmt.Println("Quit2 now...")
			return
		}
	}
}

func onExit() {
	fmt.Println("Bye, World!")
}
