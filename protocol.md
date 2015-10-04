# Etude du jeu Miko

## Précisions sur le fonctionnement

* on garde une liste d'entités, créées et supprimées sur demande du serveur.
* chaque entité est identifiée de manière unique grâce à un "entityid" (envoyable sous forme de unsigned short)
* il y a un nombre fini d'actions possibles de la part des entités, on envoit action avec entityid et byte d'action (params supplémentaires éventuellement)
* on envoit assez souvent des updates des entités (plus précisément sans doute pos+vitesse)
* ainsi, le niveau d'interaction des clients avec le serveur est pos et actions, qui est assez haut niveau, il faudra prendre des mesures de protection au niveau du serveur (par exemple vérifier si l'entité a la capacité de se déplacer là/d'utiliser telle action (cooldown, etc))
* pour l'envoi du terrain, normalement ça se fait automatiquement, mais le client peut request du terrain en particulier
* on va peut etre devoir, si besoin, en fonction du jeu, envoyer une requete explicite depuis le serveur pour supprimer du terrain pour aider le client à poubeller le terrain inutile au fur et à mesure
* niveau implémentation, il faudra trouver moyen d'associer connection/compte à chaque socket ouvert, pour gérer envois de ping et fermeture de session de manière objet/jolie, etc.

## Protocole

* Binaire, TCP
* ne pas envoyer de OK/ack systématique (on est en tcp)
* tout est envoyé en BIGENDIAN
* str: unsigned short contentlength, utf-8 bytes
* on parse tout en mode: unsigned byte d'enum puis contenu du message
* contenu du message dépend de enum
  par exemple message vide peut avoir 0 contenu: pas besoin de spécifier de content-length avant le message tout le temps


