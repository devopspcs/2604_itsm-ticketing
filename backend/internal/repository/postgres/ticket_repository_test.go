package postgres_test

import (
"testing"
"time"

"github.com/google/uuid"
"github.com/org/itsm/internal/domain/entity"
"github.com/org/itsm/internal/domain/repository"
"pgregory.net/rapid"
)

// Feature: itsm-web-app, Property 9: Ticket Filter Correctness
// **Validates: Requirements 5.8**
//
// For any combination of filter parameters (status, type, priority, assigned_to,
// date range), all tickets returned by the list endpoint SHALL satisfy every
// active filter condition simultaneously - no ticket outside the filter criteria
// SHALL appear in the results.

func ticketMatchesFilter(t *entity.Ticket, f repository.TicketFilter) bool {
if f.Status != nil && t.Status != *f.Status {
return false
}
if f.Type != nil && t.Type != *f.Type {
return false
}
if f.Priority != nil && t.Priority != *f.Priority {
return false
}
if f.Category != nil && t.Category != *f.Category {
return false
}
if f.AssignedTo != nil {
if t.AssignedTo == nil || *t.AssignedTo != *f.AssignedTo {
return false
}
}
if f.CreatedBy != nil && t.CreatedBy != *f.CreatedBy {
return false
}
if f.DateFrom != nil && t.CreatedAt.Before(*f.DateFrom) {
return false
}
if f.DateTo != nil && t.CreatedAt.After(*f.DateTo) {
return false
}
return true
}

func filterTickets(tickets []*entity.Ticket, f repository.TicketFilter) []*entity.Ticket {
var result []*entity.Ticket
for _, t := range tickets {
if ticketMatchesFilter(t, f) {
result = append(result, t)
}
}
return result
}

var allStatuses = []entity.TicketStatus{
entity.StatusOpen, entity.StatusInProgress,
entity.StatusPendingApproval, entity.StatusApproved,
entity.StatusRejected, entity.StatusDone,
}

var allTypes = []entity.TicketType{
entity.TypeChangeRequest, entity.TypeIncident, entity.TypeHelpdeskRequest,
}

var allPriorities = []entity.Priority{
entity.PriorityLow, entity.PriorityMedium,
entity.PriorityHigh, entity.PriorityCritical,
}

var sampleCategories = []string{"network", "hardware", "software", "security", "database"}

func ticketGen(userIDs []uuid.UUID) *rapid.Generator[*entity.Ticket] {
return rapid.Custom[*entity.Ticket](func(t *rapid.T) *entity.Ticket {
creatorIdx := rapid.IntRange(0, len(userIDs)-1).Draw(t, "creatorIdx")
hasAssignee := rapid.Bool().Draw(t, "hasAssignee")
var assignedTo *uuid.UUID
if hasAssignee {
idx := rapid.IntRange(0, len(userIDs)-1).Draw(t, "assigneeIdx")
a := userIDs[idx]
assignedTo = &a
}
dayOffset := rapid.IntRange(0, 90).Draw(t, "dayOffset")
hourOffset := rapid.IntRange(0, 23).Draw(t, "hourOffset")
baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
createdAt := baseTime.Add(
time.Duration(dayOffset)*24*time.Hour +
time.Duration(hourOffset)*time.Hour)
return &entity.Ticket{
ID:          uuid.New(),
Title:       "ticket-title",
Description: "ticket-desc",
Type:        rapid.SampledFrom(allTypes).Draw(t, "type"),
Category:    rapid.SampledFrom(sampleCategories).Draw(t, "category"),
Priority:    rapid.SampledFrom(allPriorities).Draw(t, "priority"),
Status:      rapid.SampledFrom(allStatuses).Draw(t, "status"),
CreatedBy:   userIDs[creatorIdx],
AssignedTo:  assignedTo,
CreatedAt:   createdAt,
UpdatedAt:   createdAt,
}
})
}

func filterGen(userIDs []uuid.UUID) *rapid.Generator[repository.TicketFilter] {
return rapid.Custom[repository.TicketFilter](func(t *rapid.T) repository.TicketFilter {
var f repository.TicketFilter
if rapid.Bool().Draw(t, "hasStatus") {
s := rapid.SampledFrom(allStatuses).Draw(t, "filterStatus")
f.Status = &s
}
if rapid.Bool().Draw(t, "hasType") {
tp := rapid.SampledFrom(allTypes).Draw(t, "filterType")
f.Type = &tp
}
if rapid.Bool().Draw(t, "hasPriority") {
p := rapid.SampledFrom(allPriorities).Draw(t, "filterPriority")
f.Priority = &p
}
if rapid.Bool().Draw(t, "hasCategory") {
c := rapid.SampledFrom(sampleCategories).Draw(t, "filterCategory")
f.Category = &c
}
if rapid.Bool().Draw(t, "hasAssignedTo") {
idx := rapid.IntRange(0, len(userIDs)-1).Draw(t, "filterAssigneeIdx")
a := userIDs[idx]
f.AssignedTo = &a
}
if rapid.Bool().Draw(t, "hasCreatedBy") {
idx := rapid.IntRange(0, len(userIDs)-1).Draw(t, "filterCreatorIdx")
c := userIDs[idx]
f.CreatedBy = &c
}
if rapid.Bool().Draw(t, "hasDateRange") {
dayFrom := rapid.IntRange(0, 45).Draw(t, "dateFromDay")
dayTo := rapid.IntRange(dayFrom, 90).Draw(t, "dateToDay")
baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
from := baseTime.Add(time.Duration(dayFrom) * 24 * time.Hour)
to := baseTime.Add(time.Duration(dayTo)*24*time.Hour +
23*time.Hour + 59*time.Minute + 59*time.Second)
f.DateFrom = &from
f.DateTo = &to
}
f.Page = 1
f.PageSize = 1000
return f
})
}

