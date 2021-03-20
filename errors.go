package main

import "errors"

var (
	errNoDefinitionsFound = errors.New("no definitions found")
	errNoSearchTerm       = errors.New("no search term provided")
	errInvalidConfig      = errors.New("config invalid. must have a valid app key and id")
)
