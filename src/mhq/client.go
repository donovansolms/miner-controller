/*
  MiningHQ Miner Controller - manages cryptocurrency miners on a user's machine.
  https://mininghq.io

	Copyright (C) 2018  Donovan Solms     <https://github.com/donovansolms>
                                        <https://github.com/mininghq>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package mhq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sethgrid/pester"
)

// Client is the API client for MiningHQ APIs
//
// It provides the basic functionality in plain JSON calls. This will be
// moved to a gRPC API in the future
//
// TODO: Implement as gRPC API client before moving out of beta
type Client struct {
	// miningKey is the user's mining key. It is used as an identification
	// token for API calls and websocket connections
	miningKey string
	// endpoint of the MiningHQ API
	endpoint string
}

// NewClient creates and returns a new MiningHQ API client
func NewClient(miningKey string, endpoint string) (*Client, error) {
	if miningKey == "" {
		return nil, errors.New(
			"You must provide a valid mining key to the MiningHQ API client, not blank")
	}
	if endpoint == "" {
		return nil, errors.New(
			"You must provide a valid endpoint to the MiningHQ API client, not blank")
	}

	client := Client{
		miningKey: miningKey,
		endpoint:  endpoint,
	}

	return &client, nil
}

// RegisterRig registers this system/rig for the user with MiningHQ
// It returns the MiningHQ RigID
func (client *Client) RegisterRig(
	registerRequest RegisterRigRequest) (string, error) {

	jsonBytes, err := json.Marshal(registerRequest)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/register-rig", client.endpoint),
		bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.miningKey))

	apiClient := pester.New()
	apiClient.MaxRetries = 5
	response, err := apiClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Unable to register rig: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to register rig: Status %s", response.Status)
	}

	var registerResponse RegisterRigResponse
	err = json.NewDecoder(response.Body).Decode(&registerResponse)
	if err != nil {
		return "", fmt.Errorf("Unable to read register response: %s", err)
	}

	if registerResponse.Status == "err" {
		return "", errors.New(registerResponse.Message)
	}

	return registerResponse.RigID, nil
}

// DeregisterRig removes this system/rig from the user with MiningHQ
func (client *Client) DeregisterRig(
	deregisterRequest DeregisterRigRequest) error {

	jsonBytes, err := json.Marshal(deregisterRequest)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/deregister-rig", client.endpoint),
		bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.miningKey))

	apiClient := pester.New()
	apiClient.MaxRetries = 5
	response, err := apiClient.Do(request)
	if err != nil {
		return fmt.Errorf("Unable to deregister rig: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to deregister rig: Status %s", response.Status)
	}

	var deregisterResponse DeregisterRigResponse
	err = json.NewDecoder(response.Body).Decode(&deregisterResponse)
	if err != nil {
		return fmt.Errorf("Unable to read deregister response: %s", err)
	}

	if deregisterResponse.Status == "err" {
		return errors.New(deregisterResponse.Message)
	}

	return nil
}

// // GetRecommendedMiners submits the rig capabilities to MiningHQ and gets
// // back a list of compatible miners, if any
// func (client *Client) GetRecommendedMiners(
// 	systemInfo caps.SystemInfo) ([]RecommendedMiner, error) {
//
// 	jsonBytes, err := json.Marshal(systemInfo)
// 	if err != nil {
// 		return []RecommendedMiner{}, err
// 	}
//
// 	apiClient := pester.New()
// 	apiClient.MaxRetries = 5
// 	response, err := apiClient.Post(
// 		fmt.Sprintf("%s/recommend-miners", client.endpoint),
// 		"application/json",
// 		bytes.NewReader(jsonBytes))
// 	if err != nil {
// 		return []RecommendedMiner{},
// 			fmt.Errorf("Unable to get recommended miners: %s", err)
// 	}
// 	defer response.Body.Close()
//
// 	if response.StatusCode != http.StatusOK {
// 		return []RecommendedMiner{},
// 			fmt.Errorf("Unable to get recommended miners: %s", err)
// 	}
//
// 	var recommendedMinerResponse RecommendedMinerResponse
// 	err = json.NewDecoder(response.Body).Decode(&recommendedMinerResponse)
// 	if err != nil {
// 		return []RecommendedMiner{},
// 			fmt.Errorf("Error reading recommended miners response: %s", err)
// 	}
// 	if recommendedMinerResponse.Status == "err" {
// 		return []RecommendedMiner{}, errors.New(recommendedMinerResponse.Message)
// 	}
//
// 	return recommendedMinerResponse.Miners, nil
// }

// DownloadMiner downloads and verifies the given miner
// func (client *Client) DownloadMiner(
// 	destination string,
// 	recommendedMiner RecommendedMiner,
// 	progressChan chan Progress) error {
//
// 	downloadClient := grab.NewClient()
// 	req, err := grab.NewRequest(destination, recommendedMiner.DownloadLink)
// 	if err != nil {
// 		return err
// 	}
//
// 	// start download
// 	resp := downloadClient.Do(req)
//
// 	// start progress loop
// 	t := time.NewTicker(500 * time.Millisecond)
// 	defer t.Stop()
//
// Loop:
// 	for {
// 		select {
// 		case <-t.C:
// 			progressChan <- Progress{
// 				BytesCompleted: resp.BytesComplete(),
// 				BytesTotal:     resp.Size,
// 			}
// 		case <-resp.Done:
// 			// Downoad completed, close progress channel
// 			close(progressChan)
// 			break Loop
// 		}
// 	}
//
// 	// check for errors
// 	err = resp.Err()
// 	if err != nil {
// 		return err
// 	}
//
// 	// Download complete, verify the download
// 	file, err := os.Open(destination)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()
//
// 	hasher := sha512.New()
// 	if _, err := io.Copy(hasher, file); err != nil {
// 		return err
// 	}
//
// 	if hex.EncodeToString(hasher.Sum(nil)) == recommendedMiner.DownloadSHA512 {
// 		return nil
// 	}
// 	return errors.New(
// 		"The downloaded miner could not be verified, verification hashes differ. Please try again later")
// }
