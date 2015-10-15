# Server

## Installation

```bash
mkdir -p $GOPATH/src/git.emersion.fr/saucisse-royale
cd $GOPATH/src/git.emersion.fr/saucisse-royale
git clone ssh://git@git.emersion.fr/saucisse-royale/miko.git
cd miko/server
go get ./...
```

## Utilisation

```bash
make # Pour compiler
make start # Pour compiler et lancer
```
