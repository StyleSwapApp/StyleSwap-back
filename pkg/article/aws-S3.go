package article

import (
	"StyleSwap/utils"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
		ACL:         aws.String("public-read"),    // Permissions de l'image
	})
	if err != nil {
		fmt.Printf("Failed to upload image to S3: %v\n", err)
		return "", fmt.Errorf("failed to upload image to S3: %v", err)
	}

	// Retourner l'URL S3 de l'image
	imageURL := fmt.Sprintf("https://%s.s3.eu-west-3.amazonaws.com/%s", "styleswapbucket", filename)
	return imageURL, nil
}

func extractS3KeyFromURL(s3URL string) (string, error) {
	// Vérifier si l'URL commence bien par le préfixe
	const s3Prefix = "https://styleswapbucket.s3.eu-west-3.amazonaws.com/"
	if !strings.HasPrefix(s3URL, s3Prefix) {
		return "", fmt.Errorf("URL S3 invalide")
	}

	key := strings.TrimPrefix(s3URL, s3Prefix)

	return key, nil
}

func (config *ArticleConfig) DeleteImageFromS3(IdArticle int) error {

	// Initialiser une session AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	})
	utils.HandleError(err, "Failed to create AWS session")

	// Récupérer url à partir de l'id de l'article
	urlImage, err:= config.ArticleRepository.FindImageByID(IdArticle)
	utils.HandleError(err, "Failed to get image URL from article ID")

	// Créer un client S3
	s3Client := s3.New(sess)
	imageKey, err := extractS3KeyFromURL(urlImage)
	utils.HandleError(err, "Failed to extract S3 key from URL")

	// Créer une requête pour supprimer l'objet
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String("styleswapbucket"),
		Key:    aws.String(imageKey),
	}

	// Supprimer l'image de S3
	_, err = s3Client.DeleteObject(deleteObjectInput)
	utils.HandleError(err, "Failed to delete image from S3")
	log.Println("Image deleted successfully from S3")
	return nil
}
