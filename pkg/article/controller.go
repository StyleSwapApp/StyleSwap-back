package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ArticleConfig struct {
	*config.Config
}

func New(configuration *config.Config) *ArticleConfig {
	return &ArticleConfig{configuration}
}

func uploadToS3(file multipart.File, filename string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"), // Région de ton bucket
	})
	if err != nil {
		fmt.Printf("Session creation failed: %v\n", err)
		return "", fmt.Errorf("session creation failed: %v", err)
	}

	// Créer un client S3
	s3Client := s3.New(sess)

	// Créer un buffer pour lire le fichier
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		fmt.Printf("Failed to read file into buffer: %v\n", err)
		return "", fmt.Errorf("failed to read file into buffer: %v", err)
	}

	// Télécharger l'image sur S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String("styleswapbucket"),
		Key:         aws.String(filename),         // Nom du fichier sur S3
		Body:        bytes.NewReader(buf.Bytes()), // Contenu du fichier
		ContentType: aws.String("image/png"),      // Type MIME de l'image
	})
	if err != nil {
		fmt.Printf("Failed to upload image to S3: %v\n", err)
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Retourner l'URL S3 de l'image
	imageURL := fmt.Sprintf("https://%s.s3.eu-west-3.amazonaws.com/%s", "styleswapbucket", filename)
	return imageURL, nil
}

func (config *ArticleConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse la requête multipart pour obtenir les données de formulaire
	err := r.ParseMultipartForm(10 << 20) // Limite de taille de 10 Mo pour l'image
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to parse multipart form"})
		return
	}

	// Récupérer les champs de texte du formulaire
	name := r.FormValue("name")
	priceStr := r.FormValue("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid price format"})
		return
	}
	description := r.FormValue("description")

	// Vérification que tous les champs sont fournis
	if name == "" || price == 0 || description == "" {
		render.JSON(w, r, map[string]string{"error": "Missing required fields"})
		return
	}

	// Récupérer le fichier de la requête (champ "image")
	file, _, err := r.FormFile("image")
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to get file from form"})
		return
	}
	defer file.Close()

	// Générer un nom unique pour le fichier sur S3
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("%s.png", uniqueID) // Utiliser un nom unique pour éviter les collisions

	// Télécharger l'image sur S3
	imageURL, err := uploadToS3(file, filename)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to upload image to S3"})
		return
	}

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		Name:        name,
		Price:       price,
		Description: description,
		ImageURL:    imageURL,
	}

	// Ajouter l'article à la base de données
	if err := config.ArticleRepository.Create(articleEntry); err != nil {
		render.JSON(w, r, map[string]string{"error": "Error while adding article to the database"})
		return
	}

	// Répondre avec un message de succès
	render.JSON(w, r, map[string]string{"message": "Article added successfully", "image_url": imageURL})
}
