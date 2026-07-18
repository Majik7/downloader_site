package main

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var jobs = newJobStore()

func main() {
	_ = godotenv.Load()
	_ = clearDownloadDir()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", welcome)
	r.GET("/hello/", hello)
	r.GET("/hello/:name", hello)
	r.GET("/app", showForm)
	r.POST("/app", submitForm)
	r.GET("/status/:id", jobStatusHandler)
	r.GET("/download/:id", downloadHandler)
	r.GET("/shelltest", shellTest)

	r.Run("localhost:8080")
}

func welcome(c *gin.Context) {
	c.String(http.StatusOK, "<p>Welcome</p>")
}

func hello(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", gin.H{
		"Person": c.Param("name"),
	})
}

func showForm(c *gin.Context) {
	c.HTML(http.StatusOK, "main_form.html", gin.H{})
}

func submitForm(c *gin.Context) {
	videoURL := c.PostForm("url")

	if !isValidURL(videoURL) {
		c.HTML(http.StatusOK, "main_form.html", gin.H{
			"URL":   videoURL,
			"Error": "Please enter a valid URL.",
		})
		return
	}

	id := uuid.New().String()
	outPath := filepath.Join(downloadDir, id+".mp4")
	jobs.create(id)

	go func() {
		if err := downloadVideo(videoURL, outPath); err != nil {
			jobs.setError(id, err.Error())
			return
		}
		jobs.setDone(id, outPath)
	}()

	c.HTML(http.StatusOK, "processing.html", gin.H{
		"JobID": id,
	})
}

func jobStatusHandler(c *gin.Context) {
	j, ok := jobs.get(c.Param("id"))
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"status": "not_found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": j.Status, "error": j.Err})
}

func downloadHandler(c *gin.Context) {
	id := c.Param("id")
	j, ok := jobs.get(id)
	if !ok || j.Status != statusDone {
		c.String(http.StatusNotFound, "Not ready")
		return
	}

	c.FileAttachment(j.Path, "video.mp4")

	// FileAttachment blocks until the response is fully written, so it's
	// safe to clean up immediately after.
	_ = os.Remove(j.Path)
	jobs.delete(id)
}

func shellTest(c *gin.Context) {
	out, err := exec.Command("python", "--version").CombinedOutput()
	if err != nil {
		out = []byte(err.Error())
	}
	c.HTML(http.StatusOK, "shelltest.html", gin.H{
		"Out": string(out),
	})
}
