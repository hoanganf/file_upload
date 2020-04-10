package persistence

import (
	"encoding/json"
	"fmt"
	"github.com/hoanganf/file_upload/src/infrastructure/client"
	"github.com/hoanganf/file_upload/src/infrastructure/client/resource"
	"github.com/hoanganf/pos_domain/entity"
	"github.com/hoanganf/pos_domain/entity/exception"
	"github.com/hoanganf/pos_domain/repository"
	"net/http"
)

type LoginRepositoryImpl struct {
	Client *client.UserClient
}

func NewLoginRepository(client *client.UserClient) repository.LoginRepository {
	return &LoginRepositoryImpl{Client: client}
}

func (r *LoginRepositoryImpl) GetUserByJwt(jwt string) (*entity.User, *exception.Error) {
	data := &resource.LoginRequestResource{JWT: jwt}
	resp, err := r.Client.PostRequest(data)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	user, err := r.GetResponse(resp)
	if err != nil {
		return nil, exception.CreateError(exception.CodeUnknown, err.Error())
	}

	return user, nil

}
func (r *LoginRepositoryImpl) GetUserByUserNameAndPassword(userName string, password string) (*entity.User, *exception.Error) {
	//not use
	return nil, exception.CreateError(exception.CodeUnknown, "method not supported.")
}

func (s *LoginRepositoryImpl) GetResponse(resp *http.Response) (*entity.User, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse exception.Error
		err := json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			return nil, fmt.Errorf("can not decode login response. [%w]", err)
		}
		if errResponse.ErrorCode == exception.CodeNotFound {
			return nil, fmt.Errorf("UserName or password is not invalid.")
		}
		return nil, fmt.Errorf("api has error. [%s]", errResponse.ErrorMessage)
	}

	var user entity.User
	err := json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("can not decode user response. [%w]", err)
	}
	return &user, nil
}
