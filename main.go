package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"video-screenshoter/screenshoter"
	"video-screenshoter/utils"
)

func main() {
	currentPath := utils.GetCurrentPath()
	logChannel := make(chan string)

	myApp := app.New()
	myWindow := myApp.NewWindow("Video Screenshoter 1.0.0")

	// Video screenshot interval
	videoScreenshotIntervalTimeLabel := widget.NewLabel("Screenshot inverval time (HH:MM:SS):")
	videoScreenshotIntervalTimeEntry := widget.NewEntry()
	videoScreenshotIntervalTimeEntry.SetText("00:01:00")
	videoScreenshotIntervalTimeContainerRow := container.New(layout.NewGridLayoutWithColumns(3),
		videoScreenshotIntervalTimeLabel,
		videoScreenshotIntervalTimeEntry)

	// Video input entry
	videoInputDirPathLabel := widget.NewLabel("Video input path:")
	videoInputDirPathEntry := widget.NewEntry()
	videoInputDirPathEntry.SetText(currentPath + "\\video-input")
	videoInputDirPathSelectButton := widget.NewButton("Choose folder", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				videoInputDirPathEntry.SetText(uri.Path())
			}
		}, myWindow)
		listableUri, err := storage.ListerForURI(storage.NewFileURI(currentPath))
		utils.CheckError(err)

		folderDialog.SetLocation(listableUri)
		folderDialog.Show()
	})
	videoInputDirPathContaine := container.New(layout.NewGridLayoutWithRows(2),
		videoInputDirPathLabel,
		container.New(layout.NewGridLayoutWithColumns(2),
			videoInputDirPathEntry,
			videoInputDirPathSelectButton))

	// Video output entry
	videoOutputDirPathLabel := widget.NewLabel("Video output path:")
	videoOutputDirPathEntry := widget.NewEntry()
	videoOutputDirPathEntry.SetText(currentPath + "\\video-output")
	videoOutputDirPathSelectButton := widget.NewButton("Choose folder", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				videoOutputDirPathEntry.SetText(uri.Path())
			}
		}, myWindow)
		listableUri, err := storage.ListerForURI(storage.NewFileURI(currentPath))
		utils.CheckError(err)

		folderDialog.SetLocation(listableUri)
		folderDialog.Show()
	})
	videoOutputDirPathContainer := container.New(layout.NewGridLayoutWithRows(2),
		videoOutputDirPathLabel,
		container.New(layout.NewGridLayoutWithColumns(2),
			videoOutputDirPathEntry,
			videoOutputDirPathSelectButton))

	// ffmpeg path entry
	ffmpegDirPathLabel := widget.NewLabel("FFMPEG file path:")
	ffmpegDirPathEntry := widget.NewEntry()
	ffmpegDirPathEntry.SetText(currentPath + "\\bin\\ffmpeg.exe")
	ffmpegDirPathSelectButton := widget.NewButton("Choose file", func() {
		folderDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				ffmpegDirPathEntry.SetText(file.URI().Path())
			}
		}, myWindow)
		listableUri, err := storage.ListerForURI(storage.NewFileURI(currentPath))
		utils.CheckError(err)

		folderDialog.SetLocation(listableUri)
		folderDialog.Show()
	})
	ffmpegDirPathContainer := container.New(layout.NewGridLayoutWithRows(2),
		ffmpegDirPathLabel,
		container.New(layout.NewGridLayoutWithColumns(2),
			ffmpegDirPathEntry,
			ffmpegDirPathSelectButton))

	// ffprobe path entry
	ffprobeDirPathLabel := widget.NewLabel("FFPROBE file path:")
	ffprobeDirPathEntry := widget.NewEntry()
	ffprobeDirPathEntry.SetText(currentPath + "\\bin\\ffprobe.exe")
	ffprobeDirPathSelectButton := widget.NewButton("Choose file", func() {
		folderDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				ffmpegDirPathEntry.SetText(file.URI().Path())
			}
		}, myWindow)
		listableUri, err := storage.ListerForURI(storage.NewFileURI(currentPath))
		utils.CheckError(err)

		folderDialog.SetLocation(listableUri)
		folderDialog.Show()
	})
	ffprobeDirPathContainer := container.New(layout.NewGridLayoutWithRows(2),
		ffprobeDirPathLabel,
		container.New(layout.NewGridLayoutWithColumns(2),
			ffprobeDirPathEntry,
			ffprobeDirPathSelectButton))

	// Start button
	var startButton *widget.Button
	startButton = widget.NewButton("Start", func() {
		startButton.Disable()

		go func() {
			screenshoter.Run(logChannel, screenshoter.ConfigData{
				InputPath:              utils.GetRelativePath(currentPath, videoInputDirPathEntry.Text),
				OutputPath:             utils.GetRelativePath(currentPath, videoOutputDirPathEntry.Text),
				ScreenshotInvervalTime: videoScreenshotIntervalTimeEntry.Text,
				FfmpegPath:             utils.GetRelativePath(currentPath, ffmpegDirPathEntry.Text),
				FfprobePath:            utils.GetRelativePath(currentPath, ffprobeDirPathEntry.Text)})
		}()
	})
	startButtonContainer := container.New(layout.NewGridLayoutWithRows(1), startButton)

	// Logger
	loggerLabel := widget.NewLabel("Logger")
	loggerContainerBox := container.New(layout.NewVBoxLayout(), loggerLabel)
	loggerContainerBoxScrollable := container.NewVScroll(loggerContainerBox)

	go func() {
		for logMsg := range logChannel {
			loggerContainerBox.Add(widget.NewLabel(logMsg))

			if logMsg == "finalized..." {
				startButton.Enable()
			}

			loggerContainerBoxScrollable.Refresh()
			loggerContainerBoxScrollable.ScrollToBottom()
		}
	}()

	// Main container
	grid := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(),
			container.New(layout.NewVBoxLayout(),
				container.NewPadded(videoScreenshotIntervalTimeContainerRow),
				container.NewPadded(videoInputDirPathContaine),
				container.NewPadded(videoOutputDirPathContainer),
				container.NewPadded(ffmpegDirPathContainer),
				container.NewPadded(ffprobeDirPathContainer),
				container.NewPadded(startButtonContainer)),
			container.New(layout.NewGridLayoutWithRows(1),
				loggerContainerBoxScrollable)))

	myWindow.SetContent(grid)
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(1080, 400))
	myWindow.ShowAndRun()
}
