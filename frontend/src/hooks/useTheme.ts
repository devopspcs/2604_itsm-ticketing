import { useCallback, useEffect, useState } from 'react'

export type ThemeColor = 'red' | 'blue' | 'green' | 'yellow' | 'purple'
export type ThemeMode = 'light' | 'dark'

export const THEMES: { id: ThemeColor; label: string; color: string }[] = [
  { id: 'red', label: 'Merah', color: '#C41E3A' },
  { id: 'blue', label: 'Biru', color: '#00307d' },
  { id: 'green', label: 'Hijau', color: '#15803d' },
  { id: 'yellow', label: 'Kuning', color: '#a16207' },
  { id: 'purple', label: 'Ungu', color: '#7c3aed' },
]

const STORAGE_KEY = 'itsm-theme'
const MODE_STORAGE_KEY = 'itsm-theme-mode'

function getInitialTheme(): ThemeColor {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved && THEMES.some(t => t.id === saved)) return saved as ThemeColor
  return 'red'
}

function getInitialMode(): ThemeMode {
  const saved = localStorage.getItem(MODE_STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') return saved
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  return 'light'
}

function applyMode(m: ThemeMode) {
  const html = document.documentElement
  if (m === 'dark') {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
}

export function useTheme() {
  const [theme, setThemeState] = useState<ThemeColor>(getInitialTheme)
  const [mode, setModeState] = useState<ThemeMode>(getInitialMode)

  const setTheme = useCallback((t: ThemeColor) => {
    setThemeState(t)
    localStorage.setItem(STORAGE_KEY, t)
    document.documentElement.setAttribute('data-theme', t)
  }, [])

  const setMode = useCallback((m: ThemeMode) => {
    setModeState(m)
    localStorage.setItem(MODE_STORAGE_KEY, m)
    applyMode(m)
  }, [])

  const toggleMode = useCallback(() => {
    const newMode = mode === 'light' ? 'dark' : 'light'
    setModeState(newMode)
    localStorage.setItem(MODE_STORAGE_KEY, newMode)
    applyMode(newMode)
  }, [mode])

  // Apply on mount
  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
    applyMode(mode)
  }, [])

  return { theme, setTheme, mode, setMode, toggleMode, themes: THEMES }
}
