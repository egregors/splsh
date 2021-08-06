package main

import (
	"fmt"
	"os/exec"
	"strconv"

	"strings"

	"github.com/egregors/splsh/icon"
	"github.com/getlantern/systray"
	"github.com/reujab/wallpaper"
	"github.com/skratchdot/open-golang/open"
)

// TODO: extract it to OS specific modules
// getScreenResolution returns the screen resolution of the current monitor (works on MacOS)
func getScreenResolution() (int, int) {
	cmd := "system_profiler SPDisplaysDataType | awk '/Resolution/{print $2, $3, $4}'"
	// FIXME: catch err
	out, _ := exec.Command("bash", "-c", cmd).Output()
	res := strings.Split(strings.Trim(string(out), "\n"), " ")
	height, _ := strconv.Atoi(res[0])
	width, _ := strconv.Atoi(res[2])
	return height, width
}

func main() {
	println("Hello, World!")
	getScreenResolution()
	systray.Run(onReady, onExit)
}

func onReady() {
	// TODO: find icon in modern MacOS style
	systray.SetIcon(icon.Data)
	systray.SetTooltip("Splsh")
	mImgSourceURL := systray.AddMenuItem("Unsplash", "All images from https://unsplash.com/")
	h, w := getScreenResolution()
	mScreenResolution := systray.AddMenuItem(fmt.Sprintf("%d x %d", h, w), "")
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

	// mbp 16" for example
	// imgURL := fmt.Sprintf("https://picsum.photos/3072/1920", h, w)
	imgURL := fmt.Sprintf("https://picsum.photos/%d/%d", h, w)

	for {
		select {
		case <-mImgSourceURL.ClickedCh:
			_ = open.Run("https://unsplash.com/")
		case <-mNext.ClickedCh:
			println("Next")
			err := wallpaper.SetFromFile(imgURL)
			if err != nil {
				println(err)
			}
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
