package bdd_test

import (
	"bdd/asserter"
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gofiber/fiber/v2"
	flag "github.com/spf13/pflag"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
	Tags:   "~pending",
	Paths:  []string{"example_service"},
}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)

}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run())
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	var dockerTestServices asserter.Asserter
	dockerTestServices.App = fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return err
		},
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		dockerTestServices.ShoutDown()
		return ctx, nil
	})

	dockerTestServices.App.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	// service up
	ctx.Step(`^есть сервис авторизации$`, dockerTestServices.ThereAreAuthorizeService)

	// request
	ctx.Step(`^я делаю (GET|POST|PUT|DELETE)? запрос на ([^"]*)$`, dockerTestServices.MakeRequest)
	ctx.Step(`^я делаю (GET|POST|PUT|DELETE)? запрос по адресу ([^"]*) и телом:$`, dockerTestServices.MakeRequestWithBody)

	// assert response
	ctx.Step(`^код ответа должен быть (\d+)$`, dockerTestServices.AssertResponseCode)
	ctx.Step(`^тело ответа должно соответствовать JSON:$`, dockerTestServices.AssertResponseBody)
	ctx.Step(`^тело ответа должно содержать JSON:$`, dockerTestServices.AssertContainBody)
}
