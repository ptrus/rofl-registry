# Oasis ROFL Registry Client

Query the ROFL registry from Oasis Sapphire. Fetches information about registered ROFL applications and their running instances.

All data is attested by Oasis nodes through remote attestation. Attestation reports are verified on-chain by Oasis consensus. Each registration includes verified enclave identities (MrEnclave, MrSigner), policies, and instance metadata attested by the Oasis network to originate from the registered ROFL application running the exact enclave binary.

**Note:** Application metadata (name, version, description, homepage, repository) is provided by the ROFL app admin and is not cryptographically attested. All other registry data is cryptographically attested by Oasis consensus and can be trusted.

## A Note on Attestations

The Oasis network automatically verifies and enforces remote attestation on the consensus layer. This covers:

- Enclave identities – trusted enclave code is verified by measurement
- Attestation quote and hardware TCB status – validated against network-defined policies
- RAK/REK binding – ensures enclave signing and encryption keys are genuine
- Freshness – attestation quotes must be recent and are re-verified periodically
- Node endorsements – binds the attestation to a specific node (for ROFL applications)
- …and other attestation conditions required by the network

Unlike manual quote verification, this ensures policy enforcement, binding, and freshness cannot be skipped or misinterpreted - unlike in other systems that present clients with vague or one-off attestation signatures, which often provide a false sense of security.

That’s why the Oasis ROFL Registry, verified directly on the Oasis consensus layer, represents a more trustworthy and complete attestation source than most systems relying on client attestations.

## A Note on Accessing Trusted Registry State

This tool uses public RPC endpoints by default (`grpc.oasis.io:443`, `testnet.grpc.oasis.io:443`). For end-to-end trustless verification (without relying on public RPC operators), run your own Oasis node or light client and point the network config to your local endpoint.

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

## Verifying ROFL Instance TLS Certificates

ROFL instances using [rofl-proxy](https://github.com/oasisprotocol/oasis-sdk/tree/main/rofl-proxy) expose their TLS public key in the instance metadata under the key `net.oasis.tls.pk`. The TLS private key is generated inside the TEE and never leaves it, ensuring that only the attested enclave can use it (see rofl-proxy implementation for details).

You can verify that a service's TLS certificate matches the attested on-chain instance:

**Example: Verify TLS Certificate Used by ROFL Helios Testnet Deployment**

[ROFL Helios](https://github.com/ptrus/rofl-helios) is an Ethereum Light Client running in ROFL.

1. First, fetch the state and query the ROFL app to see its details:
```bash
go run main.go testnet rofl1qzul9krxsnuanfqqte337utwxr47gqe4zu6rcr5z
```

This will show you the app details. Notice the reported URL in the app metadata (`net.oasis.rofl.homepage: https://ethrpc.rofl.cloud`) and the TLS public key in the instance metadata:
```
net.oasis.tls.pk: MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0gqbijVXY+svlCtWuZw5uA82AlU0gth36TAK+zi6tDZscCHSDl4fy82/DrPBdopa3N5kvB9bx+cekiTGpG2kCg==
```

2. Extract the TLS certificate's public key from the running service and ensure it matches:
```bash
echo | openssl s_client -servername ethrpc.rofl.cloud -connect ethrpc.rofl.cloud:443 2>/dev/null | \
  openssl x509 -pubkey -noout | \
  openssl ec -pubin -outform DER 2>/dev/null | \
  base64
```

Output:
```
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE0gqbijVXY+svlCtWuZw5uA82AlU0gth36TAK+zi6tDZscCHSDl4fy82/DrPBdopa3N5kvB9bx+cekiTGpG2kCg==
```

The `net.oasis.tls.pk` value from the on-chain registry matches the TLS certificate served by `https://ethrpc.rofl.cloud`, confirming the service is operated by the attested ROFL instance.

**What this proves:**
- The TLS private key exists only inside the TEE and cannot be extracted or accessed by anyone outside the enclave
- Your HTTPS connection terminates inside the enclave, and the TLS private key is sealed within it — no external entity (node operator, scheduler, proxy) can access it
- The service at `https://ethrpc.rofl.cloud` is operated by the attested ROFL instance with verified enclave identity (MrEnclave)
- The instance's code matches the cryptographic measurements recorded in the app policy

This creates a complete trust chain: Oasis consensus → ROFL registry → Enclave attestation → TLS certificate → HTTPS endpoint.

## Verifying Exact ROFL Code

You can verify that the on-chain enclave measurements correspond exactly to the code in its repository. This proves that the running instances execute the specific code version you can audit.

**Example: Verify Talos Agent Code**

Using the Talos example from above (app ID: `rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9`), we can verify that the registered enclave identities match the code in the repository.

By rebuilding the ROFL application locally, you can reproduce its enclave measurement (MrEnclave) and confirm that it matches the one registered on-chain.

1. Query the app to see its metadata and enclave measurements:
```bash
go run main.go mainnet rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9
```

This shows the repository URL and registered MrEnclave values:
```
  App Metadata:
    net.oasis.rofl.repository: https://github.com/talos-agent/talos
  Policy:
    Enclaves: 2
      Enclave #1:
        MrEnclave: 8f2a41d6a7d8876629a1741b0e0948c4cc47036c2a396a47ebc70b021a740819
      Enclave #2:
        MrEnclave: bfa37737aec498bb4a8021ae2e26baf9ac3f66d801d99c5c7c7431bb7067f9cd
```

2. Clone the repository and rebuild with verification:
```bash
git clone https://github.com/talos-agent/talos
cd talos
oasis rofl build --verify --deployment mainnet
```

Output on successful verification:
```
Building a ROFL application...
Deployment: mainnet
Network:    mainnet
ParaTime:   sapphire
App ID:     rofl1qpykfkl6ea78cyy67d35f7fmpk3pg36vashka4v9
Name:       talos
Version:    0.1.0
...
Computing enclave identity...
Built enclave identities MATCH latest manifest enclave identities.
Manifest enclave identities MATCH on-chain enclave identities.
```

**What this proves:**
- The exact source code in the repository produces the same enclave measurements (MrEnclave) as registered on-chain
- Any running instance with these measurements must be executing this specific code version
- You can now audit the full source code to understand precisely what the ROFL app executes inside the enclave
- The build is reproducible - anyone can independently verify the same result

This completes the verification chain: Source code → Deterministic build → Enclave measurement (MrEnclave) → On-chain registration → Attested running instance.

## Resources

- [Oasis ROFL Documentation](https://docs.oasis.io/build/rofl/)
