# Security boundary

The local Go server is the authoritative editor and binds to `127.0.0.1` by default.

The GitHub Pages build is public static hosting. Its SHA-256 password gate is designed to stop casual browsing, not to provide strong secrecy. The verifier hash is necessarily shipped to the browser and can be inspected or attacked offline.

Therefore the real deployment boundary is the exporter:

- `proprietary: true` => never exported;
- `hiddenFromDeploy: true` => never exported;
- `/admin` => local server only and never generated into static output.

Never put credentials, private weights, proprietary designs, private corpora, secrets, or unreleased material into a record that is eligible for static export.
