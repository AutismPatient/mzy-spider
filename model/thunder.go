package model

type Thunder struct {
	MinVersion         string `json:"min_version"`
	DownloadDir        string `json:"downloadDir"`
	RunParams          string `json:"runParams"`
	TaskGroupName      string `json:"taskGroupName"`
	ThreadCount        int    `json:"threadCount"`
	ThunderInstallPack string `json:"thunderInstallPack"`
	Tasks              []Task `json:"tasks"`
}

type Task struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Dir  string `json:"dir"`
}
