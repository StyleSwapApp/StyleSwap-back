package utils

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DeleteImageFromS3 : Fonction utilitaire pour supprimer une image de S3

func DeleteImageFromS3(imageURL string) error {
	// Création de la session S3
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	})
	if err != nil {
		return fmt.Errorf("échec de la création de la session AWS : %v", err)
	}

	// Créer un client S3
	s3Client := s3.New(sess)

	// Extraire la clé S3 à partir de l'URL de l'image
	imageKey, err := extractS3KeyFromURL(imageURL)
	if err != nil {
		return fmt.Errorf("échec de l'extraction de la clé S3 depuis l'URL : %v", err)
	}

	// Créer une requête pour supprimer l'objet S3
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String("styleswapbucket"),
		Key:    aws.String(imageKey),
	}

	// Supprimer l'image de S3
	_, err = s3Client.DeleteObject(deleteObjectInput)
	if err != nil {
		return fmt.Errorf("échec de la suppression de l'image sur S3 : %v", err)
	}

	log.Println("Image supprimée avec succès de S3")
	return nil
}

// Fonction utilitaire pour extraire la clé S3 depuis l'URL
func extractS3KeyFromURL(s3URL string) (string, error) {
	const s3Prefix = "https://styleswapbucket.s3.eu-west-3.amazonaws.com/"
	if !strings.HasPrefix(s3URL, s3Prefix) {
		return "", fmt.Errorf("URL S3 invalide")
	}

	return strings.TrimPrefix(s3URL, s3Prefix), nil
}

func UploadToS3(file multipart.File, filename string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-3"),
	})
	HandleError(err, "Error creating a new session")

	uploader := s3.New(sess)

	buffer := make([]byte, file.(sizer).Size())
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("S3_BUCKET")),
		Key:           aws.String(filename),
		Body:          fileBytes,
		ContentLength: aws.Int64(file.(sizer).Size()),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	})
	HandleError(err, "Error uploading image to S3")
	return fmt.Sprintf("https://%s.s3.eu-west-3.amazonaws.com/%s", os.Getenv("S3_BUCKET"), filename), nil

}

type sizer interface {
	Size() int64
}
