export interface Photo {
  id: number
  placeId: number
  url: string
}

export interface Place {
  id: number
  userid: number
  name?: string
  desc: string
  cover: string
  private: boolean
  group: string
  date?: Date
  photos: Photo[]
  lon: number
  lat: number
}

export interface PlaceAddForm {
  name: string
  desc: string
  lon: number
  lat: number
  cover: string
  group: string
  photos: string[]
}

export interface PlaceUpdateForm {
  name: string
  desc: string
  lon: number
  lat: number
  group: string
}
