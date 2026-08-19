import { TableState } from '../api/types'

interface Props {
  tables: TableState[] | null
}

const STATE_COLOUR: Record<string, string> = {
  ready: '#4ade80',
  synchronized: '#4ade80',
  'finished copy': '#a3e635',
  'copying data': '#facc15',
  initializing: '#6b7280',
}

export function TableSyncGrid({ tables }: Props) {
  if (!tables || tables.length === 0) {
    return <div className="card"><h2>Table Sync</h2><p className="empty">No subscription data</p></div>
  }

  return (
    <div className="card">
      <h2>Table Sync</h2>
      <div className="table-grid">
        {tables.map((t) => (
          <div key={t.table_name} className="table-cell" title={t.state}>
            <span className="table-dot" style={{ background: STATE_COLOUR[t.state] ?? '#6b7280' }} />
            <span className="table-name">{t.table_name.replace('public.', '')}</span>
            <span className="table-state">{t.state}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
