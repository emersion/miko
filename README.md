# Miko

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
S | 4 | loginresponse | uint1 loginresponsecode
C | 5 | register | str pseudo + str password
S | 6 | registerresponse | uint1 registerresponsecode
S | 7 | playermeta | uint2 entityid + uint1 metaactioncode + bytes metaactionbody
S | 8 | terrainupdate | bytes terrain
C | 9 | terrainrequest | bytes terrainhint
S | 10 | entitiesupdate | bytes entities
C | 11 | entityupdate | bytes entity
S | 12 | actions | bytes actions
C | 13 | action | bytes action
S | 14 | entitycreate | uint2 entityid + bytes entitycreate
S | 15 | entitiydestroy | uint2 entityid + bytes entitydestroy
C | 16 | chatsend | str message
S | 17 | chatreceive | uint2 entityid + str message

#### loginresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | unknownpseudo
2 | wrongpassword
3 | toomanytries
4 | alreadyconnected
5 | playerlimitreached

#### registerresponsecode

Valeur | Signification
--- | ---
0 | ok
1 | usedpseudo
2 | invalidpseudo
2 | invalidpassword
3 | toomanytries
4 | registerdisabled

#### metaactioncode

Valeur | Signification | Contenu
--- | --- | ---
0 | playerjoined | str pseudo
1 | playerleft

### Terrain

#### Description

* blocs de 256x256 : x dans [k;256+k[ ; y dans [l;256+l[
* x vers la droite; y vers le haut
* un uint1 par case pour décrire caractéristiques
* chaque point est décrit par (bx;by;x;y) avec bx et by coordonnées du bloc, x;y coordonnées de la case dans le bloc
* coordonnées d'un bloc: multiple de 256 de ses coordonnées en x et y
* coordonnées d'une case dans un bloc: différence de ses coordonnées par rapport à la case en bas à gauche du bloc

#### Envoi d'un bloc

```
sint2 bx + sint2 by
uint1 defaultvalue
uint2 size
size times:
	uint1 x + uint1 y
    uint1 value
```

* bx;by : coordonnées du bloc
* default : valeur par défaut des cases non envoyées
* size : nombre de cases envoyées
* x, y : coordonnées de la case dans le bloc
* value : valeur de la case

### Entités

```
uint2 entityid
uint1 bitfield
for bit in bitfield:
	bytes data
```

#### Bitfield

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
7 | object | objectbitfield + bytes objectdata

#### Object

* caractéristiques spécifiques au type de l'entité à update
* même principe que précédemment, mais la table de bitfield dépend du type de l'entité

### Session exemple

```
[initiation du ssl, session tcp établie]
S ping
C pong
...
C login(pseudo,password)
S loginresponse(unknownpseudo)
C register(pseudo,password)
S registerresponse(ok)
C login(pseudo,password)
S loginresponse(ok)
S (en broadcast) playermeta(id,playerjoined,pseudo)
S (en broadcast) entitiesupdate(...)
...
S ping
C pong
...
S terrainupdate(terrain)
C entityupdate(bytes)
S (en broadcast) entitiesupdate([...])
S (en broadcast) entitiesupdate([...])
C terrainrequest(hint)
S terrainupdate(terrain)
...
C chat("cc")
S (en broadcast) (id, "cc")
...
C ping
S pong
...
C action(action)
S (en broadcast) actions(...)
...
C exit(player quit)
S (broadcast) entitiesdestroy(id)
S (broadcast) playermeta(id, left)
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
* de manière générale, les commandes qui ont une response sont envoyées en attendant la réponse
* le ping peut se résumer à j'envois un ping si j'ai pas reçu de message depuis 5 secondes, et si j'ai rien pendant 10 secondes d'affilée (en tout) je exit(ping timeout)
* paramètres à faire varier: fréquence d'envoi de entitiesupdate/entityupdate, du terrain (même si il est invariant, on peut dire qu'on envoit par plus petits morceaux mais plus souvent, voir plus précisément en testant), timeout du ping, et quelques autres trucs plus précis dont il faudra parler liés à l'optimisation du multijoueur
* le serveur attend avant d'envoyer des entitiesupdate, par exemple il le fait toutes les n frames, il ne l'envoit pas dès qu'il reçoit une update
* notes d'implémentation: ça risque d'être assez difficile à implémenter: comment gérer les collisions efficacement, comment gérer les murs efficacement ? puisque l'espace mémoire n'est pas critique, il faut privilégier la facilité de programmation (contrairement à quand on envoit par réseau). je pense que ça peut être intéressant de garder le système de segmentation par bloc, parce que ça permettera de "naturellement" filtrer les murs à observer lors d'une collision/d'un affichage, vu qu'on n'aura qu'à s'intéresser à tous ceux d'un bloc donné (éventuellement, il faudrait peut-être, dans le cas où on se retrouverait avec une grande quantité de cases non standard, séparer une liste avec toutes les cases non standard, et une avec juste celles qui interviennent dans les collisions, pour itérer plus vite dans la liste de cases lorsqu'il y a une collision)
* il est très difficile de faire un joli système d'entités (et je ne parle même pas de l'implémentation, mais juste du protocole)
* chaque entité a un entityid, une position
* éventuellement une vitesse (négligeons l'accélération pour l'instant)
* la position sera envoyée en coordonnées cartésiennes, la vitesse en coordonnées _polaires_ ie (x,y,normevitesse,anglevitesse) (on a besoin d'envoyer la vitesse pour faire l'extrapolation (et l'interpolation) des coordonnées des entités entre chaque frame reçue)
* on a pas besoin de précision extrême sur l'angle de la vitesse ou sur sa norme, on enverra donc des _float_
* note d'implémentation: il vaudrait mieux stocker l'angle et la norme sous forme de double pour mieux gérer une accélération éventuelle, mais pas besoin d'autant de précision dans l'envoi puisqu'on mettra à jour suffisament souvent pour corriger les défauts de trajectoire liées aux imprécisions ie (short, short, float, float)
* notes d'implémentation: comment représenter une entité, avec son image et toutes les caractéristiques supplémentaires ? on peut pour l'instant associer un unsgined byte à chaque entité qui correspond à une image, et un unsigned byte qui correspond à son TYPE (joueur ? ennemi ? piège ? interrupteur ? ... ). clairement on aura besoin de caractéristiques supplémentaires pour certains types d'entités, on pourra utiliser éventuellement un système d'héritage avec des classes, mais ça risque de devenir assez lourd, peut etre un héritage par interface? mais ça reste lourd, il faudra aviser.
* Puisque le type d'entité restera constant tout au long de son temps de vie, on enverra son type avec ses caractéristiques supplémentaires pendant entitycreate.
* Il faudra envoyer des caractéristiques supplémentaires en fonction du type de l'entité, au fur et à mesure de son temps de vie, imagine un interrupteur qui passe en mode activé, ou un monstre qui prend une animation particulière lorsqu'il envoit un projectile, ou lorsqu'il dash, ... je propose d'ajouter un bit au bitfield, qui permettra d'ajouter des paramètres spécifiques au type de l'objet
* pour spécifier plus loin, il faudrait donner les types d'entités, et les caractéristiques supplémentaires de chaque type d'entités, et alors on spécifierait entitycreate assez facilement
* pour entitesdestroy, il faudrait peut etre garder ça assez minimal, par exemple si la mort d'un personnage crée une explosion, il faudrait peut etre faire spawn une entité explosion, en tous cas déléguer les taches complexes à d'autres choses que entitiesdestroy qui devrait vraiment être minimaliste et générique
