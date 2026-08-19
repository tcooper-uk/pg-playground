import { useState, useEffect } from 'react'
import { WeightConfig } from '../api/types'
import { setWeights, setRate } from '../api/client'

interface Props {
  currentWeights: WeightConfig | null
  currentRate: number | null
}

const SLIDERS: { key: keyof WeightConfig; label: string; colour: string }[] = [
  { key: 'rental',         label: 'Rentals', colour: '#4ade80' },
  { key: 'return',         label: 'Returns', colour: '#60a5fa' },
  { key: 'customer_churn', label: 'Churn',   colour: '#c084fc' },
  { key: 'read',           label: 'Reads',   colour: '#facc15' },
]

export function WeightsControl({ currentWeights, currentRate }: Props) {
  const [weights, setLocalWeights] = useState<WeightConfig>({ rental: 40, return: 30, customer_churn: 10, read: 20 })
  const [rate, setLocalRate] = useState(10)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => { if (currentWeights) setLocalWeights(currentWeights) }, [currentWeights == null])
  useEffect(() => { if (currentRate != null) setLocalRate(currentRate) }, [currentRate == null])

  const total = Object.values(weights).reduce((a, b) => a + b, 0)

  const handleWeightChange = (key: keyof WeightConfig, value: number) => {
    setLocalWeights((prev) => ({ ...prev, [key]: value }))
    setSaved(false)
  }

  const handleApply = async () => {
    setSaving(true)
    try {
      await Promise.all([setWeights(weights), setRate(rate)])
      setSaved(true)
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="card">
      <h2>Simulator Controls</h2>

      <h3>Rate</h3>
      <div className="weight-row">
        <span className="weight-label" style={{ color: '#e0e0e0' }}>Ops / sec</span>
        <input
          type="range"
          min={1}
          max={500}
          value={rate}
          onChange={(e) => { setLocalRate(Number(e.target.value)); setSaved(false) }}
          style={{ '--thumb-colour': '#e0e0e0' } as React.CSSProperties}
        />
        <span className="weight-value">{rate}</span>
        <span className="weight-pct" />
      </div>

      <h3>Weights</h3>
      <div className="weights-sliders">
        {SLIDERS.map(({ key, label, colour }) => {
          const pct = total > 0 ? Math.round((weights[key] / total) * 100) : 0
          return (
            <div key={key} className="weight-row">
              <span className="weight-label" style={{ color: colour }}>{label}</span>
              <input
                type="range"
                min={0}
                max={100}
                value={weights[key]}
                onChange={(e) => handleWeightChange(key, Number(e.target.value))}
                style={{ '--thumb-colour': colour } as React.CSSProperties}
              />
              <span className="weight-value">{weights[key]}</span>
              <span className="weight-pct">{pct}%</span>
            </div>
          )
        })}
      </div>

      <div className="weights-footer">
        <span className="weight-total">Total weight: {total}</span>
        <button
          className={`apply-btn ${saved ? 'saved' : ''}`}
          onClick={handleApply}
          disabled={saving || total === 0}
        >
          {saving ? 'Applying…' : saved ? 'Applied ✓' : 'Apply'}
        </button>
      </div>
    </div>
  )
}
