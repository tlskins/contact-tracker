import React, { useState, useEffect } from "react"
import {
  Alert,
  Text,
  View,
  StyleSheet,
  Animated,
  TouchableOpacity,
  Image,
} from "react-native"
import { BarCodeScanner } from "expo-barcode-scanner"
import { connect, ConnectedProps } from "react-redux"

import { RootState } from "../state"
import { instance } from "../behavior/rest"

interface BarCodeData {
  type: string
  data: any
}

type permission = true | false | null

const mapState = (state: RootState) => ({
  profile: state.profile,
})

const mapDispatch = (dispatch: any) => {
  return {
    checkIn: async (userId: string, placeId: string): Promise<boolean> => {
      try {
        console.log("before checking in...")
        const resp = await instance.post(`/check-ins`, { userId, placeId })
        console.log("resp", resp)
        const title = resp?.data?.out ? "Check Out" : "Check In"
        Alert.alert(title, "Success", [{ text: "OK" }], {
          cancelable: false,
        })
      } catch (e) {
        const errMsg = e.response?.data?.message || e.message || ""
        Alert.alert("Check In Error", errMsg, [{ text: "OK" }], {
          cancelable: false,
        })
      }
      return false
    },
  }
}

const connector = connect(mapState, mapDispatch)

const ScannerScreen = (props: ConnectedProps<typeof connector>) => {
  const { checkIn, profile } = props
  const [hasCameraPermission, setCameraPermission] = useState(
    null as permission
  )
  const [scanned, setScanned] = useState(false)
  const [animationLineHeight, setAnimationLineHeight] = useState(0)
  const [focusLineAnimation, setFocusLineAnimation] = useState(
    new Animated.Value(0)
  )
  useEffect(() => {
    getPermissionsAsync()
    animateLine()
  }, [])
  const animateLine = () => {
    Animated.sequence([
      Animated.timing(focusLineAnimation, {
        toValue: 1,
        duration: 1000,
        useNativeDriver: true,
      }),
      Animated.timing(focusLineAnimation, {
        toValue: 0,
        duration: 1000,
        useNativeDriver: true,
      }),
    ]).start(animateLine)
  }
  const getPermissionsAsync = async () => {
    const { status } = await BarCodeScanner.requestPermissionsAsync()
    const isPermissionGranted = status === "granted"
    console.log(isPermissionGranted)
    setCameraPermission(isPermissionGranted)
  }
  const handleBarCodeScanned = ({ type, data }: BarCodeData) => {
    setScanned(true)
    if (profile && data) {
      checkIn(data, profile.id)
    } else {
      Alert.alert(
        "Check In Error",
        `Bad data... ${data} ${profile.id}`,
        [{ text: "OK" }],
        {
          cancelable: false,
        }
      )
    }
  }
  if (hasCameraPermission === null) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
        <Text>Requesting for camera permission</Text>
      </View>
    )
  }
  if (hasCameraPermission === false) {
    return (
      <View style={{ flex: 1, alignItems: "center", justifyContent: "center" }}>
        <Text>No access to camera</Text>
      </View>
    )
  }
  return (
    <View style={styles.container}>
      <BarCodeScanner
        onBarCodeScanned={scanned ? () => {} : handleBarCodeScanned}
        style={StyleSheet.absoluteFillObject}
      />
      <View style={styles.overlay}>
        <View style={styles.unfocusedContainer}></View>
        <View style={styles.middleContainer}>
          <View style={styles.unfocusedContainer}></View>
          <View
            onLayout={(e) =>
              setAnimationLineHeight(e.nativeEvent.layout.height)
            }
            style={styles.focusedContainer}
          >
            {!scanned && (
              <Animated.View
                style={[
                  styles.animationLineStyle,
                  {
                    transform: [
                      {
                        translateY: focusLineAnimation.interpolate({
                          inputRange: [0, 1],
                          outputRange: [0, animationLineHeight],
                        }),
                      },
                    ],
                  },
                ]}
              />
            )}
            {scanned && (
              <TouchableOpacity
                onPress={() => setScanned(false)}
                style={styles.rescanIconContainer}
              >
                <Image
                  source={require("./rescan.png")}
                  style={{ width: 50, height: 50 }}
                />
              </TouchableOpacity>
            )}
          </View>
          <View style={styles.unfocusedContainer}></View>
        </View>
        <View style={styles.unfocusedContainer}></View>
      </View>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    position: "relative",
  },
  overlay: {
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
  },
  unfocusedContainer: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.7)",
  },
  middleContainer: {
    flexDirection: "row",
    flex: 1.5,
  },
  focusedContainer: {
    flex: 6,
  },
  animationLineStyle: {
    height: 2,
    width: "100%",
    backgroundColor: "red",
  },
  rescanIconContainer: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
  },
})

export default connector(ScannerScreen)
