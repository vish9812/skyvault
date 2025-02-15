package integration

import (
	"fmt"
	"skyvault/internal/domain/media"
	"skyvault/pkg/paging"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMediaManagementFlow tests the complete media management functionality
// through realistic user workflows. It focuses on positive cases and integration
// between different operations. Edge cases and error conditions are covered
// in unit tests (see: file_info_test.go and folder_info_test.go)
func TestMediaManagementFlow(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	// Create initial folder
	folder1 := createFolder(t, env, token, 0, "Documents")
	require.Equal(t, "Documents", folder1.Name, "created folder should have correct name")

	// Upload a file to the folder
	file1 := uploadFile(t, env, token, folder1.ID, "test.txt", media.BytesPerKB)
	require.Equal(t, "test.txt", file1.Name, "uploaded file should have correct name")
	require.Equal(t, int64(media.BytesPerKB), file1.Size, "uploaded file should have correct size")

	// Verify folder contents
	contents1 := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1.FilePage.Items, 1, "folder should contain exactly one file")
	require.Equal(t, file1.Name, contents1.FilePage.Items[0].Name, "file in folder should have correct name")
	require.Len(t, contents1.FolderPage.Items, 0, "folder should not contain any subfolders")

	// Verify root folder contents
	rootContents := getFolderContents(t, env, token, 0)
	require.Len(t, rootContents.FolderPage.Items, 1, "root should contain exactly one folder")
	require.Equal(t, folder1.Name, rootContents.FolderPage.Items[0].Name, "folder in root should have correct name")
	require.Len(t, rootContents.FilePage.Items, 0, "root should not contain any files")

	// Download the file
	buf := make([]byte, media.BytesPerKB)
	downloadFile(t, env, token, file1.ID, buf)
	require.Len(t, buf, media.BytesPerKB, "downloaded file should have correct size")

	// Rename the file
	renameFile(t, env, token, file1.ID, "renamed.txt")

	// Verify file name change
	contents1Renamed := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1Renamed.FilePage.Items, 1, "folder should contain exactly one file after rename")
	require.Equal(t, "renamed.txt", contents1Renamed.FilePage.Items[0].Name, "renamed file should have updated name")

	// Create another folder
	folder2 := createFolder(t, env, token, 0, "Archive")
	require.Equal(t, "Archive", folder2.Name, "second folder should have correct name")

	// Move file to new folder
	moveFile(t, env, token, file1.ID, folder2.ID)

	// Verify contents of both folders
	contents1Moved := getFolderContents(t, env, token, folder1.ID)
	require.Len(t, contents1Moved.FilePage.Items, 0, "source folder should be empty after move")

	contents2 := getFolderContents(t, env, token, folder2.ID)
	require.Len(t, contents2.FilePage.Items, 1, "destination folder should contain exactly one file")
	require.Equal(t, file1.ID, contents2.FilePage.Items[0].ID, "file in destination folder should be the moved file")

	// Trash both folders
	trashFolders(t, env, token, []int64{folder1.ID, folder2.ID})

	// Verify folders are trashed
	rootContentsTrashed := getFolderContents(t, env, token, 0)
	require.Len(t, rootContentsTrashed.FolderPage.Items, 0, "root should not contain any folders after trash")

	// Restore the folder which contains the file
	restoreFolder(t, env, token, folder2.ID)

	// Verify folder is restored
	rootContentsRestored := getFolderContents(t, env, token, 0)
	require.Len(t, rootContentsRestored.FolderPage.Items, 1, "root should contain exactly one folder after restore")
	require.Equal(t, folder2.ID, rootContentsRestored.FolderPage.Items[0].ID, "restored folder should be the correct one")

	// Verify the nested file is also restored
	contents2Restored := getFolderContents(t, env, token, folder2.ID)
	require.Len(t, contents2Restored.FilePage.Items, 1, "restored folder should contain exactly one file")
	require.Equal(t, file1.ID, contents2Restored.FilePage.Items[0].ID, "restored file should be the correct one")
}

