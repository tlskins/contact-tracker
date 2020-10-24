import { createStore, applyMiddleware } from "redux"
import { persistStore } from "redux-persist"
import thunk from "redux-thunk"
// import reduxLogger from "redux-logger"
import { rootReducer } from "../state"

// const logger = __DEV__ ? [reduxLogger] : []
const logger = []

const middlewares = applyMiddleware(thunk, ...logger)

const initialState = {}

const store = createStore(rootReducer, initialState, middlewares)
const persistor = persistStore(store)

export { store, persistor }
