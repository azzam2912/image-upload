package main

import (
	"fmt"
	"log"

	//"strings"

	"github.com/gofiber/fiber/v2"
	//"github.com/google/uuid"
)

func uploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		log.Println("Error in uploading Image : ", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})

	}

	// uniqueId := uuid.New()

	// filename := strings.Replace(uniqueId.String(), "-", "", -1)

	// fileExt := strings.Split(file.Filename, ".")[1]

	// image := fmt.Sprintf("%s.%s", filename, fileExt)
	image := file.Filename

	err = c.SaveFile(file, fmt.Sprintf("./images/%s", image))

	if err != nil {
		log.Println("Error in saving Image :", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	imageUrl := fmt.Sprintf("http://localhost:3000/images/%s", image)

	data := map[string]interface{}{

		"imageName": image,
		"imageUrl":  imageUrl,
		"header":    file.Header,
		"size":      file.Size,
	}

	return c.JSON(fiber.Map{"status": 201, "message": "Image uploaded successfully", "data": data})
}

func downloadImage(c *fiber.Ctx) error {
	return c.SendFile("./images/" + c.Params("imageName"))
}

func main() {
	// Create a new Fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	// Define a route for handling image uploads
	app.Post("/upload", uploadImage)
	app.Get("/images/:imageName", downloadImage)

	// Start the Fiber server on port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
