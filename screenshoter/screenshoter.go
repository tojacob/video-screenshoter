package screenshoter

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"video-screenshoter/utils"
)

type ConfigData struct {
	InputPath              string
	OutputPath             string
	ScreenshotInvervalTime string
	FfmpegPath             string
	FfprobePath            string
}

var logChannel chan string

func log(message string) {
	logChannel <- message
}

func createMainOutputDir(outputDir string) {
	_, err := os.Stat(outputDir)

	if os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0777)
		utils.CheckError(err)
	}

	log("output folder is correct")
}

func getListOfVideos(inputDirPath string) []fs.DirEntry {
	var err error
	var dirEntries []fs.DirEntry

	dirEntries, err = os.ReadDir(inputDirPath)
	utils.CheckError(err)

	log("input folder is correct")

	return dirEntries
}

func getVideoDurationInSec(ffprobePath string, videoPath string) int {
	var err error
	var durationStr string
	var durationInt int
	var cmd = exec.Command(ffprobePath, "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)

	out, err := cmd.Output()
	utils.CheckError(err)

	durationStr = strings.TrimSpace(string(out))

	durationInt, err = strconv.Atoi(strings.Split(durationStr, ".")[0])
	utils.CheckError(err)

	return durationInt
}

func createOutputDir(inputPath string, dirName string) {
	var dirPath = inputPath + "/" + dirName
	var err = os.Mkdir(dirPath, 0777)

	utils.CheckError(err)
}

func getOutputPath(inputPath string, name string, screenshotIndex int) string {
	var screenshotIndexWithFormat = fmt.Sprintf("%03d", screenshotIndex)
	var indexName = inputPath + "/" + name + "/" + screenshotIndexWithFormat + ".jpeg"

	return indexName
}

func getScreenshotTimeWithFormat(screenshotIntervalInSec int, iterations int) string {
	var totalSeconds = screenshotIntervalInSec * iterations
	var newHours = totalSeconds / 3600
	var newMinutes = (totalSeconds % 3600) / 60
	var newSeconds = totalSeconds % 60
	var result = fmt.Sprintf("%02d:%02d:%02d", newHours, newMinutes, newSeconds)

	return result
}

func getScreenshotIntervalTimeInSec(intervalTime string) int {
	var err error
	var totalSeconds int
	var parts = strings.Split(intervalTime, ":")

	hours, err := strconv.Atoi(parts[0])
	utils.CheckError(err)

	minutes, err := strconv.Atoi(parts[1])
	utils.CheckError(err)

	seconds, err := strconv.Atoi(parts[2])
	utils.CheckError(err)

	totalSeconds = (hours * 3600) + (minutes * 60) + seconds

	return int(totalSeconds)
}

func generateScreenshotCapture(ffmpegPath string, videoInputPath string, videoOutputPath string, intervalTime string) {
	var err error
	var cmd = exec.Command(ffmpegPath, "-ss", intervalTime, "-i", videoInputPath, "-frames:v", "1", videoOutputPath)

	err = cmd.Run()
	utils.CheckError(err)
}

func generateScreenshotFromVideo(videoNameWithExt string, configData ConfigData) {
	var videoOutputPath string
	var screenshotIntervalTime string
	var videoInputPath = fmt.Sprintf("%s/%s", configData.InputPath, videoNameWithExt)
	var videoDurationInSec = getVideoDurationInSec(configData.FfprobePath, videoInputPath)
	var screenshotIntervalTimeInSec = getScreenshotIntervalTimeInSec(configData.ScreenshotInvervalTime)
	var iterations = math.Floor(float64(videoDurationInSec) / float64(screenshotIntervalTimeInSec))

	for i := 0; i < int(iterations); i++ {
		videoOutputPath = getOutputPath(configData.OutputPath, videoNameWithExt, i)
		screenshotIntervalTime = getScreenshotTimeWithFormat(screenshotIntervalTimeInSec, i+1)

		generateScreenshotCapture(configData.FfmpegPath, videoInputPath, videoOutputPath, screenshotIntervalTime)
		log("generated " + videoOutputPath)
	}
}

func processListOfVideos(entries []fs.DirEntry, configData ConfigData) {
	var videoNameWithExt string

	for i := 0; i < len(entries); i++ {
		videoNameWithExt = entries[i].Name()

		log("processing " + videoNameWithExt)
		createOutputDir(configData.OutputPath, videoNameWithExt)
		generateScreenshotFromVideo(videoNameWithExt, configData)
	}
}

func Run(logChan chan string, configData ConfigData) {
	logChannel = logChan

	log("starting...")
	log("screenshot interval time: " + configData.ScreenshotInvervalTime)
	log("video input path: " + configData.InputPath)
	log("video output path: " + configData.OutputPath)
	log("ffmpeg path: " + configData.FfmpegPath)
	log("ffprobe path: " + configData.FfprobePath)

	var dirEntries = getListOfVideos(configData.InputPath)

	createMainOutputDir(configData.OutputPath)
	processListOfVideos(dirEntries, configData)

	log("finalized...")
}