import type { Place } from './place'
import type { Config } from './config'

export interface User {
  id: number
  username: string
  password?: string
  token?: string
  lastlogin?: Date
  admin: boolean
  places?: Place[]
  config?: Config
}

export interface LoginResult {
  token: string
  username: string
  admin: boolean
}

export interface LoginForm {
  username: string
  password: string
}
