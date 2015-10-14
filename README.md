﻿# Miko

An experimental minimalist multiplayer top-down adventure game (en français)

## Organisation du projet

### Client

* `/client`
* java 8 se
* non fonctionnel

### Serveur

* `/server`
* go
* non fonctionnel

## Protocole

### Généralités

* Binaire
* TCP + SSL
* Bigendian
* uint(n) : entier non signé de n bytes
* sint(n) : entier signé de n bytes
* float : nombre décimal sur 4 bytes
* str : chaine de caractères, uint2 contentlength + utf-8 bytes
* messsage : headers + contenu
* header : uint1, type du message
* le stream est fermé après réception ou envoi d'un exit

### Messages

Envoyeur | Valeur | Nom | Contenu
--- | --- | --- | ---
SC | 0 | ping |
SC | 1 | pong |
SC | 2 | exit | uint1 exitcode
C | 3 | login | str pseudo + str password
S | 4 | login_response | uint1 loginresponsecode
C | 5 | register | str pseudo + str password
S | 6 | register_response | uint1 registerresponsecode
S | 7 | meta_action | uint2 entityid + uint1 metaactioncode + bytes metaactionbody
S | 8 | terrain_update | bytes terrain
C | 9 | terrain_request | bytes terrainhint
S | 10 | entities_update | bytes entities
C | 11 | entity_update | bytes entity
S | 12 | actions | bytes actions
C | 13 | action | bytes action
S | 14 | entity_create | uint2 entityid + bytes entity_create
S | 15 | entity_destroy | uint2 entityid + bytes entity_destroy
C | 16 | chat_send | str message
S | 17 | chat_receive | uint2 entityid + str message

#### exitcode

Valeur | Signification
--- | ---
0 | client_quit
1 | server_closed
2 | network_error
3 | ping_timeout
4 | client_kicked
5 | client_banned

#### loginresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | unknown_pseudo
2 | wrong_password
3 | too_many_tries
4 | already_connected
5 | player_limit_reached

#### registerresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | used_pseudo
2 | invalid_pseudo
3 | invalid_password
4 | too_many_tries
5 | register_disabled

#### metaactioncode

Valeur | Signification | Contenu
--- | --- | ---
0 | player_joined | str pseudo
1 | player_left

### Terrain

#### Description

