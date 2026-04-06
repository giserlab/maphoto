package api

import (
	"fmt"
	"io"
	"maphoto/internal/util"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	UploadDir       = "./uploads"
	ThumbsDir       = "thumbs"
	PhotosDir       = "photos"
	MaxUploadSize   = 20 * 1024 * 1024 // 20MB
)

// FileInfo 文件信息
type FileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	URL          string    `json:"url"`
}

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// getStoragePath 获取存储路径
func getStoragePath(folder string) string {
	return filepath.Join(UploadDir, folder)
}

// @Summary  上传文件到本地存储
// @Tags storage
// @Security JWT
// @Accept   multipart/form-data
// @Param    folder formData string true "目标文件夹: thumbs 或 photos"
// @Param    files  formData file   true "要上传的文件"
// @Success  200 {object} Response
// @Router   /api/v1/storage/upload [POST]
func StorageUpload(c echo.Context) error {
	folder := c.FormValue("folder")
	if folder != ThumbsDir && folder != PhotosDir {
		return ApiFailed(c, 400, "无效的文件夹，只能是 thumbs 或 photos")
	}

	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		return ApiFailed(c, 400, "无法解析上传的文件: "+err.Error())
	}

	files := form.File["files"]
	if len(files) == 0 {
		return ApiFailed(c, 400, "没有选择文件")
	}

	// 确保目录存在
	storagePath := getStoragePath(folder)
	if err := ensureDir(storagePath); err != nil {
		return ApiFailed(c, 500, "创建目录失败: "+err.Error())
	}

	var uploadedFiles []FileInfo
	var failedFiles []string

	for _, file := range files {
		// 检查文件大小
		if file.Size > MaxUploadSize {
			failedFiles = append(failedFiles, file.Filename+" (文件过大)")
			continue
		}

		// 检查文件类型
		if !isAllowedImageType(file) {
			failedFiles = append(failedFiles, file.Filename+" (不支持的文件类型)")
			continue
		}

		// 生成唯一文件名
		ext := filepath.Ext(file.Filename)
		filename := util.ShortUID(16) + ext
		filepath := filepath.Join(storagePath, filename)

		// 保存文件
		src, err := file.Open()
		if err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}
		defer src.Close()

		dst, err := os.Create(filepath)
		if err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		// 获取文件信息
		info, _ := os.Stat(filepath)
		fileInfo := FileInfo{
			Name:         filename,
			Path:         folder + "/" + filename,
			Size:         info.Size(),
			LastModified: info.ModTime(),
			URL:          getFileURL(c, folder, filename),
		}
		uploadedFiles = append(uploadedFiles, fileInfo)
	}

	return ApiSuccess(c, map[string]interface{}{
		"success": uploadedFiles,
		"failed":  failedFiles,
		"count":   len(uploadedFiles),
	})
}

// @Summary  获取本地存储文件列表
// @Tags storage
// @Security JWT
// @Param    folder query string false "文件夹: thumbs 或 photos，不传则返回所有"
// @Success  200 {object} Response
// @Router   /api/v1/storage/files [GET]
func StorageList(c echo.Context) error {
	folder := c.QueryParam("folder")

	var folders []string
	if folder == "" {
		folders = []string{ThumbsDir, PhotosDir}
	} else if folder == ThumbsDir || folder == PhotosDir {
		folders = []string{folder}
	} else {
		return ApiFailed(c, 400, "无效的文件夹参数")
	}

	var files []FileInfo

	for _, f := range folders {
		storagePath := getStoragePath(f)

		// 检查目录是否存在
		if _, err := os.Stat(storagePath); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(storagePath)
		if err != nil {
			util.Logger.Error("读取目录失败: " + err.Error())
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, FileInfo{
				Name:         info.Name(),
				Path:         f + "/" + info.Name(),
				Size:         info.Size(),
				LastModified: info.ModTime(),
				URL:          getFileURL(c, f, info.Name()),
			})
		}
	}

	// 按时间倒序排列
	sort.Slice(files, func(i, j int) bool {
		return files[i].LastModified.After(files[j].LastModified)
	})

	return ApiSuccess(c, files)
}

