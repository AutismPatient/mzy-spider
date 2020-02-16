package model

type RunTime struct {
	ID      int32  `json:"id"`
	RunKey  string `json:"run_key"`
	IsRun   int    `json:"is_run"`
	RunTime int64  `json:"run_time"`
}
