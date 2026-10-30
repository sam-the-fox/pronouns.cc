package db

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"codeberg.org/pronounscc/pronouns.cc/backend/log"
	"emperror.dev/errors"
	"github.com/BurntSushi/toml"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioCluster struct {
	clients []minioClusterClient
}

type minioClusterClient struct {
	client *minio.Client
	bucket string
}

func NewMinioCluster() (*MinioCluster, error) {
	config, err := os.ReadFile("minio-clients.toml")
	if err != nil {
		log.Infof("could not read minio-clients.toml: %v", err)

		minioClient, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
			Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ACCESS_KEY_ID"), os.Getenv("MINIO_ACCESS_KEY_SECRET"), ""),
			Secure: os.Getenv("MINIO_SSL") == "true",
		})
		if err != nil {
			return nil, errors.Wrap(err, "creating minio client")
		}

		return &MinioCluster{
			clients: []minioClusterClient{{client: minioClient, bucket: os.Getenv("MINIO_BUCKET")}},
		}, nil
	}

	configs := minioClientsConfig{}
	err = toml.Unmarshal(config, &configs)
	if err != nil {
		return nil, fmt.Errorf("reading minio-clients.toml: %w", err)
	}

	cluster := &MinioCluster{clients: []minioClusterClient{}}

	for i := range configs.Clients {
		client, err := minio.New(configs.Clients[i].Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(configs.Clients[i].AccessKeyID, configs.Clients[i].AccessKeySecret, ""),
			Secure: configs.Clients[i].SSL,
		})
		if err != nil {
			return nil, fmt.Errorf("creating minio client for config %d: %w", i, err)
		}

		cluster.clients = append(cluster.clients, minioClusterClient{client: client, bucket: configs.Clients[i].Bucket})
	}

	return cluster, nil
}

func (c *MinioCluster) PutObject(ctx context.Context, objectName string, buf *bytes.Buffer, objectSize int64,
	opts minio.PutObjectOptions,
) (info minio.UploadInfo, err error) {
	for i, client := range c.clients {
		info, err = client.client.PutObject(ctx, client.bucket, objectName, bytes.NewReader(buf.Bytes()), objectSize, opts)
		if err != nil {
			return info, fmt.Errorf("PutObject to client %d: %w", i, err)
		}
	}
	return info, err
}

func (c *MinioCluster) RemoveObject(ctx context.Context, objectName string, opts minio.RemoveObjectOptions) (err error) {
	for i, client := range c.clients {
		err = client.client.RemoveObject(ctx, client.bucket, objectName, opts)
		if err != nil {
			return fmt.Errorf("RemoveObject to client %d: %w", i, err)
		}
	}
	return err
}

func (c *MinioCluster) GetObject(ctx context.Context, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	return c.clients[0].client.GetObject(ctx, c.clients[0].bucket, objectName, opts)
}

type minioClientsConfig struct {
	Clients []MinioClientConfig `toml:"clients"`
}

type MinioClientConfig struct {
	Endpoint        string `toml:"endpoint"`
	AccessKeyID     string `toml:"access_key_id"`
	AccessKeySecret string `toml:"access_key_secret"`
	SSL             bool   `toml:"ssl"`

	Bucket string `toml:"bucket"`

	client *minio.Client
}
