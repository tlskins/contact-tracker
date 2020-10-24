import { Profile, User, Place, CheckIn } from "./types"
import {
  ProfileActionTypes,
  UserActionTypes,
  PlaceActionTypes,
  CheckInActionTypes,
  SET_PROFILE,
  CLEAR_PROFILE,
  SET_USER,
  SET_PLACE,
  SET_CHECK_INS,
} from "./actions"

type ProfileState = Profile | null
type UserState = User | null
type PlaceState = Place | null
type CheckInsState = CheckIn[] | null

const initialProfileState: ProfileState = null
const initialUserState: UserState = null
const initialPlaceState: PlaceState = null
const initialCheckInsState: CheckInsState = []

export const profileReducer = (
  state: ProfileState = initialProfileState,
  action: ProfileActionTypes
): ProfileState => {
  switch (action.type) {
    case SET_PROFILE:
      return {
        ...state,
        ...action.payload,
      }
    case CLEAR_PROFILE:
      return null
    default:
      return state
  }
}

export const userReducer = (
  state: UserState = initialUserState,
  action: UserActionTypes
): UserState => {
  switch (action.type) {
    case SET_USER:
      return {
        ...state,
        ...action.payload,
      }
    default:
      return state
  }
}

export const placeReducer = (
  state: PlaceState = initialPlaceState,
  action: PlaceActionTypes
): PlaceState => {
  switch (action.type) {
    case SET_PLACE:
      return {
        ...state,
        ...action.payload,
      }
    default:
      return state
  }
}

export const checkInsReducer = (
  state: CheckInsState = initialCheckInsState,
  action: CheckInActionTypes
): CheckInsState => {
  switch (action.type) {
    case SET_CHECK_INS:
      return [...action.payload]
    default:
      return state
  }
}
