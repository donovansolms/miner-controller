package mhq

import "github.com/donovansolms/mininghq-miner-controller/src/caps"

// Progress holds information about the current download progress
type Progress struct {
	BytesCompleted int64
	BytesTotal     int64
}

// RecommendedMinerResponse contains the recommended miners (if any)
// from the MiningHQ API
type RecommendedMinerResponse struct {
	Status  string             `json:"Status"`
	Message string             `json:"Message"`
	Miners  []RecommendedMiner `json:"Miners"`
}

// RecommendedMiner contains the information to download a recommended miner
type RecommendedMiner struct {
	Name           string `json:"Name"`
	Version        string `json:"Version"`
	Type           string `json:"Type"`
	DownloadLink   string `json:"DownloadLink"`
	DownloadSHA512 string `json:"DownloadSHA512"`
	SizeBytes      int64  `json:"SizeBytes"`
}

// RegisterRigRequest is the request sent to MiningHQ to register a new rig
type RegisterRigRequest struct {
	// Name is a custom name for this rig,
	// if blank, it will be set to the hostname
	Name string
	// Caps is the capabilities of this rig
	Caps caps.SystemInfo
}

// RegisterRigResponse is returned after a RegisterRigRequest
type RegisterRigResponse struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	RigID   string `json:"RigID"`
}