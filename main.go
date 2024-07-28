package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
)

func main() {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Select an option:")
		fmt.Println("1. DUMP (receive ALL data in a single output)")
		fmt.Println("2. PAGINATE (receive paginated data)")
		fmt.Print("Enter your choice (1 or 2): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		originalProvider := &folders.OriginalDataProvider{}
		req := &folders.FetchFolderRequest{
				OrgId: uuid.FromStringOrNil(folders.DefaultOrgId),
		}

		switch choice {
		case "1":
				req.Paginate = false
				fmt.Println("You selected DUMP.")
				res, err := folders.FetchAllFoldersByOrgId(req, originalProvider)
				PrintData(res, err)
		case "2":
				req.Paginate = true
				fmt.Println("You selected PAGINATE.")
				// fmt.Println("Enter page token or leave empty to get the first page:")
				for {
						if req.PageToken != "" {
								fmt.Printf("Using page token: %s\n", req.PageToken)
						} else {
								fmt.Println("Fetching the first page.")
						}
			
						res, err := folders.FetchAllFoldersByOrgId(req, originalProvider)
						PrintData(res, err)
			
						if res.NextPageToken == "" {
								fmt.Println("No more pages available.")
								break
						}
						
						fmt.Println("\nOptions:")
						fmt.Println("1. Fetch next page")
						fmt.Println("2. Exit")
						fmt.Print("Enter your choice (1 or 2): ")
						
						nextChoice, _ := reader.ReadString('\n')
						nextChoice = strings.TrimSpace(nextChoice)
			
						if nextChoice != "1" {
								fmt.Println("Exiting pagination.")
								break
						}
						req.PageToken = res.NextPageToken
				}
		default:
				fmt.Println("Invalid choice. Defaulting to DUMP.")
				req.Paginate = false
				res, err := folders.FetchAllFoldersByOrgId(req, originalProvider)
				PrintData(res, err)
		}

}

func PrintData (res *folders.FetchFolderResponse, err error) { 
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		folders.PrettyPrint(res)
}
