package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	//"strings"

	"github.com/gofiber/fiber/v2"
	//"github.com/google/uuid"
)

type FileRequest struct {
	Filename string `json:"filename"`
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error in uploading Image: " + err.Error())
	}
	fileName := file.Filename
	err = c.SaveFile(file, filepath.Join("./files_repo", fileName))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal error saving file: " + err.Error())
	}

	data := map[string]interface{}{

		"fileName": fileName,
		"header":   file.Header,
		"size":     file.Size,
	}

	return c.JSON(fiber.Map{"status": 201, "message": "File has been uploaded and saved successfully", "data": data})
}

func downloadFile(c *fiber.Ctx) error {
	var req FileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing download request body: " + err.Error())
	}
	filePath := filepath.Join("./files_repo", req.Filename)
	return c.SendFile(filePath)
}

func deleteFile(c *fiber.Ctx) error {
	var req FileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error parsing delete request body: " + err.Error())
	}
	filePath := filepath.Join("./files_repo", req.Filename)
	err := os.Remove(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting file: " + err.Error())
	}
	return c.SendString(fmt.Sprintf("File %s deleted successfully", req.Filename))
}

func getAllFileNames(c *fiber.Ctx) error {
	filePath := "./files_repo"
	files, err := os.ReadDir(filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading files directory: " + err.Error())
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return c.JSON(fileNames)
}

func main() {
	app := fiber.New()

	app.Get("/", getAllFileNames)
	app.Post("/upload", uploadFile)
	app.Get("/:fileName", downloadFile)
	app.Delete("/", deleteFile)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
