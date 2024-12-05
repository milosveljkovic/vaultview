package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vaultview/pkg/config"
	"vaultview/pkg/vault"
)

type InfoListeners interface {
	UpdateInfoTable(infos Info)
}

type Info struct {
	VaultViewRev        string
	VaultRev            string
	VaultAddr           string
	Sealed              string
	TokenPolicies       string
	TokenExpirationTime string
	listeners           []InfoListeners
}

type vaultResponse struct {
	Sealed  *bool  `json:"sealed"`
	Version string `json:"version"`
}

var (
	version = "v0.0.0"
)

func NewInfo(vaultCli vault.VaultSvc, cfg *config.Config) (*Info, error) {
	info := &Info{
		VaultViewRev: version,
		VaultAddr:    cfg.VaultAddr,
		Sealed:       "",
	}
	vs, err := info.getVaultInfo()
	if err != nil {
		return info, err
	}

	info.VaultRev = vs.Version
	if vs.Sealed == nil {
		info.Sealed = ""
	} else if *vs.Sealed {
		info.Sealed = "true"
	} else {
		info.Sealed = "false"
	}
	tokenInfo, err := vaultCli.ReadTokenInfo()
	if err != nil {
		return info, err
	}
	info.TokenPolicies = tokenInfo["policies"]
	info.TokenExpirationTime = tokenInfo["expire_time"]
	return info, err
}

func (i *Info) getVaultInfo() (vaultResponse, error) {
	var vr vaultResponse

	//todo: make sure to remove last '/' from vault addr
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/sys/health", i.VaultAddr), nil)
	if err != nil {
		return vaultResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return vaultResponse{}, err
	}

	response, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(response, &vr); err != nil {
		return vaultResponse{}, err
	}
	return vr, nil
}

func (i *Info) RegisterListener(listener InfoListeners) {
	i.listeners = append(i.listeners, listener)
}

func (i *Info) TriggerInfoChange() {
	for _, l := range i.listeners {
		l.UpdateInfoTable(*i)
	}
}
