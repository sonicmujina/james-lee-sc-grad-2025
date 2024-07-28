package folders_test

import (
	"math/rand"
	"testing"
	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

type MockDataProvider struct {
	Folders []*folders.Folder
	PageToken string
}

func (m *MockDataProvider) GetSampleData() []*folders.Folder { 
	return m.Folders
}

func Test_FetchAllFoldersByOrgId(t *testing.T) {
	// Create sample OrgIDs
	orgId1 := uuid.Must(uuid.NewV4())
	orgId2 := uuid.Must(uuid.NewV4())
	orgId3 := uuid.Must(uuid.NewV4())
	unknownOrgId := uuid.Must(uuid.NewV4())
	
	deleted := rand.Int() % 2
	
	// Create sample Folder objects
	folder1 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 1", 
		OrgId: orgId1, 
		Deleted: deleted != 0,
	}

	folder2 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 2", 
		OrgId: orgId1, 
		Deleted: deleted != 0,
	}

	folder3 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 3", 
		OrgId: orgId2, 
		Deleted: deleted != 0,
	}

	
	mockProvider := &MockDataProvider{Folders: []*folders.Folder{
		folder1, 
		folder2, 
		folder3,
		},
	}

	
	var tests = []struct { 
		name string
		request *folders.FetchFolderRequest
		provider folders.DataProvider
		want *folders.FetchFolderResponse
		wantErr bool
	}{
		{
			name: "Fetch w orgId1",
			request: &folders.FetchFolderRequest{OrgId: orgId1},
			provider: folders.DataProvider(mockProvider),
			want: &folders.FetchFolderResponse{Folders: []*folders.Folder{folder1, folder2}},
		},
		{
			name: "Fetch w orgId2",
			request: &folders.FetchFolderRequest{OrgId: orgId2},
			provider: folders.DataProvider(mockProvider),
			want: &folders.FetchFolderResponse{Folders: []*folders.Folder{folder3}},
		},
		{
			name: "Fetch w orgId3 that has no folders",
			request: &folders.FetchFolderRequest{OrgId: orgId3},
			provider: folders.DataProvider(mockProvider),
			want: &folders.FetchFolderResponse{Folders: []*folders.Folder{}},
		},
		{
			name: "Fetch w unknown orgId and return err", 
			request: &folders.FetchFolderRequest{OrgId: unknownOrgId},
			provider: folders.DataProvider(mockProvider),
			want: &folders.FetchFolderResponse{Folders: []*folders.Folder{}},
			wantErr: true,
		},
		{
			name: "Fetch w orgId1 and nil provider",
			request: &folders.FetchFolderRequest{OrgId: orgId1},
			provider: nil,
			want: nil,
			wantErr: true,
		},
		{
			name: "Fetch w nil orgId",
			request: &folders.FetchFolderRequest{OrgId: uuid.Nil},
			provider: nil,
			want: nil,
			wantErr: true,
		},
		{
			name: "Fetch w nil request",
			request: nil,
			provider: nil,
			want: nil,
			wantErr: true,
		},
		
	}


	for _, tt := range tests { 

		t.Run(tt.name, func(t *testing.T) {
			result, err := folders.FetchAllFoldersByOrgId(tt.request, tt.provider)
			assert.Equal(t, result, tt.want, "They should be equal" )
			
			if tt.wantErr { // Check for error messages
				assert.Error(t, err)
			} 

		})
	}
}

func Test_Pagination(t *testing.T) { 
	// Create sample OrgIDs
	orgId1 := uuid.Must(uuid.NewV4())

	deleted := rand.Int() % 2

	// Create sample Folder objects
	folder1 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 1", 
		OrgId: orgId1, 
		Deleted: deleted != 0,
	}

	folder2 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 2", 
		OrgId: orgId1, 
		Deleted: deleted != 0,
	}

	folder3 := &folders.Folder{
		Id: uuid.Must(uuid.NewV4()), 
		Name: "Test-Folder 3", 
		OrgId: orgId1, 
		Deleted: deleted != 0,
	}

	t.Run("Correct return values for first and last page", func(t *testing.T) { 
	
			mockFolders := []*folders.Folder{folder1, folder2}
			pageToken := ""
		
			resp, _ := folders.Pagination(mockFolders, pageToken)
			
			// Check number of folders
			assert.Equal(t, 2, len(resp.Folders), "Number of folders returned should be 2")
			
			// Check folder content
			assert.Equal(t, folder1, resp.Folders[0], "First folder should match")
			assert.Equal(t, folder2, resp.Folders[1], "Second folder should match")

			// Check that NextPageToken is empty
			assert.Empty(t, resp.NextPageToken, "NextPageToken should be empty for the last page")
	})

	t.Run("Generate next page token", func(t *testing.T) { 
			mockFolders := []*folders.Folder{folder1, folder2, folder3}
			pageToken := ""
			
			resp, _ := folders.Pagination(mockFolders, pageToken)

			// Check that NextPageToken is not empty
			assert.NotEmpty(t, resp.NextPageToken, "NextPageToken should not be empty")
	})
	
	t.Run("Return correct page with valid page token", func(t *testing.T) { 
			mockFolders := []*folders.Folder{folder1, folder2, folder3}
			pageToken := ""
			
			resp, _ := folders.Pagination(mockFolders, pageToken)

			finalResp, _ := folders.Pagination(mockFolders, resp.NextPageToken)
			
			// Check that pagination decodes the pageToken to correct page
			assert.Equal(t, folder3, finalResp.Folders[0], "First folder of 2nd page should match")
	})

	t.Run("Detect invalid page token", func(t *testing.T) { 
			mockFolders := []*folders.Folder{folder1}
			invalidPageToken := "ab271ac"

			_, err := folders.Pagination(mockFolders, invalidPageToken)

			// Check that err is thrown for invalid page token
			assert.Error(t, err, "Pagination should return an error on invalid page token")
	})
}
