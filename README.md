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
- **Backend :** Node.js / Express
- **Base de données :** MongoDB
- **Authentification :** JWT
- **Paiement :** Stripe
- **Documentation API :** Swagger

## 📦 Installation et démarrage

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/votre-utilisateur/styleswap.git
   ```
2. Accédez au dossier du projet :
   ```bash
   cd styleswap
   ```
3. Installez les dépendances :
   ```bash
   npm install
   ```
4. Configurez vos variables d'environnement dans un fichier `.env`.
5. Lancez le serveur :
   ```bash
   npm start
   ```

L'API sera accessible à l'adresse : `http://localhost:3000`

## 🧪 Tests
Pour exécuter les tests :
```bash
npm test
```

## 📄 Documentation API

### 🔑 Authentification
- **POST /api/auth/register** : Crée un nouvel utilisateur.
  - **Body :**
    ```json
    {
      "username": "string",
      "email": "string",
      "password": "string"
    }
    ```
  - **Réponse :** 201 Created

- **POST /api/auth/login** : Connecte un utilisateur existant.
  - **Body :**
    ```json
    {
      "email": "string",
      "password": "string"
    }
    ```
  - **Réponse :** 200 OK, retourne un token JWT

### 👗 Gestion des annonces
- **GET /api/ads** : Récupère toutes les annonces.
  - **Réponse :** 200 OK, retourne une liste d'annonces

- **POST /api/ads** : Crée une nouvelle annonce.
  - **Body :**
    ```json
    {
      "title": "string",
      "description": "string",
      "price": "number",
      "category": "string"
    }
    ```
  - **Réponse :** 201 Created

- **PUT /api/ads/:id** : Met à jour une annonce existante.
  - **Body :**
    ```json
    {
      "title": "string",
      "description": "string",
      "price": "number"
    }
    ```
  - **Réponse :** 200 OK

- **DELETE /api/ads/:id** : Supprime une annonce.
  - **Réponse :** 204 No Content

### 💬 Messagerie
- **GET /api/messages** : Récupère les messages de l'utilisateur connecté.
  - **Réponse :** 200 OK

- **POST /api/messages** : Envoie un message.
  - **Body :**
    ```json
    {
      "recipientId": "string",
      "content": "string"
    }
    ```
  - **Réponse :** 201 Created

### 💳 Paiements
- **POST /api/payments** : Traite un paiement.
  - **Body :**
    ```json
    {
      "adId": "string",
      "paymentMethodId": "string"
    }
    ```
  - **Réponse :** 200 OK

---
**StyleSwap** – Parce que chaque vêtement mérite une seconde vie ! 👗♻️

*Happy Swapping!* ✌️

