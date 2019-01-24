package s3Downloader

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3DownloadInfo struct {
	AccessKeyId     string
	SecretAccessKey string
	Region          string
	BucketName      string
}

func FileDownloadFromS3(s3Info S3DownloadInfo, filePath, fileName string) {
	//ダウンロード保存先に、ダウンロード用のファイルを作成
	file, err := os.Create(fmt.Sprintf("%s%s", filePath, fileName))
	if err != nil {
		fmt.Println("File Create Err.")
		return
	}
	defer file.Close()

	// s3manager作成
	cre := credentials.NewStaticCredentials(
		s3Info.AccessKeyId,
		s3Info.SecretAccessKey,
		"",
	)
	sess := session.New(&aws.Config{
		Credentials: cre,
		Region:      aws.String(s3Info.Region),
	})
	downloader := s3manager.NewDownloader(sess)

	// download実行
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s3Info.BucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		fmt.Println("Download Err.")
		return
	}
}
