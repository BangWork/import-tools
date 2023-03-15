package file

import (
	"encoding/json"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "golang.org/x/image/webp"

	"github.com/bangwork/import-tools/serve/common"
	account2 "github.com/bangwork/import-tools/serve/services/account"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/utils"
)

func UploadFile(account *account2.Account, file *os.File, realFileName string) (resourceUUID string, err error) {
	if cache.SharedDiskPath != "" {
		return uploadToShareDisk(file, account, realFileName)
	}
	return upload(file, account, realFileName)
}

func uploadToShareDisk(file *os.File,
	account *account2.Account, realFileName string) (resourceUUID string, err error) {
	cacheInfo := account.Cache
	fileInfo, err := utils.GetFileInfo(file)
	if err != nil {
		return
	}
	body := &RecordRequest{
		Type:          LabelUploadAttachment,
		ReferenceType: EntityTypeUnrelatedLabel,
		Name:          realFileName,
		Hash:          fileInfo.Hash,
		Mime:          fileInfo.Mime,
		Size:          fileInfo.Size,
		ImageWidth:    fileInfo.ImageWidth,
		ImageHeight:   fileInfo.ImageHeight,
		Exif:          fileInfo.Exif,
	}
	srcPath := file.Name()
	if err = file.Close(); err != nil {
		return "", err
	}
	file, err = os.Open(srcPath)
	if err != nil {
		return "", err
	}
	dstPath := fmt.Sprintf("%s/%s/%s", cache.SharedDiskPath, common.ShareDiskPathPrivate, body.Hash)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("open %s failed, err:%v.\n", dstPath, err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return
	}

	url := common.GenApiUrl(cacheInfo.URL, fmt.Sprintf(fileRecordUri, cacheInfo.ImportTeamUUID))
	resp, err := utils.PostJSONWithHeader(url, body, account.AuthHeader)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("file record failed")
		return
	}
	recordResponse := new(RecordResponse)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, &recordResponse); err != nil {
		return
	}
	resourceUUID = recordResponse.ResourceUUID
	return
}

func upload(file *os.File, account *account2.Account, realFileName string) (resourceUUID string, err error) {
	fileUploadResponse, err := PrepareUploadInfo(realFileName, LabelUploadAttachment, EntityTypeUnrelatedLabel, account)
	if err != nil {
		return "", err
	}
	token := fileUploadResponse.Token
	uploadUrl := fileUploadResponse.UploadURL
	resp, err := utils.PostFileUpload(uploadUrl, token, file, realFileName)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != 579 {
		log.Printf("doUpload file failed")
		return
	}
	resourceUUID = fileUploadResponse.ResourceUUID
	return
}

func PrepareUploadInfo(fileName, label string, refType string, account *account2.Account) (*UploadResponse, error) {
	cacheInfo := account.Cache
	url := common.GenApiUrl(cacheInfo.URL, fmt.Sprintf(fileUploadUri, cacheInfo.ImportTeamUUID))
	body := &UploadRequest{
		Name:          fileName,
		Type:          label,
		ReferenceType: refType,
	}
	resp, err := utils.PostJSONWithHeader(url, body, account.AuthHeader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("upload file failed")
		return nil, err
	}
	fileUploadResponse := new(UploadResponse)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &fileUploadResponse); err != nil {
		return nil, err
	}
	return fileUploadResponse, nil
}
