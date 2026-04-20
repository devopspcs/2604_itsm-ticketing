import '@testing-library/jest-dom'
import { afterEach, afterAll, beforeAll, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { server } from './mocks/server'

beforeAll(() => server.listen())
afterEach(() => {
  cleanup()
  server.resetHandlers()
})
afterAll(() => server.close())

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})
