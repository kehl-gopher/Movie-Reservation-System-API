package data

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kehl-gopher/Movie-Reservation-System-API/internal/validator"
	"github.com/rs/xid"
)

// handle user image file upload
func HandleImageFile(r *http.Request, key, uploadPath string) (string, error) {
	path, header, err := r.FormFile(key)
	if err != nil {
		return "", &ValidationError{Err: err}
	}
	defer path.Close()
	im_path, err := handleImageUpload(path, header, uploadPath)
	if err != nil {
		return "", err
	}
	return im_path, err
}

func handleImageUpload(path multipart.File, header *multipart.FileHeader, uploadPath string) (string, error) {
	buf := make([]byte, 512)
	_, err := path.Read(buf)
	if err != nil {
		return "", &BadRequestError{Err: err}
	}
	path.Seek(0, 0)

	// validate mime type
	mimeType := http.DetectContentType(buf)
	if !validator.Isimage(mimeType) {
		return "", &BadRequestError{Err: errors.New("Image provided is not a supported image type")}
	}

	// validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !validator.In(ext, ".jpg", ".png", ".jpeg", ".gif", ".webp") {
		return "", &ValidationError{Err: errors.New(fmt.Sprintf("Invalid image parsed with ext not supprted %s", ext))}
	}

	// upload image :to disk...
	imageUploadDir := uploadPath
	if _, err := os.Stat(imageUploadDir); err != nil {
		os.MkdirAll(imageUploadDir, os.ModePerm)
	}

	// generate unique id and concat to image file
	guid := xid.New()
	fileName := strings.Split(header.Filename, ".")
	filePath := filepath.Join(imageUploadDir, fmt.Sprintf("%s-%s.%s", fileName[0], guid.String(), fileName[1]))

	dst, err := os.Create(filePath)

	if err != nil {
		return "", &ServerError{Err: err}
	}

	if _, err = io.Copy(dst, path); err != nil {
		return "", &ServerError{Err: err}
	}

	return filePath, nil
}
