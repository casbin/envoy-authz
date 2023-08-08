# envoy-authz

[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/casbin/envoy-authz/issues)
[![Discord](https://img.shields.io/discord/1022748306096537660?logo=discord&label=discord&color=5865F2)](https://discord.gg/S5UjpzGZjN)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

<p align="center">
    <img width="400" height="400" src="casbin-envoy-logo.png" alt="envoy-authz" />
</p>

Envoy-authz is a middleware of Envoy which performs [external authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter#arch-overview-ext-authz) through casbin. This proxy would be deployed on any type of envoy-based service meshes like Istio. 

## Installation

```
go get github.com/casbin/envoy-authz
```

## Requirements
- Envoy 1.17+ 
- Istio or any type of service mesh
- grpc dependencies

## Working
- A client would make a http request.
- Envoy proxy would send that request to grpc server.
- The grpc server would then authorize the request based on casbin policies.
- If authorized, the request would be sent through or else, it gets denied.

The grpc server is based on protocol buffer from [external_auth.proto](https://github.com/envoyproxy/envoy/blob/master/api/envoy/service/auth/v2alpha/external_auth.proto). 

## Usage
- Define the Casbin policies under config files by following this [guide](https://casbin.org/docs/how-it-works).

You can verify/test your policies on online [casbin-editor](https://casbin.org/editor/).

- Start the authorizing server by running:-
```
$ go build .
$ ./authz 
```
- Load the envoy configuration:-
```
$  envoy -c authz.yaml -l info
```
Once the envoy starts, it will start intercepting requests for the authorization process.

## Integrating to Istio
You need to send custom headers, which would contain usernames in the JWT token OF headers for this middleware to work. You can check the official [Istio docs](https://istio.io/v1.4/docs/tasks/policy-enforcement/control-headers/) to get more info on modifying `Request Headers`.

## Community

In case of any query, you can ask on our [Discord](https://discord.gg/S5UjpzGZjN).

