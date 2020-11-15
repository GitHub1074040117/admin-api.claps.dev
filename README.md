# admin-api.claps.dev
## Before running
### Create a folder in the current working path. Name the folder "config".
Create application.yml file under this folder.
Here are some examples:
```
# the port that the server running on
server:
  port: 7003

# database connection
datasource:
  driverName : postgres
  host : localhost
  port : 5432
  database : 
  username : 
  password : 
  charset : utf8

# github oauth
oauthConf:
  clientId: 37*********d6cb8ac
  clientSecret: 0d24840ca01c***********463c9681df0
  redirectUrl: http://localhost:8001/login

# your github account, role can either be "super" or "common"
superAdmin:
  userId: 1
  name: 
  account: 
  avatarUrl: 
  role: super

# mixin keystore
keystore:
  path: ./keystore-7000103520.json

# merico api configuration
merico:
  basePath: https://cloud.merico.cn
  contentType: application/json
  appId: 2*******3
  key: 9ebf****************23a9
```

### Create a minxin robot
Visit https://developers.mixin.one/dashboard to create your robot.
Create file keystore-7000103520.json under the "config" directory.
```
{
 "pin": "",
 "client_id": "",
 "session_id": "",
 "pin_token": "",
 "private_key": ""
}
```

## How to run
Compile to generate executable file
```
go build
```

Deploy to the server and run on port 7003
```
./claps-admin &
```
