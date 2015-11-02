# Server

## Installation

```bash
go get -u git.emersion.fr/saucisse-royale/miko.git/server
cd $GOPATH/src/git.emersion.fr/saucisse-royale/miko.git/server
make
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
