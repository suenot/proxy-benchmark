# Proxy Configuration

<cite>
**Referenced Files in This Document**  
- [config.go](file://config.go)
- [proxy.go](file://proxy.go)
- [config.example.json](file://config.example.json)
- [http_client.go](file://http_client.go)
- [socks5_client.go](file://socks5_client.go)
- [benchmark.go](file://benchmark.go)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [Proxies Field Structure](#proxies-field-structure)
3. [Proxy String Format and Parsing](#proxy-string-format-and-parsing)
4. [Supported Protocols](#supported-protocols)
5. [Authentication Handling](#authentication-handling)
6. [Validation and Error Handling](#validation-and-error-handling)
7. [Enabled/Disabled Status Management](#enabled-disabled-status-management)
8. [Integration with Benchmarking Engine](#integration-with-benchmarking-engine)
9. [Common Configuration Mistakes](#common-configuration-mistakes)
10. [Real-World Examples](#real-world-examples)

## Introduction
This document provides comprehensive guidance on configuring proxies for the benchmarking application. The Proxies field in the Config struct serves as the central mechanism for defining proxy servers used during performance testing. Each proxy is represented as a string following a specific format that encodes protocol, host, port, authentication credentials, and operational status. This section explains the complete configuration workflow from JSON definition to runtime usage.

## Proxies Field Structure
The Proxies field is defined within the Config struct as a JSON array of strings, where each string represents an individual proxy server configuration. This design choice enables flexible configuration through external JSON files while maintaining simplicity in parsing and validation.

```mermaid
classDiagram
class Config {
+[]string Proxies
+BenchmarkConfig Benchmark
+StatisticsConfig Statistics
}
class Proxy {
+string Protocol
+string Host
+string Port
+string Username
+string Password
+string Status
}
```

**Diagram sources**
- [config.go](file://config.go#L8-L12)
- [proxy.go](file://proxy.go#L8-L15)

**Section sources**
- [config.go](file://config.go#L8-L12)

## Proxy String Format and Parsing
Each proxy string follows the format: `protocol:host:port:username:password:status`. The ParseProxy function in proxy.go splits the input string by colons and validates that exactly six components are present. This strict format ensures consistent parsing and prevents ambiguous configurations.

The parsing process converts the string representation into an internal Proxy struct that is used throughout the application. If the format is invalid (e.g., missing components or extra colons), the parser returns an error which is handled gracefully by the benchmark engine.

```mermaid
flowchart TD
Start([Proxy String Input]) --> Split["Split by ':' delimiter"]
Split --> Count{"Exactly 6 parts?"}
Count --> |No| Error["Return 'invalid proxy format' error"]
Count --> |Yes| Construct["Construct Proxy struct"]
Construct --> Validate["Set Protocol, Host, Port, etc."]
Validate --> Output["Return *Proxy object"]
style Error fill:#f8b7bd,stroke:#333
style Output fill:#a8e4a0,stroke:#333
```

**Diagram sources**
- [proxy.go](file://proxy.go#L18-L32)
- [proxy.go](file://proxy.go#L21)

**Section sources**
- [proxy.go](file://proxy.go#L18-L32)

## Supported Protocols
The system supports three proxy protocols:
- **http**: Standard HTTP proxy connections
- **https**: Secure HTTP proxy connections
- **socks**: SOCKS5 proxy protocol (referred to as "socks" in configuration)

Protocol identification occurs at runtime when the benchmark engine processes each proxy. The protocol value is case-sensitive and must match exactly ("http", "https", or "socks"). Any other value results in the proxy being rejected during execution.

```mermaid
stateDiagram-v2
[*] --> CheckProtocol
CheckProtocol --> HTTP : protocol == "http"
CheckProtocol --> HTTPS : protocol == "https"
CheckProtocol --> SOCKS : protocol == "socks"
CheckProtocol --> Invalid : otherwise
Invalid --> Reject["Reject proxy"]
HTTP --> ConfigureHTTPClient
HTTPS --> ConfigureHTTPClient
SOCKS --> ConfigureSOCKS5Client
ConfigureHTTPClient --> Ready
ConfigureSOCKS5Client --> Ready
Ready --> [*]
```

**Diagram sources**
- [benchmark.go](file://benchmark.go#L94-L128)
- [benchmark.go](file://benchmark.go#L190-L237)

## Authentication Handling
Authentication credentials are embedded directly in the proxy string format through the username and password fields. These values are extracted during parsing and used to configure the appropriate client:

- For HTTP/HTTPS proxies: Credentials are incorporated into the proxy URL using standard basic authentication syntax (`http://user:pass@host:port`)
- For SOCKS5 proxies: Credentials are passed separately to the SOCKS5 dialer via an Auth struct

The system assumes all credentials are provided in plain text within the configuration file. Users should ensure proper file permissions to protect sensitive authentication information.

```mermaid
sequenceDiagram
participant Config as "Config File"
participant Parser as "ParseProxy()"
participant HTTPClient as "HTTPClient"
participant SOCKS5Client as "SOCKS5Client"
Config->>Parser : proxyString="http : host : port : user : pass : enabled"
Parser->>Parser : Split into components
Parser->>HTTPClient : NewHTTPClient(proxy, timeout)
HTTPClient->>HTTPClient : Build URL with auth
HTTPClient-->>HTTPClient : http : //user : pass@host : port
Config->>Parser : proxyString="socks : host : port : user : pass : enabled"
Parser->>Parser : Split into components
Parser->>SOCKS5Client : NewSOCKS5Client(proxy, timeout)
SOCKS5Client->>SOCKS5Client : Create Auth struct
SOCKS5Client-->>SOCKS5Client : User=user, Password=pass
```

**Diagram sources**
- [http_client.go](file://http_client.go#L17-L36)
- [socks5_client.go](file://socks5_client.go#L16-L40)

## Validation and Error Handling
The system implements strict validation of proxy strings during parsing. The primary validation rule checks that the input string contains exactly six colon-separated components. If this condition fails, ParseProxy returns an error with the message "invalid proxy format".

Errors are handled gracefully by the benchmark engine, which skips invalid proxies after logging a warning message. This allows partial configuration validity—when some proxies are malformed, the system continues processing valid ones rather than failing entirely.

```go
// Example of validation logic
if len(parts) != 6 {
    return nil, fmt.Errorf("invalid proxy format: %s", proxyString)
}
```

**Section sources**
- [proxy.go](file://proxy.go#L21)
- [benchmark.go](file://benchmark.go#L18-L36)

## Enabled/Disabled Status Management
Proxy activation status is controlled through the sixth component of the proxy string. Only proxies with status "enabled" are included in benchmarking operations. This status field acts as a simple toggle mechanism without requiring structural changes to the configuration.

Unlike traditional approaches that might use separate enabled/disabled arrays or boolean flags, this implementation uses inclusion-based filtering: disabled proxies can either be set to any value other than "enabled" or omitted entirely from the configuration array.

```go
// In benchmark.go - only enabled proxies are added
if proxy.IsValid() {
    proxies = append(proxies, proxy)
}
```

**Section sources**
- [proxy.go](file://proxy.go#L34-L37)
- [benchmark.go](file://benchmark.go#L18-L36)

## Integration with Benchmarking Engine
The proxy configuration integrates seamlessly with the benchmarking engine through a multi-step process:

1. Configuration loading reads the JSON file into the Config struct
2. Proxy strings are parsed into Proxy objects
3. Valid and enabled proxies are filtered and stored
4. During execution, the engine creates protocol-specific clients based on the Proxy.Protocol field

HTTP and HTTPS proxies use the standard net/http library with proxy configuration, while SOCKS5 proxies leverage the golang.org/x/net/proxy package for connection handling.

```mermaid
graph TB
A[config.json] --> B[LoadConfig]
B --> C[Config.Proxies]
C --> D{For each proxy}
D --> E[ParseProxy]
E --> F{Valid?}
F --> |No| G[Log Warning]
F --> |Yes| H{Status=enabled?}
H --> |No| I[Skip]
H --> |Yes| J[Add to active proxies]
J --> K[BenchmarkEngine]
K --> L{Make Requests}
L --> M[HTTPClient]
L --> N[SOCKS5Client]
```

**Diagram sources**
- [config.go](file://config.go#L32-L47)
- [benchmark.go](file://benchmark.go#L18-L36)
- [http_client.go](file://http_client.go#L17-L36)
- [socks5_client.go](file://socks5_client.go#L16-L40)

## Common Configuration Mistakes
Several common errors occur when configuring proxies:

- **Incorrect component count**: Missing or extra colons resulting in ≠6 parts
- **Invalid protocol names**: Using "socks5" instead of "socks", or "HTTP" with uppercase letters
- **Missing ports**: Omitting the port number even when required by the protocol
- **Special characters in credentials**: Including colons or other delimiters in usernames/passwords
- **Wrong status values**: Using "enable" instead of "enabled"

These mistakes typically result in parsing errors that prevent the proxy from being used, though the system continues processing other valid entries.

**Section sources**
- [proxy.go](file://proxy.go#L18-L32)
- [benchmark.go](file://benchmark.go#L18-L36)

## Real-World Examples
Based on config.example.json, valid proxy configurations include:

- `"socks:proxy1.example.com:1080:username:password:enabled"` - SOCKS5 proxy with authentication
- `"http:proxy2.example.com:8080:username:password:enabled"` - HTTP proxy with authentication
- `"https:proxy3.example.com:443:username:password:enabled"` - HTTPS proxy with authentication

All examples follow the exact six-part format with lowercase protocol names and explicit "enabled" status. The benchmark engine uses these configurations to create appropriate clients for performance testing against the target URL.

```json
{
  "proxies": [
    "socks:proxy1.example.com:1080:username:password:enabled",
    "http:proxy2.example.com:8080:username:password:enabled",
    "https:proxy3.example.com:443:username:password:enabled"
  ]
}
```

**Section sources**
- [config.example.json](file://config.example.json#L2-L8)