package article

import (
	"StyleSwap/config"
	"StyleSwap/database/dbmodel"
	"StyleSwap/pkg/model"
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/render"
)

type ArticleConfig struct {
	*config.Config
}

func New(configuration *config.Config) *ArticleConfig {
	return &ArticleConfig{configuration}
}

func uploadToS3(file multipart.File, filename string) (string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	})
	if err != nil {
		return "", fmt.Errorf("session creation failed: %v", err)
	}

	// Créer un client S3
	s3Client := s3.New(sess)

	// Créer un buffer pour lire le fichier
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file into buffer: %v", err)
	}

	// Télécharger l'image sur S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("mon-bucket-images"), // Remplace par ton bucket S3
		Key:    aws.String(filename),            // Nom du fichier sur S3
		Body:   bytes.NewReader(buf.Bytes()),    // Contenu du fichier
		ContentType: aws.String("image/png"),    // Type MIME du fichier
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Retourner l'URL S3 de l'image
	imageURL := fmt.Sprintf("https://%s.s3.eu-west-3.amazonaws.com/%s", "styleswapbucket", filename)
	return imageURL, nil
}


func (config *ArticleConfig) ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Initialiser la structure ArticleRequest pour lire les données du body
	req := &model.ArticleRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Parse le formulaire multipart (pour pouvoir récupérer le fichier)
	err := r.ParseMultipartForm(10 << 20) // Limite de taille de 10 Mo pour l'image
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to parse multipart form"})
		return
	}

	// Récupérer le fichier de la requête (champ "image")
	file, _, err := r.FormFile("image")
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Unable to get file from form"})
		return
	}
	defer file.Close()

	// Nom du fichier pour S3
	filename := fmt.Sprintf("%s.png", req.Name) // Utiliser le nom de l'article comme base

	// Télécharger le fichier sur S3
	imageURL, err := uploadToS3(file, filename)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to upload image to S3"})
		return
	}

	// Créer un nouvel article dans la base de données
	articleEntry := &dbmodel.ArticleEntry{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
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
