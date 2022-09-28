package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
)

type H map[string]any

func Success(c *fiber.Ctx, message string, a ...interface{}) error {
	output := H{"_benchmark": bench(c), "type": "success", "success": true, "message": message}
	if len(a) > 0 {
		output["data"] = a[0]
	}
	return c.Status(http.StatusAccepted).JSON(output)

}

func Error(c *fiber.Ctx, errors any) error {
	if reflect.TypeOf(errors).Kind() == reflect.String {
		errors = [1]string{fmt.Sprintf("%v", errors)}
	}
	return c.Status(http.StatusBadRequest).JSON(H{"error": true, "errors": errors})
}

func Template(c *fiber.Ctx, name string, data any) string {
	buffer := new(bytes.Buffer)
	if err := c.App().Config().Views.Render(buffer, name, data); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

func HTML(c *fiber.Ctx, html string) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
	return c.SendString(html)
}

func Render(c *fiber.Ctx, data any, a ...interface{}) error {
	json := H{"_benchmark": bench(c), "data": data}
	if len(a) > 0 {
		json["_meta"] = a[0]
	}
	return c.Status(http.StatusAccepted).JSON(json)
}

func bench(c *fiber.Ctx) float64 {
	benchmark := c.Locals("tonicBenchmark")
	diff := (float64(time.Now().UnixMicro() - benchmark.(int64))) / 1000000
	return diff
}
