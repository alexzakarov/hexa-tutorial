package common

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	ent "main/internal/v1/core_api/domain/entities"
	"math/rand"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	Bundle *i18n.Bundle
)

func init() {
	Bundle = i18n.NewBundle(language.Turkish)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	Bundle.MustLoadMessageFile("i18n/tr.json")
	Bundle.MustLoadMessageFile("i18n/en.json")
}

func HTTPResponser(data interface{}, status bool, message string) fiber.Map {
	return fiber.Map{
		"error":   status,
		"message": message,
		"data":    data,
	}
}

func GenNum() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(999999-100000) + 100000
}

// Get config path for local or docker
func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}

func ValueTrim(value string) string {
	result := strings.TrimSpace(value)
	return result
}

func Placeholder(data []string) string {
	var columns []string
	for i := 1; i <= len(data); i++ {
		columns = append(columns, fmt.Sprintf("$%d", i))
	}
	return strings.Join(columns, ",")
}

func Column(data []string) string {
	var columns []string
	for i := 0; i < len(data); i++ {
		columns = append(columns, data[i])
	}
	return strings.Join(columns, ",")
}

func ThrowError(defination string) error {
	return errors.New(defination)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

// CheckStringIfContains check a string if contains given param
func CheckStringIfContains(input_text string, search_text string) bool {
	CheckContains := strings.Contains(input_text, search_text)
	return CheckContains
}

func CheckNullDate(input_text string) sql.NullTime {
	var parsedDate sql.NullTime
	if input_text != "" {
		t, _ := time.Parse("2006-01-02", input_text)
		parsedDate = sql.NullTime{Time: t, Valid: true}
	} else {
		parsedDate = sql.NullTime{Valid: false}
	}

	return parsedDate
}

func GetAuthIdFromToken(c *fiber.Ctx) (int64, int8) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	AuthID := int64(claims["user_id"].(float64))
	UserType := int8(claims["user_type"].(float64))

	return AuthID, UserType
}

func GetAuthDataFromToken(c *fiber.Ctx) (dat ent.UserData) {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	mapper := ent.UserData{
		UserId:       int64(claims["user_id"].(float64)),
		CurOrgId:     int64(claims["cur_org_id"].(float64)),
		TenantId:     int64(claims["tenant_id"].(float64)),
		PersonalCode: claims["personal_code"].(string),
		AuthTitle:    claims["auth_title"].(string),
		DefLang:      claims["def_lang"].(string),
		OrgIds:       claims["org_ids"].([]interface{}),
	}

	return mapper

}

func GetRawAccessToken(c *fiber.Ctx) (token string) {
	user := c.Locals("user").(*jwt.Token)
	rawToken := user.Raw
	return rawToken
}

func RemoveBasePath(path string) string {
	pathSlice := strings.Split(path, "/")
	pathSlice = pathSlice[2:]
	newPath := fmt.Sprintf("/%s", strings.Join(pathSlice, "/"))
	return newPath
}

func getStructType(struc interface{}) reflect.Type {
	sType := reflect.TypeOf(struc)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	return sType
}

func Convert(struc interface{}) ([]string, []interface{}, error) {

	var returnMap []interface{}
	var returnJsons []string

	sType := getStructType(struc)

	if sType.Kind() != reflect.Struct {
		return returnJsons, returnMap, errors.New("variable given is not a struct or a pointer to a struct")
	}

	for i := 0; i < sType.NumField(); i++ {
		structFieldName := sType.Field(i).Name
		structJsonName := sType.Field(i).Tag.Get("json")
		structDbName := sType.Field(i).Tag.Get("db")
		if structDbName == "" {
			if !strings.Contains(structJsonName, "sub") {
				switch structJsonName {
				case "-":
				case "":
				default:
					parts := strings.Split(structJsonName, ",")
					name := parts[0]
					if name == "" {
						name = structJsonName
					}
					structJsonName = name
				}
				structVal := reflect.ValueOf(struc)
				returnMap = append(returnMap, structVal.FieldByName(structFieldName).Interface())
				returnJsons = append(returnJsons, structJsonName)
			}
		} else {
			if !strings.Contains(structJsonName, "sub") {
				switch structDbName {
				case "-":
				case "":
				default:
					parts := strings.Split(structDbName, ",")
					name := parts[0]
					if name == "" {
						name = structDbName
					}
					structDbName = name
				}
				structVal := reflect.ValueOf(struc)
				returnMap = append(returnMap, structVal.FieldByName(structFieldName).Interface())
				returnJsons = append(returnJsons, structDbName)
			}
		}

	}

	return returnJsons, returnMap, nil
}

func GetQueryIndex(struc interface{}, query_param string) (int64, error) {

	var index int64

	sType := getStructType(struc)

	if sType.Kind() != reflect.Struct {
		return index, errors.New("variable given is not a struct or a pointer to a struct")
	}

	for i := 0; i < sType.NumField(); i++ {
		structJsonName := sType.Field(i).Tag.Get("json")
		structDbName := sType.Field(i).Tag.Get("db")
		if structDbName == "" {
			if !strings.Contains(structJsonName, "sub") {
				if query_param == structJsonName {
					index = int64(i + 2)
					break
				}
			}
		} else {
			if !strings.Contains(structJsonName, "sub") {
				if query_param == structJsonName {
					index = int64(i + 2)
					break
				}
			}
		}

	}
	return index, nil
}

func MapToArr(dat map[string]interface{}) (data []interface{}) {
	for _, value := range dat {
		data = append(data, value)
	}
	return
}

func Translate(lang string, id string) (translation string) {

	localized := i18n.NewLocalizer(Bundle, lang)
	translation = localized.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	return
}

func ConvertStringToIntArray(queryString string) []int64 {
	var err error
	var intVal int64

	stringParts := strings.Split(queryString, ",")
	var intParts []int64

	for _, part := range stringParts {
		intVal, err = strconv.ParseInt(part, 10, 64)

		if intVal != 0 && err == nil {
			intParts = append(intParts, intVal)
		}
	}

	return intParts
}

func ConvertStringToArray(queryString string) []string {
	stringParts := strings.Split(queryString, ",")
	var stringValues []string

	for _, part := range stringParts {
		if part != "" {
			stringValues = append(stringValues, fmt.Sprintf("%s", part))
		}
	}

	return stringValues
}
