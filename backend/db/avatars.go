package db

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"os/exec"
	"strings"

	"emperror.dev/errors"
	"github.com/minio/minio-go/v7"
	"github.com/rs/xid"
)

var (
	webpArgs = []string{"-resize", "256x256", "-quality", "50", "webp:-"}
	jpgArgs  = []string{"-resize", "256x256", "-quality", "50", "jpg:-"}
)

const ErrInvalidDataURI = errors.Sentinel("invalid data URI")
const ErrInvalidContentType = errors.Sentinel("invalid avatar content type")

// ConvertAvatar parses an avatar from a data URI, converts it to WebP and JPEG, and returns the results.
func (db *DB) ConvertAvatar(data string) (
	webp io.Reader,
	jpg io.Reader,
	err error,
) {
	data = strings.TrimSpace(data)
	if !strings.Contains(data, ",") || !strings.Contains(data, ":") || !strings.Contains(data, ";") {
		return nil, nil, ErrInvalidDataURI
	}
	split := strings.Split(data, ",")
	rest, b64 := split[0], split[1]

	rest = strings.Split(rest, ":")[1]
	contentType := strings.Split(rest, ";")[0]

	var contentArg []string
	switch contentType {
	case "image/png":
		contentArg = []string{"png:-"}
	case "image/jpeg":
		contentArg = []string{"jpg:-"}
	case "image/gif":
		contentArg = []string{"gif:-"}
	case "image/webp":
		contentArg = []string{"webp:-"}
	default:
		return nil, nil, ErrInvalidContentType
	}

	rawData, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, nil, errors.Wrap(err, "invalid base64 data")
	}

	// create webp convert command and get its pipes
	webpConvert := exec.Command("convert", append(contentArg, webpArgs...)...)
	stdIn, err := webpConvert.StdinPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting webp stdin")
	}
	stdOut, err := webpConvert.StdoutPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting webp stdout")
	}

	// start webp command
	err = webpConvert.Start()
	if err != nil {
		return nil, nil, errors.Wrap(err, "starting webp command")
	}

	// write data
	_, err = stdIn.Write(rawData)
	if err != nil {
		return nil, nil, errors.Wrap(err, "writing webp data")
	}
	err = stdIn.Close()
	if err != nil {
		return nil, nil, errors.Wrap(err, "closing webp stdin")
	}

	// read webp output
	webpBuffer := new(bytes.Buffer)
	_, err = io.Copy(webpBuffer, stdOut)
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading webp data")
	}
	webp = webpBuffer

	// finish webp command
	err = webpConvert.Wait()
	if err != nil {
		return nil, nil, errors.Wrap(err, "running webp command")
	}

	// create jpg convert command and get its pipes
	jpgConvert := exec.Command("convert", append(contentArg, jpgArgs...)...)
	stdIn, err = jpgConvert.StdinPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting jpg stdin")
	}
	stdOut, err = jpgConvert.StdoutPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "getting jpg stdout")
	}

	// start jpg command
	err = jpgConvert.Start()
	if err != nil {
		return nil, nil, errors.Wrap(err, "starting jpg command")
	}

	// write data
	_, err = stdIn.Write(rawData)
	if err != nil {
		return nil, nil, errors.Wrap(err, "writing jpg data")
	}
	err = stdIn.Close()
	if err != nil {
		return nil, nil, errors.Wrap(err, "closing jpg stdin")
	}

	// read jpg output
	jpgBuffer := new(bytes.Buffer)
	_, err = io.Copy(jpgBuffer, stdOut)
	if err != nil {
		return nil, nil, errors.Wrap(err, "reading jpg data")
	}
	jpg = jpgBuffer

	// finish jpg command
	err = jpgConvert.Wait()
	if err != nil {
		return nil, nil, errors.Wrap(err, "running jpg command")
	}

	return webp, jpg, nil
}

func (db *DB) WriteUserAvatar(ctx context.Context,
	userID xid.ID, webp io.Reader, jpeg io.Reader,
) (
	webpLocation string,
	jpegLocation string,
	err error,
) {
	webpInfo, err := db.minio.PutObject(ctx, db.minioBucket, "/users/"+userID.String()+".webp", webp, -1, minio.PutObjectOptions{
		ContentType: "image/webp",
	})
	if err != nil {
		return "", "", errors.Wrap(err, "uploading webp avatar")
	}

	jpegInfo, err := db.minio.PutObject(ctx, db.minioBucket, "/users/"+userID.String()+".jpg", jpeg, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", "", errors.Wrap(err, "uploading jpeg avatar")
	}

	return webpInfo.Location, jpegInfo.Location, nil
}

func (db *DB) WriteMemberAvatar(ctx context.Context,
	memberID xid.ID, webp io.Reader, jpeg io.Reader,
) (
	webpLocation string,
	jpegLocation string,
	err error,
) {
	webpInfo, err := db.minio.PutObject(ctx, db.minioBucket, "/members/"+memberID.String()+".webp", webp, -1, minio.PutObjectOptions{
		ContentType: "image/webp",
	})
	if err != nil {
		return "", "", errors.Wrap(err, "uploading webp avatar")
	}

	jpegInfo, err := db.minio.PutObject(ctx, db.minioBucket, "/members/"+memberID.String()+".jpg", jpeg, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return "", "", errors.Wrap(err, "uploading jpeg avatar")
	}

	return webpInfo.Location, jpegInfo.Location, nil
}
