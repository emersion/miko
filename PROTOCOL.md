# Protocole version *9*

## Généralités

* Binaire
* TCP + SSL
* Bigendian
* uint(n) : entier non signé de n bits
* sint(n) : entier signé de n bits
* float : nombre décimal sur 4 bytes
* str : chaine de caractères, uint16 contentlength + utf-8 bytes
* messsage : headers + contenu
* header : uint8, type du message
* le stream est fermé après réception ou envoi d'un exit

## Messages

Envoyeur | Valeur | Tick | Nom | Contenu
--- | --- | --- | --- | ---
SC | 0 | Non | ping |
SC | 1 | Non | pong |
SC | 2 | Non | exit | uint8 exitcode
C | 3 | Non | login | str pseudo + str password
S | 4 | Non [1] | login_response | uint8 loginresponsecode
C | 5 | Non | register | str pseudo + str password
S | 6 | Non | register_response | uint8 registerresponsecode
S | 7 | Oui | meta_action | uint16 entityid + uint8 metaactioncode + bytes metaactionbody
S | 8 | Oui | chunk_update | bytes chunk
C | 9 | Non | terrain_request | bytes terrainhint
S | 10 | Oui | entities_update | bytes entities
C | 11 | Oui | entity_update | bytes entity
S | 12 | Oui | actions_done | bytes actions
C | 13 | Oui | action_do | bytes action
S | 14 | Oui | entity_create | uint16 entityid + bytes entity_create
S | 15 | Oui | entity_destroy | uint16 entityid
C | 16 | Non | chat_send | str message
S | 17 | Oui | chat_receive | uint16 entityid + str message
C | 18 | Non | version | uint16 versionid
S | 19 | Non | config | bytes config
S | 20 | Non | entity_id_change | uint16 oldentityid + uint16 newentityid
S | 21 | Oui | chunks_update | bytes chunks

* Si le message possède un tick, il l'envoit avant son contenu : headers puis tick puis contenu. Le tick est un uint16 et est la frame de logique actuelle sur le simulateur du jeu envoyant le message, elle revient à 0 après avoir atteint son maximum (2^16 - 1)
* [1] : Si c'est le code "ok", la réponse contient en plus : uint16 tick | uint64 timestamp

### exitcode

Valeur | Signification
--- | ---
0 | client_quit
1 | server_closed
2 | network_error
3 | ping_timeout
4 | client_kicked
5 | client_banned
6 | client_outdated
7 | server_outdated

### loginresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | unknown_pseudo
2 | wrong_password
3 | too_many_tries
4 | already_connected
5 | player_limit_reached

Un tick est envoyé après le code si celui-ci est "ok".

### registerresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | used_pseudo
2 | invalid_pseudo
3 | invalid_password
4 | too_many_tries
5 | register_disabled

### metaactioncode

Valeur | Signification | Contenu
--- | --- | ---
0 | player_joined | str pseudo
1 | player_left

### versionid

Identifiant unique de la version du protocole utilisé par le client.

### config

```
uint16 maxRollbackTicks
uint16 timeServerPort
```

* maxRollbackTicks : nombre maximum de ticks où l'on peut revenir dans le passé pour appliquer des actions
* timeServerPort : le port du serveur de temps

## Terrain

### Description

* blocs de 256x256 : x dans [k;256+k[ ; y dans [l;256+l[
* x vers la droite; y vers le haut
* un uint8 par case pour décrire caractéristiques
* chaque point est décrit par (bx;by;x;y) avec bx et by coordonnées du bloc, x;y coordonnées de la case dans le bloc
* coordonnées d'un bloc: multiple de 256 de ses coordonnées en x et y
* coordonnées d'une case dans un bloc: différence de ses coordonnées par rapport à la case en bas à gauche du bloc

### chunk

```
sint16 bx + sint16 by
uint8 defaultvalue
uint16 size
size times:
	uint8 x + uint8 y
	uint8 value
```

* bx, by : coordonnées du bloc
* default : valeur par défaut des cases non envoyées
* size : nombre de cases envoyées
* x, y : coordonnées de la case dans le bloc
* value : valeur de la case

### chunks

```
uint16 size
size times:
	chunk
```

### terrain_hint

Le client envoit une liste de blocs dont il veut recevoir le terrain.

```
uint8 size
size times:
	sint16 bx + sint16 by
```

* bx, by : coordonnées du bloc

## Entités

Le client peut envoyer plusieurs entity_update si interaction d'objets différents de lui.

### entityid

Un identifiant unique de l'entité, représenté par un uint16. Les 1000 dernières valeurs possibles (de 64536 à 65535 inclus) sont réservées pour des ids temporaires générés par le client et ne doivent jamais être générées par le serveur.

### entity_id_change

Un action peut provoquer la création d'une entité. Si tel est le cas, le client va générer un entityid temporaire. Lorsque le serveur va créer l'entité, un entityid permanent sera attribué. Le serveur doit alors envoyer un entity_id_change au client qui a provoqué l'action.

### entity_update

```
uint16 entityid
uint8 bitfield
for bit in bitfield:
	bytes data
```

### entity_create

```
entity_update
```

### entities_update

```
uint16 size
size times:
	entity_create
```

### bitfield

Les bits sont lus de gauche à droite : 01234567.

Bit | Signification | Contenu
--- | --- | ---
0 | position | sint16 bx + sint16 by + uint8 x + uint8 y
1 | speedangle | float angle
2 | speednorm | float norm
3 |
4 |
5 | entitytype | uint16 entitytype
6 | sprite | uint16 sprite
7 | object | bytes objectattributes

### entitytype

Un identifiant unique du type de l'entité.

### sprite

Un identifiant unique de l'animation (graphique) d'une entité. Définit aussi une hitbox pour cette entité.

### objectattributes

* Paires de (type;valeur) correspondants à des attributs spécifiques à des entités
* Un message object ne va que mettre à jour la paire qu'il spécifie

```
uint16 size
size times:
    uint16 attribute
    bytes value
```

### hitbox

Définit le contour d'une entité. Le barycentre du contour est toujours centré sur la position de l'entité.

* null: aucun contour
* circle(float radius): un cercle de rayon radius
* rectangle(float width, float height): un rectangle de longueur width et de hauteur height

## Actions

### Description

Action caractérisée par un id et des paramètres optionnels.

### action_do

```
uint16 actionid
bytes params
```

### actions_done

```
uint16 size
size times:
	uint16 entityid
	action_do
```

## Serveur de temps

* Permet au client de synchroniser son horloge avec celle du serveur
* UDP, écoute sur le port spécifié dans `config`
* Lorsqu'il reçoit un paquet UDP vide, il  envoie régulièrement et pendant un certain temps des `uint64` qui représentent le timestamp Unix actuel en microsecondes

## Session exemple (OUTDATED, TODO)

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

## Notes de fonctionnement

* le ping peut se résumer à j'envois un ping si j'ai pas reçu de message depuis 5 secondes, et si j'ai rien pendant 10 secondes d'affilée (en tout) je exit(ping timeout)
* spells à channel: envoyer un action pour début channel et un pour fin channel
