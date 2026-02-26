package result

type CheckUploadStatusResult struct {
	Uploaded       bool     `json:"uploaded"`
	UploadedChunks []string `json:"uploadedChunks"`
}
