
### Install git and basic workstation stuff

sudo apt-get install git
git clone https://github.com/tcotav/workstation-init.git
cd workstation-init
ln -s ~/workstation-init/vim ~/.vim && ln -s ~/workstation-init/tmux.conf ~/.tmux.conf && ln -s ~/workstation-init/gitignore ~/.gitignore 
git clone https://github.com/fatih/vim-go.git ~/.vim/bundle/vim-go


### install golang and dev env


wget https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz
tar -xzvf go1.5.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
export GOROOT=/usr/local/go
echo "export GOROOT=/usr/local/go" >> ~/.bashrc
echo "export GOPATH=~/go" >> ~/.bashrc
export GOPATH=~/go
export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:$GOROOT/bin
echo "export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:$GOROOT/bin" >> ~/.bashrc
mkdir -p ~/go/src
#mkdir -p ~/go/bin
cd ~/go

# get the etcd and any other packages you need
go get github.com/coreos/go-etcd/etcd


### Set up etcd

curl -L  https://github.com/coreos/etcd/releases/download/v2.1.2/etcd-v2.1.2-linux-amd64.tar.gz -o etcd-v2.1.2-linux-amd64.tar.gz
tar xzvf etcd-v2.1.2-linux-amd64.tar.gz
cd etcd-v2.1.2-linux-amd64
sudo ./etcd &


