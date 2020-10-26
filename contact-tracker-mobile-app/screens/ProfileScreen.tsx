import * as React from "react"
import { StyleSheet, TouchableOpacity } from "react-native"
import tailwind from "tailwind-rn"
import { connect, ConnectedProps } from "react-redux"
import SvgQRCode from "react-native-qrcode-svg"

import { Text, View } from "../components/Themed"
import { RootState } from "../state"
import { clearProfile } from "../state/actions"

const mapState = (state: RootState) => ({
  profile: state.profile,
})

const mapDispatch = (dispatch: any) => {
  return {
    signOut: async (): Promise<void> => {
      dispatch(clearProfile())
    },
  }
}

const connector = connect(mapState, mapDispatch)

const ProfileStackScreen = (props: ConnectedProps<typeof connector>) => {
  const { profile, signOut } = props
  const name = profile?.name || "Not Logged In"

  console.log("profile", profile)

  return (
    <View style={tailwind("h-full items-center bg-gray-500 p-12 pt-20")}>
      <Text style={styles.title}>{name}</Text>
      <Text style={styles.title}>{ profile?.email }</Text>
      <View
        style={styles.separator}
        lightColor="#eee"
        darkColor="rgba(255,255,255,0.1)"
      />
      <TouchableOpacity
        style={tailwind(
          "flex justify-center bg-white rounded-lg border border-black p-2 m-4"
        )}
        onPress={signOut}
      >
        <Text>Logout</Text>
      </TouchableOpacity>
      <SvgQRCode value={profile?.id} />
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

export default connector(ProfileStackScreen)
