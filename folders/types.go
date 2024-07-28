package folders

import "github.com/gofrs/uuid"

const FIXED_PAGE_SIZE = 2

type FetchFolderRequest struct {
	OrgId uuid.UUID
	PageToken string
	Paginate bool
}

type FetchFolderResponse struct {
	Folders []*Folder
	NextPageToken string
}

/* 
		Using DataProvider interface abstraction to decouple data retriaval
		from folder filtering logic. 
*/
type DataProvider interface { 
	GetSampleData() []*Folder
}

/*
		Added OriginalDataProvider to allow for different datasets to 
		be filtered. Mainly so that I can use mock data to test and original
		data for actual execution.
*/
type OriginalDataProvider struct {}

