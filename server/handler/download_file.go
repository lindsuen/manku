// manku - download_file.go
// Copyright (C) 2025 LindSuen <lindsuen@foxmail.com>
//
// Use of this source code is governed by a BSD 2-Clause license that can be
// found in the LICENSE file.

package handler

import (
	"encoding/base64"

	"github.com/bytedance/sonic"
	"github.com/labstack/echo/v4"
	"github.com/lindsuen/manku/internal/db"
	"github.com/lindsuen/manku/server/core"
)

func DownloadFile(c echo.Context) error {
	fileId := c.QueryParam("fileid")
	file := new(core.File)
	value, _ := base64.RawURLEncoding.DecodeString(string(db.Get([]byte(fileId))))
	err := sonic.Unmarshal(value, &file)
	if err != nil {
		return err
	}
	return c.Attachment(file.Path, file.Name)
}
