# Oasis ROFL Registry Client

Query the ROFL registry from Oasis Sapphire. Fetches information about registered ROFL applications and their running instances.

All data is attested by Oasis nodes through remote attestation. Each ROFL app registration includes cryptographically verified enclave identities (MrEnclave/MrSigner) and node endorsements.

## Note on Trust

This tool uses public RPC endpoints by default (`grpc.oasis.io:443`, `testnet.grpc.oasis.io:443`). For trustless verification, run your own Oasis node or light client and point the network config to your local endpoint.

## Usage

```bash
go run main.go <testnet|mainnet> [app_id]
```

Without an app ID, lists all active ROFL apps. With an app ID, queries that specific app.

**Examples:**

```bash
# List active apps on testnet
go run main.go testnet

# Query specific app
go run main.go testnet rofl1qqg3qrpk4484gm8dcayfmnrkcwcg5v3nnusp5d0h

# List active apps on mainnet
go run main.go mainnet
```

**Example output:**

```
Connecting to mainnet...

Latest Block Information:
  Height: 26976318
  Hash: c9934b45ef6471088409838ca0c00ba117d86c20f43933b10e0e45d61fbb9236
  Time: 2025-10-23 15:18:30 +0200 CEST

--- Querying Sapphire ROFL Apps on mainnet ---

Querying specific app: rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9

ROFL App #1:
  ID: rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9
  Admin: oasis1qz2lty9v4glt5ts8ljhfpnd05dy3cwmtnyshws8q
  Stake: 100000000000000000000 <native>
  Policy:
    Fee Policy: 2
    Max Expiration: 3
    Enclaves: 2
      Enclave #1:
        MrEnclave: 8f2a41d6a7d8876629a1741b0e0948c4cc47036c2a396a47ebc70b021a740819
        MrSigner: 0000000000000000000000000000000000000000000000000000000000000000
      Enclave #2:
        MrEnclave: bfa37737aec498bb4a8021ae2e26baf9ac3f66d801d99c5c7c7431bb7067f9cd
        MrSigner: 0000000000000000000000000000000000000000000000000000000000000000
    Endorsements: 1
      Endorsement #1:
        Type: And (2 policies)
  App Metadata:
    net.oasis.rofl.description: Talos is.
    net.oasis.rofl.name: talos
    net.oasis.rofl.version: 0.1.0
    net.oasis.rofl.homepage: https://talos.is
    net.oasis.rofl.repository: https://github.com/talos-agent/talos
  Secrets: 1
    OPENAI_API_KEY: <encrypted>
  Instances: 1
    Instance #1:
      App: rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9
      Node ID: 89HoWT4NuI0cXQOcjXqQUwHjVtICw04sTjPVSH4qb4o=
      Entity ID: <none>
      RAK: fsNl407dnzQggP6eSgKhWzsA4U3GeJcr8uVvb3aexpw=
      Expiration: 44949
      Extra Keys: 1
        Key #1: A55nPX6H4GIaGXJiN/qhEwjmWTx73buKTcCThy28fkI2
      Instance Metadata:
        net.oasis.tls.pk: MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAETd0R2py5Xzu1f7V2GaxL8I+M2dZ+/czz3INYZEHIvyS4kODXe57pTDa3LrGMD0SW0NedSt00l+lHnrayucKGPg==
        net.oasis.policy.provider: omlzaWduYXR1cmVYQEIeDSBx+gSTYLcwkog5WXJHQB9QC9iK+UhZIHSvxeZm56gft/TbJSMaFWm6003sOi1pZGO8Jf9upzMj6/1NlgBxbGFiZWxfYXR0ZXN0YXRpb25YiKJjcmFrWCB+w2XjTt2fNCCA/p5KAqFbOwDhTcZ4lyvy5W9vdp7GnGZsYWJlbHOhcm5ldC5vYXNpcy5wcm92aWRlcnhEb21ocGJuTjBZVzVqWlVnQUFBQUFBQUFBSG1od2NtOTJhV1JsY2xVQXNIRDlyTm5adUxlTWFPV291UXl6bzFJWG1FND0=
```

## Resources

- [Oasis ROFL Documentation](https://docs.oasis.io/build/rofl/)
