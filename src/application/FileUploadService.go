package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hoanganf/pos_domain/entity/exception"
	"github.com/hoanganf/pos_domain/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileUploadService struct {
	LoginUrl        string
	TokenName       string
	Directory       string
	MaxFileSize     int64
	FileType        []string
	LoginRepository repository.LoginRepository
}

func NewFileUploadService(loginUrl string, tokenName string, directory string, maxFileSize int64, loginRepository repository.LoginRepository) *FileUploadService {

	return &FileUploadService{LoginUrl: loginUrl,
		TokenName:       tokenName,
		Directory:       directory,
		MaxFileSize:     maxFileSize * 1000000,
		FileType:        []string{".jpg", ".png", ".jpeg", ".gif"},
		LoginRepository: loginRepository}
}

func (s *FileUploadService) Get(c *gin.Context) {
	/*if !s.IsLogin(c) {
		log.Print("host %q", c.Request.URL.Host)
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s?origin=%s", s.LoginUrl, c.Request.URL))
		return
	}*/
	c.HTML(http.StatusOK, "index.tmpl", nil)
	return
}

func (s *FileUploadService) Post(c *gin.Context) {
	/*if !s.IsLogin(c) {
		c.JSON(http.StatusUnauthorized, exception.CreateError(exception.CodeSignatureInvalid, "Access denied."))
		return
	}*/

	directory := c.PostForm("directory")
	if directory != "" {
		directory = filepath.Base(directory)
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "file required."))
		return
	}

	//check file Size
	if file.Size > s.MaxFileSize {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, fmt.Sprintf("Sorry, your file is too large. < %fMB are allowed: ", float64(s.MaxFileSize/1000000))))
		return
	}

	//check file name

	baseFileName := filepath.Base(file.Filename)
	if s.IsInvalidFileExtension(baseFileName) {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeSignatureInvalid, "Sorry, only JPG, JPEG, PNG & GIF files are allowed."))
		return
	}

	filename := filepath.Join(s.Directory, directory, baseFileName)
	err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeUnknown, "can not create file folder."))
		return
	}

	//check file is exits
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, exception.CreateError(exception.CodeValueInvalid, "Sorry, file already exists."))
		return
	}

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, exception.CreateError(exception.CodeUnknown, fmt.Sprintf("upload file err: %s", err.Error())))
		return
	}
	c.JSON(http.StatusOK, gin.H{"filePath": c.Request.URL.String()+"/view/" + directory + "/" + baseFileName, "message": fmt.Sprintf("'%s' uploaded!", file.Filename)})
}

func (s *FileUploadService) IsLogin(c *gin.Context) bool {
	cookie, err := c.Cookie(s.TokenName)
	if err != nil {
		log.Print(err)
		return false
	}
	user, cErr := s.LoginRepository.GetUserByJwt(cookie)
	if cErr != nil {
		log.Print(cErr.ErrorMessage)
		return false
	}

	if user.JWT != "" {
		return true
	}
	return false
}

func (s *FileUploadService) IsInvalidFileExtension(fileName string) bool {
	ext := filepath.Ext(fileName)
	for _, v := range s.FileType {
		if v == ext {
			return false
		}
	}
	return true
}
