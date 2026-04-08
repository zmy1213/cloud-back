package types

type ClusterAppDetail struct {
	ID                 uint64 `json:"id"`
	ClusterUuid        string `json:"clusterUuid"`
	AppName            string `json:"appName"`
	AppCode            string `json:"appCode"`
	AppType            int64  `json:"appType"`
	IsDefault          int64  `json:"isDefault"`
	AppUrl             string `json:"appUrl"`
	Port               int64  `json:"port"`
	Protocol           string `json:"protocol"`
	AuthEnabled        int64  `json:"authEnabled"`
	AuthType           string `json:"authType"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Token              string `json:"token"`
	AccessKey          string `json:"accessKey"`
	AccessSecret       string `json:"accessSecret"`
	TlsEnabled         int64  `json:"tlsEnabled"`
	CaFile             string `json:"caFile"`
	CaKey              string `json:"caKey"`
	CaCert             string `json:"caCert"`
	ClientCert         string `json:"clientCert"`
	ClientKey          string `json:"clientKey"`
	InsecureSkipVerify int64  `json:"insecureSkipVerify"`
	Status             int64  `json:"status"`
	CreatedBy          string `json:"createdBy"`
	UpdatedBy          string `json:"updatedBy"`
	CreatedAt          int64  `json:"createdAt"`
	UpdatedAt          int64  `json:"updatedAt"`
}

type AddClusterAppRequest struct {
	ClusterUuid        string `json:"clusterUuid"`
	AppName            string `json:"appName"`
	AppCode            string `json:"appCode"`
	AppType            int64  `json:"appType"`
	IsDefault          int64  `json:"isDefault"`
	AppUrl             string `json:"appUrl"`
	Port               int64  `json:"port"`
	Protocol           string `json:"protocol"`
	AuthEnabled        int64  `json:"authEnabled"`
	AuthType           string `json:"authType"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Token              string `json:"token"`
	AccessKey          string `json:"accessKey"`
	AccessSecret       string `json:"accessSecret"`
	TlsEnabled         int64  `json:"tlsEnabled"`
	CaFile             string `json:"caFile"`
	CaKey              string `json:"caKey"`
	CaCert             string `json:"caCert"`
	ClientCert         string `json:"clientCert"`
	ClientKey          string `json:"clientKey"`
	InsecureSkipVerify int64  `json:"insecureSkipVerify"`
	UpdatedBy          string `json:"updatedBy"`
}

type ClusterAppDetailRequest struct {
	ID uint64 `path:"id"`
}

type ClusterAppValidateRequest struct {
	ID uint64 `path:"id"`
}

type ClusterAppListRequest struct {
	ClusterUuid string `form:"clusterUuid"`
}
