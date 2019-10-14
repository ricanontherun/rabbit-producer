package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type RequiredFlagSet struct {
	*flag.FlagSet
	requiredFields []string
}

func NewRequiredFlagSet(required []string) *RequiredFlagSet {
	return &RequiredFlagSet{
		FlagSet:        flag.NewFlagSet("", flag.ContinueOnError),
		requiredFields: required,
	}
}

func (fp *RequiredFlagSet) Parse() error {
	return fp.parseArguments(os.Args[1:])
}

func (fp *RequiredFlagSet) parseArguments(arguments []string) error {
	// Invoke the embedded FlagSet object.
	if parseErr := fp.FlagSet.Parse(arguments); parseErr != nil {
		return parseErr
	}

	seen := make(map[string]bool)
	fp.FlagSet.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})

	// Validate required fields, returning an error as soon as one is found to be missing.
	for _, requiredField := range fp.requiredFields {
		if _, requiredFieldSet := seen[requiredField]; !requiredFieldSet {
			return errors.New(fmt.Sprintf("%s is a required flag", requiredField))
		}
	}

	return nil
}
