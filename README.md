# StyleSwap

Bienvenue sur **StyleSwap**, une API conçue pour faciliter le commerce entre particuliers de vêtements de seconde main. StyleSwap vise à offrir une plateforme intuitive et sécurisée pour donner une seconde vie à vos habits tout en contribuant à une mode plus durable.

## 🚀 Objectif du projet
L'objectif de StyleSwap est de permettre aux utilisateurs de vendre facilement leurs anciens vêtements qu'ils ne portent plus et d'acheter des pièces uniques directement auprès d'autres particuliers.

## 📚 Fonctionnalités de l'API
- **Gestion des utilisateurs :** Authentification et gestion sécurisée des comptes utilisateurs.
- **Ajout et gestion des annonces :** Les utilisateurs peuvent ajouter, modifier ou supprimer leurs annonces de vêtements.
- **Système de recherche et de filtres :** Trouvez rapidement les vêtements qui vous intéressent grâce à des filtres avancés.
- **Messagerie intégrée :** Communication directe entre acheteurs et vendeurs.
- **Paiements sécurisés :** Intégration de solutions de paiement sécurisées pour finaliser les transactions.

## 🛠️ Technologies utilisées
- **Backend :** Golang
- **Base de données :** MySQL
- **Authentification :** JWT
- **Documentation API :** Swagger

## A NOTER: PRENDRE LA BRANCHE DEV ET NON MAIN

## 📄 Documentation API

Doc Postman: 

https://documenter.getpostman.com/view/31025137/2sAYQWKtN7

## Dev non répertoriés dans la doc
Fonction messagerie utilisant un websocket
**Pour ce qui est du serveur websocket, voici comment procéder:**

https://styleswap-7059.postman.co/workspace/StyleSwap-Workspace~4e439639-a50e-4301-ad0f-b7a9ec0f8e53/collection/6776a59dfbeb7916c7e6f45a?action=share&creator=31025137&active-environment=31025137-6585c1b9-6f66-4dff-b063-8dde845fdab7

Vous pourrez vous connecter au serveur à cette addresse, 
l'API vous crée une session entre vous(id récupérer grâce au token JWT) et la personne avec qui vous souhaitez parler (idReceveur)

  ws://localhost:8080/api/v1/chat/ws/{idReceveur}

```
{
    "content":"bonjour"
}
```

l'API récupère la conversation passée (stocké dans la BDD) et vous l'envoie


## Devs en cours: 

Fonction de payement
Fonction de recherche 
  - recherche par User
  - recherche par nom d'article
---
**StyleSwap** – Parce que chaque vêtement mérite une seconde vie ! 👗♻️

*Happy Swapping!* ✌️

