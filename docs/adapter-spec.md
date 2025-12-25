# Adapter Specification

## Overview

Adapters are declarative bundles that define how to interact with specific database systems for PK/FK discovery. Adapters are NOT arbitrary code - they contain only declarative files.

## Bundle Structure

An adapter bundle is a tarball/zip containing:

```
adapter-name/
├── manifest.yaml          # Adapter metadata
├── capabilities.yaml      # Feature flags
├── type_mapping.yaml      # DB type -> canonical type mappings
├── sql/                   # SQL templates
│   ├── list_tables.sql
│   ├── list_columns.sql
│   ├── list_indexes.sql
│   ├── list_constraints.sql
│   ├── profile_column_sample.sql
│   ├── profile_column_full.sql
│   ├── fk_inclusion_sample.sql
│   ├── fk_inclusion_exact.sql
│   └── explain.sql
├── tests/                 # Validation test definitions
│   └── *.yaml
├── bundle.json           # Checksums
└── signature             # HMAC signature
```

## Manifest Schema

```yaml
id: postgres-adapter
name: PostgreSQL Adapter
vendor: pkfk-discovery
db_family: postgresql
supported_db_versions:
  - "12+"
adapter_schema_version: "1.0"
engine_compat: ">=1.0"
maturity_level: L2
created_at: "2025-01-01T00:00:00Z"
updated_at: "2025-01-01T00:00:00Z"
```

## Maturity Levels

### L0: Metadata Introspection
- `list_tables.sql` passes
- `list_columns.sql` passes
- `list_indexes.sql` passes
- `list_constraints.sql` passes

### L1: Profiling
- All L0 tests pass
- `profile_column_sample.sql` passes
- `profile_column_full.sql` passes (optional)

### L2: FK Evidence
- All L1 tests pass
- `fk_inclusion_sample.sql` passes
- `fk_inclusion_exact.sql` passes

### L3: Performance Mode
- All L2 tests pass
- Sampling works correctly
- `explain.sql` integration (if supported)

### L4: Enterprise Mode
- All L3 tests pass
- Permissions matrix validated
- Large-table safeguards tested

## SQL Template Requirements

### Safety Rules

1. **SELECT/EXPLAIN only**: No INSERT, UPDATE, DELETE, DROP, ALTER, CREATE, TRUNCATE, GRANT, REVOKE, COPY, CALL, EXEC, MERGE, VACUUM, ANALYZE
2. **Exception**: EXPLAIN ANALYZE allowed if capability flag set and policy allows
3. **Single statement**: One statement per template
4. **Parameterized**: Use placeholders ($1, $2, etc. for Postgres, ? for MySQL)
5. **Sample limits**: Sample-mode templates MUST include LIMIT clause

### Template Parameters

Templates accept parameters:
- `schema`: Schema name (optional)
- `table`: Table name
- `column`: Column name
- `limit`: Row limit for sampling
- `offset`: Offset for pagination

### Example Template

```sql
-- list_tables.sql
SELECT 
    table_schema,
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema = COALESCE($1, 'public')
ORDER BY table_schema, table_name;
```

## Type Mapping

```yaml
# type_mapping.yaml
mappings:
  - db_type: "VARCHAR"
    canonical_type: "string"
    normalization:
      max_length: true
      nullable: true
  
  - db_type: "INTEGER"
    canonical_type: "integer"
    normalization:
      signed: true
      nullable: true
```

## Bundle Packaging

1. Collect all files
2. Generate `bundle.json` with SHA256 checksums
3. Generate HMAC signature using master key
4. Create tarball/zip
5. Upload to MinIO

## Signing

- **MVP**: HMAC-SHA256 signature
- **Future**: Ed25519 signature support

## Validation

Before publishing, adapters must pass all tests for the selected maturity level. Publishing is blocked if tests fail.

