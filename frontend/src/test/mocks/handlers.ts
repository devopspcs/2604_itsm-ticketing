import { http, HttpResponse } from 'msw'

export const handlers = [
  // Auth endpoints
  http.post('/api/auth/login', () => {
    return HttpResponse.json({
      access_token: 'test-access-token',
      refresh_token: 'test-refresh-token',
      user: {
        id: '123e4567-e89b-12d3-a456-426614174000',
        email: 'test@example.com',
        name: 'Test User',
        role: 'user',
      },
    })
  }),

  // Ticket endpoints
  http.get('/api/tickets', () => {
    return HttpResponse.json({
      data: [
        {
          id: '123e4567-e89b-12d3-a456-426614174001',
          title: 'Test Ticket 1',
          description: 'Test Description 1',
          status: 'open',
          type: 'incident',
          priority: 'high',
          category: 'network',
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ],
      pagination: {
        page: 1,
        page_size: 10,
        total: 1,
      },
    })
  }),

  http.get('/api/tickets/:id', () => {
    return HttpResponse.json({
      id: '123e4567-e89b-12d3-a456-426614174001',
      title: 'Test Ticket',
      description: 'Test Description',
      status: 'open',
      type: 'incident',
      priority: 'high',
      category: 'network',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    })
  }),

  // Dashboard endpoints
  http.get('/api/dashboard/stats', () => {
    return HttpResponse.json({
      total_tickets: 10,
      by_status: {
        open: 5,
        in_progress: 3,
        pending_approval: 1,
        approved: 1,
        done: 0,
      },
      by_type: {
        incident: 5,
        change_request: 3,
        helpdesk_request: 2,
      },
      by_priority: {
        critical: 1,
        high: 3,
        medium: 4,
        low: 2,
      },
      recent_tickets: [],
      sla_compliance_rate: 95.5,
      avg_resolution_hours: 4.2,
      on_time_count: 19,
      breached_count: 1,
    })
  }),
]
