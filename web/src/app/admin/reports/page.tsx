'use client'

const mockReports = [
  {
    id: 'r1',
    targetType: 'Video',
    targetName: 'How to hack banking apps',
    reportedBy: 'user_good123',
    reason: 'Dangerous acts',
    date: '2 hours ago',
    description: 'This video promotes illegal hacking tools and shows sensitive information.',
  },
  {
    id: 'r2',
    targetType: 'User Profile',
    targetName: 'toxic_guy',
    reportedBy: 'anime_fan99',
    reason: 'Harassment',
    date: '5 hours ago',
    description: 'This user is constantly insulting people in the live chat of multiple streams.',
  },
  {
    id: 'r3',
    targetType: 'Video',
    targetName: 'Free Robux Generator 2026',
    reportedBy: 'system_automod',
    reason: 'Spam / Scam',
    date: '1 day ago',
    description: 'Automated flag: Suspicious link in description matching known phishing database.',
  },
]

export default function AdminReportsPage() {
  return (
    <div className="flex flex-col gap-4 animate-in fade-in duration-500 max-w-4xl">
      {/* Tabs / Filters for reports */}
      <div className="flex items-center gap-3 mb-2">
        <button className="px-4 py-2 bg-neutral-100 text-foreground-inverse-subtle font-bold rounded-xl text-sm">
          Open (3)
        </button>
        <button
          className="px-4 py-2 bg-[#0A0A0A] hover:bg-neutral-800 text-foreground-muted border border-neutral-800/60
            rounded-xl text-sm font-medium transition-colors"
        >
          Resolved
        </button>
      </div>

      {/* Reports Feed */}
      {mockReports.map(report => (
        <div key={report.id} className="bg-[#0A0A0A] border border-neutral-800/60 rounded-2xl p-5 flex flex-col gap-4">
          {/* Header */}
          <div className="flex items-start justify-between">
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-2">
                <span
                  className="px-2 py-0.5 bg-neutral-800 text-foreground-subtle rounded text-xs font-bold uppercase
                    tracking-wider"
                >
                  {report.targetType}
                </span>
                <span className="text-red-400 text-sm font-bold">{report.reason}</span>
                <span className="text-foreground-disabled text-xs ml-2">{report.date}</span>
              </div>
              <h4 className="text-lg font-bold text-foreground-secondary mt-1">Target: {report.targetName}</h4>
            </div>
            <div className="text-sm text-foreground0">
              Reported by:{' '}
              <span className="text-emerald-500 font-medium cursor-pointer hover:underline">@{report.reportedBy}</span>
            </div>
          </div>

          {/* Body */}
          <div className="bg-neutral-900/50 rounded-xl p-4 text-foreground-subtle text-sm border-l-4 border-red-500/50">
            &quot;{report.description}&quot;
          </div>

          {/* Actions */}
          <div className="flex items-center gap-3 mt-2 border-t border-neutral-800/60 pt-4">
            {report.targetType === 'Video' && (
              <button
                className="px-5 py-2 bg-red-500 hover:bg-red-600 text-foreground-strong font-semibold rounded-xl text-sm
                  transition-colors shadow-lg shadow-red-500/20"
              >
                Delete Content
              </button>
            )}
            {report.targetType === 'User Profile' && (
              <button
                className="px-5 py-2 bg-red-500 hover:bg-red-600 text-foreground-strong font-semibold rounded-xl text-sm
                  transition-colors shadow-lg shadow-red-500/20"
              >
                Ban User
              </button>
            )}
            <button
              className="px-5 py-2 bg-neutral-800 hover:bg-neutral-700 text-foreground-subtle font-semibold rounded-xl
                text-sm transition-colors"
            >
              Dismiss Report
            </button>
            <button
              className="px-5 py-2 text-foreground0 hover:text-foreground-subtle font-medium rounded-xl text-sm
                transition-colors ml-auto"
            >
              View Details
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
