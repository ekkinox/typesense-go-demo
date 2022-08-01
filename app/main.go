package main

import (
	"embed"
	"log"
	"net/http"
	"strconv"

	"github.com/ekkinox/typesense-go-demo/app/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/spf13/viper"
)

//go:embed views/*
var vfs embed.FS

func main() {

	// viper
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	// fiber
	views := html.NewFileSystem(http.FS(vfs), ".html")
	views.Delims("[[", "]]")

	fiberApp := fiber.New(fiber.Config{
		Views: views,
	})

	fiberApp.Use(logger.New())

	// typesense
	typesenseClient := client.InitTypesenseClient()

	// main endpoint
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.Render("views/index", fiber.Map{})
	})

	// search endpoint
	fiberApp.Get("/search", func(c *fiber.Ctx) error {

		expression := c.Query("expression", "")
		page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 0)
		size, _ := strconv.ParseInt(c.Query("size", "20"), 10, 0)
		min, _ := strconv.ParseInt(c.Query("min", "1"), 10, 0)
		max, _ := strconv.ParseInt(c.Query("max", "10000"), 10, 0)
		countries := c.Query("countries", "")

		res, err := typesenseClient.Search(expression, int(page), int(size), countries, int(min), int(max))
		if err != nil {
			return err
		}

		return c.JSON(res)
	})

	// populate endpoint
	fiberApp.Get("/populate", func(c *fiber.Ctx) error {

		wipe := c.Query("wipe", "false")
		if wipe == "true" {
			err := typesenseClient.WipeAndCreateSchema()
			if err != nil {
				return err
			}
		}

		count, _ := strconv.ParseInt(c.Query("count", "500"), 10, 0)
		err := typesenseClient.PopulateSchema(int(count))
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"count": count,
		})

	})

	log.Fatal(fiberApp.Listen(":8080"))

}