* blocs de 256x256 : x dans [k;256+k[ ; y dans [l;256+l[
* x vers la droite; y vers le haut
* un uint1 par case pour décrire caractéristiques
* chaque point est décrit par (bx;by;x;y) avec bx et by coordonnées du bloc, x;y coordonnées de la case dans le bloc
* coordonnées d'un bloc: multiple de 256 de ses coordonnées en x et y
* coordonnées d'une case dans un bloc: différence de ses coordonnées par rapport à la case en bas à gauche du bloc

#### terrain

```
sint2 bx + sint2 by
uint1 defaultvalue
uint2 size
size times:
	uint1 x + uint1 y
    uint1 value
```

* bx, by : coordonnées du bloc
* default : valeur par défaut des cases non envoyées
* size : nombre de cases envoyées
* x, y : coordonnées de la case dans le bloc
* value : valeur de la case

#### terrain_hint

* le client envoit une liste de blocs dont il veut recevoir le terrain

```
uint1 size
size times:
	sint2 bx + sint2 by
```

* bx, by : coordonnées du bloc

### Entités

* le client peut envoyer plusieurs entity_update si interaction d'objets différents de lui

#### entity_update

```
uint2 entityid
uint1 bitfield
for bit in bitfield:
	bytes data
```

#### entities_update

```
uint2 size
size times:
	uint2 entityid
	uint1 bitfield
	for bit in bitfield:
		bytes data
```

#### bitfield

* les bits sont lus de gauche à droite : 01234567

Bit | Signification | Contenu
--- | --- | ---
0 | position | sint2 bx + sint2 by + uint1 x + uint1 y
1 | speedangle | float angle
2 | speednorm | float norm
3 |
4 |
5 |
6 |
7 | object | uint1 size + bytes objectupdatetypes

#### objectupdatetype

* liste de uint1 indépendants entre eux
* pas de paramètres supplémentaires associés aux bytes

Valeur | Signification | Détail
--- | --- | ---
0 | disabled | objet désactivé (interrupteur, ...)
1 | enabled | object activé (inteerrupteur, ...)

#### entitycreate

#### entitydestroy

### Actions

#### Description

* action caractérisée par id et a paramètres optionnels
* paramètres spécifiés par type de action, type de action<==>type de params
* on n'envoit pas l'actiontype parce que le peer sait le type de chaque action (fixée)
* spells à channel: envoyer un action pour début channel et un pour fin channel

#### action

```
uint2 actionid
bytes params
```

#### actions

```
uint2 size
size times:
	uint2 actionid
	bytes params
```

#### actionid

Valeur | Signification | actiontype
--- | --- | ---

* exemple possible: 0 | heal | entitytarget

#### actiontype

Valeur | Signification | Types de params
--- | --- | ---
0 | paramless |
1 | onefloat | float value
2 | entitytarget | uint2 targetentityid
3 | terraintarget | sint2 bx sint2 by uint1 x uint1 y

### Session exemple

```
[initiation du ssl, session tcp établie]
S ping
C pong
...
C login(pseudo,password)
S login_response(unknownpseudo)
C register(pseudo,password)
S register_response(ok)
C login(pseudo,password)
S login_response(ok)
S (en broadcast) player_meta(id, player_joined, pseudo)
S (en broadcast) entities_update(...)
...
S ping
C pong
...
S terrain_update(terrain)
C entity_update(bytes)
S (en broadcast) entities_update([...])
S (en broadcast) entities_update([...])
C terrain_request(hint)
S terrain_update(terrain)
...
C chat_send("cc")
S (en broadcast) chat_receive(id, "cc")
...
C ping
S pong
...
C action(action)
S (en broadcast) actions(...)
...
C exit(player quit)
S (broadcast) entities_destroy(id)
S (broadcast) player_meta(id, left)
```

## Fonctionnement (todo)

* on garde une liste d'entités, créées et supprimées sur demande du serveur.
* chaque entité est identifiée de manière unique grâce à un "entityid" (envoyable sous forme de unsigned short)
* il y a un nombre fini d'actions possibles de la part des entités, on envoit action avec entityid et byte d'action (params supplémentaires éventuellement)
* on envoit assez souvent des updates des entités (plus précisément sans doute pos+vitesse)
* ainsi, le niveau d'interaction des clients avec le serveur est pos et actions, qui est assez haut niveau, il faudra prendre des mesures de protection au niveau du serveur (par exemple vérifier si l'entité a la capacité de se déplacer là/d'utiliser telle action (cooldown, etc))
* pour l'envoi du terrain, normalement ça se fait automatiquement, mais le client peut request du terrain en particulier
* on va peut etre devoir, si besoin, en fonction du jeu, envoyer une requete explicite depuis le serveur pour supprimer du terrain pour aider le client à poubeller le terrain inutile au fur et à mesure
* niveau implémentation, il faudra trouver moyen d'associer connection/compte à chaque socket ouvert, pour gérer envois de ping et fermeture de session de manière objet/jolie, etc.
* c'est assez général comme protocole, mais il faudra assez spécifier les bytes de action/terrain/entité
* le ping peut se résumer à j'envois un ping si j'ai pas reçu de message depuis 5 secondes, et si j'ai rien pendant 10 secondes d'affilée (en tout) je exit(ping timeout)
* paramètres à faire varier: fréquence d'envoi de entitiesupdate/entityupdate, du terrain (même si il est invariant, on peut dire qu'on envoit par plus petits morceaux mais plus souvent, voir plus précisément en testant), timeout du ping, et quelques autres trucs plus précis dont il faudra parler liés à l'optimisation du multijoueur
* le serveur attend avant d'envoyer des entitiesupdate, par exemple il le fait toutes les n frames, il ne l'envoit pas dès qu'il reçoit une update
* chaque entité a un entityid, une position
* éventuellement une vitesse (négligeons l'accélération pour l'instant)
* la position sera envoyée en coordonnées cartésiennes, la vitesse en coordonnées _polaires_ ie (x,y,normevitesse,anglevitesse) (on a besoin d'envoyer la vitesse pour faire l'extrapolation (et l'interpolation) des coordonnées des entités entre chaque frame reçue)
* on a pas besoin de précision extrême sur l'angle de la vitesse ou sur sa norme, on enverra donc des _float_
* note d'implémentation: il vaudrait mieux stocker l'angle et la norme sous forme de double pour mieux gérer une accélération éventuelle, mais pas besoin d'autant de précision dans l'envoi puisqu'on mettra à jour suffisament souvent pour corriger les défauts de trajectoire liées aux imprécisions ie (short, short, float, float)
* notes d'implémentation: comment représenter une entité, avec son image et toutes les caractéristiques supplémentaires ? on peut pour l'instant associer un unsgined byte à chaque entité qui correspond à une image, et un unsigned byte qui correspond à son TYPE (joueur ? ennemi ? piège ? interrupteur ? ... ). clairement on aura besoin de caractéristiques supplémentaires pour certains types d'entités, on pourra utiliser éventuellement un système d'héritage avec des classes, mais ça risque de devenir assez lourd, peut etre un héritage par interface? mais ça reste lourd, il faudra aviser.
* Puisque le type d'entité restera constant tout au long de son temps de vie, on enverra son type avec ses caractéristiques supplémentaires pendant entitycreate.
* pour spécifier plus loin, il faudrait donner les types d'entités, et les caractéristiques supplémentaires de chaque type d'entités, et alors on spécifierait entitycreate assez facilement
* pour entitesdestroy, il faudrait peut etre garder ça assez minimal, par exemple si la mort d'un personnage crée une explosion, il faudrait peut etre faire spawn une entité explosion, en tous cas déléguer les taches complexes à d'autres choses que entitiesdestroy qui devrait vraiment être minimaliste et générique
* les objets liront un set de objectupdatetype pour s'update et ne s'intéresseront qu'aux bytes qui les concerne, en les interprétant comme ils le veulent