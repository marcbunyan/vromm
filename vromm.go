package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type LoginData struct {
	Username   string `json:"username"`
	AuthSource string `json:"authSource"`
	Password   string `json:"password"`
}

type ObjectResponse struct {
	ResourceList []struct {
		Identifier string `json:"identifier"`
	} `json:"resourceList"`
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run yourname.go <vROPS_Hostname> <VM_Name> <start/end>")
		return
	}

	vropsHostname := os.Args[1]
	vmName := os.Args[2]
	action := os.Args[3]

	// Login to vROPS API and get auth token
	loginData := LoginData{
		Username:   "james.bond",
		AuthSource: "vIDMAuthSource",
		Password:   "SuperSecretPassword,
	}

	loginDataBytes, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println("Error marshaling login data:", err)
		return
	}

	loginUri := fmt.Sprintf("https://%s/suite-api/api/auth/token/acquire?_no_links=true", vropsHostname)
	req, err := http.NewRequest("POST", loginUri, bytes.NewBuffer(loginDataBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json") // Add Accept header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}

	token := tokenResponse["token"].(string)

	// Show the auth token on the console
	fmt.Println("Auth Token:", token)

	// Get vROPS object ID by VM name
	objectUri := fmt.Sprintf("https://%s/suite-api/api/resources?resourceKind=VirtualMachine&name=%s", vropsHostname, vmName)
	req, err = http.NewRequest("GET", objectUri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "application/json") // Add Accept header
	req.Header.Set("Authorization", "vRealizeOpsToken "+token)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	var objectResponse ObjectResponse
	err = json.Unmarshal(body, &objectResponse)
	if err != nil {
		fmt.Println("Error unmarshaling response:", err)
		return
	}

	// Check if the resourceList array has at least one item
	if len(objectResponse.ResourceList) == 0 {
		fmt.Println("Error: resourceList is empty")
		return
	}

	objectID := objectResponse.ResourceList[0].Identifier

	// Show object ID on the console
	fmt.Println("Object ID:", objectID)

	// Define the maintenance mode endpoint
	maintenanceUri := fmt.Sprintf("https://%s/suite-api/api/resources/%s/maintained?_no_links=true", vropsHostname, objectID)

	// Check the value of action and perform the corresponding operation
	if action == "start" {
		// Send a PUT request to the maintenance mode endpoint to start maintenance
		req, err := http.NewRequest("PUT", maintenanceUri, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Authorization", "vRealizeOpsToken "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer resp.Body.Close()

		fmt.Println("The maintenance mode for object ID", objectID, "has been started.")
	} else if action == "end" {
		// Send a DELETE request to the maintenance mode endpoint to end maintenance
		req, err := http.NewRequest("DELETE", maintenanceUri, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Authorization", "vRealizeOpsToken "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		defer resp.Body.Close()

		fmt.Println("The maintenance mode for object ID", objectID, "has been ended.")
	} else {
		fmt.Println("Invalid action:", action, ". Expected 'start' or 'end'.")
	}
}
