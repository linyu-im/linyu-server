package utils

import (
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
)

func ResolvePath(filePath, fileKey string) (string, error) {
	absBasePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(absBasePath, fileKey)
	rel, err := filepath.Rel(absBasePath, fullPath)
	if err != nil {
		return "", errors.New("unsafe file path detected")
	}

	if strings.HasPrefix(rel, "..") || rel == ".." {
		return "", errors.New("unsafe file path detected")
	}

	return fullPath, nil
}

func GetFileSystemPath(filePath, fileKey string) (string, error) {
	fullPath, err := ResolvePath(filePath, fileKey)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	}
	return fullPath, nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string, perm ...fs.FileMode) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var mode os.FileMode = 0o750
	if len(perm) > 0 {
		mode = perm[0]
	}
	dir := filepath.Dir(dst)
	if err = os.MkdirAll(dir, mode); err != nil {
		return err
	}
	if err = os.Chmod(dir, mode); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

var extCategoryMap = map[string]string{
	"jpg": constant.FileCategory.Image, "jpeg": constant.FileCategory.Image,
	"png": constant.FileCategory.Image, "gif": constant.FileCategory.Image,
	"bmp": constant.FileCategory.Image, "webp": constant.FileCategory.Image,
	"svg": constant.FileCategory.Image, "ico": constant.FileCategory.Image,
	"tiff": constant.FileCategory.Image,

	"mp4": constant.FileCategory.Video, "avi": constant.FileCategory.Video,
	"mkv": constant.FileCategory.Video, "mov": constant.FileCategory.Video,
	"wmv": constant.FileCategory.Video, "flv": constant.FileCategory.Video,
	"webm": constant.FileCategory.Video, "m4v": constant.FileCategory.Video,
	"3gp": constant.FileCategory.Video,

	"mp3": constant.FileCategory.Audio, "wav": constant.FileCategory.Audio,
	"flac": constant.FileCategory.Audio, "aac": constant.FileCategory.Audio,
	"ogg": constant.FileCategory.Audio, "wma": constant.FileCategory.Audio,
	"m4a": constant.FileCategory.Audio, "opus": constant.FileCategory.Audio,

	"doc": constant.FileCategory.Document, "docx": constant.FileCategory.Document,
	"xls": constant.FileCategory.Document, "xlsx": constant.FileCategory.Document,
	"ppt": constant.FileCategory.Document, "pptx": constant.FileCategory.Document,
	"pdf": constant.FileCategory.Document, "txt": constant.FileCategory.Document,
	"csv": constant.FileCategory.Document, "md": constant.FileCategory.Document,
	"rtf": constant.FileCategory.Document, "odt": constant.FileCategory.Document,
	"ods": constant.FileCategory.Document, "odp": constant.FileCategory.Document,

	"zip": constant.FileCategory.Archive, "rar": constant.FileCategory.Archive,
	"7z": constant.FileCategory.Archive, "tar": constant.FileCategory.Archive,
	"gz": constant.FileCategory.Archive, "bz2": constant.FileCategory.Archive,
	"xz": constant.FileCategory.Archive, "zst": constant.FileCategory.Archive,
}

// FileCategoryFromExt 根据文件后缀推断文件分类
func FileCategoryFromExt(ext string) string {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	if category, ok := extCategoryMap[ext]; ok {
		return category
	}
	return constant.FileCategory.Other
}
