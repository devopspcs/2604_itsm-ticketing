# Phase 8: Integration & Testing Checklist
## Jira-like Project Board Upgrade

---

## 8.1 End-to-End Testing Checklist

### Sprint Workflow
- [ ] Create sprint with name, goal, dates
- [ ] Verify sprint status = "Planned"
- [ ] Verify sprint appears in sprint list
- [ ] Assign 5 records to sprint
- [ ] Verify records appear in sprint board
- [ ] Verify records removed from backlog
- [ ] Start sprint
- [ ] Verify sprint status = "Active"
- [ ] Verify actual_start_date set
- [ ] Transition record: To Do → In Progress
- [ ] Verify status change in board
- [ ] Verify activity log entry created
- [ ] Transition record: In Progress → Done
- [ ] Verify completion reflected in metrics
- [ ] Complete sprint
- [ ] Verify sprint status = "Completed"
- [ ] Verify actual_end_date set
- [ ] Verify sprint metrics calculated:
  - [ ] Total records = 5
  - [ ] Completed records = 3
  - [ ] Completion % = 60%
- [ ] Verify incomplete records moved to backlog
- [ ] Verify completed records in sprint history

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Backlog Workflow
- [ ] Create 10 records without sprint
- [ ] Verify all records in backlog
- [ ] Verify records ordered by priority
- [ ] Drag record #3 to position #1
- [ ] Verify priority updated
- [ ] Verify order persisted
- [ ] Drag record #5 to position #2
- [ ] Verify new ordering maintained
- [ ] Select records #1, #2, #3
- [ ] Bulk assign to sprint
- [ ] Verify records removed from backlog
- [ ] Verify records in sprint board
- [ ] Verify priority maintained in sprint
- [ ] Drag record from sprint to backlog
- [ ] Verify record removed from sprint
- [ ] Verify record added to backlog
- [ ] Complete sprint with incomplete records
- [ ] Verify incomplete records at top of backlog
- [ ] Verify priority ordering maintained

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Comment Workflow
- [ ] Open record detail
- [ ] Add comment "This needs review"
- [ ] Verify comment appears with author, timestamp
- [ ] Verify comment text displayed
- [ ] Add comment with "@john.doe"
- [ ] Verify mention dropdown appears
- [ ] Select user from dropdown
- [ ] Verify mention inserted
- [ ] Submit comment
- [ ] Verify mentions parsed
- [ ] Verify CommentMention record created
- [ ] Verify mention highlighted
- [ ] Verify notification sent to john.doe
- [ ] Verify notification contains record link
- [ ] Verify notification contains comment text
- [ ] Verify notification marked unread
- [ ] Edit comment to add "@jane.smith"
- [ ] Verify comment updated
- [ ] Verify new mention parsed
- [ ] Verify new notification sent
- [ ] Verify "edited" indicator shown
- [ ] Delete comment
- [ ] Verify comment removed
- [ ] Verify CommentMention records deleted
- [ ] Verify activity log records deletion

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Attachment Workflow
- [ ] Open record detail
- [ ] Click "Add Attachment"
- [ ] Select file (test.pdf, 2MB)
- [ ] Verify file uploaded
- [ ] Verify attachment appears with:
  - [ ] File name: test.pdf
  - [ ] File size: 2MB
  - [ ] Uploader: current user
  - [ ] Upload timestamp
- [ ] Upload image.jpg (1MB)
- [ ] Upload document.docx (3MB)
- [ ] Verify all 3 files in list
- [ ] Verify attachment count = 3 in header
- [ ] Attempt to upload .exe file
- [ ] Verify error: "File type not supported"
- [ ] Attempt to upload 60MB file
- [ ] Verify error: "File exceeds 50MB limit"
- [ ] Click on attachment
- [ ] Verify file downloads correctly
- [ ] Verify file content matches original
- [ ] Hover over image.jpg
- [ ] Verify thumbnail preview displayed
- [ ] Click on image
- [ ] Verify image opens in viewer
- [ ] Delete test.pdf
- [ ] Verify attachment removed
- [ ] Verify attachment count = 2
- [ ] Verify file removed from storage
- [ ] Delete record with attachments
- [ ] Verify all attachments deleted
- [ ] Verify files removed from storage

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

## 8.2 Performance Testing Checklist

### Search Performance (10,000 Records)
- [ ] Create project with 10,000 records
- [ ] Distribute across 5 sprints
- [ ] Add various labels, custom fields, assignees
- [ ] Search for "urgent"
- [ ] Measure response time: < 500ms ✓
- [ ] Verify results accurate
- [ ] Search for "bug fix"
- [ ] Measure response time: < 500ms ✓
- [ ] Filter by "Bug" (2,000 records)
- [ ] Measure response time: < 300ms ✓
- [ ] Verify correct records returned
- [ ] Filter by "In Progress" (3,000 records)
- [ ] Measure response time: < 300ms ✓
- [ ] Filter by single assignee (500 records)
- [ ] Measure response time: < 200ms ✓
- [ ] Filter: Issue Type = "Bug" AND Status = "In Progress" AND Assignee = "john.doe"
- [ ] Measure response time: < 500ms ✓
- [ ] Verify results accurate
- [ ] Save complex filter
- [ ] Load saved filter
- [ ] Measure response time: < 200ms ✓

