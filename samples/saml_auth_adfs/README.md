# SAML ADFS example
Package main is an example how to use Active Directory Federation Services as SAML IdP for vCD.
main() function has an example how to setup vCD client with SAML auth
To run this command please supply parameters as below
```
go run main.go --username test@test-forest.net --password my-password --org my-org --endpoint https://192.168.1.160/api
```

Results should look similar to:
```
Found 1 Edge Gateways
my-edge-gw
```