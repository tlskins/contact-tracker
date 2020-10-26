import React, { useEffect } from "react"
import { Alert, StyleSheet } from "react-native"
import tailwind from "tailwind-rn"
import { connect, ConnectedProps } from "react-redux"

import { Text, View } from "../components/Themed"
import { RootState } from "../state"
import { setCheckIns } from "../state/actions"
import { instance } from "../behavior/rest"

const mapState = (state: RootState) => ({
  profile: state.profile,
  checkIns: state.checkIns,
})

const mapDispatch = (dispatch: any) => {
  return {
    loadCheckIns: async (userId: string): Promise<boolean> => {
      try {
        console.log("before loading checkins...")
        const resp = await instance.get(`/check-ins?userId=${userId}`)
        console.log("resp", resp)
        const data = resp?.data || []
        dispatch(setCheckIns(data))
        return true
      } catch (e) {
        console.log("loading checkins err", e)
        const errMsg = e.response?.data?.message || e.message || ""
        Alert.alert("Loading Error", errMsg, [{ text: "OK" }], {
          cancelable: false,
        })
      }
      return false
    },
  }
}

const connector = connect(mapState, mapDispatch)

const CheckInsStackScreen = (props: ConnectedProps<typeof connector>) => {
  const { checkIns, loadCheckIns, profile } = props
  const name = profile?.name || "Not Logged In"

  console.log("checkIns", checkIns)

  useEffect(() => {
    if (profile) {
      console.log("loading checkins...")
      loadCheckIns(profile.id)
    }
  }, [])
  
  return (
    <View style={tailwind("h-full items-center bg-gray-500 p-12 pt-20")}>
      <Text style={styles.title}>{name}</Text>
      <View
        style={styles.separator}
        lightColor="#eee"
        darkColor="rgba(255,255,255,0.1)"
      />
      { (checkIns || []).length === 0 &&
        <Text style={styles.title}>No Check-Ins</Text>
      }
      {(checkIns || []).map((check) => (
        <View key={check.id} style={tailwind("rounded-lg p-6 m-2")}>
          <View
            style={tailwind(
              "rounded-full text-center items-center content-center bg-blue-200 m-2 p-2"
            )}
          >
            <Text style={tailwind("text-blue-800 font-semibold")}>
              {check.place.name}
            </Text>
          </View>
          <Text style={tailwind("text-blue-800 font-semibold")}>
            In: {check.in}
          </Text>
          <Text style={tailwind("text-blue-800 font-semibold")}>
            Out: {check.out}
          </Text>
        </View>
      ))}
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

export default connector(CheckInsStackScreen)
