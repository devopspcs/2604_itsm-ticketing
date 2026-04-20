import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { configureStore } from '@reduxjs/toolkit'
import { DashboardPage } from '../DashboardPage'
import authReducer from '../../store/authSlice'

// Mock dashboard service
vi.mock('../../services/dashboard.service', () => ({
  dashboardService: {
    getStats: vi.fn().mockResolvedValue({
      data: {
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
      },
    }),
  },
}))

const createMockStore = () => {
  return configureStore({
    reducer: {
      auth: authReducer,
    },
  })
}

const renderWithProviders = (component: React.ReactElement) => {
  const store = createMockStore()
  return render(
    <Provider store={store}>
      <BrowserRouter>
        {component}
      </BrowserRouter>
    </Provider>
  )
}

describe('DashboardPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders loading state initially', () => {
    renderWithProviders(<DashboardPage />)
    
    expect(screen.getByText(/loading dashboard/i)).toBeInTheDocument()
  })

  it('displays dashboard stats after loading', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      expect(screen.getByText(/total tickets/i)).toBeInTheDocument()
    })
    
    expect(screen.getByText(/pending approvals/i)).toBeInTheDocument()
    expect(screen.getByText(/active incidents/i)).toBeInTheDocument()
    expect(screen.getByText(/completed/i)).toBeInTheDocument()
  })

  it('displays SLA compliance metrics', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      expect(screen.getByText(/sla compliance/i)).toBeInTheDocument()
    })
    
    expect(screen.getAllByText(/compliance/i).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/on time/i).length).toBeGreaterThan(0)
    expect(screen.getByText(/breached/i)).toBeInTheDocument()
  })

  it('displays average resolution hours', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      expect(screen.getByText(/average resolution time/i)).toBeInTheDocument()
    })
    
    expect(screen.getByText(/hours/i)).toBeInTheDocument()
  })

  it('displays recent tickets section', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      expect(screen.getByText(/recent tickets/i)).toBeInTheDocument()
    })
  })

  it('displays create new ticket button', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      const buttons = screen.getAllByRole('link', { name: /create new ticket/i })
      expect(buttons.length).toBeGreaterThan(0)
    })
  })

  it('displays chart with ticket categories', async () => {
    renderWithProviders(<DashboardPage />)
    
    await waitFor(() => {
      expect(screen.getByText(/tickets per category/i)).toBeInTheDocument()
    })
    
    expect(screen.getByText(/helpdesk/i)).toBeInTheDocument()
    expect(screen.getAllByText(/incident/i).length).toBeGreaterThan(0)
    expect(screen.getByText(/change/i)).toBeInTheDocument()
  })
})
