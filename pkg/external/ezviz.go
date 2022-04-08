package external

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/duclmse/fengine/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type EzvizHttpRequester struct {
	ezvizHost  string
	appKey     string
	appSecret  string
	httpClient http.Client
}

type ResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type ResponseGetToken struct {
	Code    string        `json:"code"`
	Message string        `json:"msg"`
	Data    TokenResponse `json:"data"`
}
type ResponseCreateSubAccount struct {
	Code    string           `json:"code"`
	Message string           `json:"msg"`
	Data    CreateSubAccount `json:"data"`
}

type TokenResponse struct {
	AccessToken string  `json:"accessToken"`
	AreaDomain  string  `json:"areaDomain"`
	ExpireTime  float64 `json:"expireTime"`
}

type CreateSubAccount struct {
	AccountId string `json:"accountId"`
}

func NewClient(ezvizHost, appKey, appSecret string) EzvizHttpRequester {
	return EzvizHttpRequester{
		ezvizHost:  ezvizHost,
		appKey:     appKey,
		appSecret:  appSecret,
		httpClient: http.Client{},
	}
}

func (h EzvizHttpRequester) GetAppToken() (TokenResponse, error) {
	body := url.Values{}
	body.Set("appKey", h.appKey)
	body.Set("appSecret", h.appSecret)
	endpoint := h.ezvizHost + "/api/lapp/token/get"

	fmt.Printf("request get app token [%s], [%s], [%s] \n", endpoint, body.Get("appKey"), body.Get("appSecret"))
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return TokenResponse{}, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}
	newStr := buf.String()

	var responseBody ResponseGetToken
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return TokenResponse{}, err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return TokenResponse{}, errors.New(responseBody.Message)
	}

	return responseBody.Data, nil
}

func (h EzvizHttpRequester) CreateSubAccount(token, areaDomain, accountName, password string) (string, error) {
	passEncode := md5.Sum([]byte(fmt.Sprintf("%s#%s", h.appKey, password)))

	body := url.Values{}
	body.Set("accessToken", token)
	body.Set("accountName", accountName)
	body.Set("password", strings.ToLower(hex.EncodeToString(passEncode[:])))

	fmt.Printf("request create sub account [%s] [%s], [%s], [%s]\n", areaDomain, token, accountName, strings.ToLower(hex.EncodeToString(passEncode[:])))

	endpoint := areaDomain + "/api/lapp/ram/account/create"
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	newStr := buf.String()

	var responseBody ResponseCreateSubAccount
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return "", err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return "", errors.New(responseBody.Message)
	}

	return responseBody.Data.AccountId, nil
}

func (h EzvizHttpRequester) GetSubAccountToken(token, areaDomain, accountId string) (TokenResponse, error) {
	body := url.Values{}
	body.Set("accessToken", token)
	body.Set("accountId", accountId)

	endpoint := areaDomain + "/api/lapp/ram/token/get"
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return TokenResponse{}, err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return TokenResponse{}, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return TokenResponse{}, err
	}
	newStr := buf.String()

	var responseBody ResponseGetToken
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return TokenResponse{}, err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return TokenResponse{}, errors.New(responseBody.Message)
	}

	return responseBody.Data, nil
}

func (h EzvizHttpRequester) CreateDevice(token, areaDomain, deviceSerial, validateCode string) error {
	body := url.Values{}
	body.Set("accessToken", token)
	body.Set("deviceSerial", deviceSerial)
	body.Set("validateCode", validateCode)

	endpoint := areaDomain + "/api/lapp/device/add"
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	newStr := buf.String()

	var responseBody ResponseBody
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return errors.New(responseBody.Message)
	}

	return nil
}

func (h EzvizHttpRequester) AddPermission(token, areaDomain, accountId, deviceSerial string) error {
	type Statement struct {
		Permission string   `json:"Permission"`
		Resource   []string `json:"Service"`
	}
	var resource []string
	resource = append(resource, fmt.Sprintf("dev:%s", deviceSerial))
	resource = append(resource, fmt.Sprintf("cam:%s:1", deviceSerial))
	statement := Statement{Permission: "Get,Update,DevCtrl", Resource: resource}

	statementBody, _ := json.Marshal(statement)

	body := url.Values{}
	body.Set("accessToken", token)
	body.Set("accountId", accountId)
	body.Set("statement", string(statementBody))

	endpoint := areaDomain + "/api/lapp/ram/statement/add"
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	newStr := buf.String()

	var responseBody ResponseBody
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return errors.New(responseBody.Message)
	}

	return nil
}

func (h EzvizHttpRequester) DeleteDevice(token, areaDomain, deviceSerial string) error {
	body := url.Values{}
	body.Set("accessToken", token)
	body.Set("deviceSerial", deviceSerial)

	endpoint := areaDomain + "/api/lapp/device/delete"
	reqHttp, err := http.NewRequest("POST", endpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	reqHttp.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(reqHttp)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	newStr := buf.String()

	var responseBody ResponseBody
	if err := json.Unmarshal([]byte(newStr), &responseBody); err != nil {
		return err
	}
	code, _ := strconv.Atoi(responseBody.Code)
	if code != 200 {
		return errors.New(responseBody.Message)
	}

	return nil
}
