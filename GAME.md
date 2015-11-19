# boulemagique version *6*

## config

Tous les champs spécifiés dans le protocole, suivis de :

```
float defaultPlayerSpeed
uint16 playerBallCooldown
float defaultBallSpeed
uint16 defaultBallLifespan
```

* defaultPlayerSpeed : vitesse par défaut des joueurs
* playerBallCooldown : cooldown de lancer de boules par les joueurs
* defaultBallSpeed : vitesse par défaut des boules
* defaultBallLifespan : temps de vie des boules

## entity

### entitytype

Valeur | Description du type
--- | ---
0 | player
1 | ball

### sprite

Valeur | Description de l'animation
--- | ---
0 | placeholder
1 | player
2 | ball

### objectattributes

Valeur | Type | Signification | Détail
--- | --- | --- | ---
0 | oneshort | ticks_left | Temps en ticks avant la destruction d'un objet éphémère
1 | oneshort | health | Nombre de points de vie d'une entité
2 | oneentity | sender | entity_id du créateur de l'entité
30000 | oneshort | cooldown_one | Cooldown en ticks (numéro 1)

## action

### actionid

Valeur | Type | Signification
--- | --- | ---
0 | onefloat | Envoit une boule dans la direction spécifiée
