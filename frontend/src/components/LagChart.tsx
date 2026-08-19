import { LagSample } from '../hooks/useLagHistory'
import {
  LineChart, Line, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer, CartesianGrid,
} from 'recharts'

interface Props {
  history: LagSample[]
}

export function LagChart({ history }: Props) {
  if (history.length === 0) {
    return <div className="card"><h2>Replication Lag</h2><p className="empty">No WAL sender connected</p></div>
  }

  const data = history.map((s) => ({
    t: new Date(s.t).toLocaleTimeString(),
    lagBytes: s.lagBytes,
    replayMs: s.replayLagMs,
  }))

  return (
    <div className="card">
      <h2>Replication Lag</h2>
      <ResponsiveContainer width="100%" height={220}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#333" />
          <XAxis dataKey="t" tick={{ fill: '#aaa', fontSize: 11 }} interval="preserveStartEnd" />
          <YAxis yAxisId="bytes" tick={{ fill: '#aaa', fontSize: 11 }} />
          <YAxis yAxisId="ms" orientation="right" tick={{ fill: '#aaa', fontSize: 11 }} />
          <Tooltip contentStyle={{ background: '#1e1e1e', border: '1px solid #444' }} />
          <Legend />
          <Line yAxisId="bytes" type="monotone" dataKey="lagBytes" stroke="#4ade80" dot={false} name="Lag (bytes)" isAnimationActive={false} />
          <Line yAxisId="ms" type="monotone" dataKey="replayMs" stroke="#60a5fa" dot={false} name="Replay lag (ms)" isAnimationActive={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
