import { useCallback, useEffect, useState } from 'react'

export type ThemeColor = 'red' | 'blue' | 'green' | 'yellow' | 'purple'

export const THEMES: { id: ThemeColor; label: string; color: string }[] = [
  { id: 'red', label: 'Merah', color: '#C41E3A' },
  { id: 'blue', label: 'Biru', color: '#00307d' },
  { id: 'green', label: 'Hijau', color: '#15803d' },
  { id: 'yellow', label: 'Kuning', color: '#a16207' },
  { id: 'purple', label: 'Ungu', color: '#7c3aed' },
]

const STORAGE_KEY = 'itsm-theme'

function getInitialTheme(): ThemeColor {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved && THEMES.some(t => t.id === saved)) return saved as ThemeColor
  return 'blue'
}

export function useTheme() {
  const [theme, setThemeState] = useState<ThemeColor>(getInitialTheme)

  const setTheme = useCallback((t: ThemeColor) => {
    setThemeState(t)
    localStorage.setItem(STORAGE_KEY, t)
    document.documentElement.setAttribute('data-theme', t)
  }, [])

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme)
  }, [theme])

  return { theme, setTheme, themes: THEMES }
}
