sudo apt install golang-go
sudo apt install cobra
go mod init cactus_cli
cobra init .
go get github.com/spf13/viper

sudo rm -rf /usr/local/go
which go
whereis go
sudo rm -rf /usr/local/bin/go
sudo rm -rf /usr/bin/go
wget https://go.dev/dl/go1.20.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
source ~/.profile  # ou source ~/.bashrc, ou source ~/.zshrc

go mod tidy
go build -o cactus_cli
