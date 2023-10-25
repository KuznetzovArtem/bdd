package asserter

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/cucumber/godog"
	"github.com/gofiber/fiber/v2"
	"github.com/kinbiko/jsonassert"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Asserter struct {
	App      *fiber.App
	CloseFns []func() error
	Resp     *http.Response
	Err      error
}

func (a *Asserter) ShoutDown() {
	for _, fn := range a.CloseFns {
		fn()
	}
}

func (a *Asserter) Errorf(format string, args ...interface{}) {
	a.Err = fmt.Errorf(format, args...)
}

func (a *Asserter) AssertResponseCode(code int) error {
	assert.Equal(a, code, a.Resp.StatusCode)
	return a.Err
}

func (a *Asserter) AssertResponseBody(expectedBody *godog.DocString) error {
	var actualBody []byte

	actualBody, a.Err = io.ReadAll(a.Resp.Body)
	if a.Err != nil {
		return a.Err
	}

	expectedBody.Content = strings.ReplaceAll(expectedBody.Content, "{{number}}", "\"«PRESENCE»\"")
	expectedBody.Content = strings.ReplaceAll(expectedBody.Content, "{{string}}", "\"«PRESENCE»\"")

	jsonassert.New(a).Assertf(string(actualBody), expectedBody.Content)
	if a.Err != nil {
		a.Err = fmt.Errorf("Err: %v\n\nwant: %s\n\ngot: %s", a.Err, expectedBody.Content, string(actualBody))
	}
	return a.Err
}

type ContainElement struct {
	RepeatCount int         `json:"repeat_count"`
	SearchKey   string      `json:"search_key"`
	SearchValue interface{} `json:"search_value"`
}

func (a *Asserter) AssertContainBody(expectedBody *godog.DocString) error {
	var actualBody []byte

	actualBody, a.Err = ioutil.ReadAll(a.Resp.Body)
	if a.Err != nil {
		return a.Err
	}

	var errs []error

	_, err := jsonparser.ArrayEach([]byte(expectedBody.Content), func(searchValue []byte, dataTypeOfSearchElement jsonparser.ValueType, offsetSearch int, err error) {
		var f ContainElement

		if err := json.Unmarshal(searchValue, &f); err != nil {
			errs = append(errs, err)
			return
		}
		val := extractValue(string(actualBody), f.SearchKey)
		if len(val) != f.RepeatCount {
			errs = append(errs, errors.Errorf("expecred count %v of value %s: actual count %v", f.RepeatCount, f.SearchValue, len(val)))
			return
		}
		searchValueExist := false
		for _, result := range val {
			keyValMatch := strings.Split(string(result), ":")
			value := strings.ReplaceAll(keyValMatch[1], "\"", "")
			switch f.SearchValue.(type) {
			case int64:
				val, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					continue
				}
				if val == f.SearchValue {
					searchValueExist = true
				}
			case float64:
				val, err := strconv.ParseFloat(value, 64)
				if err != nil {
					continue
				}
				if val == f.SearchValue {
					searchValueExist = true
				}
			case string:
				if value == f.SearchValue {
					searchValueExist = true
				}
			default:
				errs = append(errs, errors.Errorf("implement determenation for this case assertContainBody support float64 int64 sting %s \n", f.SearchValue))
				return
			}
		}
		if !searchValueExist {
			errs = append(errs, errors.Errorf("value %v with key %v not found in %s \n", f.SearchValue, f.SearchKey, actualBody))
			return
		}
	})
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		a.Err = fmt.Errorf("Err: %s", errs)
	}

	return a.Err
}

func extractValue(body string, key string) [][]byte {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindAll([]byte(body), len(body))
	return match
}
