package folders

import (
	"fmt"
	"github.com/gofrs/uuid"
)

/*
		Component 1:
		Improvements that can be made:
		1. Error checking
		2. Removing a few unnecessary steps and repeated code
		3. Could remove or use any unused variables

		Retrospective 
		After making the changes and trying to come up with unit tests
		I have realised that the current code I have is tightly coupled 
		and will refactor to separate concerns. 
		
		FetchAllFoldersByOrgID should only filter folders by orgId 
		for a given set of data instead of internally using GetSampleData 
		to retrieve the data.
		
		Will do a dependancy injection by creating an interface DataProvider
		in types.go and having struct types that implement it to achieve 
		loosely coupled functions and abide to the Dependancy Inversion Principle.
*/


/*  
		This function retrieves all folders from the sample data that have 
		matching OrgID, it then returns a slice of those folders.
*/
func FetchAllFoldersByOrgId(req *FetchFolderRequest, provider DataProvider) (*FetchFolderResponse, error) {
	var (
		ffr *FetchFolderResponse
		err error
	)
	if provider == nil || req == nil || req.OrgId == uuid.Nil { 
		return nil, fmt.Errorf("one or more parameters contains nil input")
	}	
	
	folders := provider.GetSampleData()

	filteredFolder := []*Folder{}	
	for _, folder := range folders {
		if folder.OrgId == req.OrgId {
			filteredFolder = append(filteredFolder, folder)  
		}
	}

	if req.Paginate {	// Check if pagination is true 
		ffr, err = Pagination(filteredFolder, req.PageToken)
		if err != nil {
				return nil, fmt.Errorf("pagination error: %v", err)
		}
	} else { 
		ffr = &FetchFolderResponse{Folders: filteredFolder}
	}

	if len(ffr.Folders) == 0 { 
		return ffr, fmt.Errorf("no folders found for organization ID %s", req.OrgId)
	}
	return ffr, nil
}

// Seems like this function returns a slice of pointers to folders that match the req.OrgID
// func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
// 	var (
// 		ffr *FetchFolderResponse
// 		// f1  Folder
// 		// fs  []*Folder
// 	)

// 	// f := []Folder{}    // Slice of folder objects
// 	// r, err := FetchAllFoldersByOrgID(req.OrgID)		// returns all folders with associated OrgID
// 	if r, err := FetchAllFoldersByOrgID(req.OrgId); err != nil {
// 		fmt.Println("FetchAllFoldersByOrgID failed:", err)
// 		return nil, err
// 	} else {
// 		ffr = &FetchFolderResponse{Folders: r}
// 		return ffr, nil
// 	}

// 	// for k, v := range r {   // k(index)-v(pointer to folder object), this goes over the folders and creates a new copy into f
// 	// 	f = append(f, *v)
// 	// }

// 	// var fp []*Folder
// 	// for _, v1 := range r {		// adds pointers to folder objects into fp
// 	// 	fp = append(fp, v1)
// 	// }
// 	// ffr = &FetchFolderResponse{Folders: r}	// ffr is of type FetchFolderResponse containing a slice of folder pointers
// 	// return ffr, nil
// }