**Performance Targets Met**: ⏳ Pending  
**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Bulk Operations Performance (1,000 Records)
- [ ] Select 1,000 records
- [ ] Change status to "In Progress"
- [ ] Measure response time: < 2 seconds ✓
- [ ] Verify all records updated
- [ ] Select 500 records
- [ ] Assign to user "john.doe"
- [ ] Measure response time: < 1 second ✓
- [ ] Verify all records assigned
- [ ] Select 800 records
- [ ] Add label "urgent"
- [ ] Measure response time: < 1.5 seconds ✓
- [ ] Verify all records labeled
- [ ] Select 200 records
- [ ] Delete all
- [ ] Measure response time: < 1 second ✓
- [ ] Verify all records deleted

**Performance Targets Met**: ⏳ Pending  
**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Sprint Board Rendering (500 Records)
- [ ] Load sprint board with 500 records
- [ ] Measure page load time: < 2 seconds ✓
- [ ] Verify all records rendered
- [ ] Drag record between columns
- [ ] Measure response time: < 300ms ✓
- [ ] Verify status updated
- [ ] Apply filter (assignee = "john.doe")
- [ ] Measure response time: < 500ms ✓
- [ ] Verify filtered records displayed
- [ ] Scroll through 500 records
- [ ] Verify smooth scrolling (60 FPS) ✓
- [ ] No lag or jank observed
- [ ] Search for "urgent" in board
- [ ] Measure response time: < 300ms ✓
- [ ] Verify results highlighted

**Performance Targets Met**: ⏳ Pending  
**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

## 8.3 Backward Compatibility Testing Checklist

### Existing Projects
- [ ] Open project created before Jira upgrade
- [ ] Verify project loads without errors
- [ ] Verify all existing records display
- [ ] Verify default issue type scheme created
- [ ] Verify default workflow created
- [ ] Verify default custom fields created
- [ ] Verify all existing records have:
  - [ ] Default issue type = "Task"
  - [ ] Default status = "To Do"
  - [ ] All existing fields preserved
- [ ] Verify existing columns still display
- [ ] Verify column names match workflow statuses
- [ ] Verify drag-and-drop still works
- [ ] Verify existing filters still work
- [ ] Verify search still works
- [ ] Verify sorting still works

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Existing Records Display
- [ ] Verify record title displays
- [ ] Verify record description displays
- [ ] Verify assignee avatar displays
- [ ] Verify due date displays
- [ ] Verify issue type icon displays (default Task)
- [ ] Open record detail
- [ ] Verify all fields display
- [ ] Verify custom fields section displays
- [ ] Verify comments section displays
- [ ] Verify attachments section displays
- [ ] Filter by assignee
- [ ] Verify correct records displayed
- [ ] Filter by status
- [ ] Verify correct records displayed
- [ ] Filter by label
- [ ] Verify correct records displayed

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Existing Drag-and-Drop
- [ ] Drag record from "To Do" to "In Progress"
- [ ] Verify record moves to new column
- [ ] Verify status updated in database
- [ ] Verify activity log records change
- [ ] Drag record within same column
- [ ] Verify record reorders
- [ ] Verify order persisted
- [ ] Drag record to backlog
- [ ] Verify record removed from sprint
- [ ] Verify record added to backlog
- [ ] Drag record from backlog to sprint
- [ ] Verify record added to sprint
- [ ] Verify record removed from backlog

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

## 8.4 Migration Testing Checklist

### Database Migration
- [ ] Verify existing tables exist
- [ ] Verify existing data intact
- [ ] Verify no new tables exist
- [ ] Execute migration 000009_jira_features.up.sql
- [ ] Verify migration completes without errors
- [ ] Verify all new tables created
- [ ] Verify issue_types table created with correct columns
- [ ] Verify issue_type_schemes table created
- [ ] Verify custom_fields table created
- [ ] Verify workflows table created
- [ ] Verify sprints table created
- [ ] Verify comments table created
- [ ] Verify attachments table created
- [ ] Verify labels table created
- [ ] Verify all indexes created
- [ ] Verify existing project_records data intact
- [ ] Verify foreign key constraints working
- [ ] Verify cascade deletes configured
- [ ] Execute migration rollback
- [ ] Verify all new tables removed
- [ ] Verify existing data intact
- [ ] Verify no errors during rollback

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Default Configuration Creation
- [ ] For each existing project:
  - [ ] Create default issue type scheme
  - [ ] Create default workflow
  - [ ] Create default custom fields
