package config

import (
	"os"
	"strconv"

	"github.com/apex/log"
)

type EnvVar string

const (
	ChatbotCommandInputQueue  EnvVar = "CHATBOT_COMMAND_INPUT_QUEUE"
	ChatbotCommandOutputQueue EnvVar = "CHATBOT_COMMAND_OUTPUT_QUEUE"
	ChatbotMaxWorkers         EnvVar = "CHATBOT_MAX_WORKERS"

	RabbitMQUser EnvVar = "RABBITMQ_USER"
	RabbitMQPass EnvVar = "RABBITMQ_PASS"
	RabbitMQHost EnvVar = "RABBITMQ_HOST"

	MySqlPass EnvVar = "MYSQL_PASSWORD"
	MySqlUser EnvVar = "MYSQL_USER"
	MySqlDB   EnvVar = "MYSQL_DATABASE"
	MySqlHost EnvVar = "MYSQL_HOST"
)

func GetStingEnvVarOrPanic(env EnvVar) string {
	v := os.Getenv(string(env))

	if v == "" {
		log.Fatalf("var %s is required", v)
	}

	return v
}

func GetIntEnvVarOrDefault(env EnvVar, d int) int {
	v := os.Getenv(string(env))

	intv, err := strconv.Atoi(v)
	if err != nil {
		return d
	}

	return intv
}
