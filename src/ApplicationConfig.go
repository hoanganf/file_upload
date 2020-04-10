package src

import (
	"errors"
	"github.com/hoanganf/file_upload/src/application"
	"github.com/hoanganf/file_upload/src/infrastructure/client"
	"github.com/hoanganf/file_upload/src/infrastructure/persistence"
	"os"
	"strconv"
)

type Bean struct {
	FileUploadService *application.FileUploadService
}

func InitBean() (*Bean, error) {
	userTimeout, err := strconv.Atoi(getEnvWithDefault("POS_USER_TIMEOUT", "5000"))
	if err != nil {
		return nil, err
	}
	userClient := client.NewUserClient(
		getEnvWithDefault("POS_USER_HOST", "http://localhost:8080"),
		userTimeout,
	)

	loginRepository := persistence.NewLoginRepository(userClient)
	maxFileSize, err := strconv.ParseInt(getEnvWithDefault("POS_FILE_UPLOAD_SIZE", "10"), 0, 64)
	if err != nil {
		return nil, err
	}
	fileUploadService := application.NewFileUploadService(
		getEnvWithDefault("POS_LOGIN_URL", "http://localhost:8080/login"),
		getEnvWithDefault("POS_LOGIN_TOKEN", "pos_access_token"),
		getEnvWithDefault("POS_FILE_UPLOAD_DIRECTORY", "./files"),
		maxFileSize,
		loginRepository)
	return &Bean{FileUploadService: fileUploadService}, nil
}

func getEnvWithDefault(name, def string) string {
	env := os.Getenv(name)
	if len(env) != 0 {
		return env
	}
	return def
}

func getEnvRequired(name string) (string, error) {
	env := os.Getenv(name)
	if len(env) != 0 {
		return env, nil
	}
	return "", errors.New("not found env: " + name)
}
