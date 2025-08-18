# SkyVault Web TODO - Sharing Feature

## Epic 1: Contact Management System

### 1.1 API Client Layer
- [ ] Create `/web/src/apis/sharing/`
  - [ ] `models.ts` - TypeScript interfaces for all sharing entities
    - [ ] Contact interface
    - [ ] ContactGroup interface  
    - [ ] ShareConfig interface
    - [ ] ShareRecipient interface
    - [ ] Request/Response DTOs
  - [ ] `index.ts` - API client functions
    - [ ] Contact CRUD operations
    - [ ] Contact group CRUD operations
    - [ ] Contact group membership operations

### 1.2 Contact Management Components
- [ ] Create `/web/src/components/sharing/contacts/`
  - [ ] `contactList.tsx` - Display all user contacts with search
  - [ ] `contactForm.tsx` - Add/edit contact modal
  - [ ] `contactItem.tsx` - Individual contact display with actions
  - [ ] `contactDeleteDialog.tsx` - Confirmation dialog for contact deletion

### 1.3 Contact Group Management
- [ ] Create `/web/src/components/sharing/groups/`
  - [ ] `groupList.tsx` - Display all contact groups
  - [ ] `groupForm.tsx` - Create/edit contact group modal
  - [ ] `groupMembers.tsx` - View and manage group membership
  - [ ] `memberPicker.tsx` - Multi-select contact picker for groups

### 1.4 Contact Manager Page
- [ ] Create `/web/src/pages/contacts.tsx`
  - [ ] Tabbed interface (Contacts / Groups)
  - [ ] Search and filter functionality
  - [ ] Import contacts functionality
  - [ ] Bulk operations (delete, group assignment)

## Epic 2: Core File Sharing

### 2.1 Share Dialog Component
- [ ] Create `/web/src/components/sharing/shareDialog.tsx`
  - [ ] Recipient selection interface
    - [ ] Contact picker with search/autocomplete
    - [ ] Contact group selection
    - [ ] Manual email input with validation
    - [ ] "Save as contact" toggle for new emails
  - [ ] Share settings panel
    - [ ] Password protection toggle + input
    - [ ] Expiration date picker (date-fns integration)
    - [ ] Download limit input with validation
  - [ ] Share preview and confirmation

### 2.2 Recipient Selection Components  
- [ ] Create `/web/src/components/sharing/recipients/`
  - [ ] `recipientPicker.tsx` - Multi-select interface for recipients
  - [ ] `recipientChip.tsx` - Individual recipient display chip
  - [ ] `recipientSearch.tsx` - Search across contacts, groups, and emails
  - [ ] `recipientValidation.tsx` - Email format and duplicate validation

### 2.3 Share Settings Components
- [ ] Create `/web/src/components/sharing/settings/`
  - [ ] `passwordSettings.tsx` - Password protection controls
  - [ ] `expirySettings.tsx` - Date picker for share expiration
  - [ ] `downloadLimitSettings.tsx` - Download count limits
  - [ ] `sharePreview.tsx` - Preview of share configuration

### 2.4 Integration with File/Folder Items
- [ ] Update `/web/src/components/folderContent/gridItem.tsx`
  - [ ] Add share button to hover actions
  - [ ] Add sharing indicator icon for shared items
  - [ ] Context menu integration for share option
- [ ] Update `/web/src/components/folderContent/listItem.tsx` 
  - [ ] Add share button to item actions
  - [ ] Add sharing status column
  - [ ] Context menu share option

## Epic 3: Shared Content Management

### 3.1 Shared Page Implementation
- [ ] Replace `/web/src/pages/shared.tsx` with full implementation
  - [ ] Tabbed interface: "Shared with Me" / "Shared by Me"
  - [ ] List/Grid view toggle (consistent with Drive)
  - [ ] Search functionality for shared items
  - [ ] Filter by share status (active, expired, password-protected)

### 3.2 Shared Content Components
- [ ] Create `/web/src/components/sharing/content/`
  - [ ] `sharedWithMe.tsx` - Items others shared with user
  - [ ] `sharedByMe.tsx` - Items user shared with others
  - [ ] `sharedItem.tsx` - Individual shared item display
  - [ ] `shareStatus.tsx` - Share status indicators and badges

### 3.3 Share Management Components
- [ ] Create `/web/src/components/sharing/management/`
  - [ ] `shareSettings.tsx` - Edit share settings panel
  - [ ] `shareRecipients.tsx` - View and manage share recipients
  - [ ] `shareAnalytics.tsx` - Download stats and access logs
  - [ ] `shareActions.tsx` - Share action buttons (edit, revoke, copy link)

