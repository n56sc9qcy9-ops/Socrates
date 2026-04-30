# Knowledge Data

This directory organizes knowledge data by purpose.

## Layout

```
knowledge/
├── reference/       # Non-runtime source/reference material
│   ├── concepts.yaml
│   ├── forms.yaml
│   └── relations.yaml
└── README.md
```

## Active Runtime

Active runtime knowledge is embedded in `internal/knowledge/*.yaml` and loaded at build time.

**Do not edit `knowledge/reference/` thinking it is the active runtime.**

## Reference Material

Files in `knowledge/reference/` are source/reference/seed material.
They are not loaded at runtime and are kept for curation reference.

## Adding Knowledge

1. Edit `internal/knowledge/*.yaml` (active runtime)
2. Run `socrates knowledge validate` to check changes
3. Commit and push changes

The `knowledge/reference/` directory is for archive/seed material only.