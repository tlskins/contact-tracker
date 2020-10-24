import { Moment } from "moment"

export interface Profile {
  id: string
  name: string
  email: string
  profileType: ProfileType
}

export enum ProfileType {
  User = "User",
  Place = "Place",
}
export interface User {
  id: string
  name: string
  email: string
}

export interface Place {
  id: string
  name: string
  email: string
}

export interface CheckIn {
  id: string
  in: Moment
  out: Moment
  place: Place
  user: User
}

export interface CheckInHistory {
  id: string
  in: Moment
  out: Moment
  place: Place
  user: User
  contacts: CheckIn[]
}

export interface Place {
  id: string
  name: string
}
