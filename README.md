# go-auth-api

## 1. Pull the docker image :
```
  docker pull siddheshk02/auth-api
```

## 2. Run the image :

```
  docker run -p 3000:3000 auth-api
```
## Output :
![image](https://github.com/Siddheshk02/go-auth-api/assets/90148705/1e54d705-c7bf-48d9-a69a-36b530231767)

3. Use Postman or any other API testing tool and type this url ``` http://127.0.0.1:3000/auth ``` Method = POST, Pass the email and password as shown below through the API's request Body
``` 
   {
    "email": "<YOUR_EMAIL>",
    "password": "<YOUR_PASSWORD>"
   }
```

This will Authenticate the User if Already registered else it will create new User and generate the token.

4. For the Second Endpoint type this url in place of previous one ``` http://127.0.0.1:3000/user ``` Method = POST, Pass the data as shown below through the API's request Body :
``` 
   {
    "email": "<YOUR_EMAIL>",
    "name" : "<YOUR_NAME",
    "phone" : "<YOUR_PHONE>"
   }
```
 
