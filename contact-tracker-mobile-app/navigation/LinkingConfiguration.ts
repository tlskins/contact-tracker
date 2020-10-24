import * as Linking from "expo-linking"

export default {
  prefixes: [Linking.makeUrl("/")],
  config: {
    screens: {
      Root: {
        screens: {
          Login: {
            screens: {
              LoginScreen: "login",
            },
          },
          CheckIns: {
            screens: {
              CheckInsScreen: "CheckIns",
            },
          },
          Scanner: {
            screens: {
              ScannerScreen: "Scanner",
            },
          },
          Profile: {
            screens: {
              ProfileScreen: "Profile",
            },
          },
        },
      },
      NotFound: "*",
    },
  },
}
