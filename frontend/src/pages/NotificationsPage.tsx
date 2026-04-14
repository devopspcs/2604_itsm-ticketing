import { useNotifications } from '../hooks/useNotifications'
import { Link } from 'react-router-dom'

export function NotificationsPage() {
  const { notifications, unreadCount, markAsRead } = useNotifications()

  return (
    <div className="max-w-4xl mx-auto p-8">
      <div className="flex items-end justify-between mb-10">
        <div>
          <h1 className="text-4xl font-extrabold tracking-tight text-on-surface mb-2">Notifications</h1>
          <p className="text-on-surface-variant">
            {unreadCount > 0
              ? <><span className="font-bold text-primary">{unreadCount} unread</span> notifications require your attention.</>
              : 'All caught up — no unread notifications.'}
          </p>
        </div>
        {unreadCount > 0 && (
          <div className="bg-surface-container-low px-4 py-2 rounded-full border border-outline-variant/10 flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-amber-500" />
            <span className="text-xs font-bold text-on-surface-variant uppercase tracking-tighter">{unreadCount} Unread</span>
          </div>
        )}
      </div>

      {notifications.length === 0 ? (
        <div className="bg-surface-container-lowest rounded-xl p-16 text-center shadow-sm">
          <span className="material-symbols-outlined text-5xl text-on-surface-variant/30 block mb-3">notifications_none</span>
          <p className="text-on-surface-variant font-medium">No notifications yet</p>
        </div>
      ) : (
        <div className="bg-surface-container-lowest rounded-xl overflow-hidden shadow-sm border border-outline-variant/10">
          <div className="divide-y divide-surface-container">
            {notifications.map((n) => (
              <div
                key={n.id}
                className={`flex items-start gap-4 px-6 py-5 hover:bg-slate-50 transition-colors ${!n.is_read ? 'bg-blue-50/30' : ''}`}
              >
                <div className={`flex-shrink-0 w-10 h-10 rounded-xl flex items-center justify-center ${
                  !n.is_read ? 'bg-primary-fixed text-primary' : 'bg-surface-container-high text-on-surface-variant'
                }`}>
                  <span className="material-symbols-outlined text-lg">notifications</span>
                </div>
                <div className="flex-grow min-w-0">
                  <p className={`text-sm ${!n.is_read ? 'font-bold text-on-surface' : 'font-medium text-on-surface-variant'}`}>
                    {n.message}
                  </p>
                  <div className="flex items-center gap-3 mt-1">
                    <span className="text-[10px] text-on-surface-variant">{new Date(n.created_at).toLocaleString()}</span>
                    <Link to={`/tickets/${n.ticket_id}`} className="text-[10px] font-bold text-primary hover:underline">
                      View Ticket →
                    </Link>
                  </div>
                </div>
                {!n.is_read && (
                  <button
                    onClick={() => markAsRead(n.id)}
                    className="flex-shrink-0 px-3 py-1.5 text-xs font-bold bg-primary text-white rounded-xl hover:opacity-90 transition-opacity"
                  >
                    Mark Read
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
