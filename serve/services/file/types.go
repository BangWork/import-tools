package file

type UploadRequest struct {
	Type             string `json:"type"`
	Source           string `json:"source"`
	Name             string `json:"name"`
	Hash             string `json:"hash"`
	ReferenceType    string `json:"ref_type"`
	ReferenceID      string `json:"ref_id"`
	IgnoreNotice     bool   `json:"ignore_notice"`
	Description      string `json:"description"`
	ResourceId       string `json:"-"`
	Skip             bool
	CreateIfNotExist bool   `json:"-"`
	Ctype            string `json:"ctype"`
	ImageWidth       int    `json:"image_width"`
	ImageHeight      int    `json:"image_height"`
}

type UploadResponse struct {
	NeedUpload   bool                   `json:"need_upload"`
	BaseURL      string                 `json:"base_url"`
	UploadURL    string                 `json:"upload_url"`
	Token        string                 `json:"token,omitempty"`
	ResourceUUID string                 `json:"resource_uuid,omitempty"`
	File         *UploadCallbackPayload `json:"file,omitempty"`
	SizeLimit    int64                  `json:"size_limit"`
	NeedFormdata bool                   `json:"need_form_data"`
	NeedCallback bool                   `json:"need_callback"`
	CallbackURL  string                 `json:"callback_url"`
	ExcInfo      string                 `json:"exc_info"`
	PostFormdata map[string]string      `json:"post_form_data"`
}

type UploadCallbackPayload struct {
	Hash    string `json:"hash"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	MIME    string `json:"mime"`
	Version int    `json:"version"`
}

type RecordRequest struct {
	Type          string `json:"type"`
	Source        string `json:"source"`
	Name          string `json:"name"`
	ReferenceType string `json:"ref_type"`
	ReferenceID   string `json:"ref_id"`
	ProjectUUID   string `json:"project_uuid"`
	Description   string `json:"description"`

	Hash        string `json:"hash"`
	Mime        string `json:"mime"`
	Size        int64  `json:"size"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
	Exif        string `json:"exif"`
}

type RecordResponse struct {
	ResourceUUID string `json:"resource_uuid"`
}
