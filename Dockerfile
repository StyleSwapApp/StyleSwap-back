# Étape 1: Construction de l'application Go
FROM golang:1.23.0-alpine AS builder

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers go mod et go sum pour télécharger les dépendances
COPY go.mod go.sum ./
RUN go mod tidy

# Copier l'ensemble du code source dans le conteneur
COPY . .

# Construire l'application Go
RUN go build -o /app/app .

# Étape 2: Construction de l'image de production
FROM alpine:latest

# Définir le répertoire de travail dans l'image finale
WORKDIR /app

# Copier l'exécutable Go du conteneur builder vers le conteneur de production
COPY --from=builder /app/app .

# Copier également le fichier .env (si tu en as un pour les variables d'environnement)
COPY --from=builder /app/.env ./.env

# Exposer le port de l'API (ici 8080, à adapter si besoin)
EXPOSE 8080

# Commande d'exécution de l'application Go
CMD ["/app/app"]
