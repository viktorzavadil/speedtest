package service

func IsLoadedStatus(status TestStatus) bool {
	return status == DownloadPhase || status == UploadPhase
}
