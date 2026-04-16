package usecase

import (
	"context"
	"time"

	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
)

var slaTargets = map[entity.Priority]time.Duration{
	entity.PriorityCritical: 4 * time.Hour,
	entity.PriorityHigh:     8 * time.Hour,
	entity.PriorityMedium:   24 * time.Hour,
	entity.PriorityLow:      72 * time.Hour,
}

type dashboardUseCase struct {
	ticketRepo repository.TicketRepository
}

func NewDashboardUseCase(ticketRepo repository.TicketRepository) domainUC.DashboardUseCase {
	return &dashboardUseCase{ticketRepo: ticketRepo}
}

func (uc *dashboardUseCase) GetStats(ctx context.Context, filter domainUC.DashboardFilter, requester domainUC.UserClaims) (*domainUC.DashboardStats, error) {
	baseFilter := repository.TicketFilter{Page: 1, PageSize: 100000}

	if requester.Role == entity.RoleUser {
		// For user: we'll use the ticket usecase logic via direct repo calls
		// Fetch created + assigned tickets
		createdFilter := repository.TicketFilter{Page: 1, PageSize: 100000, CreatedBy: &requester.UserID}
		assignedFilter := repository.TicketFilter{Page: 1, PageSize: 100000, AssignedTo: &requester.UserID}

		createdResult, err := uc.ticketRepo.List(ctx, createdFilter)
		if err != nil {
			return nil, err
		}
		assignedResult, err := uc.ticketRepo.List(ctx, assignedFilter)
		if err != nil {
			return nil, err
		}

		seen := map[string]bool{}
		var allTickets []*entity.Ticket
		for _, t := range append(createdResult.Tickets, assignedResult.Tickets...) {
			if !seen[t.ID.String()] {
				seen[t.ID.String()] = true
				allTickets = append(allTickets, t)
			}
		}

		stats := &domainUC.DashboardStats{
			TotalTickets:  int64(len(allTickets)),
			ByStatus:      make(map[entity.TicketStatus]int64),
			ByType:        make(map[entity.TicketType]int64),
			ByPriority:    make(map[entity.Priority]int64),
			RecentTickets: []*entity.Ticket{},
		}
		for _, t := range allTickets {
			stats.ByStatus[t.Status]++
			stats.ByType[t.Type]++
			stats.ByPriority[t.Priority]++
		}
		recent := allTickets
		if len(recent) > 10 {
			recent = recent[:10]
		}
		stats.RecentTickets = recent

		// Calculate SLA metrics for user path
		complianceRate, avgHours, onTime, breached := calculateSLAMetrics(allTickets)
		stats.SLAComplianceRate = complianceRate
		stats.AvgResolutionHours = avgHours
		stats.OnTimeCount = onTime
		stats.BreachedCount = breached

		return stats, nil
	}
	if filter.DateFrom != nil {
		t, err := time.Parse(time.RFC3339, *filter.DateFrom)
		if err == nil {
			baseFilter.DateFrom = &t
		}
	}
	if filter.DateTo != nil {
		t, err := time.Parse(time.RFC3339, *filter.DateTo)
		if err == nil {
			baseFilter.DateTo = &t
		}
	}

	result, err := uc.ticketRepo.List(ctx, baseFilter)
	if err != nil {
		return nil, err
	}

	stats := &domainUC.DashboardStats{
		TotalTickets:  result.Total,
		ByStatus:      make(map[entity.TicketStatus]int64),
		ByType:        make(map[entity.TicketType]int64),
		ByPriority:    make(map[entity.Priority]int64),
		RecentTickets: []*entity.Ticket{},
	}

	for _, t := range result.Tickets {
		stats.ByStatus[t.Status]++
		stats.ByType[t.Type]++
		stats.ByPriority[t.Priority]++
	}

	// Recent tickets: up to 10
	recent := result.Tickets
	if len(recent) > 10 {
		recent = recent[:10]
	}
	stats.RecentTickets = recent

	// Calculate SLA metrics for admin path
	complianceRate, avgHours, onTime, breached := calculateSLAMetrics(result.Tickets)
	stats.SLAComplianceRate = complianceRate
	stats.AvgResolutionHours = avgHours
	stats.OnTimeCount = onTime
	stats.BreachedCount = breached

	return stats, nil
}

func calculateSLAMetrics(tickets []*entity.Ticket) (complianceRate float64, avgHours float64, onTime int64, breached int64) {
	var totalResolutionHours float64
	var resolvedCount int64

	for _, t := range tickets {
		if t.Status != entity.StatusDone || t.ResolvedAt == nil {
			continue
		}
		duration := t.ResolvedAt.Sub(t.CreatedAt)
		totalResolutionHours += duration.Hours()
		resolvedCount++

		target, ok := slaTargets[t.Priority]
		if !ok {
			target = 72 * time.Hour // default to low priority
		}
		if duration <= target {
			onTime++
		} else {
			breached++
		}
	}

	if resolvedCount > 0 {
		avgHours = totalResolutionHours / float64(resolvedCount)
		complianceRate = float64(onTime) / float64(resolvedCount) * 100
	}

	return complianceRate, avgHours, onTime, breached
}
