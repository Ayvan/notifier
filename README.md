**Install**

Ubuntu 14.04

Required Golang ver. >=1.2.1
apt-get install golang-go

apt-get install redis-server

mkdir /home/USER/go

Add GOPATH and GOROOT:
add to ~/.bashrc and ~/.profile and start this command in console:
export GOPATH=/home/USER/go
export GOROOT=/usr/lib/go

Get BeeGo, bee Tool and Redis for Golang:
go get github.com/astaxie/beego
go get github.com/garyburd/redigo/redis
go get github.com/beego/bee
sudo ln -s /home/USER/go/bin/bee /usr/bin


Get project from github:
cd /home/USER/go/src
git clone https://github.com/fintech-fab/iforgetgo.git

Start:
cd /home/USER/go/src/iforgetgo
bee run