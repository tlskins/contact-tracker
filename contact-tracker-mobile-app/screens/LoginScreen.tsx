import * as React from "react"
import { Alert, TextInput, TouchableOpacity, StyleSheet } from "react-native"
import tailwind from "tailwind-rn"
import { connect, ConnectedProps } from "react-redux"
import { Picker } from "@react-native-community/picker"

import { Text, View } from "../components/Themed"
import { useState } from "react"
import { instance, setAuthToken } from "../behavior/rest"
import { setProfile, setUser, setPlace } from "../state/actions"
import { RootState } from "../state"
import { ProfileType } from "../state/types"

interface newAccountRequest {
  email: string
  password: string
  name: string
}
interface loginRequest {
  email: string
  password: string
}

const connector = connect(
  (state: RootState) => ({ profile: state.profile }),
  (dispatch: any) => {
    return {
      createUser: async ({
        email,
        password,
        name,
      }: newAccountRequest): Promise<boolean> => {
        try {
          console.log("before new user...")
          const resp = await instance.post("/users", { email, password, name })
          console.log("create resp", resp)
          if (resp?.data) {
            const { data } = resp
            setAuthToken(data.authToken)
            dispatch(
              setProfile({
                ...data,
                profileType: ProfileType.User,
              })
            )
            dispatch(setUser(data))
            return true
          }
        } catch (e) {
          console.log("Register error", JSON.parse(JSON.stringify(e)))
          const errMsg = e.response?.data?.message || e.message || ""
          Alert.alert("Register Error", errMsg, [{ text: "OK" }], {
            cancelable: false,
          })
        }
        return false
      },
      userLogin: async ({
        email,
        password,
      }: loginRequest): Promise<boolean> => {
        try {
          console.log("before user login...")
          const resp = await instance.post("/users/login", { email, password })
          console.log("login resp", resp)
          if (resp?.data) {
            const { data } = resp
            setAuthToken(data.authToken)
            dispatch(
              setProfile({
                ...data,
                profileType: ProfileType.User,
              })
            )
            dispatch(setUser(data))
            return true
          }
        } catch (e) {
          console.log("login error", e)
          const errMsg = e.response?.data?.message || e.message || ""
          Alert.alert("Login Error", errMsg, [{ text: "OK" }], {
            cancelable: false,
          })
        }
        return false
      },
      placeLogin: async ({
        email,
        password,
      }: loginRequest): Promise<boolean> => {
        try {
          console.log("before place login...")
          const resp = await instance.post("/places/login", { email, password })
          console.log("place login resp", resp)
          if (resp?.data) {
            const { data } = resp
            setAuthToken(data.authToken)
            dispatch(
              setProfile({ ...data, profileType: ProfileType.Place })
            )
            dispatch(setPlace(data))
            return true
          }
        } catch (e) {
          const errMsg = e.response?.data?.message || e.message || ""
          Alert.alert("Login Error", errMsg, [{ text: "OK" }], {
            cancelable: false,
          })
        }
        return false
      },
      createPlace: async ({
        email,
        password,
        name,
      }: newAccountRequest): Promise<boolean> => {
        try {
          console.log("before place create...")
          const resp = await instance.post("/places", { email, password, name })
          console.log("place create resp", resp)
          if (resp?.data) {
            const { data } = resp
            dispatch(
              setProfile({
                ...data,
                profileType: ProfileType.Place,
              })
            )
            dispatch(setPlace(data))
            return true
          }
        } catch (e) {
          const errMsg = e.response?.data?.message || e.message || ""
          Alert.alert("Login Error", errMsg, [{ text: "OK" }], {
            cancelable: false,
          })
        }
        return false
      },
    }
  }
)

