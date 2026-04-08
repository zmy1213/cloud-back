package apprepo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

type ClusterApp struct {
	ID                 uint64
	ClusterUuid        string
	AppName            string
	AppCode            string
	AppType            int64
	IsDefault          int64
	AppUrl             string
	Port               int64
	Protocol           string
	AuthEnabled        int64
	AuthType           string
	Username           string
	Password           string
	Token              string
	AccessKey          string
	AccessSecret       string
	TlsEnabled         int64
	CaFile             string
	CaKey              string
	CaCert             string
	ClientCert         string
	ClientKey          string
	InsecureSkipVerify int64
	Status             int64
	CreatedBy          string
	UpdatedBy          string
	CreatedAt          int64
	UpdatedAt          int64
}

type UpsertParams struct {
	ClusterUuid        string
	AppName            string
	AppCode            string
	AppType            int64
	IsDefault          int64
	AppUrl             string
	Port               int64
	Protocol           string
	AuthEnabled        int64
	AuthType           string
	Username           string
	Password           string
	Token              string
	AccessKey          string
	AccessSecret       string
	TlsEnabled         int64
	CaFile             string
	CaKey              string
	CaCert             string
	ClientCert         string
	ClientKey          string
	InsecureSkipVerify int64
	UpdatedBy          string
}

type Service struct {
	mu     sync.RWMutex
	apps   []ClusterApp
	nextID uint64

	db    *sql.DB
	useDB bool
}

func NewService(mysqlCfg appcfg.MysqlConfig) *Service {
	s := &Service{apps: make([]ClusterApp, 0), nextID: 1}
	if !mysqlCfg.Enabled || strings.TrimSpace(mysqlCfg.DataSource) == "" {
		return s
	}

	db, err := sql.Open("mysql", mysqlCfg.DataSource)
	if err != nil {
		return s
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return s
	}
	s.db = db
	s.useDB = true
	return s
}

