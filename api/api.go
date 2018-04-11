package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ggiamarchi/linuxboot-ci-api/logger"
	"github.com/ggiamarchi/linuxboot-ci-api/utils"
	"gopkg.in/gin-gonic/gin.v1"
)

// Run runs the API server
func Run(port int) {
	logger.Info("Starting API server...")

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      api(),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
	}
	s.ListenAndServe()
}

func api() *gin.Engine {
	api := gin.New()
	api.Use(logger.APILogger(), gin.Recovery())

	v1 := api.Group("/v1")

	healthcheck(v1)
	submitJob(v1)
	getJob(v1)
	getJobLog(v1)

	return api
}

func authenticated(c *gin.Context) bool {
	if c.GetHeader("X-Auth-Secret") == "J4Cpcbq35BGLRm9vilV6Wg4zAUpFDiwd" {
		return true
	}
	return false
}

func healthcheck(api *gin.RouterGroup) {
	api.GET("/healthcheck", func(c *gin.Context) {
		if !authenticated(c) {
			c.Writer.WriteHeader(403)
			return
		}
		c.Writer.WriteHeader(204)
	})
}

func submitJob(api *gin.RouterGroup) {
	api.POST("/jobs", func(c *gin.Context) {
		if !authenticated(c) {
			c.Writer.WriteHeader(403)
			return
		}

		var job Job

		if c.BindJSON(&job) != nil {
			c.Writer.WriteHeader(400)
			return
		}

		if job.Repository.URL == "" {
			c.Writer.WriteHeader(400)
			return
		}

		if job.Repository.Branch != nil && *job.Repository.Branch == "" {
			job.Repository.Branch = nil
		}

		branch := ""
		if job.Repository.Branch != nil {
			branch = *job.Repository.Branch
		}

		stdout, _, err := utils.ExecCommand("/usr/bin/sbatch kvmjob %s %s", job.Repository.URL, branch)

		if err != nil {
			logger.Error("%s", err)
			c.Writer.WriteHeader(500)
			return
		}

		job.SubmitDate = time.Now()
		job.Status = "PENDING"
		job.ID, err = strconv.ParseInt(strings.TrimSuffix(strings.Split(stdout, " ")[3], "\n"), 10, 64)

		if err != nil {
			logger.Error("%s", err)
			c.Writer.WriteHeader(500)
			return
		}

		if os.Mkdir(fmt.Sprintf("/var/lib/ci/%d", job.ID), 0744) != nil {
			logger.Error("%s", err)
			c.Writer.WriteHeader(500)
			return
		}

		c.JSON(202, job)
	})
}

func getJob(api *gin.RouterGroup) {
	api.GET("/jobs/:id", func(c *gin.Context) {
		if !authenticated(c) {
			c.Writer.WriteHeader(403)
			return
		}

		job := Job{}

		var err error
		job.ID, err = strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			logger.Error("%s", err)
			c.Writer.WriteHeader(400)
			return
		}

		jobDir := fmt.Sprintf("/var/lib/ci/%d", job.ID)
		statusFile := fmt.Sprintf("%s/status", jobDir)

		if _, err := os.Stat(statusFile); os.IsNotExist(err) {
			job.Status = "RUNNING"
		} else {
			status, err := ioutil.ReadFile(statusFile)
			if err != nil {
				logger.Error("%s", err)
				c.Writer.WriteHeader(400)
				return
			}

			if strings.TrimSuffix(string(status), "\n") == "0" {
				job.Status = "SUCCESS"
			} else {
				job.Status = "FAILURE"
			}
		}

		c.JSON(200, job)
	})
}

func getJobLog(api *gin.RouterGroup) {
	api.GET("/jobs/:id/logs", func(c *gin.Context) {
		if !authenticated(c) {
			c.Writer.WriteHeader(403)
			return
		}

		job := Job{}

		var err error
		job.ID, err = strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			logger.Error("%s", err)
			c.Writer.WriteHeader(400)
			return
		}

		jobDir := fmt.Sprintf("/var/lib/ci/%d", job.ID)

		userLogFile := fmt.Sprintf("%s/log", jobDir)
		var userLogArray []byte
		logger.Error("1")
		if _, err := os.Stat(userLogFile); err == nil {
			logger.Error("2")
			userLogArray, err = ioutil.ReadFile(userLogFile)
			if err != nil {
				logger.Error("%s", err)
				c.Writer.WriteHeader(400)
				return
			}
		}
		userLogString := strings.TrimSuffix(string(userLogArray), "\n")

		// TODO remove occurrence of 'compute01'
		slurmLogFile := fmt.Sprintf("%s/../compute01.job.%d.out", jobDir, job.ID)
		var slurmLogArray []byte
		if _, err := os.Stat(slurmLogFile); err == nil {
			logger.Error("3")
			slurmLogArray, err = ioutil.ReadFile(slurmLogFile)
			if err != nil {
				logger.Error("%s", err)
				c.Writer.WriteHeader(400)
				return
			}
		}
		slurmLogString := string(slurmLogArray)

		logger.Error("4")
		logString := strings.Replace(slurmLogString, "<% USER LOG PLACEHOLDER %>", userLogString, 1)

		raw := c.Query("raw")

		if raw == "true" {
			c.String(200, logString)
			return
		}

		c.JSON(200, Log{
			Log: logString,
		})
	})
}
