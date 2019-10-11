package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type project struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ProjectLocation	string	`json: "plocation"`
        ProjectName       string `json:"pname"`    //the fieldtags are needed to keep case from bouncing around
	ProjectId       int `json:"pid"`
	ProjectSurveyNumber      int `json:"psnum"`
	ProjectNocResponse string `json:pnres"`
	ProjectLakeAuthResponse string `json:plares"`
	ProjectForestAuthResponse string `json:pfares"`
	ProjectCityDevelopmentStatus string `"json:pcdstat"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "NocRequest" {
		return t.NocRequest(stub, args)
	} else if function == "BDAProjectStatus" { 
		return t.BDAProjectStatus(stub, args)
	} else if function == "delete" { 
		return t.delete(stub, args)
	} else if function == "readProject" { 
		return t.readProject(stub, args)
	} else if function == "LakeAuthRequest" { 
		return t.LakeAuthRequest(stub, args)
	} else if function == "ForestAuthRequest" { 
		return t.ForestAuthRequest(stub, args)
	} else if function == "getHistoryForProject" { 
		return t.getHistoryForProject(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initBDA - create a new project, store into chaincode state
// ============================================================
func (t *SimpleChaincode) NocRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error


	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init project")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return shim.Error("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return shim.Error("8th argument must be a non-empty string")
	}
	ProjectLocation := args[0]
	ProjectName := args[1]
	ProjectId, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	ProjectSurveyNumber, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	ProjectNocResponse := args[4]
	ProjectLakeAuthResponse := args[5]
	ProjectForestAuthResponse := args[6]
	ProjectCityDevelopmentStatus := args[7]

	
	ProjectAsBytes, err := stub.GetState(ProjectName)
	if err != nil {
		return shim.Error("Failed to get Project: " + err.Error())
	} else if ProjectAsBytes != nil {
		fmt.Println("This project already exists: " + ProjectName)
		return shim.Error("This project already exists: " + ProjectName)
	}


	objectType := "project"
	project := &project{objectType, ProjectLocation , ProjectName , ProjectId, ProjectSurveyNumber , ProjectNocResponse , ProjectLakeAuthResponse , ProjectForestAuthResponse , ProjectCityDevelopmentStatus}
	projectJSONasBytes, err := json.Marshal(project)
	if err != nil {
		return shim.Error(err.Error())
	}


	err = stub.PutState(ProjectName, projectJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "project~name"
	ProjectNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{project.ProjectLocation, project.ProjectName})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutState(ProjectNameIndexKey, value)

	fmt.Println("- end init project")
	return shim.Success(nil)
}

// ===============================================
// readProject - read a project from chaincode state
// ===============================================
func (t *SimpleChaincode) readProject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the BDA to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the project from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Project does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a project key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var projectJSON project
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	ProjectName := args[0]

	valAsbytes, err := stub.GetState(ProjectName) //get the project from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ProjectName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"BDA does not exist: " + ProjectName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &projectJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + ProjectName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(ProjectName) //remove the project from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	indexName := "project~name"
	ProjectNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{projectJSON.ProjectLocation, projectJSON.ProjectName})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(ProjectNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// ===========================================================
// Get project status by setting a getting approvals
// ===========================================================
func (t *SimpleChaincode) BDAProjectStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProjectName := args[0]
	fmt.Println("- start NocRequest ", ProjectName)

	projectAsBytes, err := stub.GetState(ProjectName)
	if err != nil {
		return shim.Error("Failed to get project:" + err.Error())
	} else if projectAsBytes == nil {
		return shim.Error("Project does not exist")
	}
	projecttoapprove := project{}
	err = json.Unmarshal(projectAsBytes, &projecttoapprove) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	lakeresponse := t.LakeAuthRequest(stub, []string{ProjectName})
	Forestresponse := t.ForestAuthRequest(stub, []string{ProjectName})
	if lakeresponse.Status != shim.OK {
		return shim.Error("project failed to acquire approval: " + lakeresponse.Message)
	}
	if Forestresponse.Status != shim.OK {
		return shim.Error("project failed to acquire approval: " + Forestresponse.Message)
	}
	if lakeresponse.Status == shim.OK && Forestresponse.Status == shim.OK {
		projecttoapprove.ProjectNocResponse = "Yes"
		projecttoapprove.ProjectCityDevelopmentStatus = "Yes"
	}
	

	projectJSONasBytes, _ := json.Marshal(projecttoapprove)
	err = stub.PutState(ProjectName, projectJSONasBytes) //rewrite the BDA
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end NocRequest (success)")
	return shim.Success(nil)
}


func (t *SimpleChaincode) LakeAuthRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProjectName := args[0]
	fmt.Println("- start LakeAuthRequest ", ProjectName)

	projectAsBytes, err := stub.GetState(ProjectName)
	if err != nil {
		return shim.Error("Failed to get project:" + err.Error())
	} else if projectAsBytes == nil {
		return shim.Error("Project does not exist")
	}
	
	projecttoapprove := project{}
	err = json.Unmarshal(projectAsBytes, &projecttoapprove) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	projecttoapprove.ProjectLakeAuthResponse = "Yes" //change the owner

	projectJSONasBytes, _ := json.Marshal(projecttoapprove)
	err = stub.PutState(ProjectName, projectJSONasBytes) //rewrite the BDA
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end LakeAuthRequest (success)")
	return shim.Success(nil)
}



func (t *SimpleChaincode) ForestAuthRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProjectName := args[0]
	fmt.Println("- start ForestAuthRequest ", ProjectName)

	projectAsBytes, err := stub.GetState(ProjectName)
	if err != nil {
		return shim.Error("Failed to get project:" + err.Error())
	} else if projectAsBytes == nil {
		return shim.Error("Project does not exist")
	}
	
	projecttoapprove := project{}
	err = json.Unmarshal(projectAsBytes, &projecttoapprove) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	projecttoapprove.ProjectForestAuthResponse = "Yes" //change the owner

	projectJSONasBytes, _ := json.Marshal(projecttoapprove)
	err = stub.PutState(ProjectName, projectJSONasBytes) //rewrite the BDA
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end ForestAuthRequest (success)")
	return shim.Success(nil)
}


func (t *SimpleChaincode) getHistoryForProject(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ProjectName := args[0]

	fmt.Printf("- start getHistoryForBDA: %s\n", ProjectName)

	resultsIterator, err := stub.GetHistoryForKey(ProjectName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the BDA
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON BDA)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForProject returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
