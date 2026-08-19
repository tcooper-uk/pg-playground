export interface Lag {
  application_name: string
  state: string
  lag_bytes: number | null
  write_lag_ns: number | null
  flush_lag_ns: number | null
  replay_lag_ns: number | null
}

export interface TableState {
  table_name: string
  state: 'initializing' | 'copying data' | 'finished copy' | 'synchronized' | 'ready'
}

export interface SlotHealth {
  slot_name: string
  active: boolean
  retained_bytes: number | null
}

export interface SenderInfo {
  application_name: string
  state: string
  sent_lsn: string
  replay_lsn: string
}

export interface ReceiverInfo {
  sub_name: string
  received_lsn: string
  latest_end_lsn: string
  latest_end_time: string
}

export interface WorkerInfo {
  senders: SenderInfo[] | null
  receivers: ReceiverInfo[] | null
}

export interface SimulatorSnapshot {
  rental: number
  return: number
  churn: number
  read: number
  errors: number
}

export interface WeightConfig {
  rental: number
  return: number
  customer_churn: number
  read: number
}

export interface SubscriptionStatus {
  name: string
  enabled: boolean
}

export interface Stats {
  lag: Lag[] | null
  tables: TableState[] | null
  slots: SlotHealth[] | null
  workers: WorkerInfo
  simulator: SimulatorSnapshot
  weights: WeightConfig
  rate: number
  subscriptions: SubscriptionStatus[] | null
}
