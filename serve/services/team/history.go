package team

import (
	"github.com/bangwork/import-tools/serve/models/ones"
)

func GetImportHistory(orgUUID, url string, header map[string]string) ([]*ones.ImportHistory, error) {
	history, err := ones.GetImportHistory(orgUUID, url, header)
	if err != nil {
		return nil, err
	}
	return history, nil
}
