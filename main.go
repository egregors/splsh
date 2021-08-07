package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"strings"

	"github.com/egregors/splsh/icon"
	"github.com/getlantern/systray"
	"github.com/reujab/wallpaper"
	"github.com/skratchdot/open-golang/open"
)

// TODO: extract it to OS specific modules
// getScreenResolution returns the screen resolution of the current monitor (works on MacOS)
func getScreenResolution() (height, width int) {
	cmd := "system_profiler SPDisplaysDataType | awk '/Resolution/{print $2, $3, $4}'"
	// FIXME: catch err
	out, _ := exec.Command("bash", "-c", cmd).Output()
	res := strings.Split(strings.Trim(string(out), "\n"), " ")
	height, _ = strconv.Atoi(res[0])
	width, _ = strconv.Atoi(res[2])
	return height, width
}

// TODO: add revision
func main() {
	fmt.Println(">>> Splsh")
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
	mGrayscaleMode.Disable()
	mBlurMode := systray.AddMenuItemCheckbox("Blur", "", false)
	mBlurMode.Disable()
	systray.AddSeparator()
	mClearCache := systray.AddMenuItem("Clear Cache", "")
	mClearCache.Disable()
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")

	// mbp 16" for example
	// imgURL := fmt.Sprintf("https://picsum.photos/3072/1920", h, w)
	imgURL := fmt.Sprintf("https://picsum.photos/%d/%d", h, w)

	for {
		select {
		case <-mImgSourceURL.ClickedCh:
			_ = open.Run("https://unsplash.com/")
		// TODO: extract all this mess to a separate func
		case <-mNext.ClickedCh:
			println("Next")
			err := wallpaper.SetFromFile(imgURL)
			if err != nil {
				println(err)
			}
			imgURL := fmt.Sprintf("https://picsum.photos/%d/%d", h, w)
			l, err := downloadImage(imgURL)
			if err != nil {
				fmt.Println(err)
			}
			// FIXME: for some reason on MacOS 11.4 background image changes through default one
			err = wallpaper.SetFromFile(l)
			if err != nil {
				fmt.Println(err)
			}
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
	// TODO: clean cache here maybe
	fmt.Println("Bye!")
}

func downloadImage(url string) (string, error) {
	// nolint:gosec // it's ok for this goal
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	cacheDir := getCacheDir()

	file, err := os.Create(filepath.Join(cacheDir, fmt.Sprintf("%s.jpg", randStringRunes(16))))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

// TODO: get cache dir from OS
func getCacheDir() string {
	_ = os.MkdirAll(filepath.Join("tmp", "splsh"), os.ModePerm)
	return "/tmp/splsh"
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		// nolint:gosec // it's ok for this goal
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
