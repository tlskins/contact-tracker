## Contact Tracker Mobile App
Contact tracking mobile app for store owners.

# Mobile App
Iphone and Android compatible mobile app written using expo, react-native, and typescript. Expo is a framework and platform for creating and distributing react native applications. To test an Expo app the user will want to decide what host(s) and client(s) you will want to use which can allow you to test different distributions on different mobile devices.

# App Host
After cloning this project a user can run their own Expo server via Expo CLI or you can test a staging distrubtion of this project via the url below. Running your own Expo server will allow you to test changes you make to the codebase as well as let you test the app on emulators you have installed on your computer. The staging distribution will be a stable distribution of this project that is publicly available. The staging distribution's backend api host will be pointing to serverless host that I will maintain for a limited period of time.

Staging URL
https://expo.io/@tlskins/projects/contact-tracker-mobile-app?release-channel=staging

# Client
The client can be tested via an Iphone or Android mobile device or an emulator. To test on your mobile device you will need to download the app Expo client. If you are testing the staging distribution, clicking the link for the distribution app host will invite you to test the app in Expo client. 

Test on Android simulator using Android Studio 3.0+
https://docs.expo.io/workflow/android-studio-emulator/

Test on Iphone simulator using Xcode
https://docs.expo.io/workflow/ios-simulator/


## Contact Tracker API
The API for this app is a go serverless microservice architecture with MongoDB Atlas and AWS as the cloud provider. This backend has been designed using the clean architecture in mind so the use cases have no dependencies on the database, cloud provider, or delivery method. This repo has both http and serverless deliveries available so you can run your own Go microservice web servers or use your own AWS Lambdas.


# AWS Serverless
After cloning the repo you will need to create a config.prod.yml similar to the one below:

```
USERS_HOST: "https://abcde.execute-api.us-east-1.amazonaws.com/dev"
PLACES_HOST: "https://abcde.execute-api.us-east-1.amazonaws.com/dev"
CHECK_INS_HOST: "https://abcde.execute-api.us-east-1.amazonaws.com/dev"
MONGO_DB_NAME: test
MONGO_HOST: cluster0-abcd.mongodb.net
MONGO_USER: admin
MONGO_PWD: PWD123467890
JWT_KEY_PATH: ./bin/id_rsa.pub
JWT_SECRET_PATH: ./bin/id_rsa
JWT_ACCESS_EXPIR: 20160
JWT_REFRESH_EXPIR: 20160
AWS_SES_ACCESS_KEY: KEY1234567890
AWS_SES_ACCESS_SECRET: SECRET1234567890
AWS_SES_REGION: us-east-1
SENDER_EMAIL: contact.tracker@email.provider.com
RPC_AUTH_PWD: PWD123!@#
```

To deploy to AWS you can just execute the make file `deploy`
```
make deploy
```


