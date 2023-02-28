package ones

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginONES(url, email, password string) (*http.Response, error) {
	if len(url) == 0 || len(email) == 0 || len(password) == 0 {
		log.Println("login param missing", url, email, password)
		return nil, common.Errors(common.ParameterMissingError, nil)
	}
	body := new(LoginRequest)
	body.Email = email
	body.Password = password
	url = common.GenApiUrl(url, loginUri)
	resp, err := utils.PostJSON(url, body)
	if err != nil {
		return nil, common.Errors(common.NetworkError, nil)
	}
	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		respBody := new(loginErrorResponse)
		if err = json.Unmarshal(data, &respBody); err != nil {
			return nil, err
		}
		return nil, common.Errors(common.AccountError, respBody)
	}
	return resp, nil
}
