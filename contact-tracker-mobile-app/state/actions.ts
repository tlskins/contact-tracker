import { Profile, User, Place, CheckIn } from "./types"

export const SET_PROFILE = "SET_PROFILE"
export const CLEAR_PROFILE = "CLEAR_PROFILE"
export const SET_USER = "SET_USER"
export const SET_PLACE = "SET_PLACE"
export const SET_CHECK_INS = "SET_CHECK_INS"

interface SetProfileAction {
  type: typeof SET_PROFILE
  payload: Profile
}
interface ClearProfileAction {
  type: typeof CLEAR_PROFILE
}
interface SetUserAction {
  type: typeof SET_USER
  payload: User
}
interface SetPlaceAction {
  type: typeof SET_PLACE
  payload: Place
}
interface SetCheckInsAction {
  type: typeof SET_CHECK_INS
  payload: CheckIn[]
}

export type ProfileActionTypes = SetProfileAction | ClearProfileAction
export type UserActionTypes = SetUserAction
export type PlaceActionTypes = SetPlaceAction
export type CheckInActionTypes = SetCheckInsAction

export function setProfile(profile: Profile): ProfileActionTypes {
  return {
    type: SET_PROFILE,
    payload: profile,
  }
}
export function clearProfile(): ProfileActionTypes {
  return {
    type: CLEAR_PROFILE,
  }
}
export function setUser(user: User): UserActionTypes {
  return {
    type: SET_USER,
    payload: user,
  }
}
export function setPlace(place: Place): PlaceActionTypes {
  return {
    type: SET_PLACE,
    payload: place,
  }
}
export function setCheckIns(checkIns: CheckIn[]): CheckInActionTypes {
  return {
    type: SET_CHECK_INS,
    payload: checkIns,
  }
}
