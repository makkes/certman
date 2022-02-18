# certman

certman is an opinionated tool for managing TLS certificates. It creates certificate authorities (CAs) and certificates signed by these CAs. All just with a handful of simple CLI commands.

## Creating a CA

```sh
certman create ca --cert-out ca.pem --key-out ca.key --config ca.yaml
```

This commands stores the CA certificate in `ca.pem` and its private key in `ca.key`. The values of the certificate are derived from the configuration given in `ca.yaml` which looks similar to this example:

```yaml
apiVersion: v1
commonName: My awesome CA
organization: Me
country: US
province: California
locality: San Francisco
NotBefore: 2022-02-18T15:00:00+01:00
NotAfter: 2032-02-18T15:00:00+01:00
```

## Creating a certificate

After creating a CA you use it to create certificates that you can directly inject into your favourite web server:

```sh
certman create cert --csr-config test.yaml --ca-key ca.key --ca-cert ca.pem --out test.pem --privkey-out test.key
```

The given `test.yaml` looks similar to this example:

```yaml
apiVersion: v1
commonName: example.org
organization: Me
country: US
province: California
locality: San Francisco
```

The certificate has the following properties:

- Valid for 90 days
- Valid for "TLS Web Server Authentication" (this is the only x509 extended key usage)
- The subject alternative name (SAN) is set to the same value as the common name in the cnofiguration file.
