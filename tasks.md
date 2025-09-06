# Implementation Tasks

## Media Options Feature

### File Operations

#### 1. Download File
**API Endpoint:** `GET /api/v1/media/files/{file-id}/download`
- **Implementation:** `DownloadFile` handler in media_api.go:346
- **Method:** Uses `http.ServeContent` to serve file with proper headers
- **Frontend Tasks:**
  - Add download button/option in file context menu or file actions
  - Make GET request to download endpoint
  - Handle file download in browser (triggers browser download)
  - Add appropriate icons and loading states

#### 2. Rename File  
**API Endpoint:** `PATCH /api/v1/media/files/{file-id}/rename`
- **Implementation:** `RenameFile` handler in media_api.go:403
- **Request Body:** `{"name": "new-filename.ext"}`
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add rename option in file context menu
  - Create rename modal/dialog with input field
  - Validate file name input (non-empty, reasonable length)
  - Make PATCH request with new name
  - Update file name in UI after success
  - Handle validation errors from backend

#### 3. Move File
**API Endpoint:** `PATCH /api/v1/media/files/{file-id}/move`
- **Implementation:** `MoveFile` handler in media_api.go:436
- **Request Body:** `{"folderId": "target-folder-uuid"}` (empty string for root)
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add move option in file context menu
  - Create folder selection dialog/tree view
  - Allow selecting target folder (including root)
  - Make PATCH request with target folder ID
  - Remove file from current view after success
  - Handle folder hierarchy validation errors

#### 4. Trash File
**API Endpoint:** `DELETE /api/v1/media/files/`
- **Implementation:** `TrashFiles` handler in media_api.go:371 (supports bulk)
- **Request Body:** `{"fileIds": ["file-uuid-1"]}`
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add delete/trash option in file context menu
  - Show confirmation dialog before deletion
  - Make DELETE request with file ID array
  - Remove file from current view after success
  - Show success feedback to user
  - Handle invalid file ID errors

### Folder Operations

#### 1. Rename Folder
**API Endpoint:** `PATCH /api/v1/media/folders/{folder-id}/rename`
- **Implementation:** `RenameFolder` handler in media_api.go:569
- **Request Body:** `{"name": "new-folder-name"}`
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add rename option in folder context menu
  - Create rename modal/dialog with input field
  - Validate folder name input (non-empty, reasonable length)
  - Make PATCH request with new name
  - Update folder name in UI after success
  - Handle validation errors from backend

#### 2. Move Folder
**API Endpoint:** `PATCH /api/v1/media/folders/{folder-id}/move`
- **Implementation:** `MoveFolder` handler in media_api.go:602
- **Request Body:** `{"folderId": "target-parent-folder-uuid"}` (empty string for root)
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add move option in folder context menu
  - Create folder selection dialog (exclude current folder and descendants)
  - Allow selecting target parent folder (including root)
  - Make PATCH request with target parent folder ID
  - Remove folder from current view after success
  - Handle circular dependency validation errors

#### 3. Trash Folder
**API Endpoint:** `DELETE /api/v1/media/folders/`
- **Implementation:** `TrashFolders` handler in media_api.go:536 (supports bulk)
- **Request Body:** `{"folderIds": ["folder-uuid-1"]}`
- **Response:** 204 No Content on success
- **Frontend Tasks:**
  - Add delete/trash option in folder context menu  
  - Show confirmation dialog with warning about contents
  - Make DELETE request with folder ID array
  - Remove folder from current view after success
  - Show success feedback to user
  - Handle invalid folder ID errors

### UI/UX Implementation Details

#### Context Menu System
- Implement right-click context menus for files and folders
- Include relevant options based on item type
- Use consistent styling with existing UI components

#### Modals and Dialogs
- Create reusable modal components for rename operations
- Implement folder picker component for move operations
- Add confirmation dialogs for destructive operations (trash)
- Ensure proper keyboard navigation and accessibility

#### Error Handling
- Display user-friendly error messages for API failures
- Handle network errors gracefully with retry options
- Validate input on frontend before API calls
- Show loading states during operations

#### State Management
- Update local state after successful operations
- Refresh folder contents when items are moved/deleted
- Handle optimistic updates where appropriate
- Maintain consistent state across components

### API Integration Notes
- All endpoints require authentication (JWT token)
- File and folder IDs must be valid UUIDs
- Error responses follow AppError format with metadata
- Use appropriate HTTP status codes for different scenarios