func (s *Service) Upsert(params UpsertParams) error {
	params.ClusterUuid = strings.TrimSpace(params.ClusterUuid)
	params.AppName = strings.TrimSpace(params.AppName)
	params.AppCode = strings.TrimSpace(params.AppCode)
	params.AppUrl = strings.TrimSpace(params.AppUrl)
	params.Protocol = strings.ToLower(strings.TrimSpace(params.Protocol))
	params.AuthType = strings.TrimSpace(params.AuthType)
	params.UpdatedBy = strings.TrimSpace(params.UpdatedBy)
	if params.UpdatedBy == "" {
		params.UpdatedBy = "system"
	}

	if s.useDB {
		err := s.upsertToDB(params)
		if err != nil {
			log.Printf("[apprepo] op=Upsert source=db cluster_uuid=%q app_code=%q app_type=%d error=%v", params.ClusterUuid, params.AppCode, params.AppType, err)
			return err
		}
		log.Printf("[apprepo] op=Upsert source=db cluster_uuid=%q app_code=%q app_type=%d", params.ClusterUuid, params.AppCode, params.AppType)
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	for idx := range s.apps {
		if s.apps[idx].ClusterUuid == params.ClusterUuid && s.apps[idx].AppCode == params.AppCode && s.apps[idx].AppType == params.AppType {
			app := &s.apps[idx]
			app.AppName = params.AppName
			app.IsDefault = params.IsDefault
			app.AppUrl = params.AppUrl
			app.Port = params.Port
			app.Protocol = params.Protocol
			app.AuthEnabled = params.AuthEnabled
			app.AuthType = params.AuthType
			app.Username = params.Username
			app.Password = params.Password
			app.Token = params.Token
			app.AccessKey = params.AccessKey
			app.AccessSecret = params.AccessSecret
			app.TlsEnabled = params.TlsEnabled
			app.CaFile = params.CaFile
			app.CaKey = params.CaKey
			app.CaCert = params.CaCert
			app.ClientCert = params.ClientCert
			app.ClientKey = params.ClientKey
			app.InsecureSkipVerify = params.InsecureSkipVerify
			app.UpdatedBy = params.UpdatedBy
			app.UpdatedAt = now
			log.Printf("[apprepo] op=Upsert source=memory mode=update cluster_uuid=%q app_code=%q app_type=%d", params.ClusterUuid, params.AppCode, params.AppType)
			return nil
		}
	}

	s.apps = append(s.apps, ClusterApp{
		ID:                 s.nextID,
		ClusterUuid:        params.ClusterUuid,
		AppName:            params.AppName,
		AppCode:            params.AppCode,
		AppType:            params.AppType,
		IsDefault:          params.IsDefault,
		AppUrl:             params.AppUrl,
		Port:               params.Port,
		Protocol:           params.Protocol,
		AuthEnabled:        params.AuthEnabled,
		AuthType:           params.AuthType,
		Username:           params.Username,
		Password:           params.Password,
		Token:              params.Token,
		AccessKey:          params.AccessKey,
		AccessSecret:       params.AccessSecret,
		TlsEnabled:         params.TlsEnabled,
		CaFile:             params.CaFile,
		CaKey:              params.CaKey,
		CaCert:             params.CaCert,
		ClientCert:         params.ClientCert,
		ClientKey:          params.ClientKey,
		InsecureSkipVerify: params.InsecureSkipVerify,
		Status:             1,
		CreatedBy:          params.UpdatedBy,
		UpdatedBy:          params.UpdatedBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	})
	s.nextID++
	log.Printf("[apprepo] op=Upsert source=memory mode=insert cluster_uuid=%q app_code=%q app_type=%d", params.ClusterUuid, params.AppCode, params.AppType)
	return nil
}

func (s *Service) ListByClusterUUID(clusterUUID string) ([]ClusterApp, error) {
	clusterUUID = strings.TrimSpace(clusterUUID)
	if s.useDB {
		items, err := s.listByClusterUUIDFromDB(clusterUUID)
		if err != nil {
			log.Printf("[apprepo] op=ListByClusterUUID source=db cluster_uuid=%q error=%v", clusterUUID, err)
			return nil, err
		}
		return items, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]ClusterApp, 0)
	for _, item := range s.apps {
		if item.ClusterUuid != clusterUUID {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Service) GetByID(id uint64) (ClusterApp, bool, error) {
	if id == 0 {
		return ClusterApp{}, false, nil
	}

	if s.useDB {
		item, ok, err := s.getByIDFromDB(id)
		if err != nil {
			log.Printf("[apprepo] op=GetByID source=db id=%d error=%v", id, err)
			return ClusterApp{}, false, err
		}
		return item, ok, nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.apps {
		if item.ID == id {
			return item, true, nil
		}
	}
	return ClusterApp{}, false, nil
}

func (s *Service) UpdateStatus(id uint64, status int64, updatedBy string) error {
	if id == 0 {
		return errors.New("invalid app id")
	}
	updatedBy = strings.TrimSpace(updatedBy)
	if updatedBy == "" {
		updatedBy = "system"
	}

	if s.useDB {
		_, err := s.db.Exec(`
UPDATE onec_cluster_app
SET status = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0
`, status, updatedBy, id)
		if err != nil {
			log.Printf("[apprepo] op=UpdateStatus source=db id=%d status=%d error=%v", id, status, err)
			return err
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for idx := range s.apps {
		if s.apps[idx].ID == id {
			s.apps[idx].Status = status
			s.apps[idx].UpdatedBy = updatedBy
			s.apps[idx].UpdatedAt = time.Now().Unix()
			return nil
		}
	}
	return errors.New("app not found")
}

func (s *Service) upsertToDB(params UpsertParams) error {
	query := `
INSERT INTO onec_cluster_app (
	cluster_uuid, app_name, app_code, app_type, is_default,
	app_url, port, protocol,
	auth_enabled, auth_type, username, password, token,
	access_key, access_secret,
	tls_enabled, ca_file, ca_key, ca_cert, client_cert, client_key,
	insecure_skip_verify, created_by, updated_by, is_deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
ON DUPLICATE KEY UPDATE
	app_name = VALUES(app_name),
	is_default = VALUES(is_default),
	app_url = VALUES(app_url),
	port = VALUES(port),
	protocol = VALUES(protocol),
	auth_enabled = VALUES(auth_enabled),
	auth_type = VALUES(auth_type),
	username = VALUES(username),
	password = VALUES(password),
	token = VALUES(token),
	access_key = VALUES(access_key),
	access_secret = VALUES(access_secret),
	tls_enabled = VALUES(tls_enabled),
	ca_file = VALUES(ca_file),
	ca_key = VALUES(ca_key),
	ca_cert = VALUES(ca_cert),
	client_cert = VALUES(client_cert),
	client_key = VALUES(client_key),
	insecure_skip_verify = VALUES(insecure_skip_verify),
	updated_by = VALUES(updated_by),
	is_deleted = 0,
	updated_at = CURRENT_TIMESTAMP
`

	_, err := s.db.Exec(
		query,
		params.ClusterUuid,
		params.AppName,
		params.AppCode,
		params.AppType,
		params.IsDefault,
		params.AppUrl,
		params.Port,
		params.Protocol,
		params.AuthEnabled,
		params.AuthType,
		params.Username,
		params.Password,
		params.Token,
		params.AccessKey,
		params.AccessSecret,
		params.TlsEnabled,
		params.CaFile,
		params.CaKey,
		params.CaCert,
		params.ClientCert,
		params.ClientKey,
		params.InsecureSkipVerify,
		params.UpdatedBy,
		params.UpdatedBy,
	)
	return err
}

func (s *Service) listByClusterUUIDFromDB(clusterUUID string) ([]ClusterApp, error) {
	rows, err := s.db.Query(`
SELECT id, cluster_uuid, app_name, app_code, app_type, is_default,
       app_url, port, protocol,
       auth_enabled, auth_type, username, password, token,
       access_key, access_secret,
       tls_enabled, ca_file, ca_key, ca_cert, client_cert, client_key,
       insecure_skip_verify, status,
       created_by, updated_by, created_at, updated_at
FROM onec_cluster_app
WHERE cluster_uuid = ? AND is_deleted = 0
ORDER BY id DESC
`, clusterUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ClusterApp, 0)
	for rows.Next() {
		item, scanErr := scanClusterApp(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Service) getByIDFromDB(id uint64) (ClusterApp, bool, error) {
	row := s.db.QueryRow(`
SELECT id, cluster_uuid, app_name, app_code, app_type, is_default,
       app_url, port, protocol,
       auth_enabled, auth_type, username, password, token,
       access_key, access_secret,
       tls_enabled, ca_file, ca_key, ca_cert, client_cert, client_key,
       insecure_skip_verify, status,
       created_by, updated_by, created_at, updated_at
FROM onec_cluster_app
WHERE id = ? AND is_deleted = 0
LIMIT 1
`, id)

	item, err := scanClusterApp(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ClusterApp{}, false, nil
		}
		return ClusterApp{}, false, err
	}
	return item, true, nil
}

type appScanner interface {
	Scan(dest ...any) error
}

func scanClusterApp(scanner appScanner) (ClusterApp, error) {
	var item ClusterApp
	var token sql.NullString
	var caFile sql.NullString
	var caKey sql.NullString
	var caCert sql.NullString
	var clientCert sql.NullString
	var clientKey sql.NullString
	var createdAt time.Time
	var updatedAt time.Time

	err := scanner.Scan(
		&item.ID,
		&item.ClusterUuid,
		&item.AppName,
		&item.AppCode,
		&item.AppType,
		&item.IsDefault,
		&item.AppUrl,
		&item.Port,
		&item.Protocol,
		&item.AuthEnabled,
		&item.AuthType,
		&item.Username,
		&item.Password,
		&token,
		&item.AccessKey,
		&item.AccessSecret,
		&item.TlsEnabled,
		&caFile,
		&caKey,
		&caCert,
		&clientCert,
		&clientKey,
		&item.InsecureSkipVerify,
		&item.Status,
		&item.CreatedBy,
		&item.UpdatedBy,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return ClusterApp{}, err
	}

	item.Token = token.String
	item.CaFile = caFile.String
	item.CaKey = caKey.String
	item.CaCert = caCert.String
	item.ClientCert = clientCert.String
	item.ClientKey = clientKey.String
	item.CreatedAt = createdAt.Unix()
	item.UpdatedAt = updatedAt.Unix()

	return item, nil
}
