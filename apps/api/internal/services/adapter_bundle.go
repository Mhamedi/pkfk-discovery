package services

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/pkfk-discovery/api/internal/integrations/minio"
)

type BundleFile struct {
	Path string
	Data []byte
}

type BundleManifest struct {
	Files map[string]string `json:"files"` // path -> SHA256 checksum
}

type AdapterBundleService struct {
	minioClient *minio.Client
	signingKey  []byte
}

func NewAdapterBundleService(minioClient *minio.Client, signingKey string) *AdapterBundleService {
	return &AdapterBundleService{
		minioClient: minioClient,
		signingKey: []byte(signingKey),
	}
}

func (s *AdapterBundleService) PackageBundle(files []BundleFile) ([]byte, string, error) {
	// Generate checksums
	manifest := BundleManifest{
		Files: make(map[string]string),
	}

	for _, file := range files {
		hash := sha256.Sum256(file.Data)
		manifest.Files[file.Path] = hex.EncodeToString(hash[:])
	}

	// Create tarball
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	// Add manifest
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return nil, "", err
	}

	if err := s.addFileToTar(tw, "bundle.json", manifestJSON); err != nil {
		return nil, "", err
	}

	// Add all files
	for _, file := range files {
		if err := s.addFileToTar(tw, file.Path, file.Data); err != nil {
			return nil, "", err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, "", err
	}
	if err := gzw.Close(); err != nil {
		return nil, "", err
	}

	bundleData := buf.Bytes()

	// Generate HMAC signature
	mac := hmac.New(sha256.New, s.signingKey)
	mac.Write(bundleData)
	signature := hex.EncodeToString(mac.Sum(nil))

	return bundleData, signature, nil
}

func (s *AdapterBundleService) addFileToTar(tw *tar.Writer, path string, data []byte) error {
	hdr := &tar.Header{
		Name: path,
		Mode: 0644,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(data)
	return err
}

func (s *AdapterBundleService) UploadBundle(ctx context.Context, objectName string, bundleData []byte) error {
	reader := bytes.NewReader(bundleData)
	return s.minioClient.UploadAdapterBundle(ctx, objectName, reader, int64(len(bundleData)), "application/gzip")
}

func (s *AdapterBundleService) DownloadBundle(ctx context.Context, objectName string) ([]byte, error) {
	reader, err := s.minioClient.DownloadAdapterBundle(ctx, objectName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func (s *AdapterBundleService) VerifySignature(bundleData []byte, signature string) (bool, error) {
	mac := hmac.New(sha256.New, s.signingKey)
	mac.Write(bundleData)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedSignature)), nil
}

