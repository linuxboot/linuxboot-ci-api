package api

import (
	"fmt"
	"time"
)

type Repository struct {
	URL    string  `json:"url"`
	Branch *string `json:"branch"`
}

func (r *Repository) String() string {
	return fmt.Sprintf("%+v", *r)
}

type Job struct {
	ID         int64       `json:"id"`
	Repository *Repository `json:"repository"`
	SubmitDate time.Time   `json:"submitDate"`
	Status     string      `json:"status"`
}

func (j *Job) String() string {
	return fmt.Sprintf("%+v", *j)
}

type Log struct {
	Log string `json:"log"`
}

func (l *Log) String() string {
	return fmt.Sprintf("%+v", *l)
}