- [ ] Verify scheme created with all 5 issue types
- [ ] Verify scheme assigned to project
- [ ] Verify scheme items created for each type
- [ ] Verify workflow created with name "Default Workflow"
- [ ] Verify statuses created: To Do, In Progress, In Review, Done
- [ ] Verify transitions created for all valid paths
- [ ] Verify initial status = "To Do"
- [ ] Verify no custom fields created (empty by default)
- [ ] Verify custom field table accessible
- [ ] Verify no labels created (empty by default)
- [ ] Verify label table accessible

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Existing Data Preservation
- [ ] Count existing projects: 50
- [ ] Count existing records: 5,000
- [ ] Count existing users: 100
- [ ] Count existing comments: 1,000
- [ ] Execute migration
- [ ] Verify project count = 50
- [ ] Verify record count = 5,000
- [ ] Verify user count = 100
- [ ] Verify comment count = 1,000
- [ ] Verify all data values unchanged
- [ ] Verify record titles unchanged
- [ ] Verify record descriptions unchanged
- [ ] Verify record assignees unchanged
- [ ] Verify record due dates unchanged
- [ ] Verify record created_at unchanged
- [ ] Verify project-record relationships intact
- [ ] Verify user-record relationships intact
- [ ] Verify all foreign keys valid

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

## 8.5 User Acceptance Testing Checklist

### Project Manager UAT
- [ ] Create sprint "Sprint 1" with goal "Implement user auth"
- [ ] Set dates: 2024-01-15 to 2024-01-29
- [ ] Verify sprint created
- [ ] View backlog with 50 records
- [ ] Select top 10 priority records
- [ ] Bulk assign to sprint
- [ ] Verify records in sprint
- [ ] Start sprint
- [ ] View sprint board
- [ ] Verify all 10 records displayed
- [ ] Verify sprint metrics: 10 total, 0 completed, 0%
- [ ] View sprint board daily
- [ ] See records transitioning through statuses
- [ ] View updated metrics
- [ ] Verify completion % increasing
- [ ] Complete sprint
- [ ] View sprint metrics: 10 total, 8 completed, 80%
- [ ] View velocity: 8 points
- [ ] View incomplete records moved to backlog
- [ ] Create Sprint 2
- [ ] Assign incomplete records + new records
- [ ] Start Sprint 2

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Team Member UAT
- [ ] Open sprint board
- [ ] See all records assigned to sprint
- [ ] See records organized by status
- [ ] Click on "To Do" record
- [ ] View record details
- [ ] Update custom fields
- [ ] Add comment "Started working on this"
- [ ] Drag record to "In Progress"
- [ ] Add comment "@john.doe please review"
- [ ] Verify mention notification sent
- [ ] See john.doe's response comment
- [ ] Add attachment (design.pdf)
- [ ] Add label "design-review"
- [ ] Update custom fields (% complete = 100%)
- [ ] Drag record to "Done"
- [ ] Verify record marked complete
- [ ] View backlog records
- [ ] See priority ordering
- [ ] See unassigned records
- [ ] Search for "urgent"
- [ ] Filter by "assigned to me"
- [ ] Filter by "due this week"
- [ ] Verify correct records displayed

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

### Project Owner UAT
- [ ] Open project settings
- [ ] View configuration tabs
- [ ] View issue type scheme
- [ ] Verify all 5 types available
- [ ] Verify types assigned to project
- [ ] View workflow
- [ ] See statuses: To Do, In Progress, In Review, Done
- [ ] See transitions between statuses
- [ ] Verify workflow applied to all records
- [ ] Create custom field "Priority" (dropdown)
- [ ] Add options: Low, Medium, High, Critical
- [ ] Mark as required
- [ ] Verify field appears on all new records
- [ ] Create label "urgent" (red)
- [ ] Create label "design-review" (blue)
- [ ] Create label "blocked" (yellow)
- [ ] Verify labels available for all records
- [ ] View project members
- [ ] Add new member
- [ ] Remove member
- [ ] Verify permissions updated
- [ ] View activity log
- [ ] See record changes
- [ ] See configuration changes
- [ ] See user actions

**Status**: ⏳ Not Started  
**Pass/Fail**: -

---

## Test Execution Summary

### Overall Status
- **8.1 End-to-End Testing**: ⏳ Not Started
- **8.2 Performance Testing**: ⏳ Not Started
- **8.3 Backward Compatibility Testing**: ⏳ Not Started
- **8.4 Migration Testing**: ⏳ Not Started
- **8.5 User Acceptance Testing**: ⏳ Not Started

### Pass Rate
- **Total Tests**: 150+
- **Passed**: 0
- **Failed**: 0
- **Pending**: 150+
- **Pass Rate**: 0%

### Issues Found
- **Critical**: 0
- **High**: 0
- **Medium**: 0
- **Low**: 0

### Blockers
- None identified

---

## Sign-Off

| Role | Name | Date | Status |
|------|------|------|--------|
| QA Lead | - | - | ⏳ Pending |
| Dev Lead | - | - | ⏳ Pending |
| PM | - | - | ⏳ Pending |

---

## Notes

- All tests should be executed in order
- Each test should be documented with pass/fail status
- Any failures should be investigated and documented
- Performance targets must be met for performance tests
- All backward compatibility tests must pass
- Migration must be reversible
- UAT must be completed with actual users

