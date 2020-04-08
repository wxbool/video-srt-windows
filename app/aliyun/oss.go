package aliyun

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strconv"
	"strings"
	"time"
)

type AliyunOss struct {
	Endpoint string
	AccessKeyId string
	AccessKeySecret string
	BucketName string //yourBucketName
	BucketDomain string //Bucket 域名
}


//获取Buckets列表
func (c AliyunOss) GetListBuckets() ([]string , error) {
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return nil,err
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		return nil,err
	}

	result := []string{}
	for _, bucket := range lsRes.Buckets {
		result = append(result , bucket.Name)
	}

	return result,nil
}


//上传本地文件
//localFileName:本地文件
//objectName:oss文件名称
func (c AliyunOss) UploadFile(localFileName string , objectName string) (string , error) {
	// 创建OSSClient实例
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return "",err
	}
	// 获取存储空间
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return "",err
	}

	//分日期存储
	date := time.Now()
	year := date.Year()
	month := date.Month()
	day  := date.Day()
	objectName = strconv.Itoa(year) + "/" + strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/" + objectName

	// 上传文件
	err = bucket.PutObjectFromFile(objectName , localFileName)
	if err != nil {
		return "",err
	}

	return objectName , nil
}


//删除oss文件
func (c AliyunOss) DeleteFile(objectName string) error {
	// 创建OSSClient实例
	client, err := oss.New(c.Endpoint , c.AccessKeyId , c.AccessKeySecret)
	if err != nil {
		return err
	}
	// 获取存储空间
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		return err
	}

	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}


//获取文件 url link
func (c AliyunOss) GetObjectFileUrl(objectFile string) string {
	if strings.Index(c.BucketDomain, "http://") == -1 && strings.Index(c.BucketDomain, "https://") == -1 {
		return "http://" + c.BucketDomain + "/" +  objectFile
	} else {
		return c.BucketDomain + "/" +  objectFile
	}
}