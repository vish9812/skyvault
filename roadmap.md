# SkyVault Roadmap

This roadmap replaces the older TODO files. It is ordered by user value and implementation dependency: finish the core drive experience first, then build sharing in usable slices, then add advanced collaboration and analytics.

## Current Baseline

Completed:
- User authentication with JWT.
- Folder creation and navigation.
- File upload, including chunked uploads.
- Storage quota tracking and upload quota validation.
- Server APIs for file download, rename, move, and trash.
- Server APIs for folder rename, move, and trash.

In progress:
- Web UI for file and folder operations.
- Sharing domain model, commands, and queries exist, but HTTP APIs and web UI are not complete.

Tracked GitHub issues:
- [#5 Expire the shared link after t time](https://github.com/vish9812/skyvault/issues/5)
- [#6 Expire the shared link after n no. of downloads](https://github.com/vish9812/skyvault/issues/6)
- [#7 Expire the shared link](https://github.com/vish9812/skyvault/issues/7)
- [#8 Allow files versioning upto n numbers](https://github.com/vish9812/skyvault/issues/8)
- [#10 Add OIDC Auth](https://github.com/vish9812/skyvault/issues/10)

## 1. Complete Core Drive Operations

Goal: users can manage their files and folders without leaving the web app.

Server status:
- [x] `GET /api/v1/media/files/{file-id}/download`
- [x] `PATCH /api/v1/media/files/{file-id}/rename`
- [x] `PATCH /api/v1/media/files/{file-id}/move`
- [x] `DELETE /api/v1/media/files/`
- [x] `PATCH /api/v1/media/folders/{folder-id}/rename`
- [x] `PATCH /api/v1/media/folders/{folder-id}/move`
- [x] `DELETE /api/v1/media/folders/`

Web work:
- [x] Add media API client methods for download, rename, move, and trash.
- [x] Add file actions in grid and list views: download, rename, move, trash.
- [x] Add folder actions in grid and list views: rename, move, trash.
- [x] Add a consistent action menu for mouse, keyboard, and touch users.
- [x] Add reusable rename dialogs with validation for files and folders.
- [x] Add a folder picker for move operations, including root selection.
- [x] Exclude invalid move targets, including moving a folder into itself or descendants.
- [x] Add destructive confirmation dialogs for trash operations.
- [x] Show loading states while actions are running.
- [x] Show user-friendly errors for validation, permission, network, and server failures.
- [x] Refresh or update local state after successful rename, move, and trash operations.
- [x] Keep grid and list behavior consistent.
- [ ] Verify keyboard navigation and accessible dialog behavior.

## 2. Finish Trash Lifecycle

Goal: trash is safe, reversible, and understandable.

Server work:
- [ ] Confirm whether restore and permanent delete APIs exist.
- [ ] Add restore APIs for trashed files and folders if missing.
- [ ] Add permanent delete APIs for files and folders if missing.
- [ ] Add bulk restore and bulk permanent delete where practical.
- [ ] Ensure ownership checks and quota updates are correct.

Web work:
- [ ] Complete the trash page with trashed files and folders.
- [ ] Add restore actions.
- [ ] Add permanent delete actions with stronger confirmation.
- [ ] Add empty, loading, and error states.
- [ ] Refresh drive views after restore or permanent deletion.

## 3. Add Search and Browsing Polish

Goal: users can find content quickly once their drive has real data.

Server work:
- [ ] Add or confirm search queries for files and folders by name.
- [ ] Support paging and reasonable ordering for search results.
- [ ] Ensure search respects ownership and shared access permissions.

Web work:
- [ ] Wire the existing search UI to real API results.
- [ ] Search across files and folders.
- [ ] Debounce search input.
- [ ] Add empty, loading, and error states.
- [ ] Support opening result locations.
- [ ] Keep search usable on mobile.

## 4. Improve Upload Reliability

Goal: uploads behave predictably under real-world conditions.

Server work:
- [ ] Validate interrupted chunked-upload recovery behavior.
- [ ] Confirm abandoned upload cleanup.
- [ ] Confirm concurrent upload quota reservation behavior.
- [ ] Add tests for quota, interruption, duplicate chunks, and cleanup edge cases.

Web work:
- [ ] Improve upload progress states.
- [ ] Show clear quota exceeded errors.
- [ ] Add retry behavior where the server supports it.
- [ ] Handle duplicate names clearly.
- [ ] Make interrupted or failed uploads easy to understand.

## 5. Add File Versioning

Goal: users can recover earlier versions of files without creating manual copies.

Tracked issue:
- [#8 Allow files versioning upto n numbers](https://github.com/vish9812/skyvault/issues/8)

Server work:
- [ ] Define versioning behavior for uploads with the same name or same file ID.
- [ ] Add a configurable maximum versions setting per file or globally.
- [ ] Store file version metadata, including version number, size, checksum if available, creator, and created time.
- [ ] Keep quota accounting correct across retained versions.
- [ ] Add APIs to list versions, download a specific version, restore a version, and delete old versions.
- [ ] Add cleanup behavior when version count exceeds the configured limit.
- [ ] Add tests for version retention, restore, quota accounting, and deletion.

Web work:
- [ ] Show version history from file actions.
- [ ] Add download, restore, and delete actions for individual versions.
- [ ] Show clear metadata for each version.
- [ ] Add confirmation for deleting versions.
- [ ] Explain when old versions are pruned by the configured limit.

## 6. Build Sharing API Foundation

Goal: expose the existing sharing domain through stable HTTP APIs.

Tracked issues:
- [#5 Expire the shared link after t time](https://github.com/vish9812/skyvault/issues/5)
- [#6 Expire the shared link after n no. of downloads](https://github.com/vish9812/skyvault/issues/6)
- [#7 Expire the shared link](https://github.com/vish9812/skyvault/issues/7)

Server API work:
- [ ] Create `server/internal/api/sharing_api.go`.
- [ ] Create sharing DTOs under `server/internal/api/helper/dtos/`.
- [ ] Register sharing routes in `server/internal/api/api.go`.
- [ ] Add JWT-protected contact endpoints under `/api/v1/sharing/contacts`.
- [ ] Add JWT-protected contact group endpoints under `/api/v1/sharing/contact-groups`.
- [ ] Add contact group membership endpoints under `/api/v1/sharing/contact-groups/{id}/members`.
- [ ] Add share configuration endpoints:
  - [ ] `POST /api/v1/sharing/shares`
  - [ ] `GET /api/v1/sharing/shares/{id}`
  - [ ] `PUT /api/v1/sharing/shares/{id}/expiry`
  - [ ] `PUT /api/v1/sharing/shares/{id}/password`
  - [ ] `DELETE /api/v1/sharing/shares/{id}`
- [ ] Add recipient endpoints:
  - [ ] `POST /api/v1/sharing/shares/{id}/recipients`
  - [ ] `GET /api/v1/sharing/shares/{id}/recipients`
  - [ ] `DELETE /api/v1/sharing/shares/{id}/recipients/{recipientId}`
- [ ] Add public share access endpoints:
  - [ ] `GET /api/v1/public/shares/{id}`
  - [ ] `POST /api/v1/public/shares/{id}/validate`
  - [ ] `GET /api/v1/public/shares/{id}/download`
- [ ] Enforce shared-link expiration by time.
- [ ] Enforce shared-link expiration by download count.
- [ ] Return clear expired, revoked, and max-downloads-reached errors.

Server quality requirements:
- [ ] Return consistent AppError responses.
- [ ] Validate all request DTOs.
- [ ] Validate file and folder ownership before creating shares.
- [ ] Enforce access control for private and public share access.
- [ ] Hash share passwords with bcrypt.
- [ ] Use transactions for multi-table operations.
- [ ] Add indexes for share lookup and recipient queries.
- [ ] Add paging for potentially large lists.
- [ ] Add integration tests for API endpoints and permission boundaries.

## 7. Ship Basic Sharing UI

Goal: users can share a file or folder end to end.

Web API client:
- [ ] Create `web/src/apis/sharing/models.ts`.
- [ ] Create `web/src/apis/sharing/index.ts`.
- [ ] Add TypeScript models for contacts, contact groups, share configs, recipients, and request/response DTOs.
- [ ] Add API functions for contacts, groups, group membership, shares, recipients, and public share access.

Share dialog:
- [ ] Create `web/src/components/sharing/shareDialog.tsx`.
- [ ] Add recipient selection from contacts and groups.
- [ ] Add manual email entry with validation.
- [ ] Add optional "save as contact" behavior.
- [ ] Add password protection controls.
- [ ] Add expiration controls.
- [ ] Add download limit controls.
- [ ] Show link expiration status for time-based and download-count-based limits.
- [ ] Add share preview and confirmation.

Drive integration:
- [ ] Add share actions to `gridItem.tsx`.
- [ ] Add share actions to `listItem.tsx`.
- [ ] Add sharing indicators for shared files and folders.
- [ ] Add sharing status to file and folder responses if needed.
- [ ] Keep sharing actions consistent with the core action menu.

## 8. Add Contact Management

Goal: make repeated sharing practical.

Web components:
- [ ] Create contact list, contact item, contact form, and contact delete dialog components.
- [ ] Create contact group list, group form, group members, and member picker components.
- [ ] Create `web/src/pages/contacts.tsx`.
- [ ] Add tabs for contacts and groups.
- [ ] Add contact/group search and filtering.
- [ ] Add bulk delete and group assignment if the API supports it.
- [ ] Add import contacts only after the manual contact flow is stable.

Navigation:
- [ ] Add a contacts route.
- [ ] Add a contacts entry in navigation or the user menu.

## 9. Build Shared Content Management

Goal: users can see and manage what they shared and what was shared with them.

Server work:
- [ ] Add `GET /api/v1/sharing/shared-with-me`.
- [ ] Add `GET /api/v1/sharing/shared-by-me`.
- [ ] Add `GET /api/v1/sharing/shares` with filters.
- [ ] Ensure shared content access works through the media domain.

Web work:
- [ ] Replace `web/src/pages/shared.tsx` with a full implementation.
- [ ] Add tabs for "Shared with Me" and "Shared by Me".
- [ ] Add list/grid view toggle consistent with Drive.
- [ ] Add search and filters by active, expired, and password-protected status.
- [ ] Add shared item, share status, shared-with-me, and shared-by-me components.
- [ ] Add share settings, recipients, analytics summary, and revoke/copy-link actions.

## 10. Public Share Experience

Goal: recipients can access shared content without a SkyVault account where allowed.

Tracked issues:
- [#5 Expire the shared link after t time](https://github.com/vish9812/skyvault/issues/5)
- [#6 Expire the shared link after n no. of downloads](https://github.com/vish9812/skyvault/issues/6)
- [#7 Expire the shared link](https://github.com/vish9812/skyvault/issues/7)

Web work:
- [ ] Add `/share/:shareId` route.
- [ ] Create a public share page that does not require authentication.
- [ ] Show share information before access where appropriate.
- [ ] Add password entry for protected shares.
- [ ] Add email verification for recipient-specific shares.
- [ ] Add download action after validation.
- [ ] Add clear expired, revoked, unauthorized, and max-downloads-reached states.

## 11. Add OIDC Authentication

Goal: allow deployments to authenticate users through an external identity provider.

Tracked issue:
- [#10 Add OIDC Auth](https://github.com/vish9812/skyvault/issues/10)

Server work:
- [ ] Define supported OIDC configuration: issuer URL, client ID, client secret, redirect URL, scopes, and allowed domains if needed.
- [ ] Add OIDC login, callback, and logout flow.
- [ ] Map OIDC identities to SkyVault users.
- [ ] Decide account-linking behavior for existing email/password users.
- [ ] Preserve JWT/session behavior after successful OIDC login.
- [ ] Add configuration validation and clear startup errors.
- [ ] Add tests for callback validation, user provisioning, account linking, and rejected providers.

Web work:
- [ ] Add OIDC sign-in option when enabled.
- [ ] Handle callback errors cleanly.
- [ ] Keep email/password auth available unless explicitly disabled by configuration.
- [ ] Add account/profile indication for OIDC-backed users if useful.

## 12. Sharing State, Notifications, and Mobile UX

Goal: make sharing feel integrated with the app.

State management:
- [ ] Create sharing state/context as needed.
- [ ] Cache contacts and groups.
- [ ] Track active shares where useful.
- [ ] Add notifications for share operations.

Navigation and app integration:
- [ ] Add notification badges for newly shared items if supported.
- [ ] Add route guards for authenticated and public sharing routes.
- [ ] Add consistent share indicators across app views.

Mobile:
- [ ] Make recipient selection touch-friendly.
- [ ] Optimize share management actions for small screens.
- [ ] Consider native mobile share sheet integration where available.
- [ ] Verify dialogs and pickers fit on mobile viewports.

## 13. Advanced Sharing

Goal: add power-user features after the basic sharing loop is stable.

Bulk operations:
- [ ] Add bulk share API.
- [ ] Add bulk recipient management API.
- [ ] Add bulk revoke API.
- [ ] Add bulk share dialog and selection summary UI.

Share links:
- [ ] Add public link generation and revocation.
- [ ] Add share link display and copy behavior.
- [ ] Add link preview.
- [ ] Add QR code generation only after link sharing is stable.

Analytics:
- [ ] Add share download stats endpoint.
- [ ] Add sharing summary endpoint.
- [ ] Add download history endpoint.
- [ ] Add access logs endpoint.
- [ ] Add recipient activity endpoint.
- [ ] Add dashboard, download stats, access logs, and recipient activity UI.

Security and operations:
- [ ] Rate limit sharing endpoints, especially public access.
- [ ] Add audit logging for sharing operations.
- [ ] Add share access monitoring.
- [ ] Add email notifications for share events.

## 14. UI Components, Performance, and Testing

Reusable UI:
- [ ] Add recipient avatar, share status badge, permission icon, and sharing tooltip components.
- [ ] Add date/time picker for expiration.
- [ ] Add password strength indicator.
- [ ] Add email validation suggestions.
- [ ] Add multi-select with custom options.
- [ ] Add sharing-specific loading, success, failure, and retry states.

Performance:
- [ ] Debounce search inputs.
- [ ] Cache contact and group data.
- [ ] Lazy-load sharing-heavy components where useful.
- [ ] Add virtual scrolling for large contact/share lists if needed.

Testing:
- [ ] Add server unit tests for sharing command and query handlers.
- [ ] Add server integration tests for sharing API endpoints.
- [ ] Add permission and security boundary tests.
- [ ] Add web component tests for core drive operations.
- [ ] Add web integration tests for sharing flows.
- [ ] Add end-to-end tests for upload, file management, trash, and sharing.
- [ ] Add accessibility checks for menus, dialogs, and public share flows.
