package asserter

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-retryablehttp"
	"strings"
)

func (a *Asserter) MakeRequest(method, uri string) error {
	client := retryablehttp.NewClient()

	if client == nil {
		return fmt.Errorf("error got error in client initilize")
	}

	request, err := retryablehttp.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("error cant initilize http client: %v", err)
	}

	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := client.Do(request)

	a.Resp = response
	a.Err = err

	return a.Err
}

func (a *Asserter) MakeRequestWithBody(method, uri string, body *godog.DocString) error {
	client := retryablehttp.NewClient()
	if client == nil {
		return fmt.Errorf("error got error in client initilize")
	}

	reader := strings.NewReader(body.Content)

	request, err := retryablehttp.NewRequest(method, uri, reader)
	if err != nil {
		return fmt.Errorf("error cant initilize http client: %v", err)
	}

	request.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := client.Do(request)

	a.Resp = response
	a.Err = err

	return a.Err
}
