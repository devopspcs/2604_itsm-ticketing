interface Props {
  page: number
  total: number
  pageSize: number
  onPageChange: (page: number) => void
}

export function Pagination({ page, total, pageSize, onPageChange }: Props) {
  const totalPages = Math.ceil(total / pageSize)
  if (totalPages <= 1) return null

  const pages = Array.from({ length: Math.min(totalPages, 5) }, (_, i) => i + 1)

  return (
    <div className="flex items-center gap-2">
      <button
        disabled={page <= 1}
        onClick={() => onPageChange(page - 1)}
        className="p-1 rounded-lg hover:bg-surface-container-highest transition-colors disabled:opacity-30"
      >
        <span className="material-symbols-outlined">chevron_left</span>
      </button>
      <div className="flex gap-1">
        {pages.map((p) => (
          <button
            key={p}
            onClick={() => onPageChange(p)}
            className={`w-8 h-8 rounded-lg text-xs font-medium transition-colors ${
              p === page ? 'bg-primary text-white font-bold' : 'hover:bg-surface-container-highest'
            }`}
          >
            {p}
          </button>
        ))}
        {totalPages > 5 && <span className="px-1 text-outline">...</span>}
      </div>
      <button
        disabled={page >= totalPages}
        onClick={() => onPageChange(page + 1)}
        className="p-1 rounded-lg hover:bg-surface-container-highest transition-colors disabled:opacity-30"
      >
        <span className="material-symbols-outlined">chevron_right</span>
      </button>
    </div>
  )
}
