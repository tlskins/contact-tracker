import { Ionicons } from "@expo/vector-icons"
import { MaterialCommunityIcons } from "@expo/vector-icons"
import { createBottomTabNavigator } from "@react-navigation/bottom-tabs"
import { createStackNavigator } from "@react-navigation/stack"
import * as React from "react"
import { connect, ConnectedProps } from "react-redux"

import Colors from "../constants/Colors"
import useColorScheme from "../hooks/useColorScheme"
import LoginScreen from "../screens/LoginScreen"
import CheckInsScreen from "../screens/CheckInsScreen"
import ScannerScreen from "../screens/ScannerScreen"
import ProfileScreen from "../screens/ProfileScreen"
import HistoryScreen from "../screens/HistoryScreen"
import { RootState } from "../state"
import { ProfileType, User } from "../state/types"

import {
  BottomTabParamList,
  LoginParamList,
  CheckInsParamList,
  ScannerParamList,
  ProfileParamList,
} from "../types"

const BottomTab = createBottomTabNavigator<BottomTabParamList>()

const mapState = (state: RootState) => ({ profile: state.profile })

const connector = connect(mapState, undefined)

const BottomTabNavigator = (props: ConnectedProps<typeof connector>) => {
  const { profile } = props
  const colorScheme = useColorScheme()

  if (profile === null) {
    return (
      <BottomTab.Navigator
        initialRouteName="Login"
        tabBarOptions={{ activeTintColor: Colors[colorScheme].tint }}
      >
        <BottomTab.Screen
          name="Login"
          component={LoginNavigator}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarMaterialIcon name="login" color={color} />
            ),
          }}
        />
      </BottomTab.Navigator>
    )
  }

  if (profile?.profileType === ProfileType.User) {
    return (
      <BottomTab.Navigator
        initialRouteName="Profile"
        tabBarOptions={{ activeTintColor: Colors[colorScheme].tint }}
      >
        <BottomTab.Screen
          name="Profile"
          component={ProfileNavigator}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarMaterialIcon name="account-circle" color={color} />
            ),
          }}
        />
        <BottomTab.Screen
          name="CheckIns"
          component={CheckInsNavigator}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarIcon name="md-checkmark-circle-outline" color={color} />
            ),
          }}
        />
      </BottomTab.Navigator>
    )
  } else {
    return (
      <BottomTab.Navigator
        initialRouteName="Profile"
        tabBarOptions={{ activeTintColor: Colors[colorScheme].tint }}
      >
        {/* <BottomTab.Screen
          name="History"
          component={HistoryScreen}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarMaterialIcon name="account-circle" color={color} />
            ),
          }}
        /> */}
        <BottomTab.Screen
          name="Profile"
          component={ProfileNavigator}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarMaterialIcon name="account-circle" color={color} />
            ),
          }}
        />
        <BottomTab.Screen
          name="Scanner"
          component={ScannerNavigator}
          options={{
            tabBarIcon: ({ color }) => (
              <TabBarIcon name="md-qr-scanner" color={color} />
            ),
          }}
        />
      </BottomTab.Navigator>
    )
  }
}

// You can explore the built-in icon families and icons on the web at:
// https://icons.expo.fyi/
function TabBarIcon(props: { name: string; color: string }) {
  return <Ionicons size={30} style={{ marginBottom: -3 }} {...props} />
}

function TabBarMaterialIcon(props: { name: string; color: string }) {
  return (
    <MaterialCommunityIcons size={30} style={{ marginBottom: -3 }} {...props} />
  )
}

// Each tab has its own navigation stack, you can read more about this pattern here:
// https://reactnavigation.org/docs/tab-based-navigation#a-stack-navigator-for-each-tab
const LoginStack = createStackNavigator<LoginParamList>()

function LoginNavigator() {
  return (
    <LoginStack.Navigator>
      <LoginStack.Screen
        name="LoginScreen"
        component={LoginScreen}
        options={{ headerTitle: "Login" }}
      />
    </LoginStack.Navigator>
  )
}

const CheckInsStack = createStackNavigator<CheckInsParamList>()

function CheckInsNavigator() {
  return (
    <CheckInsStack.Navigator>
      <CheckInsStack.Screen
        name="CheckInsScreen"
        component={CheckInsScreen}
        options={{ headerTitle: "Check Ins" }}
      />
    </CheckInsStack.Navigator>
  )
}

const ScannerStack = createStackNavigator<ScannerParamList>()

function ScannerNavigator() {
  return (
    <ScannerStack.Navigator>
      <ScannerStack.Screen
        name="ScannerScreen"
        component={ScannerScreen}
        options={{ headerTitle: "Check In Scanner" }}
      />
    </ScannerStack.Navigator>
  )
}

const ProfileStack = createStackNavigator<ProfileParamList>()

function ProfileNavigator() {
  return (
    <ProfileStack.Navigator>
      <ProfileStack.Screen
        name="ProfileScreen"
        component={ProfileScreen}
        options={{ headerTitle: "Profile" }}
      />
    </ProfileStack.Navigator>
  )
}

export default connector(BottomTabNavigator)
