package file

import (
	"fmt"
	"gohub/pkg/app"
	"gohub/pkg/auth"
	"gohub/pkg/helpers"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileNameWithoutExtension 获取不含后缀的文件名
func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func SaveUploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {
	var avatar string
	// 确保目录存在，不存在创建
	publicPath := "public"
	dirName := fmt.Sprintf(
		"/uploads/avatars/%s/%s/",
		app.TimenowInTimezone().Format("2006/01/02"),
		auth.CurrentUID(c),
	)
	os.MkdirAll(publicPath+dirName, 0755)

	// 保存文件
	fileName := randomNameFromUploadFile(file)
	avatarPath := publicPath + dirName + fileName
	if err := c.SaveUploadedFile(file, avatarPath); err != nil {
		return avatar, err
	}

	return avatarPath, nil
}

func randomNameFromUploadFile(file *multipart.FileHeader) string {
	return helpers.RandomString(16) + filepath.Ext(file.Filename)
}