// @Summary  删除本地存储文件
// @Tags storage
// @Security JWT
// @Param    folder path string true "文件夹: thumbs 或 photos"
// @Param    filename path string true "文件名"
// @Success  200 {object} Response
// @Router   /api/v1/storage/files/{folder}/{filename} [DELETE]
func StorageDelete(c echo.Context) error {
	folder := c.Param("folder")
	filename := c.Param("filename")

	if folder != ThumbsDir && folder != PhotosDir {
		return ApiFailed(c, 400, "无效的文件夹")
	}

	// 安全检查：防止目录遍历
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return ApiFailed(c, 400, "无效的文件名")
	}

	filepath := filepath.Join(getStoragePath(folder), filename)

	// 检查文件是否存在
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return ApiFailed(c, 404, "文件不存在")
	}

	// 删除文件
	if err := os.Remove(filepath); err != nil {
		return ApiFailed(c, 500, "删除文件失败: "+err.Error())
	}

	return ApiSuccess(c, nil)
}

// @Summary  重命名本地存储文件
// @Tags storage
// @Security JWT
// @Param    folder path string true "文件夹: thumbs 或 photos"
// @Param    filename path string true "原文件名"
// @Param    newName formData string true "新文件名"
// @Success  200 {object} Response
// @Router   /api/v1/storage/files/{folder}/{filename}/rename [POST]
func StorageRename(c echo.Context) error {
	folder := c.Param("folder")
	filename := c.Param("filename")
	newName := c.FormValue("newName")

	if folder != ThumbsDir && folder != PhotosDir {
		return ApiFailed(c, 400, "无效的文件夹")
	}

	// 安全检查
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return ApiFailed(c, 400, "无效的原文件名")
	}

	if strings.Contains(newName, "..") || strings.Contains(newName, "/") || strings.Contains(newName, "\\") {
		return ApiFailed(c, 400, "无效的新文件名")
	}

	oldPath := filepath.Join(getStoragePath(folder), filename)
	newPath := filepath.Join(getStoragePath(folder), newName)

	// 检查原文件是否存在
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return ApiFailed(c, 404, "文件不存在")
	}

	// 检查新文件名是否已存在
	if _, err := os.Stat(newPath); err == nil {
		return ApiFailed(c, 400, "新文件名已存在")
	}

	// 重命名文件
	if err := os.Rename(oldPath, newPath); err != nil {
		return ApiFailed(c, 500, "重命名失败: "+err.Error())
	}

	info, _ := os.Stat(newPath)
	return ApiSuccess(c, FileInfo{
		Name:         info.Name(),
		Path:         folder + "/" + info.Name(),
		Size:         info.Size(),
		LastModified: info.ModTime(),
		URL:          getFileURL(c, folder, info.Name()),
	})
}

// isAllowedImageType 检查是否允许的图片类型
func isAllowedImageType(file *multipart.FileHeader) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/bmp":  true,
	}

	contentType := file.Header.Get("Content-Type")
	if allowedTypes[contentType] {
		return true
	}

	// 检查后缀
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".bmp":  true,
	}
	return allowedExts[ext]
}

// getFileURL 获取文件访问URL
func getFileURL(c echo.Context, folder, filename string) string {
	scheme := c.Scheme()
	host := c.Request().Host
	return fmt.Sprintf("%s://%s/uploads/%s/%s", scheme, host, folder, filename)
}

// InitStorage 初始化存储目录
func InitStorage() error {
	// 检查并创建上传目录
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
			return err
		}
	}

	// 创建子目录
	dirs := []string{ThumbsDir, PhotosDir}
	for _, dir := range dirs {
		path := filepath.Join(UploadDir, dir)
		if err := ensureDir(path); err != nil {
			return err
		}
	}

	return nil
}
