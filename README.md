# Overview

TODO

# License

MIT

# Installation

    $ go get -u github.com/leafi/blunt

(-u forces go get to hit the network, even if you already have the package.)

go get will teach you how to set up $GOPATH if you haven't done it already.

# Update Built Assets & Run

Assets are in $GOPATH/src/github.com/leafi/blunt/assets.

    $ cd $GOPATH/src/github.com/leafi/blunt
    $ go generate; go build
    $ ./blunt

go generate will install go-bindata if you don't already have it. Do 'go generate -n' instead to see what shell commands would be run.

(TODO: vvvvv THIS IS IN PROGRESS vvvvv)
The final binary can be distributed standalone, thanks to the magic of Go. If run on a system where blunt can't find itself in $GOPATH - or where there isn't a $GOPATH - blunt will offer to either write an assets folder to disk or save nothing. Try it!

# Forking

First, fork the repository on GitHub.

    $ go get -u github.com/leafi/blunt
    $ cd $GOPATH/github.com/leafi/blunt
    $ git remote add fork https://github.com/YOUR_USER/blunt.git

    $ git push fork  # to push changes

If this isn't satisfactory because you want a long-term fork, you'll need to fix all the imports & asset folder detection code in the codebase to point to your fork yourself. 

Sorry. I don't know of a better solution.