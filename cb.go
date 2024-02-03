/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "main/coldbrew/cmd"

type Config struct {
	tags          map[string]string
	vars          map[string]string
	extraFiles    map[string]string `yaml:"extra_files,omitempty"`
	restartUnless string            `yaml:"restart_unless,omitempty"`
	addons        map[string]string
}

func main() {
	cmd.Execute()
}
