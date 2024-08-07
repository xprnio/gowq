# gowq

A terminal-based todo list application written in go.

### Building from source

```sh
# clone the repository to $HOME/.wq
git clone "https://github.com/xprnio/gowq.git" --depth=1 "$HOME/.wq"
cd "$HOME/.wq"

# assuming you have `make` installed, you can simply run
make -B

# otherwise you can also build directly with go
go build -o bin/gowq cmd/main.go

# optionally you can also link the binary to $PATH
ln -s "$(pwd)/bin/gowq" "$HOME/.local/bin/gowq"
```

### Database

The application uses SQLite as a data store.

The application checks each of the following paths (in order) looking for a database file:

1. `$XDG_CONFIG_HOME/wq.sqlite`
2. `$XDG_CONFIG_HOME/wq/database.sqlite`
3. `$HOME/.wq.sqlite`
4. `$HOME/.wq/database.sqlite`
5. `$HOME/.config/wq.sqlite`
6. `$HOME/.config/wq/database.sqlite`
7. `$(pwd)/wq.sqlite`

If none of these paths resolves to a file, `$(pwd)/wq.sqlite` will be used.
