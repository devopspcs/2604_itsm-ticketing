import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BrowserRouter } from 'react-router-dom'
import { Provider } from 'react-redux'
import { configureStore } from '@reduxjs/toolkit'
import { LoginPage } from '../LoginPage'
import authReducer from '../../store/authSlice'

// Mock useAuth hook
vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({
    login: vi.fn().mockResolvedValue(undefined),
  }),
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

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders login form with email and password fields', () => {
    renderWithProviders(<LoginPage />)
    
    expect(screen.getByLabelText(/email or username/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('allows user to enter email and password', async () => {
    const user = userEvent.setup()
    renderWithProviders(<LoginPage />)
    
    const emailInput = screen.getByLabelText(/email or username/i) as HTMLInputElement
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement
    
    await user.type(emailInput, 'test@example.com')
    await user.type(passwordInput, 'password123')
    
    expect(emailInput.value).toBe('test@example.com')
    expect(passwordInput.value).toBe('password123')
  })

  it('has SSO login link', () => {
    renderWithProviders(<LoginPage />)
    
    const ssoLink = screen.getByRole('link', { name: /login with keycloak sso/i })
    expect(ssoLink).toBeInTheDocument()
    expect(ssoLink).toHaveAttribute('href', '/api/v1/auth/sso/redirect')
  })

  it('displays form title and description', () => {
    renderWithProviders(<LoginPage />)
    
    expect(screen.getByText(/welcome back/i)).toBeInTheDocument()
    expect(screen.getByText(/please enter your credentials/i)).toBeInTheDocument()
  })

  it('displays PCS Ticketing System header', () => {
    renderWithProviders(<LoginPage />)
    
    expect(screen.getByText(/pcs ticketing system/i)).toBeInTheDocument()
    expect(screen.getByText(/it service management/i)).toBeInTheDocument()
  })
})
