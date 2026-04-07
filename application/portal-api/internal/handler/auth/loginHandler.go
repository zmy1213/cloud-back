package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	authservice "github.com/yanshicheng/cloud-back/application/portal-api/internal/auth"
	authlogic "github.com/yanshicheng/cloud-back/application/portal-api/internal/logic/auth"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		l := authlogic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			if errors.Is(err, authservice.ErrInvalidCredentials) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
