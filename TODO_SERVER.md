# SkyVault Server TODO - Sharing Feature

## Epic 1: Contact Management System

### 1.1 API Layer Implementation
- [ ] Create `/server/internal/api/sharing_api.go`
  - [ ] Contact CRUD endpoints (`/api/v1/sharing/contacts`)
  - [ ] Contact group CRUD endpoints (`/api/v1/sharing/contact-groups`) 
  - [ ] Contact group membership endpoints (`/api/v1/sharing/contact-groups/{id}/members`)
  - [ ] JWT middleware integration for all endpoints

### 1.2 DTOs Implementation  
- [ ] Create `/server/internal/api/helper/dtos/sharing_dtos.go`
  - [ ] Contact request/response DTOs
  - [ ] Contact group request/response DTOs
  - [ ] Contact group member DTOs
  - [ ] Input validation structures

### 1.3 API Registration
- [ ] Update `/server/internal/api/api.go` to register sharing routes
- [ ] Add sharing API to main router configuration

## Epic 2: Core File Sharing

### 2.1 Sharing API Endpoints
- [ ] Add sharing endpoints to `sharing_api.go`:
  - [ ] `POST /api/v1/sharing/shares` - Create share
  - [ ] `GET /api/v1/sharing/shares/{id}` - Get share config  
  - [ ] `PUT /api/v1/sharing/shares/{id}/expiry` - Update expiry
  - [ ] `PUT /api/v1/sharing/shares/{id}/password` - Update password
  - [ ] `DELETE /api/v1/sharing/shares/{id}` - Delete share

### 2.2 Share Recipients API
- [ ] Add recipient management endpoints:
  - [ ] `POST /api/v1/sharing/shares/{id}/recipients` - Add recipient
  - [ ] `DELETE /api/v1/sharing/shares/{id}/recipients/{recipientId}` - Remove recipient
  - [ ] `GET /api/v1/sharing/shares/{id}/recipients` - List recipients

### 2.3 Share DTOs
- [ ] Add to `sharing_dtos.go`:
  - [ ] ShareConfig request/response DTOs
  - [ ] ShareRecipient request/response DTOs
  - [ ] CreateShare request DTO with validation
  - [ ] Share settings update DTOs

### 2.4 Public Share Access
- [ ] Create public endpoints (no JWT required):
  - [ ] `GET /api/v1/public/shares/{id}` - Get share info
  - [ ] `POST /api/v1/public/shares/{id}/validate` - Validate access
  - [ ] `GET /api/v1/public/shares/{id}/download` - Download shared content

## Epic 3: Shared Content Management

### 3.1 Sharing Queries API
- [ ] Add query endpoints to `sharing_api.go`:
  - [ ] `GET /api/v1/sharing/shared-with-me` - Files shared with current user
  - [ ] `GET /api/v1/sharing/shared-by-me` - Files shared by current user
  - [ ] `GET /api/v1/sharing/shares` - List user's shares with filters

### 3.2 Share Analytics
- [ ] Add analytics endpoints:
  - [ ] `GET /api/v1/sharing/shares/{id}/analytics` - Share download stats
  - [ ] `GET /api/v1/sharing/analytics/summary` - User sharing summary

### 3.3 Integration with Media Domain
- [ ] Update media DTOs to include sharing status
- [ ] Add sharing indicators to file/folder responses
- [ ] Ensure proper access control for shared content

## Epic 4: Advanced Sharing Features

### 4.1 Bulk Operations
- [ ] Add bulk sharing endpoints:
  - [ ] `POST /api/v1/sharing/bulk-share` - Share multiple files/folders
  - [ ] `PUT /api/v1/sharing/bulk-recipients` - Add recipients to multiple shares
  - [ ] `DELETE /api/v1/sharing/bulk-revoke` - Revoke multiple shares

### 4.2 Enhanced Analytics
- [ ] Add advanced analytics endpoints:
  - [ ] `GET /api/v1/sharing/analytics/downloads` - Download history
  - [ ] `GET /api/v1/sharing/analytics/access-logs` - Access logs per share
  - [ ] `GET /api/v1/sharing/analytics/recipients` - Recipient activity stats

### 4.3 Share Link Features
- [ ] Add share link management:
  - [ ] `POST /api/v1/sharing/shares/{id}/link` - Generate public link
  - [ ] `DELETE /api/v1/sharing/shares/{id}/link` - Revoke public link
  - [ ] `GET /api/v1/sharing/shares/{id}/qr` - Generate QR code

### 4.4 Security Enhancements
- [ ] Implement rate limiting for sharing endpoints
- [ ] Add audit logging for sharing operations
- [ ] Implement share access monitoring and alerts
- [ ] Add email notification system for share events

## Implementation Notes

### Error Handling
- All endpoints must return consistent AppError responses
- Implement proper validation for all input DTOs
- Add comprehensive error documentation

### Security
- Validate file/folder ownership before sharing
- Implement proper access control for shared content
- Hash passwords using bcrypt for share protection
- Rate limit public share access endpoints

### Database
- Leverage existing Jet SQL patterns
- Use transactions for multi-table operations
- Implement proper indexing for share queries
- Consider pagination for large result sets

### Testing  
- Unit tests for all command/query handlers
- Integration tests for API endpoints
- Test sharing permission scenarios
- Validate security boundaries