* 0:SC ping
* 1:SC pong
* 2:SC exit (un unsigned byte exit code ? (server déco/client déco/whatever))
* 3:C login (str pseudo, str password)
* 4:S login response 1 unsigned byte responsecode (ok, usr inconnu, mauvais pwd, too many tries, déjà connecté, limite de joueurs atteinte, ...)
* 5:C register (str pseudo, str password)
* 6:S register response 1 unsigned byte responsecode (ok, usr déjà utilisé, usr impossible, pwd impossible, too many tries, création désactivée, ...)
* 7:S playermeta (unsigned short entityid, unsigned byte metaaction (joined(str playername), left)
* 8:S terrainupdate (bytes terrain (à voir plus précisément))
* 9:C terrainrequest (éventuellement un hint pour où on a besoin de terrain)
* 10:S entitiesupdate (bytes entities (pour chaque entité, bytes de position?))
* 11:C entityupdate(bytes entité (celle du joueur))
* 12:S actions(bytes actions(eventuellement des parametres supplémentaires en fonction de l'action))
* 13:C action(unsigned byte action(idem))
* 14:S entitescreate (unsigned short entityid, bytes des params)
* 15:S entitiesdestroy (unsigned short entityid, peutetre bytes des params)
* 16:C chat(str message)
* 17:S chat(unsigned short entityid, str message)

* pour faire un event global, on pourra faire actions avec comme senderid un nombre spécial (unsignedshort.max_value?)
* pour faire parler un monstre, il suffit de chat avec comme entityid son id

exemple de jeu:

```
[initiation du ssl, session tcp établie]
S ping
C ping
...
C login(usr,pwd)
S login_response (usr inconnu)
C register(usr,pwd)
S register_response(ok)
C login(usr,pwd)
S login_response(ok)
S (en broadcast) playermeta(id,joined(usr))
S (en broadcast) entitiesupdate([id, bytes de pos/vitesse])
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
S exit(server exit/player bannedmdr) OU C exit(player quit) S (broadcast) entitiesdestroy(id) S (broadcast) playermeta(id, left)
[celui qui a envoyé exit peut fermer direct le stream, celui qui le reçoit peut fermer direct le stream]
```

* c'est assez général comme protocole, mais il faudra assez spécifier les bytes de action/terrain/entité
* de manière générale, les commandes qui ont une response sont envoyées en attendant la réponse
* le ping peut se résumer à j'envois un ping si j'ai pas reçu de message depuis 5 secondes, et si j'ai rien pendant 10 secondes d'affilée (en tout) je exit(ping timeout)
* paramètres à faire varier: fréquence d'envoi de entitiesupdate/entityupdate, du terrain (même si il est invariant, on peut dire qu'on envoit par plus petits morceaux mais plus souvent, voir plus précisément en testant), timeout du ping, et quelques autres trucs plus précis dont il faudra parler liés à l'optimisation du multijoueur
* le serveur attend avant d'envoyer des entitiesupdate, par exemple il le fait toutes les n frames, il ne l'envoit pas dès qu'il reçoit une update

## Terrain

exemple :

```
 y
 ^
B|AB
C|DB
-+--> x
A|DC
```

avec A, B, C, D des types de cases

* on peut l'imaginer comme un tableau infini de bytes, un byte en chaque case, définissant le type de case
* on peut l'imaginer comme une liste infinie de cases de murs (où on stocke leurs coordonnées)
* terrainupdate: on envoit un bloc de 256x256 cases
* 2 short : coordonnées du coin en bas à gauche du bloc (multiple de 256 en X et en Y)
* unsigned byte : type de case par défaut du bloc (si l'on envoit pas de case à des coordonnées, elle aura cette valeur)
* unsigned short : nombre de cases à envoyer
(nombre de cases à envoyer) fois:
	* 2 unsigned byte : coordonnées de la case dans le bloc par rapport au coin en bas à gauche

notes d'implémentation: ça risque d'être assez difficile à implémenter: comment gérer les collisions efficacement, comment gérer les murs efficacement ? puisque l'espace mémoire n'est pas critique, il faut privilégier la facilité de programmation (contrairement à quand on envoit par réseau). je pense que ça peut être intéressant de garder le système de segmentation par bloc, parce que ça permettera de "naturellement" filtrer les murs à observer lors d'une collision/d'un affichage, vu qu'on n'aura qu'à s'intéresser à tous ceux d'un bloc donné (éventuellement, il faudrait peut-être, dans le cas où on se retrouverait avec une grande quantité de cases non standard, séparer une liste avec toutes les cases non standard, et une avec juste celles qui interviennent dans les collisions, pour itérer plus vite dans la liste de cases lorsqu'il y a une collision)

## Entités

* il est très difficile de faire un joli système d'entités (et je ne parle même pas de l'implémentation, mais juste du protocole)
* chaque entité a un entityid, une position
* éventuellement une vitesse (négligeons l'accélération pour l'instant)
* la position sera envoyée en coordonnées cartésiennes, la vitesse en coordonnées _polaires_
  ie, (x,y,normevitesse,anglevitesse)
  (on a besoin d'envoyer la vitesse pour faire l'extrapolation (et l'interpolation) des coordonnées des entités entre chaque frame reçue)
* on a pas besoin de précision extrême sur l'angle de la vitesse ou sur sa norme, on enverra donc des _float_
* note d'implémentation: il vaudrait mieux stocker l'angle et la norme sous forme de double pour mieux gérer une accélération éventuelle, mais pas besoin d'autant de précision dans l'envoi puisqu'on mettra à jour suffisament souvent pour corriger les défauts de trajectoire liées aux imprécisions
  ie (short, short, float, float)

je propose le système suivant:
pour chaque entité à update: les entités qui ne sont pas updatées continuent à agir normalement
* 1 unsigned short : entityid
* 1 byte de BITFIELD qui spécifie les données à update à ce tour ci sur cette entité
* pour chaque bit de bitfield, de gauche à droite:
  * lire depuis le stream les choses à update pour ce bit
    exemple: pour le bit position on lira deux short, x puis y
* bitfield : position ; anglevitesse ; normevitesse ; unused ; unused ; unused ; unused ; unused
* il est extensible au sens où on pourra ajouter des valeurs au bitfield
* il est efficace parce qu'on enverra à chaque fois que ce qui a besoin d'être update pour chaque entité
* pour les caractéristiques des entités qui sont les mêmes pendant toute la durée de vie de l'entité, on préfèrera les envoyer par entitycreate au moment de la création

notes d'implémentation: comment représenter une entité, avec son image et toutes les caractéristiques supplémentaires ? on peut pour l'instant associer un unsgined byte à chaque entité qui correspond à une image, et un unsigned byte qui correspond à son TYPE (joueur ? ennemi ? piège ? interrupteur ? ... ). clairement on aura besoin de caractéristiques supplémentaires pour certains types d'entités, on pourra utiliser éventuellement un système d'héritage avec des classes, mais ça risque de devenir assez lourd, peut etre un héritage par interface? mais ça reste lourd, il faudra aviser.

Puisque le type d'entité restera constant tout au long de son temps de vie, on enverra son type avec ses caractéristiques supplémentaires pendant entitycreate.

Il faudra envoyer des caractéristiques supplémentaires en fonction du type de l'entité, au fur et à mesure de son temps de vie, imagine un interrupteur qui passe en mode activé, ou un monstre qui prend une animation particulière lorsqu'il envoit un projectile, ou lorsqu'il dash, ... je propose d'ajouter un bit au bitfield, qui permettra d'ajouter des paramètres spécifiques au type de l'objet

D'où : bitfield : position ; anglevitesse ; normevitesse ; unused ; unused ; unused ; unused ; caracssupplémentaires
et les bytes de caracssupplémentaires seront :
* 1 byte de BITFIELD spécifiant les choses à update
* dépendent du type de l'entité (le serveur et le client connaissent déjà le type de l'entité donc pas besoin de le renvoyer
* par exemple pour un interrupteur on peut imaginer le bitfield : activé ; unused ; unused ; unused ; unused ; unused ; unused ; unused
* pour chaque bit de bitfield, lire des bytes du stream
* pour l'interrupteur par exemple on lira un bit (en fait paddé à un byte puisqu'on lit byte par byte) qui correspond à son état d'activation

pour spécifier plus loin, il faudrait donner les types d'entités, et les caractéristiques supplémentaires de chaque type d'entités, et alors on spécifierait entitycreate assez facilement

pour entitesdestroy, il faudrait peut etre garder ça assez minimal, par exemple si la mort d'un personnage crée une explosion, il faudrait peut etre faire spawn une entité explosion, en tous cas déléguer les taches complexes à d'autres choses que entitiesdestroy qui devrait vraiment être minimaliste et générique
