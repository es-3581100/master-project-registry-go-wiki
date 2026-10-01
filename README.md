# Master Project Registry

A local-first Go wiki for managing a growing project ecosystem and publishing a deliberately reduced, read-only GitHub Pages projection.

Copyright © 2026 Eric T. Sawtelle.

## Design

```text
LOCAL AUTHORITATIVE STATE
  data/projects.json
        │
        ├── Go local wiki + /admin (writable)
        │
        └── static exporter
              │
              ├── omits proprietary records
              ├── omits hiddenFromDeploy records
              ├── generates project wiki pages
              └── adds SHA-256 client-side access gate
                    ↓
                 GitHub Pages
```

The project intentionally follows the catalog → generated-static-site pattern used successfully in `master-ash-catalog`, while the UI uses its own pixel-sunset identity.

## Run locally

```bash
./scripts/dev.sh
# http://127.0.0.1:8788
# admin: http://127.0.0.1:8788/admin
```

The local server binds to loopback by default. Local password gating is **off** by default and can be enabled in `config.json`.

To make a local password hash:

```bash
go run ./cmd/registry hash --value 'your-password'
```

Paste the result into `localPasswordSHA256` and set `localAuthEnabled` to `true`.

## Admin capabilities

The local-only `/admin` page can:

- create records;
- edit title, family, status, kind, paths, URLs, tags, notes, and source;
- retire a project;
- toggle `hiddenFromDeploy`;
- mark a record proprietary.

Writes are atomic: the server writes `data/projects.json.tmp` and renames it over the authoritative JSON only after serialization succeeds.

## Static GitHub Pages export

The deployed build has the SHA-256 access gate **enabled by default**.

```bash
export REGISTRY_DEPLOY_PASSWORD='choose-a-password'
./scripts/export.sh
```

or:

```bash
go run ./cmd/registry export --out docs --password 'choose-a-password'
```

The generated site contains only records that are both:

```text
proprietary == false
hiddenFromDeploy == false
```

The static build never contains `/admin`.

### Important static-auth boundary

GitHub Pages is public static hosting. A SHA-256 browser password gate is a **casual access barrier, not a confidentiality boundary**: its hash and generated public files are downloadable by anyone who knows how to inspect a static site. Do not rely on it to protect secrets.

That is why this project enforces the more important rule at export time: private/proprietary/hidden records are not written into the deployment at all.

## GitHub Actions

`.github/workflows/pages.yml` tests and exports the site, then deploys it with GitHub Pages. Add this repository secret before enabling Pages:

```text
REGISTRY_DEPLOY_PASSWORD
```

The workflow refuses to deploy if the secret is absent.

## Seed inventory

`data/projects.json` is seeded from the supplied ChatGPT project snapshot plus the supplied GitHub repository snapshots. `data/source-manifest.json` records SHA-256 hashes of the repository snapshot inputs used for the initial import.

`midas-matrix-training` is seeded as proprietary + hidden from deployment.

## UI identity

- uploaded pixel-art sunset used essentially unchanged as the fixed page background;
- deep plum/pink/coral/sunset-cream/teal palette sampled conceptually from the image;
- Offworld-inspired typography rules: book serif for prose, monospace for labels/data/navigation;
- hard square corners and translucent glass panels;
- compact structural eyebrows and restrained status color.

The supplied Offworld style note explicitly separates readable serif text from mono UI/data text and uses zero-radius components; this implementation preserves those text-role and geometry rules without copying the original page wholesale.

## Data contract

Each project is a JSON object with stable ID + slug and fields for family, status, kind, summary, path, repo/site links, tags, notes, provenance source, deploy visibility, and proprietary status.

Status is intentionally editable rather than inferred. Suggested values are:

```text
active
incubation
paused
idea
retired
archived
```

## License

No source-code license has been selected by this starter. Add the license you intend before publishing the source repository. Project records and linked repositories retain their own licenses and rights.
