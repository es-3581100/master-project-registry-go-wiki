# Project record schema

```json
{
  "id": "stable source/local identifier",
  "slug": "stable URL slug",
  "title": "display name",
  "family": "project family",
  "status": "active | incubation | paused | idea | retired | archived",
  "kind": "project | workspace | repository | ...",
  "summary": "short human-readable description",
  "path": "local path or source locator",
  "repoURL": "repository URL",
  "siteURL": "deployed site URL",
  "tags": ["..."],
  "notes": "longer notes",
  "source": "provenance of the registry record",
  "hiddenFromDeploy": false,
  "proprietary": false
}
```

`proprietary` is stronger than `hiddenFromDeploy`: both are excluded by the static exporter, but `proprietary` is intended as a semantic boundary, not merely a presentation toggle.
