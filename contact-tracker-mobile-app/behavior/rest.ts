import axios from "axios"
import Constants from "expo-constants"

export let authToken = ""

export const setAuthToken = (value: string) => {
  authToken = value
}

// Set config defaults when creating the instance
export const instance = axios.create({
  baseURL: Constants.manifest.extra.apiHost,
  timeout: 10000,
  withCredentials: true,
  responseType: "json",
})

instance.interceptors.request.use(
  (config) => {
    // add jwt accessToken to auth header if present in localstorage
    // const accessToken = window.localStorage.getItem("authToken")
    console.log('authtoken:', authToken)
    if (authToken !== "") {
      // config.headers["Cookie"] = `accessToken=${accessToken};`
      config.headers["Authorization"] = `Bearer ${authToken}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// instance.interceptors.response.use((response) => {
//   // extract jwt from response data and store to localstorage
//   console.log("instance.interceptors.response", response)
//   if (response && response.data) {
//     const { authToken } = response.data
//     if (authToken) {
//       window.localStorage.setItem("authToken", authToken)
//       delete response.data.accessToken
//     }
//   }

//   return response
// })
