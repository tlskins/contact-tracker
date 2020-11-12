import React, { useEffect, useState } from "react"
import { Alert, StyleSheet, TouchableOpacity } from "react-native"
import tailwind from "tailwind-rn"
import { connect, ConnectedProps } from "react-redux"
import SvgQRCode from "react-native-qrcode-svg"

import { Text, View } from "../components/Themed"
import { RootState } from "../state"
import { CheckInHistory } from "../state/types"
import { instance } from "../behavior/rest"

const mapState = (state: RootState) => ({
  profile: state.profile,
})

const mapDispatch = (dispatch: any) => {
  return {
    getHistory: async (placeId: string): Promise<CheckInHistory[]> => {
      try {
        console.log("before get history...")
        const resp = await instance.get(`/check-ins/history/${placeId}`)
        console.log("resp", resp)
        if (resp?.data) {
          return resp?.data
        }
      } catch (e) {
        console.log("loading history err", e)
        const errMsg = e.response?.data?.message || e.message || ""
        Alert.alert("Loading Error", errMsg, [{ text: "OK" }], {
          cancelable: false,
        })
      }
      return []
    },
  }
}

const connector = connect(mapState, mapDispatch)

const HistoryStackScreen = (props: ConnectedProps<typeof connector>) => {
  const { profile, getHistory } = props
  const name = profile?.name || "Not Logged In"
  const [history, setHistory] = useState<CheckInHistory[]>([])

  useEffect(() => {
    if (profile) {
      console.log("loading history...")
      getHistory(profile.id)
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
      {(history || []).map((check) => (
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
          {(check.contacts || []).map((contact) => (
            <View key={check.id} style={tailwind("rounded-lg p-6 m-2")}>
              <View
                style={tailwind(
                  "rounded-full text-center items-center content-center bg-blue-200 m-2 p-2"
                )}
              >
                <Text style={tailwind("text-blue-800 font-semibold")}>
                  {contact.user.name}
                </Text>
              </View>
              <Text style={tailwind("text-blue-800 font-semibold")}>
                In: {contact.in}
              </Text>
              <Text style={tailwind("text-blue-800 font-semibold")}>
                Out: {contact.out}
              </Text>
            </View>
          ))}
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

export default connector(HistoryStackScreen)
