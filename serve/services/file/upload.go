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

	"github.com/bangwork/import-tools/serve/models/ones"

	_ "golang.org/x/image/webp"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

func UploadFile(cookie, teamUUID string, file *os.File, realFileName string) (resourceUUID string, err error) {
	cookieValue, err := ones.DecryptCookieValueByCookie(cookie)
	if err != nil {
		return "", err
	}

	if common.SharedDiskPath != "" {
		return uploadToShareDisk(file, cookieValue.URL, teamUUID, realFileName, cookieValue.GenAuthHeader())
	}
	return upload(file, cookieValue.URL, teamUUID, realFileName, cookieValue.GenAuthHeader())
}

func uploadToShareDisk(file *os.File,
	url, teamUUID, realFileName string, header map[string]string) (resourceUUID string, err error) {
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
	dstPath := fmt.Sprintf("%s/%s/%s", common.SharedDiskPath, common.ShareDiskPathPrivate, body.Hash)
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

	url = common.GenApiUrl(url, fmt.Sprintf(fileRecordUri, teamUUID))
	resp, err := utils.PostJSONWithHeader(url, body, header)
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

func upload(file *os.File, url, teamUUID, realFileName string, header map[string]string) (resourceUUID string, err error) {
	fileUploadResponse, err := PrepareUploadInfo(realFileName, LabelUploadAttachment, EntityTypeUnrelatedLabel, url, teamUUID, header)
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

func PrepareUploadInfo(fileName, label, refType, url, teamUUID string, header map[string]string) (*UploadResponse, error) {
	url = common.GenApiUrl(url, fmt.Sprintf(fileUploadUri, teamUUID))
	body := &UploadRequest{
		Name:          fileName,
		Type:          label,
		ReferenceType: refType,
	}
	resp, err := utils.PostJSONWithHeader(url, body, header)
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
