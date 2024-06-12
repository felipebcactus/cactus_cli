package main

import (
    "cactus_cli/cmd"
    "github.com/spf13/viper"
    "log"
)

func main() {
    viper.SetConfigName("config")
    viper.SetConfigType("ini")
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.cactus_cli")
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %s", err)
    }
    cmd.Execute()
}
