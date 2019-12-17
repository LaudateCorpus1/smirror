package gcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"smirror/auth"
	"smirror/secret/kms"
)

type service struct {
	afs.Service
}

func (s *service) downloadBase64(ctx context.Context, URL string) (string, error) {
	reader, err := s.Service.DownloadWithURL(ctx, URL)
	if err != nil {
		return "", err
	}
	defer func() { _ = reader.Close() }()
	data, err := ioutil.ReadAll(reader)
	_, err = base64.StdEncoding.DecodeString(string(data))
	if err == nil {
		return string(data), nil
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

//Decrypt decrypts plainText with supplied key
func (s *service) Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error) {
	plainText, err := s.downloadBase64(ctx, secret.URL)
	if err != nil {
		return nil, err
	}
	kmsService, err := cloudkms.NewService(ctx, option.WithScopes(cloudkms.CloudPlatformScope, cloudkms.CloudkmsScope))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create kmsService server for key %v", secret.Key))
	}
	service := cloudkms.NewProjectsLocationsKeyRingsCryptoKeysService(kmsService)
	response, err := service.Decrypt(secret.Key, &cloudkms.DecryptRequest{Ciphertext: plainText}).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decrypt with key %v", secret.Key))
	}
	return []byte(response.Plaintext), nil
}

//New creates GCP kms service
func New(storageService afs.Service) kms.Service {
	return &service{Service: storageService}
}
