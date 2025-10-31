package settings

type Minio struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	BucketName      string `json:"bucket_name"`
	Secure          bool   `json:"secure"`
}

var MinioSettings = &Minio{}
