package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"vaultview/pkg/constants"

	"github.com/hashicorp/vault-client-go"
)

type VaultSvc interface {
	ReadSecretEngines() ([]string, error)
	ListKvSecrets(mountPath, secretPath string) ([]string, error)
	ReadTokenInfo() (map[string]string, error)
	ReadKvSecret(mountPath, secretPath string) (map[string]string, map[string]string, error)
	IsErrorStatus(err error, status int) bool
}

type Vault struct {
	cli *vault.Client
}

func NewVault(addr, token string) (VaultSvc, error) {
	client, err := vault.New(
		vault.WithAddress(addr),
		vault.WithRequestTimeout(30*time.Second),
	)

	if err != nil {
		return &Vault{
			cli: client,
		}, err
	}

	if err := client.SetToken(token); err != nil {
		return &Vault{
			cli: client,
		}, err
	}

	return &Vault{
		cli: client,
	}, nil
}

func (v Vault) ReadSecretEngines() ([]string, error) {
	secretEnignesNames := []string{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	secretEngines, err := v.cli.System.MountsListSecretsEngines(ctx)

	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("%d", http.StatusForbidden)) {
			return nil, nil
		}
		return nil, err
	}

	if len(secretEngines.Data) == 0 {
		return nil, nil
	}
	for engine, _ := range secretEngines.Data {
		secretEnignesNames = append(secretEnignesNames, engine[:len(engine)-1])
	}
	return secretEnignesNames, nil
}

func (v Vault) ListKvSecrets(mountPath, secretPath string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s, err := v.cli.Secrets.KvV2List(ctx, secretPath, vault.WithMountPath(mountPath))
	if err != nil {
		return nil, err
	}
	if len(s.Data.Keys) == 0 {
		return nil, nil
	}
	return s.Data.Keys, nil
}

func (v Vault) ReadKvSecret(mountPath, secretPath string) (map[string]string, map[string]string, error) {
	sm := make(map[string]string)
	metadata := make(map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s, err := v.cli.Secrets.KvV2Read(ctx, secretPath, vault.WithMountPath(mountPath))
	if err != nil {
		return nil, nil, err
	}
	for i, s := range s.Data.Data {
		json, err := json.Marshal(s)
		if err != nil {
			return nil, nil, fmt.Errorf("Error marshaling Data: %v", err)
		}
		sm[i] = string(json)
	}
	for i, s := range s.Data.Metadata {
		switch v := s.(type) {
		case string:
			metadata[i] = v
		case bool:
			metadata[i] = fmt.Sprintf("%v", v)
		case json.Number:
			if intVal, err := v.Int64(); err == nil {
				metadata[i] = fmt.Sprintf("%d", intVal)
			}
		default:
			metadata[i] = constants.NAValue
		}
	}
	return sm, metadata, nil
}

// todo write "write func"
// func (v Vault) WriteKvSecret(mountPath, secretPath string) (map[string]string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
// 	defer cancel()
// 	// s, err := v.cli.Secrets.KvV2Write(ctx,mountPath,)
// }

func (v Vault) ReadTokenInfo() (map[string]string, error) {
	var tokenInfos = make(map[string]string)
	var plcs []string
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	s, err := v.cli.Auth.TokenLookUpSelf(ctx)
	if err != nil {
		return nil, err
	}
	policies := s.Data["policies"].([]interface{})
	expire_time := s.Data["expire_time"]
	for _, p := range policies {
		plcs = append(plcs, p.(string))
	}
	tokenInfos["policies"] = strings.Join(plcs, " ")
	if expire_time == nil && tokenInfos["policies"] == "root" {
		tokenInfos["expire_time"] = "eternal"
	} else {
		tokenInfos["expire_time"] = expire_time.(string)
	}
	return tokenInfos, nil
}

func (v Vault) IsErrorStatus(err error, status int) bool {
	return vault.IsErrorStatus(err, status)
}
