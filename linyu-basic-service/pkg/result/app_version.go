package result

type AppVersionCheckResult struct {
	Platform              string `json:"platform"`
	LatestVersion         string `json:"latestVersion"`
	LatestVersionCode     int    `json:"latestVersionCode"`
	MinSupportVersion     string `json:"minSupportVersion"`
	MinSupportVersionCode int    `json:"minSupportVersionCode"`
	DownloadUrl           string `json:"downloadUrl"`
	UpdateDesc            string `json:"updateDesc"`
}
