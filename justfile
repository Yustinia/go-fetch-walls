# list recipes
default:
    just --list

# runner
run:
    go run .

# builder
build:
    go build .

# install
install: build
    mkdir -pv $HOME/.config/go-fetch-walls
    cp -rv $(pwd)/configs/settings.json $HOME/.config/go-fetch-walls/
