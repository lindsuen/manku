// manku - upload_file.go
// Copyright (C) 2025 LindSuen <lindsuen@foxmail.com>
//
// Use of this source code is governed by a BSD 2-Clause license that can be
// found in the LICENSE file.

package handler

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	cfg "github.com/lindsuen/manku/internal/config"
	"github.com/lindsuen/manku/internal/db"
	"github.com/lindsuen/manku/server/core"
)

var content struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Path        string `json:"path"`
	CreatedTime int64  `json:"createdTime"`
	Hash        string `json:"hash"`
}

func UploadFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	formFiles := form.File["files"]
	for _, fileHeader := range formFiles {
		multiFile, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer multiFile.Close()
		fileName := fileHeader.Filename
		fileSize := fileHeader.Size

		if fileSize > int64(parseMaxLength(cfg.Config.MaxLength)) {
			return c.String(http.StatusForbidden, fileName+" is too large.")
		}

		cFile := new(core.File)
		cFile.SetFileID()
		cFile.SetFileName(fileName)
		cFile.SetFileSize(fileSize)
		cFile.SetFileCreatedTime()

		storagePath := createDateDir(cfg.Config.StoragePath) + "/" + setLocalFileName(fileName, cFile.CreatedTime)
		file, err := os.Create(storagePath)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()

		_, err = io.Copy(file, multiFile)
		if err != nil {
			return err
		}
		cFile.SetFilePath(storagePath)
		cFile.SetFileHash(file)

		content.Id = cFile.ID
		content.Name = cFile.Name
		content.Size = cFile.Size
		content.Path = cFile.Path
		content.CreatedTime = cFile.CreatedTime
		content.Hash = cFile.Hash

		key := []byte(content.Id)
		value, _ := json.Marshal(content)
		db.Set(key, []byte(base64.RawURLEncoding.EncodeToString(value)))
	}

	return c.JSON(http.StatusOK, &content)
}

func createDateDir(basePath string) string {
	subFolderName := time.Now().Format("20060102")
	folderPath := fmt.Sprint(basePath + "/" + subFolderName)
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	return folderPath
}

func setLocalFileName(name string, timestamp int64) string {
	nameByte := []byte(name)
	dataPrefix := fmt.Appendf(nil, "%x", sha1.Sum(nameByte))
	return string(dataPrefix[:29]) + strconv.FormatInt(timestamp, 10)
}

func parseMaxLength(s string) int {
	maxlength, _ := strconv.Atoi(s)
	return maxlength
}
