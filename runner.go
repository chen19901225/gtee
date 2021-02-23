package main

import (
	"fmt"
	"gtee/pkg/util"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

type fLog func(string)

func MustparseByte(byteStr string) int64 {
	var unit int64 = 1
	n := len(byteStr)
	if byteStr[n-1] == 'm' || byteStr[n-1] == 'M' {
		unit = 1024 * 1024
		byteStr = byteStr[:n-1]
	} else if byteStr[n-1] == 'k' || byteStr[n-1] == 'K' {
		unit = 1024
		byteStr = byteStr[:n-1]
	} else if byteStr[n-1] == 'g' || byteStr[n-1] == 'G' {
		unit = 1024 * 1024 * 1024
		byteStr = byteStr[:n-1]
	}
	byteN, err := strconv.Atoi(byteStr)
	util.PanicIfNotNull(err)
	return int64(byteN) * unit

}
func getSizeOrZero(pathStr string) int64 {
	stat, err := os.Stat(pathStr)
	if os.IsNotExist(err) {
		return 0
	}
	return stat.Size()
}

func isSizeSmall(size int64, pathStr string) int {
	stat, err := os.Stat(pathStr)
	if os.IsNotExist(err) {
		return 1
	}
	statSize := stat.Size()
	if statSize < size {
		return 1
	}
	return 0
}
func sizePretty(size int64) string {
	unit := ""
	arr := []string{"k", "m", "g"}
	for i := 0; i < len(arr); i++ {
		if size > 1024 {
			size = size / 1024
			unit = arr[i]
		}
	}
	return fmt.Sprintf("%d%s", size, unit)
}

func run(maxByteStr string, backupCount int, pathStr string, isDebug int) error {
	dir := path.Dir(pathStr)
	var fucLog fLog = func(text string) {
		if isDebug == 1 {
			log.Println(text)
		}
	}
	fucLog(fmt.Sprintf("dir:%s", dir))
	fucLog(fmt.Sprintf("maxByteStr:%s", maxByteStr))
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fucLog(fmt.Sprintf("create dir %s", dir))
		err = os.MkdirAll(dir, os.ModePerm)
		util.PanicIfNotNull(err)
	}
	// fileStat, err := os.Stat(pathStr)
	f, err := os.OpenFile(pathStr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend|os.ModePerm)
	util.PanicIfNotNull(err)
	maxByte := MustparseByte(maxByteStr)
	message := make([]byte, 1024)
	for {

		// n, err := io.ReadFull(os.Stdin, message)
		n, err := os.Stdin.Read(message)
		// 可能是EOF
		if err == io.EOF {
			fucLog("Meet EOF so complete")
			return nil
		}
		if err != nil {
			log.Fatal(err)
		}

		if n > 0 {
			_, err := os.Stdout.Write(message[:n])
			util.PanicIfNotNull(err)
			// should rename
			fileSize := getSizeOrZero(pathStr)
			fileSizeStr := sizePretty(fileSize)
			isSmallSize := fileSize < maxByte
			fucLog(
				fmt.Sprintf(
					"fileSize:%d,maxByte:%d, isSmall:%t, fileSizeStr:%s, maxByteStr:%s", fileSize, maxByte, isSmallSize,
					fileSizeStr,
					maxByteStr,
				),
			)
			if !isSmallSize {
				fucLog("rename begin")
				for i := backupCount - 1; i >= 0; i-- {
					srcPath := fmt.Sprintf("%s.%d", pathStr, i)
					if i == 0 {
						srcPath = pathStr
					}
					DestPath := fmt.Sprintf("%s.%d", pathStr, i+1)
					_, err := os.Stat(srcPath)
					isExist := err == nil || os.IsExist(err)
					fucLog(
						fmt.Sprintf("path:%s, isExists:%t", srcPath, isExist),
					)
					if isExist {
						fucLog(
							fmt.Sprintf("rename %s=>%s", srcPath, DestPath),
						)
						err = os.Rename(srcPath, DestPath)
						util.PanicIfNotNull(err)
					}
				}

				fucLog(
					fmt.Sprintf("reopen file %s", pathStr),
				)

				f, err = os.OpenFile(pathStr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModeAppend|os.ModePerm)
				util.PanicIfNotNull(err)

			}
			_, err = f.Write(message[:n])
			util.PanicIfNotNull(err)
		}
		time.Sleep(time.Microsecond * 10)
	}

}
