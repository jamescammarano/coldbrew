package utils

import (
	"errors"
	"math/rand"

	"github.com/spaceweasel/promptui"
)

func RandomString(n int) string {
	alphanumeric := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	random := make([]rune, n)
	for r := range random {
		random[r] = alphanumeric[rand.Int63()%int64(len(alphanumeric))]
	}
	return string(random)
}

func MergeMaps(mapArr []map[string]string) map[string]string {
	merged := map[string]string{}

	for _, t := range mapArr {
		for k, v := range t {
			merged[k] = v

		}
	}

	return merged
}

func PromptInput(label string) string {
	validator := func(input string) error {
		if input == "" {
			return errors.New("no input")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}

	res, err := prompt.Run()

	for err != nil {
		res, err = prompt.Run()
	}

	return res
}
