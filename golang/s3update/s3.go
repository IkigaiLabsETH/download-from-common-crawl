package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"example.com/s3update/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetObject(
	ctx context.Context,
	s3Client *s3.Client,
	s3BucketId string,
	key string,
) ([]byte, error) {
	log.Printf("GetObject: %s", key)
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s3BucketId,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}

	return getBytes(resp)
}

func GetWarc(
	ctx context.Context,
	client *s3.Client,
	warc *types.Warc,
) ([]byte, error) {
	s3BucketId := "commoncrawl"
	s3Range := fmt.Sprintf("bytes=%d-%d", warc.Offset, warc.Offset+warc.Length-1)
	resp, err := client.GetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: &s3BucketId,
			Key:    &warc.Filename,
			Range:  &s3Range,
		})
	if err != nil {
		return nil, err
	}

	return getBytes(resp)
}

func getBytes(resp *s3.GetObjectOutput) ([]byte, error) {
	defer resp.Body.Close()
	log.Printf("ContentLength %d", resp.ContentLength)
	var bbuffer bytes.Buffer
	buffer := make([]byte, resp.ContentLength)
	for {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			_, werr := bbuffer.Write(buffer[:num])
			if werr != nil {
				return nil, fmt.Errorf("error writing response to a buffer: %s", werr)
			}
		} else if rerr == io.EOF {
			break
		} else if rerr != nil {
			return nil, fmt.Errorf("error reading response: %s", rerr)
		}
	}
	return bbuffer.Bytes(), nil
}
