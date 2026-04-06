export interface Config {
  id: number
  userId: number
  title: string
  link: string
  iconSize: number
  lon: number
  lat: number
  zoom: number
  maxZoom: number
  minZoom: number
  tolorance: number
  autoCenter: boolean
  note: string
}

export interface ConfigUpdateForm {
  title: string
  link: string
  iconSize: number
  lon: number
  lat: number
  zoom: number
  maxZoom: number
  minZoom: number
  tolorance: number
  autoCenter: boolean
  note: string
}
