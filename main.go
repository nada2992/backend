package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Project struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	TechStack   []string `json:"tech_stack"`
	Link        string   `json:"link"`
}

func getProjectsFromFile() []Project {
	file, err := os.ReadFile("projects.json")
	if err != nil {
		return []Project{}
	}

	var projects []Project
	json.Unmarshal(file, &projects)

	return projects
}

func saveProjectsToFile(projects []Project) error {
	data, _ := json.MarshalIndent(projects, "", "  ")
	return os.WriteFile("projects.json", data, 0644)
}

func main() {
	app := gin.Default()

	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Admin-Password")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	app.GET("/api/projects", func(c *gin.Context) {
		c.JSON(http.StatusOK, getProjectsFromFile())
	})

	app.POST("/api/projects", func(c *gin.Context) {
		adminPassword := c.GetHeader("X-Admin-Password")
		if adminPassword != "my_secret_123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Wrong Password"})
			c.Abort()
			return
		}

		var newProject Project
		if err := c.ShouldBindJSON(&newProject); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		projects := getProjectsFromFile()
		projects = append(projects, newProject)
		saveProjectsToFile(projects)

		c.JSON(http.StatusOK, gin.H{"message": "Project added successfully!"})
	})

	app.DELETE("/api/projects/:id", func(c *gin.Context) {
		if c.GetHeader("X-Admin-Password") != "my_secret_123" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		id := c.Param("id")
		projects := getProjectsFromFile()

		var updatedProjects []Project
		for _, p := range projects {
			if p.ID != id {
				updatedProjects = append(updatedProjects, p)
			}
		}

		saveProjectsToFile(updatedProjects)
		c.JSON(http.StatusOK, gin.H{"message": "Deleted!"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Run(":" + port)
}
