# boulemagique version *7*

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

Valeur | Description de l'animation | Hitbox
--- | --- | ---
0 | placeholder | null
1 | player | circle(10)
2 | ball | circle(10)

### objectattributes

Valeur | Type | Signification | Détail
--- | --- | --- | ---
0 | uint16 | ticks_left | Temps en ticks avant la destruction d'un objet éphémère
1 | uint16 | health | Nombre de points de vie d'une entité
2 | entityid | sender | Créateur de l'entité
30000 | uint16 | cooldown_one | Cooldown en ticks (numéro 1)

## action

### actionid

Valeur | Type | Signification
--- | --- | ---
0 | float angle + uint16 entityid | Envoie une boule dans la direction spécifiée, avec un id temporaire
