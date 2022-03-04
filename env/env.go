package env

import (
	"os"
	"strconv"
	"time"
)

func GetEnvOrDefault(env, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		value = defaultValue
	}
	return value
}

func GetEnvAsIntOrDefault(env, defaultValue string) int {
	value := GetEnvOrDefault(env, defaultValue)
	intValue, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		panic(err)
	}
	return int(intValue)
}

func GetEnvAsBoolOrDefault(env, defaultValue string) bool {
	value := GetEnvOrDefault(env, defaultValue)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}
	return boolValue
}

func GetEnvAsDurationOrDefault(env, defaultValue string) time.Duration {
	value := GetEnvOrDefault(env, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}
	return duration
}

func GetEnvAsFloatOrDefault(env, defaultValue string) float64 {
	value := GetEnvOrDefault(env, defaultValue)
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return floatValue
}
