package common

import (
	"crypto/rand"
	"fmt"
	"github.com/duclmse/fengine/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"strings"
)

var (
	ErrGenerateOTP = errors.New("Generate OTP failed")
)

var otpLetters = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
var tokenLetters = [...]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

type ExtraUserInfoReq struct {
	UserId          string
	Owner           string
	RoleId          string
	GroupPermission string
	CurrentOrg      string
	ProjectId       string
}

func RandomString(max int, otp bool) (string, error) {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if err != nil || n != max {
		return "", err
	}
	if otp {
		for i := 0; i < len(b); i++ {
			b[i] = otpLetters[int(b[i])%len(otpLetters)]
		}
	} else {
		for i := 0; i < len(b); i++ {
			b[i] = tokenLetters[int(b[i])%len(tokenLetters)]
		}
	}
	return string(b), nil
}

//CheckGroupPermission - check groupPermission for a resource
//groupPermission format: "uuid1, uuid2, uuid3, ..."
func CheckGroupPermission(resourceGroup, groupPermission string) bool {
	if groupPermission == "" {
		return true
	}
	if groupPermission != "" && resourceGroup == "" {
		return false
	}
	arrGroupId := strings.Split(groupPermission, ", ")
	for _, groupId := range arrGroupId {
		if resourceGroup == groupId {
			return true
		}
	}
	return false
}

func PrepareSQLConditionGrantedGroups(groupPermission string) string {
	arrGroupId := strings.Split(groupPermission, ", ")
	var grantedGroupArr []string
	for _, groupId := range arrGroupId {
		grantedGroupArr = append(grantedGroupArr, fmt.Sprintf("'%s'", groupId))
	}
	return fmt.Sprintf("(%s)", strings.Join(grantedGroupArr, ", "))
}

func PrepareMongoConditionGrantedGroups(groupPermission string) bson.A {
	arrGroupId := strings.Split(groupPermission, ", ")
	var grantedGroupArr bson.A
	for _, groupId := range arrGroupId {
		grantedGroupArr = append(grantedGroupArr, groupId)
	}
	return grantedGroupArr
}

type Attribute struct {
	AttributeType string      `json:"attribute_type"`
	AttributeKey  string      `json:"attribute_key"`
	Logged        bool        `json:"logged"`
	Value         interface{} `json:"value"`
	ValueAsString string      `json:"value_as_string"`
	LastUpdateTs  int64       `json:"last_update_ts"`
	ValueType     string      `json:"value_type"`
}

type AttrValue struct {
	AttributeType string `json:"attribute_type"`
	AttributeKey  string `json:"attribute_key"`
	Logged        bool   `json:"logged"`
	Value         string `json:"value"`
	ValueType     string `json:"value_t"`
}

type AttributeQuery struct {
	QueryType  string        `json:"query_type"`
	EntityType string        `json:"entity_type"`
	Queries    []AttributeKV `json:"queries"`
}

type AttributeKV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
