package util

import (
	"fmt"
	"os"
	"product-service/src/common/log"
	"regexp"
	"time"
)

func CreateUnixFileName(filename string) string {
	re := regexp.MustCompile("[ %?#&=]")
	encodedName := re.ReplaceAllString(filename, "-")
	epochMillis := time.Now().UnixMilli()

	filename = fmt.Sprintf("%d-%s", epochMillis, encodedName)
	return filename
}

func CheckExistDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
		return err
	}

	return nil
}

func DeleteFile(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			log.Logger.Error(err.Error())
		}
	}
}
