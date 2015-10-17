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
# ou
make start # Pour compiler et lancer
```

La compilation va produire un exécutable `miko`.

> Attention ! Le serveur utilise les certificats situés dans `crypto/`.
> Il s'agit de certificats de test, ne pas les utiliser en production !
> Pour en générer de nouveaux, lancer `make crypto`.
