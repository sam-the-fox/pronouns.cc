package db

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"emperror.dev/errors"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/rs/xid"

	"github.com/chai2010/webp"
)

var (
	webpArgs = []string{"-thumbnail", "512x512^", "-gravity", "center", "-background", "none", "-extent", "512x512", "-quality", "90", "webp:-"}
	jpgArgs  = []string{"-thumbnail", "512x512^", "-gravity", "center", "-extent", "512x512", "-quality", "80", "jpg:-"}
)

const ErrInvalidDataURI = errors.Sentinel("invalid data URI")
const ErrInvalidContentType = errors.Sentinel("invalid avatar content type")

// ConvertAvatar parses an avatar from a data URI, converts it to WebP and JPEG, and returns the results.
func (db *DB) ConvertAvatar(data string) (
	webpOut *bytes.Buffer,
	jpgOut *bytes.Buffer,
	err error,
) {
	data = strings.TrimSpace(data)
	if !strings.Contains(data, ",") || !strings.Contains(data, ":") || !strings.Contains(data, ";") {
		return nil, nil, ErrInvalidDataURI
	}
	split := strings.Split(data, ",")

	rawData, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid base64 data")
	}

	img, _, err := image.Decode(bytes.NewReader(rawData))
	if err != nil {
		return nil, nil, errors.Wrap(err, "decodign image")
	}

	resized := imaging.Fill(img, 512, 512, imaging.Center, imaging.Linear)

	webpOut = new(bytes.Buffer)
	err = webp.Encode(webpOut, resized, &webp.Options{
		Quality: 90,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "encoding WebP image")
	}

	jpgOut = new(bytes.Buffer)
	err = jpeg.Encode(jpgOut, resized, &jpeg.Options{
		Quality: 80,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "encoding JPEG image")
	}

	return webpOut, jpgOut, nil
}

func (db *DB) WriteUserAvatar(ctx context.Context,
	userID xid.ID, webp *bytes.Buffer, jpeg *bytes.Buffer,
) (
	hash string, err error,
) {
	hasher := sha256.New()
	_, err = hasher.Write(webp.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "hashing webp avatar")
	}
	hash = hex.EncodeToString(hasher.Sum(nil))

	_, err = db.minio.PutObject(ctx, db.minioBucket, "/users/"+userID.String()+"/"+hash+".webp", webp, -1, minio.PutObjectOptions{
		ContentType: "image/webp",
	})
	if err != nil {
		return "", errors.Wrap(err, "uploading webp avatar")
	}

	_, err = db.minio.PutObject(ctx, db.minioBucket, "/users/"+userID.String()+"/"+hash+".jpg", jpeg, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", errors.Wrap(err, "uploading jpeg avatar")
	}

	return hash, nil
}

func (db *DB) WriteMemberAvatar(ctx context.Context,
	memberID xid.ID, webp *bytes.Buffer, jpeg *bytes.Buffer,
) (
	hash string, err error,
) {
	hasher := sha256.New()
	_, err = hasher.Write(webp.Bytes())
	if err != nil {
		return "", errors.Wrap(err, "hashing webp avatar")
	}
	hash = hex.EncodeToString(hasher.Sum(nil))

	_, err = db.minio.PutObject(ctx, db.minioBucket, "/members/"+memberID.String()+"/"+hash+".webp", webp, -1, minio.PutObjectOptions{
		ContentType: "image/webp",
	})
	if err != nil {
		return "", errors.Wrap(err, "uploading webp avatar")
	}

	_, err = db.minio.PutObject(ctx, db.minioBucket, "/members/"+memberID.String()+"/"+hash+".jpg", jpeg, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", errors.Wrap(err, "uploading jpeg avatar")
	}

	return hash, nil
}

func (db *DB) DeleteUserAvatar(ctx context.Context, userID xid.ID, hash string) error {
	err := db.minio.RemoveObject(ctx, db.minioBucket, "/users/"+userID.String()+"/"+hash+".webp", minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting webp avatar")
	}

	err = db.minio.RemoveObject(ctx, db.minioBucket, "/users/"+userID.String()+"/"+hash+".jpg", minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting jpeg avatar")
	}

	return nil
}

func (db *DB) DeleteMemberAvatar(ctx context.Context, memberID xid.ID, hash string) error {
	err := db.minio.RemoveObject(ctx, db.minioBucket, "/members/"+memberID.String()+"/"+hash+".webp", minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting webp avatar")
	}

	err = db.minio.RemoveObject(ctx, db.minioBucket, "/members/"+memberID.String()+"/"+hash+".jpg", minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "deleting jpeg avatar")
	}

	return nil
}

func (db *DB) UserAvatar(ctx context.Context, userID xid.ID, hash string) (io.ReadCloser, error) {
	obj, err := db.minio.GetObject(ctx, db.minioBucket, "/users/"+userID.String()+"/"+hash+".webp", minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "getting object")
	}
	return obj, nil
}

func (db *DB) MemberAvatar(ctx context.Context, memberID xid.ID, hash string) (io.ReadCloser, error) {
	obj, err := db.minio.GetObject(ctx, db.minioBucket, "/members/"+memberID.String()+"/"+hash+".webp", minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "getting object")
	}
	return obj, nil
}
