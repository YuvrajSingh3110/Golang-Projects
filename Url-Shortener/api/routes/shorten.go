package routes

import (
	"time"

	"github.com/YuvrajSingh3110/Url_Shortener/helpers"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	Url            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateReset     time.Duration `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx) error {
	body := new(Request)
	if err := c.BodyParser(&body); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse JSON"})
	}

	//rate limiting

	//check if input is an url
	if !govalidator.IsURL(body.Url){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid URL"})
	}

	//checking for domain error to prevent infinite loop
	if !helpers.RemoveDomainError(body.Url){
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error":"you cannot hack this system :)"})
	}

	//enforce https, SSL
	body.Url = helpers.EnforceHTTP(body.Url)
	return nil
}