func TestPagination(t *testing.T) {
	t.Parallel()
	env := setupTestEnv(t)
	_, token := createTestUser(t, env)

	fileName := func(i string) string {
		return fmt.Sprintf("file_%s.txt", i)
	}

	testCases := []struct {
		name     string
		sort     string
		expected []string // file names in expected order
	}{
		{
			name:     "ASC",
			sort:     paging.SortAsc,
			expected: []string{"01", "02", "03", "04", "05", "06", "07", "08"},
		},
		{
			name:     "DESC",
			sort:     "desc",
			expected: []string{"08", "07", "06", "05", "04", "03", "02", "01"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pagingOpt := func(direction, nextCursor, prevCursor string) *paging.Options {
				return &paging.Options{
					Direction:  direction,
					NextCursor: nextCursor,
					PrevCursor: prevCursor,
					Limit:      3,
					Sort:       tc.sort,
					SortBy:     paging.SortByName,
				}
			}

			folder := createFolder(t, env, token, 0, tc.name)

			// Upload 8 files to root folder
			for _, num := range tc.expected {
				uploadFile(t, env, token, folder.ID, fileName(num), media.BytesPerKB)
			}

			// Get first page (3 items) going forward
			contents := getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(paging.DirectionForward, "", ""))
			require.Len(t, contents.FilePage.Items, 3, "first page should contain exactly 3 files")
			require.True(t, contents.FilePage.HasMore, "should have more pages")
			require.NotEmpty(t, contents.FilePage.NextCursor, "first page should have next cursor")

			// Verify first page files
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i]), contents.FilePage.Items[i].Name)
			}

			// Get second page (3 items) going forward
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				paging.DirectionForward,
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 3, "second page should contain exactly 3 files")
			require.True(t, contents.FilePage.HasMore, "should have more pages")
			require.NotEmpty(t, contents.FilePage.PrevCursor, "second page should have prev cursor")
			require.NotEmpty(t, contents.FilePage.NextCursor, "second page should have next cursor")

			// Verify second page files
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i+3]), contents.FilePage.Items[i].Name)
			}

			// Get last page (2 items) going forward
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				paging.DirectionForward,
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 2, "last page should contain exactly 2 files")
			require.False(t, contents.FilePage.HasMore, "should not have more pages")
			require.NotEmpty(t, contents.FilePage.PrevCursor, "last page should have prev cursor")

			// Verify last page files
			for i := 0; i < 2; i++ {
				assert.Equal(t, fileName(tc.expected[i+6]), contents.FilePage.Items[i].Name)
			}

			// Go back one page using prev cursor going backward
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				"backward",
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 3, "previous page should contain exactly 3 files")
			require.True(t, contents.FilePage.HasMore, "should have more pages")
			require.NotEmpty(t, contents.FilePage.PrevCursor, "second page should have prev cursor")
			require.NotEmpty(t, contents.FilePage.NextCursor, "second page should have next cursor")

			// Verify we got back to the previous page
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i+3]), contents.FilePage.Items[i].Name)
			}

			// Reached first page using prev cursor going backward
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				"backward",
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 3, "previous page should contain exactly 3 files")
			require.False(t, contents.FilePage.HasMore, "should not have more pages")
			require.NotEmpty(t, contents.FilePage.NextCursor, "first page should have next cursor")

			// Verify we got back to the first page
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i]), contents.FilePage.Items[i].Name)
			}

			// Reached to the beginning of the list
			// Now, test forward and backward without reaching the end of the list

			// Again go forward from the first page to the second page
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				paging.DirectionForward,
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 3, "again second page should contain exactly 3 files")
			require.True(t, contents.FilePage.HasMore, "should have more pages")
			require.NotEmpty(t, contents.FilePage.PrevCursor, "second page should have prev cursor")
			require.NotEmpty(t, contents.FilePage.NextCursor, "second page should have next cursor")

			// Verify we got back to the second page
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i+3]), contents.FilePage.Items[i].Name)
			}

			// Come back from the second page to the first page
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, pagingOpt(
				"backward",
				contents.FilePage.NextCursor,
				contents.FilePage.PrevCursor),
			)
			require.Len(t, contents.FilePage.Items, 3, "again first page should contain exactly 3 files")
			require.False(t, contents.FilePage.HasMore, "should not have more pages")
			require.NotEmpty(t, contents.FilePage.NextCursor, "first page should have next cursor")

			// Verify we got back to the first page
			for i := 0; i < 3; i++ {
				assert.Equal(t, fileName(tc.expected[i]), contents.FilePage.Items[i].Name)
			}

			// Reached to the beginning of the list
			// Now test forward and backward in one jump with max. limit

			// Get all items with forward direction
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, &paging.Options{
				Limit:     8,
				Direction: paging.DirectionForward,
				Sort:      tc.sort,
				SortBy:    paging.SortByName,
			})
			require.Len(t, contents.FilePage.Items, 8, "all items should be returned with forward direction")
			require.False(t, contents.FilePage.HasMore, "should not have more pages")
			require.NotEmpty(t, contents.FilePage.NextCursor, "all items page should have next cursor")
			require.NotEmpty(t, contents.FilePage.PrevCursor, "all items page should have prev cursor")

			// Verify all items
			for i := 0; i < 8; i++ {
				assert.Equal(t, fileName(tc.expected[i]), contents.FilePage.Items[i].Name)
			}

			// In practical scenario, the following case should not be possible.
			// Going backwards,
			// when all items have been fetched without any pagination using the max. limit,
			// using the prev cursor would return empty items.
			// Since, the prev cursor already points to the beginning of the list.
			contents = getFolderContentsWithPaging(t, env, token, folder.ID, &paging.Options{
				Limit:      8,
				Direction:  "backward",
				NextCursor: contents.FilePage.NextCursor,
				PrevCursor: contents.FilePage.PrevCursor,
				Sort:       tc.sort,
				SortBy:     paging.SortByName,
			})
			require.Len(t, contents.FilePage.Items, 0, "no item is returned with backward direction")
			require.False(t, contents.FilePage.HasMore, "should not have more pages")
			require.Empty(t, contents.FilePage.NextCursor, "empty items will have empty next cursor")
			require.Empty(t, contents.FilePage.PrevCursor, "empty items will have empty prev cursor")
		})
	}
}
