package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main(){
	clipSizes, clipPaths := scanClipDir()
	generateAllChunks(clipSizes, clipPaths)
	fmt.Println("Finished")
}

func getAbsolutePath(path string) string{
	rootDir, rootErr := filepath.Abs(filepath.Dir(os.Args[0]))

	check(rootErr)

	return filepath.Join(rootDir, path)
}

func scanClipDir() ([]int, []string){
	var fileSizes []int
	var filePaths []string

	fmt.Println("Scanning clip directory")

	clipDir := getAbsolutePath(os.Args[1])

	if _, err := os.Stat(clipDir); os.IsNotExist(err) {
		fmt.Println("ERROR: The source clip directory that was provided does not exists. Please check and try again")
		os.Exit(1)
	}

	//get digit files
	for i := 0; i < 10; i++ {
		fileName := strconv.Itoa(i) + ".mp4"
		path := filepath.Join(clipDir, fileName)
		fi, err := os.Stat(path)

		if os.IsNotExist(err) {
			fmt.Println("ERROR: No file " + fileName + " found at location " + path)
		} else {
			fileSizes = append(fileSizes, int(fi.Size()))
			filePaths = append(filePaths, path)
		}
	}

	//get gap file
	gapPath := filepath.Join(clipDir, "gap.mp4")
	fi, err := os.Stat(gapPath)

	if os.IsNotExist(err) {
		fmt.Println("ERROR: No file gap.mp4 found at location " + gapPath)
	} else {
		fileSizes = append(fileSizes, int(fi.Size()))
		filePaths = append(filePaths, gapPath)
	}

	return fileSizes, filePaths
}

func generateAllChunks(sizeArr []int, pathArr []string){
	fmt.Println("Building chunks")

	const CHUNKSIZE = 1300000000
	digitLen := 1000000
	outDir := getAbsolutePath(os.Args[2]) + "\\"
	i := getResumeDigit()

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModeDir)
	}

	if _, err := os.Stat(outDir + "temp.mp4"); !os.IsNotExist(err) {
		fmt.Println("Clearing previous temp file")
		_ = os.Remove(outDir + "temp.mp4")
	}

	for i <= digitLen {
		totalBytes := 0
		startDigit := i
		digitChunk := ""
		argClips := ""
		argFilter := ""

		//calculate and display percentage
		perc := float32(i) / float32(digitLen) * 100
		fmt.Printf("%f%% Current Digit:%d\n", perc, i)

		//build chunk
		for totalBytes < CHUNKSIZE && i <= digitLen {
			curNum := strconv.Itoa(i)

			for d := 0; d < len(curNum); d++ {
				digitStr := string(curNum[d])
				digitInt, _ := strconv.Atoi(digitStr)
				digitChunk += digitStr
				totalBytes += sizeArr[digitInt]
			}

			digitChunk += "_"
			totalBytes += sizeArr[len(sizeArr) - 1]

			i++
		}

		chunkLen := len(digitChunk)
		endDigit := i - 1

		//fill in commands
		for d := 0; d < chunkLen; d++ {
			var clipPath string
			climpNum := strconv.Itoa(d)
			curStrDigit := string(digitChunk[d])

			if curStrDigit == "_" {
				clipPath = string(pathArr[len(pathArr) - 1])
			} else {
				curDigit, _ := strconv.Atoi(curStrDigit)
				clipPath = string(pathArr[curDigit])
			}

			argClips += "-i " + clipPath + " "
			argFilter += "[" + climpNum + ":v][" + climpNum + ":a] "
		}

		//execute command
		if len(argClips) > 0 && len(argFilter) > 0 {
			outFileName := strconv.Itoa(startDigit) + "-" + strconv.Itoa(endDigit)
			argStr := "ffmpeg " + argClips
			argStr += "-filter_complex \"" + argFilter
			argStr += "concat=n=" + strconv.Itoa(chunkLen)
			argStr += ":v=1:a=1 [v] [a]\" -map \"[v]\" -map \"[a]\"" + outDir + "temp.mp4"
			pws := exec.Command("powershell", "/c", argStr)
			//pws.Stdout = os.Stdout
			//pws.Stderr = os.Stderr
			err := pws.Run()

			if err != nil {
				fmt.Println("ERROR: Error combining chunk " + outFileName +  " with FFmpeg. Check source clips for corrupted files. If that does not resolve the issue, remove any temp.mp4 files and try again")
				os.Exit(1)
			} else {
				os.Rename(outDir + "temp.mp4", outDir + outFileName + ".mp4")
			}
		}

		fmt.Println(i)
	}
}

func getResumeDigit() int{
	largestNum := 0
	outPath := getAbsolutePath(os.Args[2])

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		fmt.Println("No existing clips found, starting at digit 0")
		return 0
	}

	clipPathErr := filepath.Walk(outPath, func(path string, info os.FileInfo, err error) error {
		fileName := path[len(outPath):]
		fileName = regexp.MustCompile(`\\|\.\w+`).ReplaceAllString(fileName, "")
		nums := strings.Split(fileName, "-")

		if len(nums) > 1 {
			bigNum, err := strconv.Atoi(nums[1])

			if err == nil && bigNum > largestNum{
				largestNum = bigNum
			}
		}

		return nil
	})

	check(clipPathErr)

	return largestNum
}

func check(e error){
	if (e != nil) {
		panic(e)
	}
}