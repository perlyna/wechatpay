package model

type Bill struct {
	DownloadURL string `json:"download_url"`
	HashType    string `json:"hash_type"`
	HashValue   string `json:"hash_value"`
	TarType     string `json:"tar_type"`
}
