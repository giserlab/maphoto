import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Place } from '@/types'
import * as placeApi from '@/api/place'

export const usePlacesStore = defineStore('places', () => {
  // State
  const places = ref<Place[]>([])
  const currentPlace = ref<Place | null>(null)
  const loading = ref(false)
  const groups = ref<string[]>([])

  // Getters
  const placesByGroup = computed(() => {
    const map = new Map<string, Place[]>()
    places.value.forEach((place) => {
      const group = place.group || '未分组'
      if (!map.has(group)) map.set(group, [])
      map.get(group)!.push(place)
    })
    return map
  })

  // Actions
  async function fetchPlaces() {
    loading.value = true
    try {
      const res = await placeApi.getAll()
      places.value = res.data || []
      // Extract all groups
      const groupSet = new Set(places.value.map((p) => p.group).filter(Boolean) as string[])
      groups.value = Array.from(groupSet)
    } finally {
      loading.value = false
    }
  }

  async function createPlace(data: Partial<Place>) {
    const res = await placeApi.create(data as Parameters<typeof placeApi.create>[0])
    places.value.push(res.data)
    return res.data
  }

  async function updatePlace(id: number, data: Partial<Place>) {
    const res = await placeApi.update(id, data as Parameters<typeof placeApi.update>[1])
    const index = places.value.findIndex((p) => p.id === id)
    if (index !== -1) {
      places.value[index] = { ...places.value[index], ...res.data }
    }
    return res.data
  }

  async function deletePlace(id: number) {
    await placeApi.remove(id)
    places.value = places.value.filter((p) => p.id !== id)
  }

  async function addPlacePhoto(id: number, url: string) {
    await placeApi.addPhoto(id, url)
    const place = places.value.find((p) => p.id === id)
    if (place) {
      place.photos.push({ id: Date.now(), placeId: id, url })
    }
  }

  async function removePlacePhoto(id: number, url: string) {
    await placeApi.removePhoto(id, url)
    const place = places.value.find((p) => p.id === id)
    if (place) {
      place.photos = place.photos.filter((p) => p.url !== url)
    }
  }

  async function updatePlaceCover(id: number, url: string) {
    await placeApi.updateCover(id, url)
    const place = places.value.find((p) => p.id === id)
    if (place) {
      place.cover = url
    }
  }

  function setCurrentPlace(place: Place | null) {
    currentPlace.value = place
  }

  return {
    places,
    currentPlace,
    loading,
    groups,
    placesByGroup,
    fetchPlaces,
    createPlace,
    updatePlace,
    deletePlace,
    addPlacePhoto,
    removePlacePhoto,
    updatePlaceCover,
    setCurrentPlace,
  }
})
