package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const AwesomeChatGETPromptsApiUrl = "https://datasets-server.huggingface.co/rows?dataset=fka%2Fawesome-chatgpt-prompts&config=default&split=train&offset=0&limit=100"

type huggingFaceResponse struct {
	Dataset  string `json:"dataset"`
	Config   string `json:"config"`
	Split    string `json:"split"`
	Features []struct {
		FeatureIdx int    `json:"feature_idx"`
		Name       string `json:"name"`
		Type       struct {
			Dtype string `json:"dtype"`
			Type  string `json:"_type"`
		} `json:"type"`
	} `json:"features"`
	Rows []struct {
		RowIdx int `json:"row_idx"`
		Row    struct {
			Act    string `json:"act"`
			Prompt string `json:"prompt"`
		} `json:"row"`
		TruncatedCells []interface{} `json:"truncated_cells"`
	} `json:"rows"`
}

type Prompts struct {
	Act    string `json:"act"`
	Prompt string `json:"prompt"`
}

func GetPrompts() ([]Prompts, error) {
	if err := updateCheck(); err != nil {
		fmt.Println("Error updating check:", err)
		return nil, err
	}

	jsonBytes, err := ioutil.ReadFile("awesome-chatgpt-prompts.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return nil, err
	}

	var hfResp huggingFaceResponse
	err = json.Unmarshal(jsonBytes, &hfResp)

	var prompts []Prompts
	for _, row := range hfResp.Rows {
		prompts = append(prompts, Prompts{
			Act:    row.Row.Act,
			Prompt: row.Row.Prompt,
		})
	}

	return prompts, nil
}

func updateCheck() error {
	if !fileExists("awesome-chatgpt-prompts.json") {
		err := downloadJSONFile()
		if err != nil {
			fmt.Println("Error downloading JSON file:", err)
			return err
		}
	}
	jsonBytes, err := ioutil.ReadFile("awesome-chatgpt-prompts.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	var hfResp huggingFaceResponse
	err = json.Unmarshal(jsonBytes, &hfResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	fileModTime, err := getFileModTime("awesome-chatgpt-prompts.json")
	if err != nil {
		fmt.Println("Error getting file mod time:", err)
		return err
	}

	if time.Since(fileModTime) > 24*time.Hour {
		err := downloadJSONFile()
		if err != nil {
			fmt.Println("Error downloading JSON file:", err)
			return err
		}
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getFileModTime(filename string) (time.Time, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

func downloadJSONFile() error {
	resp, err := http.Get(AwesomeChatGETPromptsApiUrl)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
			os.Exit(1)
		}
	}(resp.Body)
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("awesome-chatgpt-prompts-%d.json", time.Now().Unix())
	err = ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return os.Rename(filename, "awesome-chatgpt-prompts.json")

}
