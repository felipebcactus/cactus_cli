package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "encoding/json"
    "strings"
)

var rootCmd = &cobra.Command{
    Use:   "cactus_cli",
    Short: "CLI para interface com API",
    Long:  `CLI flexível e dinâmico para interface com API configurável via arquivo INI.`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) < 1 {
            fmt.Println("Comando não fornecido")
            return
        }
        command := args[0]
        ExecuteAPICommand(command, args[1:])
    },
}

func ExecuteAPICommand(command string, args []string) {
    baseURL := viper.GetString("endpoint.baseurl")
    commandsEndpoint := viper.GetString("endpoint.commands")

    // Função para buscar comandos disponíveis
    fetchCommands := func() ([]string, error) {
        resp, err := http.Get(fmt.Sprintf("%s%s", baseURL, commandsEndpoint))
        if err != nil {
            return nil, fmt.Errorf("erro ao buscar comandos: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return nil, fmt.Errorf("erro na resposta da API: %s", resp.Status)
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
        }

        var availableCommands []string
        if err := json.Unmarshal(body, &availableCommands); err != nil {
            return nil, fmt.Errorf("erro ao parsear resposta da API: %w", err)
        }

        return availableCommands, nil
    }

    // Busca os comandos disponíveis
    availableCommands, err := fetchCommands()
    if err != nil {
        log.Fatalf(err.Error())
    }

    // Verifica se o comando existe
    if !contains(availableCommands, command) {
        fmt.Printf("Comando '%s' não existe. Tentando atualizar comandos...\n", command)
        // Tenta buscar os comandos novamente
        availableCommands, err = fetchCommands()
        if err != nil {
            log.Fatalf(err.Error())
        }
        // Verifica novamente se o comando existe
        if !contains(availableCommands, command) {
            fmt.Printf("Comando '%s' não existe após atualização\n", command)
            return
        }
    }

    // Construir a URL com parâmetros
    params := strings.Join(args, "&")
    url := fmt.Sprintf("%s/cli/%s?%s", baseURL, command, params)

    // Executa o comando via API
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("Erro ao executar comando: %s", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Erro na resposta da API: %s", resp.Status)
    }

    result, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Erro ao ler resposta da API: %s", err)
    }

    fmt.Printf("Resultado do comando '%s':\\n%s\\n", command, string(result))
}

func contains(slice []string, item string) bool {
    for _, a := range slice {
        if a == item {
            return true
        }
    }
    return false
}

// Execute adiciona todos os comandos filho ao comando raiz e seta flags apropriadas.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
}

func initConfig() {
    // Se você deseja sobrescrever as configurações via flags de linha de comando,
    // adicione as flags aqui, e ligue-as com a configuração viper.
    viper.AutomaticEnv() // ler nas variáveis de ambiente que coincidem
}