### 3.4 Share Link Components  
- [ ] Create `/web/src/components/sharing/links/`
  - [ ] `shareLink.tsx` - Display and copy share link
  - [ ] `linkPreview.tsx` - Preview share link with access requirements
  - [ ] `qrCode.tsx` - QR code generation for mobile sharing
  - [ ] `linkSettings.tsx` - Public link specific settings

## Epic 4: Advanced Sharing Features

### 4.1 Public Share Access Flow
- [ ] Create `/web/src/pages/publicShare.tsx`
  - [ ] Public share access page (no auth required)
  - [ ] Password entry form for protected shares
  - [ ] Email verification for recipient-specific shares
  - [ ] Download interface for validated access

### 4.2 Access Validation Components
- [ ] Create `/web/src/components/sharing/access/`
  - [ ] `accessForm.tsx` - Password/email entry form
  - [ ] `accessValidation.tsx` - Validation logic and error handling
  - [ ] `downloadButton.tsx` - Download action after validation
  - [ ] `shareInfo.tsx` - Share information display

### 4.3 Bulk Sharing Operations
- [ ] Create `/web/src/components/sharing/bulk/`
  - [ ] `bulkShareDialog.tsx` - Share multiple files/folders at once
  - [ ] `bulkRecipientManager.tsx` - Add recipients to multiple shares
  - [ ] `bulkActions.tsx` - Bulk revoke, update settings
  - [ ] `selectionSummary.tsx` - Summary of selected items for sharing

### 4.4 Advanced Analytics
- [ ] Create `/web/src/components/sharing/analytics/`
  - [ ] `sharingDashboard.tsx` - User sharing overview
  - [ ] `downloadStats.tsx` - Download statistics charts
  - [ ] `accessLogs.tsx` - Share access history
  - [ ] `recipientActivity.tsx` - Per-recipient activity tracking

### 4.5 Mobile Optimization
- [ ] Mobile-specific sharing components
  - [ ] Touch-friendly recipient selection
  - [ ] Mobile share sheet integration  
  - [ ] Responsive contact management
  - [ ] Swipe actions for share management

## State Management & Integration

### 5.1 Sharing State Management
- [ ] Create `/web/src/store/sharing/`
  - [ ] `sharingCtx.ts` - Sharing context and state
  - [ ] `contactsStore.ts` - Contacts cache and management
  - [ ] `sharesStore.ts` - Active shares tracking
  - [ ] `notificationsStore.ts` - Share notifications

### 5.2 Navigation and Routing
- [ ] Update `/web/src/routes.tsx`
  - [ ] Add public share routes (`/share/:shareId`)
  - [ ] Add contact management route (`/contacts`)
  - [ ] Route guards for authenticated vs public access

### 5.3 App Integration
- [ ] Update `/web/src/pages/appNavigation.tsx`
  - [ ] Add notification badges for new shared items
  - [ ] Consider adding contacts link in user menu
- [ ] Update main app components for sharing indicators

## UI/UX Components

### 6.1 Specialized UI Components
- [ ] Create `/web/src/components/sharing/ui/`
  - [ ] `recipientAvatar.tsx` - Contact/group avatar display
  - [ ] `shareStatusBadge.tsx` - Visual share status indicators  
  - [ ] `permissionIcon.tsx` - Permission level icons
  - [ ] `sharingTooltip.tsx` - Contextual sharing help

### 6.2 Form Components
- [ ] Sharing-specific form controls
  - [ ] Date/time picker for expiration
  - [ ] Password strength indicator
  - [ ] Email validation with suggestions
  - [ ] Multi-select with custom options

### 6.3 Loading and Error States
- [ ] Sharing-specific loading states
- [ ] Error handling for sharing operations
- [ ] Success/failure notifications
- [ ] Retry mechanisms for failed operations

## Implementation Notes

### Styling
- Follow existing CSS class patterns from `/web/src/index.css`
- Use semantic classes (`text-primary`, `btn-primary`, etc.)
- Maintain responsive design principles
- Ensure accessibility compliance

### Performance
- Implement virtual scrolling for large contact/share lists
- Debounce search inputs
- Cache contact and group data
- Lazy load sharing components

### Security
- Sanitize all user inputs
- Validate email formats
- Handle authentication states properly
- Secure public share access flows

### Testing
- Unit tests for sharing components
- Integration tests for sharing flows
- E2E tests for complete sharing scenarios
- Accessibility testing for all components