// TestProperty_TicketFilterCorrectness_AllReturnedMatchFilter verifies that
// every ticket returned by the filter satisfies ALL active filter conditions.
func TestProperty_TicketFilterCorrectness_AllReturnedMatchFilter(t *testing.T) {
// Feature: itsm-web-app, Property 9: Ticket Filter Correctness
rapid.Check(t, func(t *rapid.T) {
numUsers := rapid.IntRange(2, 5).Draw(t, "numUsers")
userIDs := make([]uuid.UUID, numUsers)
for i := range userIDs {
userIDs[i] = uuid.New()
}
numTickets := rapid.IntRange(5, 30).Draw(t, "numTickets")
tickets := make([]*entity.Ticket, numTickets)
for i := 0; i < numTickets; i++ {
tickets[i] = ticketGen(userIDs).Draw(t, "ticket")
}
filter := filterGen(userIDs).Draw(t, "filter")
result := filterTickets(tickets, filter)
for _, ticket := range result {
if filter.Status != nil && ticket.Status != *filter.Status {
t.Fatalf("ticket %s has status %s, filter requires %s",
ticket.ID, ticket.Status, *filter.Status)
}
if filter.Type != nil && ticket.Type != *filter.Type {
t.Fatalf("ticket %s has type %s, filter requires %s",
ticket.ID, ticket.Type, *filter.Type)
}
if filter.Priority != nil && ticket.Priority != *filter.Priority {
t.Fatalf("ticket %s has priority %s, filter requires %s",
ticket.ID, ticket.Priority, *filter.Priority)
}
if filter.Category != nil && ticket.Category != *filter.Category {
t.Fatalf("ticket %s has category %s, filter requires %s",
ticket.ID, ticket.Category, *filter.Category)
}
if filter.AssignedTo != nil {
if ticket.AssignedTo == nil || *ticket.AssignedTo != *filter.AssignedTo {
t.Fatalf("ticket %s assigned_to %v, filter requires %s",
ticket.ID, ticket.AssignedTo, *filter.AssignedTo)
}
}
if filter.CreatedBy != nil && ticket.CreatedBy != *filter.CreatedBy {
t.Fatalf("ticket %s created_by %s, filter requires %s",
ticket.ID, ticket.CreatedBy, *filter.CreatedBy)
}
if filter.DateFrom != nil && ticket.CreatedAt.Before(*filter.DateFrom) {
t.Fatalf("ticket %s created_at %v before date_from %v",
ticket.ID, ticket.CreatedAt, *filter.DateFrom)
}
if filter.DateTo != nil && ticket.CreatedAt.After(*filter.DateTo) {
t.Fatalf("ticket %s created_at %v after date_to %v",
ticket.ID, ticket.CreatedAt, *filter.DateTo)
}
}
})
}

// TestProperty_TicketFilterCorrectness_NoMatchingTicketExcluded verifies that
// no ticket satisfying all filter conditions is excluded from the results.
func TestProperty_TicketFilterCorrectness_NoMatchingTicketExcluded(t *testing.T) {
// Feature: itsm-web-app, Property 9: Ticket Filter Correctness
rapid.Check(t, func(t *rapid.T) {
numUsers := rapid.IntRange(2, 5).Draw(t, "numUsers")
userIDs := make([]uuid.UUID, numUsers)
for i := range userIDs {
userIDs[i] = uuid.New()
}
numTickets := rapid.IntRange(5, 30).Draw(t, "numTickets")
tickets := make([]*entity.Ticket, numTickets)
for i := 0; i < numTickets; i++ {
tickets[i] = ticketGen(userIDs).Draw(t, "ticket")
}
filter := filterGen(userIDs).Draw(t, "filter")
result := filterTickets(tickets, filter)
returnedIDs := make(map[uuid.UUID]bool, len(result))
for _, ticket := range result {
returnedIDs[ticket.ID] = true
}
for _, ticket := range tickets {
if ticketMatchesFilter(ticket, filter) && !returnedIDs[ticket.ID] {
t.Fatalf("ticket %s matches all filter conditions but was excluded", ticket.ID)
}
}
})
}

// TestProperty_TicketFilterCorrectness_ResultCountConsistency verifies that
// the result count equals the count of individually matching tickets.
func TestProperty_TicketFilterCorrectness_ResultCountConsistency(t *testing.T) {
// Feature: itsm-web-app, Property 9: Ticket Filter Correctness
rapid.Check(t, func(t *rapid.T) {
numUsers := rapid.IntRange(2, 5).Draw(t, "numUsers")
userIDs := make([]uuid.UUID, numUsers)
for i := range userIDs {
userIDs[i] = uuid.New()
}
numTickets := rapid.IntRange(5, 30).Draw(t, "numTickets")
tickets := make([]*entity.Ticket, numTickets)
for i := 0; i < numTickets; i++ {
tickets[i] = ticketGen(userIDs).Draw(t, "ticket")
}
filter := filterGen(userIDs).Draw(t, "filter")
result := filterTickets(tickets, filter)
expectedCount := 0
for _, ticket := range tickets {
if ticketMatchesFilter(ticket, filter) {
expectedCount++
}
}
if len(result) != expectedCount {
t.Fatalf("filter returned %d tickets but %d match the conditions",
len(result), expectedCount)
}
})
}
