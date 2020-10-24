import { profileReducer, userReducer, checkInsReducer } from "./reducers"
import { combineReducers } from "redux"

export const rootReducer = combineReducers({
  checkIns: checkInsReducer,
  profile: profileReducer,
  user: userReducer,
})

export type RootState = ReturnType<typeof rootReducer>
