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
  return 'blue'
}

function getInitialMode(): ThemeMode {
  const saved = localStorage.getItem(MODE_STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') return saved
  // Default: ikuti system preference
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  return 'light'
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
    if (m === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [])

  const toggleMode = useCallback(() => {
    setMode(mode === 'light' ? 'dark' : 'light')
  }, [mode, setMode])

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
    if (mode === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [theme, mode])

  return { theme, setTheme, mode, setMode, toggleMode, themes: THEMES }
}
