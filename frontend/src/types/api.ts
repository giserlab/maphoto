import type { Config } from './config'

export interface ApiResponse<T = unknown> {
  status: boolean
  code: number
  message: string
  data: T
}

export interface ShareResult {
  config: Config
  features: GeoJSON.FeatureCollection
}
