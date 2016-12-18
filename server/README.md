# Server

[![GoDoc](https://godoc.org/git.emersion.fr/saucisse-royale/miko.git/server?status.svg)](https://godoc.org/git.emersion.fr/saucisse-royale/miko.git/server)
[![Build Status](https://travis-ci.org/emersion/miko.svg?branch=master)](https://travis-ci.org/emersion/miko)

## Installation

```bash
go get -u git.emersion.fr/saucisse-royale/miko.git/server
cd $GOPATH/src/git.emersion.fr/saucisse-royale/miko.git/server
```

## Utilisation

### En local

```bash
make # Pour compiler
# ou
make start # Pour compiler et lancer
```

La compilation va produire un exécutable `miko`.

> Attention ! Le serveur utilise les certificats situés dans `crypto/`.
> Il s'agit de certificats de test, ne pas les utiliser en production !
> Pour en générer de nouveaux, lancer `make crypto`.

### Avec docker

```bash
make docker # Pour construire l'image docker
make start-docker # Pour lancer le serveur dans le conteneur
```
