package usecase

import (
	"context"
	"time"

	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
)

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

	return stats, nil
}