const LoginScreen = (props: ConnectedProps<typeof connector>) => {
  const { userLogin, placeLogin, createPlace, createUser } = props
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [name, setName] = useState("")
  const [profileType, setProfileType] = useState(ProfileType.User)
  const [isRegister, setIsRegister] = useState(false)
  const isUser = profileType === ProfileType.User
  const bgColor = isUser ? "bg-white" : "bg-gray-200"

  const onLogin = async () => {
    const login = isUser ? userLogin : placeLogin
    await login({ email, password })
  }

  const onRegister = async () => {
    if ( password !== confirmPassword ) {
      Alert.alert("Register Error", "Passwords do not match", [{ text: "OK" }], {
        cancelable: false,
      })
      return
    }
    const create = isUser ? createUser : createPlace
    await create({ email, password, name })
  }

  return (
    <View style={tailwind(`flex-1 items-center justify-center ${bgColor}`)}>
      <View style={tailwind(`flex flex-col items-center ${bgColor}`)}>
        { isRegister &&
          <>
            <Text style={tailwind("text-gray-600 font-semibold")}>Name</Text>
            <TextInput
              style={tailwind("rounded-lg border border-black h-8 w-40 p-2 m-2")}
              onChangeText={(text) => setName(text)}
              value={name}
            />
          </>
        }
        
        { isUser &&
          <>
            <Text style={tailwind("text-gray-600 font-semibold")}>Email</Text>
            <TextInput
              style={tailwind("rounded-lg border border-black h-8 w-40 p-2 m-2")}
              onChangeText={(text) => setEmail(text)}
              value={email}
            />
          </>
        }
        
        <View style={tailwind(`items-center ${bgColor}`)}>
          <Text style={tailwind("text-gray-600 font-semibold")}>Password</Text>
          <TextInput
            style={tailwind("rounded-lg border border-black w-40 h-8 p-2 m-2")}
            onChangeText={(text) => setPassword(text)}
            value={password}
            secureTextEntry={true}
          />

          { isRegister &&
            <>
              <Text style={tailwind("text-gray-600 font-semibold")}>Confirm Password</Text>
              <TextInput
                style={tailwind("rounded-lg border border-black w-40 h-8 p-2 m-2")}
                onChangeText={(text) => setConfirmPassword(text)}
                value={confirmPassword}
                secureTextEntry={true}
              />
            </>
          }
        </View>
      </View>

      <View
        style={styles.separator}
        lightColor="#eee"
        darkColor="rgba(255,255,255,0.1)"
      />

      { !isRegister &&
          <TouchableOpacity
            style={tailwind("flex justify-center rounded-lg border border-black p-2 m-2")}
            onPress={onLogin}
          >
            <Text style={tailwind("text-black")}>Login</Text>
          </TouchableOpacity>
      }

      { (!isRegister && isUser) &&
          <TouchableOpacity
            style={tailwind("flex justify-center underline p-2 m-2")}
            onPress={() => setIsRegister(true)}
          >
            <Text style={tailwind("text-black")}>Register?</Text>
          </TouchableOpacity>
      }

      { isRegister &&
        <>
          <TouchableOpacity
            style={tailwind("flex justify-center rounded-lg border border-black p-2 mt-2")}
            onPress={onRegister}
          >
            <Text style={tailwind("text-black")}>Register</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={tailwind("flex justify-center underline p-2 m-2")}
            onPress={() => setIsRegister(false)}
          >
            <Text style={tailwind("text-black")}>Login?</Text>
          </TouchableOpacity>
        </>
      }

      <Picker
        selectedValue={profileType.toString()}
        style={tailwind("h-6 w-24")}
        // itemStyle={tailwind("mt-2")}
        onValueChange={(type) => setProfileType(type as ProfileType)}
      >
        <Picker.Item
          label={ProfileType.User.toString()}
          value={ProfileType.User.toString()}
        />
        <Picker.Item
          label={ProfileType.Store.toString()}
          value={ProfileType.Store.toString()}
        />
      </Picker>
    </View>
  )
}

const styles = StyleSheet.create({
  title: {
    fontSize: 20,
    fontWeight: "bold",
  },
  separator: {
    marginVertical: 30,
    height: 1,
    width: "80%",
  },
})

export default connector(LoginScreen)
