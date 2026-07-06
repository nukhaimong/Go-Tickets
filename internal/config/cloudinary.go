package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"regexp"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
	ctx context.Context
}

func NewCloudinaryService() (*CloudinaryService, error) {
	// Cloudinary will automatically read from CLOUDINARY_URL env var
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{
		cld: cld,
		ctx: context.Background(),
	}, nil
}

// UploadEventImage uploads an event photo to Cloudinary
func (s *CloudinaryService) UploadEventImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate a unique public ID with timestamp
	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("events/%d_%s", timestamp, fileHeader.Filename)

	// Remove file extension from public ID (optional)
	// publicID = strings.TrimSuffix(publicID, filepath.Ext(fileHeader.Filename))

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		PublicID: publicID,
		Folder:   "events", // Organize in a folder
		Tags:     []string{"event", "photo"},
	}

	resp, err := s.cld.Upload.Upload(s.ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

// extractPublicIDFromURL extracts the public ID from a Cloudinary URL using regex
func (s *CloudinaryService) extractPublicIDFromURL(photoURL string) (string, error) {
	if photoURL == "" {
		return "", fmt.Errorf("photo URL is empty")
	}

	// Regex to extract public ID
	// Handles: /upload/v1234567890/folder/image.jpg or /upload/folder/image.jpg
	re := regexp.MustCompile(`/upload/(?:v\d+/)?([^?]+)`)
	matches := re.FindStringSubmatch(photoURL)

	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract public ID from URL: %s", photoURL)
	}

	// Remove file extension
	publicID := strings.TrimSuffix(matches[1], path.Ext(matches[1]))

	if publicID == "" {
		return "", fmt.Errorf("extracted public ID is empty")
	}

	return publicID, nil
}

func (s *CloudinaryService) DeleteEventImage(photoURL string) error {
	// Extract public ID from the URL
	publicID, err := s.extractPublicIDFromURL(photoURL)
	if err != nil {
		return fmt.Errorf("failed to extract public ID: %w", err)
	}

	// Delete from Cloudinary
	result, err := s.cld.Upload.Destroy(s.ctx, uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	// Check if the deletion was successful
	if result.Result != "ok" {
		return fmt.Errorf("cloudinary deletion failed: %s", result.Result)
	}

	return nil
}
