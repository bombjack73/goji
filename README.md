#goji

##Overview
Goji is an auto build system for GO

Goji is in early develompment stage.
For now, it's just an oversimplified autobuilding system.
As soon as you save your files, Goji will rebuild and install them.

##Installation
    > go get github.com/bombjack73/goji

You currently need to create manually the log and pid directory
    sudo mkdir /var/log/goji
    sudo chmod 777 /var/log/goji
    sudo mkdir /var/run/goji
    sudo chmod 777 /var/run/goji
    
##Usage
For, now the only usable option is the auto build system that calls *go install* whenever a file is modified in the supplied package directory.
Note: the package dir has to be in the $GOPATH/src/ directory

To run it interactively:
    > goji autobuild github.com/myusername/mypackage

To stop it just press Ctrl+C

To run it in the background:
    > goji autobuild github.com/myusername/mypackage > /var/log/goji/goji.log 2>&1 &

To stop it:
    > kill -2 `cat /var/run/goji/goji.pid`
    
To monitor:
    > tail -f /var/log/goji/goji.log
    