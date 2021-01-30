package model

import "strings"

type WorkerWrapper struct {
	Index  int    `json:"index"`
	Url    string `json:"url"`
	Result bool   `json:"result"`
}

var keyWords = []string{
	"login",
	"log in",
	"pass",
	"password",
	"name",
	"email",
	"username",
	"sign in",
	"sign up",
}

func IsLoginForm(formFields []string) bool {
	var matchedCount int
	for _, field := range formFields {
		for _, word := range keyWords {
			if strings.Contains(
				strings.ToLower(field), word) {
				matchedCount++
			}
		}
	}
	return matchedCount > 0
}
