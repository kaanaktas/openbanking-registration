
![openbanking-register-ci](https://github.com/kaanaktas/openbanking-registration/workflows/openbanking-register-ci/badge.svg)

Introduction
------------

openbanking-registration is a complete application to process dynamic registration journey of Open Banking UK. 
To be able run this journey, user needs;

1- to get registered with [Open Banking Portal](https://www.openbanking.org.uk).

2- to create key pairs for signing and transport and export them to local

3- to create a configuration file for an aspsp same as already created for **OzoneBank** (official test platform of Open Banking) and **Danske** 


Installation and usage
----------------------
You can either use an .env file or export env variables to pass below values. It is recommended to keep keys in a secure environment rather than public access. 

`.env` file:

```dotenv
PORT=8080
#this key is used to sign JWT. It can be obtained from OpenBanking Portal  
OB_SIGN_KEY=<client_signing.key>
#KID value will be the same with your sign public cert. This can be obtained from 
KID=<CLIENT_SIGNING_KID>
#This can be the issuer of Open Banking certificate chain
CLIENT_CA_CERT_PEM=<ob_issuer.cer>
#Key pairs of TPP's transport and used for TLS-MA. This can be obtained from Open Banking portal. 
CLIENT_CERT_PEM=<client_transport.pem>
CLIENT_KEY_PEM=<client_transport.key>
```

Example
-------

### **`Request`**

`http://localhost:8080/<aspsp_id>/register?ssa=<SSA>`


 **`aspsp_id`**: This parameter should point out aspsp which we want to get registered our application.
For testing purposes, Ozone Bank and Danske were added and their configuration can be found under the aspsp folder. 

 **`ssa`**: Software Assertion is a JWT based token which is specfici for TPP. It can be retrieved from Open Banking portal. 
 Different aspsp might have different requirements for expiry time such as `not issued more than 10 mins ago`. 
### **`Expected Response`**
