package argparser

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

const (
	InvalidArgument = "Invalid arguments, see --help for more info."
)

type Args map[string]string

func (args Args) GetInt(key string, defaultValue int) int {
	val, ok := args[key]
	if !ok {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func Parse() (Args, error) {

	// +1 for program name
	if len(os.Args) < 2 || os.Args[1] == "" {
		return nil, errors.New(InvalidArgument)
	}

	args := make(Args)

	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") {
			if len(os.Args) < i+1 {
				return nil, errors.New(InvalidArgument)
			}

			key := strings.TrimLeft(os.Args[i], "-")
			if len(key) <= 0 {
				return nil, errors.New(InvalidArgument)
			}
			args[key] = os.Args[i+1]
			i++
		} else {
			args[os.Args[i]] = ""
		}
	}

	return args, nil
}
