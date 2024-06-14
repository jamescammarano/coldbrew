package utils

import (
	"errors"
	"math/rand"
	"strconv"

	"coldbrew.go/cb/common/types"
	"github.com/sirupsen/logrus"
	"github.com/spaceweasel/promptui"
)

func RandomString(n string) string {
	var err error
	num := 128

	if n != "" {
		num, err = strconv.Atoi(n)

		for err != nil {
			logrus.Error(err)
		}

	}

	alphanumeric := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	random := make([]rune, num)

	for r := range random {
		random[r] = alphanumeric[rand.Int63()%int64(len(alphanumeric))]
	}
	return string(random)
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

func GenerateVariables(variables map[string]types.Variable) map[string]string {
	generated := map[string]string{}

	for attr, val := range variables {
		if val.Util == "PromptInput" {
			generated[attr] = PromptInput(attr)
		} else if val.Util == "RandomString" {
			if len(val.Vars) > 0 {
				generated[attr] = RandomString(val.Vars[0])
			} else {
				generated[attr] = RandomString("")
			}
		} else {
			generated[attr] = val.Value
		}
	}
	return generated
}
