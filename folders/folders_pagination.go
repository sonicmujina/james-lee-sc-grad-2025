package folders

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

/*
	Component 2: 
	My solution for pagination mainly addresses the need 
	to break down large datasets into smaller chunks. 
	
	I have set the fixed paged size to 2 entries for the sake 
	of testing with smaller datasets since writing big mock test 
	data is outside of the scope. As well as just achieving a 
	more digestible amount of output data. 

	For generating and reading page tokens I used base64 encoding 
	and decoding for token creation and use. 

	I have also changed main.go to provide a more user friendly CLI 
	which provides the option to DUMP all data or paginate data.
*/


/* 
	Breaks down large dataset into smaller chunks and returns chunks
 	associated to a page token
*/
func Pagination(folders []*Folder, pageToken string) (*FetchFolderResponse, error) { 
	var (
		pffr *FetchFolderResponse
		paginatedFolders []*Folder
		nextPageToken string
		endIndex int
		err error
	) 

	if len(pageToken) > 0 { 
		startingIndex, err := DecodePageToken(pageToken)
		if err != nil { 
			return pffr, err
		}

		endIndex = startingIndex + FIXED_PAGE_SIZE
    
		// Ensure we don't go out of bounds
		if endIndex > len(folders) {
			endIndex = len(folders)
		}
		paginatedFolders = folders[startingIndex:endIndex]
	} else { // no page token given, return first page 
		endIndex = FIXED_PAGE_SIZE
		if endIndex > len(folders) { 
			endIndex = len(folders)
		}
		paginatedFolders = folders[0:endIndex]
		nextPageToken = ""
	}
    
	// Generate next page token if more folders exist
	if endIndex < len(folders) {
		nextPageToken, err = EncodeNewPageToken(endIndex)
		if err != nil { 
			return pffr, err
		}
	} else { // No more pages
		nextPageToken = "" 
	}

	pffr = 	&FetchFolderResponse{Folders: paginatedFolders, NextPageToken: nextPageToken}
	return pffr, nil 
}

// Encodes index into base64 encoding
func EncodeNewPageToken (index int) (string, error) {
	if index < 0 {
        return "", fmt.Errorf("invalid index: %d (must be non-negative)", index)
    }

	indexStr := strconv.Itoa(index)
	encodedToken := base64.StdEncoding.EncodeToString([]byte(indexStr))
	return encodedToken, nil
}

// Decodes pageToken to index
func DecodePageToken (pageToken string) (int, error) { 
	decodedBytes, err := base64.StdEncoding.DecodeString(pageToken)
	
	if err != nil {
        return 0, fmt.Errorf("invalid page token: %v", err)
    }

    index, err := strconv.Atoi(string(decodedBytes))
    
	if err != nil {
        return 0, fmt.Errorf("invalid page token: %v", err)
    }
    return index, nil
}