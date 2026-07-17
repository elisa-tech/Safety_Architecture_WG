# Gap Analysis: BASIL SPDX Export vs. SPDX 3.0.1

## Context

The file `mem-2026-05-08T150437.688Z.jsonld` is a BASIL-generated SPDX 3.0.1 JSON-LD export containing 42 software requirements for `/drivers/char/mem.c` (functions: `memory_lseek`, `memory_open`, `mmap_mem_prepare`, `open_port`, `read_mem`, `write_mem`), 6 technical annex documents, and associated code snippets.

This document covers gaps between BASIL's current export and the **SPDX 3.0.1 specification** — i.e., issues that can be fixed today without waiting for SPDX 3.1. A separate document (`spdx_31_functional_safety_gap_analysis.md`) covers the remaining gaps that require the SPDX 3.1-dev FunctionalSafety profile.

## Current BASIL Export Structure Summary

| Element Type | Count | Used For |
|---|---|---|
| `software_File` | 51 | Requirements (42), documents (6), API ref, library root |
| `software_Snippet` | 6 | Code sections mapped to requirements |
| `Annotation` (type: `other`) | 57 | Embedding all BASIL metadata as serialized JSON |
| `Relationship` | 98 | `contains` (49), `hasRequirement` (42), `hasDocumentation` (7) |
| `CreationInfo` | 57 | Per-element creation metadata |
| `Person` | 2 | Authors (`gpaoloni`, `lpellecc`) |
| `SpdxDocument` | 1 | Root document |

## Gap 1: Duplicate/Redundant Relationships

**Current**: For each snippet-to-requirement mapping, there is a `contains` relationship (API file → snippet) AND a separate `hasRequirement` relationship (snippet → requirement). When a snippet maps to multiple requirements, the `contains` relationship is repeated identically. For example, `spdx:relationship:2`, `4`, `6`, `8`, `10` are all identical `contains` from `spdx:file:basil:api:352` to `spdx:snippet:api:352:1`.

**SPDX 3.0.1 best practice**: The `contains` relationship should appear once per snippet, and `hasRequirement` should reference multiple targets in its `to` array rather than creating separate relationships per requirement.

**Impact**: The export produces inflated, redundant JSON-LD that obscures the actual traceability structure. Any SPDX-aware tooling processing the output must deduplicate relationships to build an accurate graph.

**Fix**: Deduplicate `contains` relationships (one per snippet) and consolidate `hasRequirement` targets into a single relationship per snippet with multiple entries in the `to` array.

## Gap 2: Invalid and Missing SPDX 3.0.1 Relationship Types

**Current**: BASIL's export uses only a subset of the available SPDX 3.0.1 relationship types, and one type used is invalid:

### 2a: Invalid `hasDocument` relationship type

The export uses `hasDocument` for document-to-nested-document relationships (in `spdx_manager.py` line 1093). However, `hasDocument` is **not a valid SPDX 3.0.1 relationship type**. This should be `hasDocumentation` (which is used elsewhere for similar document relationships) or `contains`.

### 2b: Missing `describes` relationship

The standard SPDX practice is for an `SpdxDocument` to use a `describes` relationship to point to its root elements. BASIL's export uses `rootElement` as a property on the document but never creates an explicit `describes` relationship. Adding this would improve conformance with SPDX tooling expectations.

### 2c: Missing `hasTestCase` relationship

BASIL currently uses the generic `hasTest` for all relationships pointing to test cases (from snippets, from software requirements, and from test specifications). SPDX 3.0.1 provides the more specific `hasTestCase` which would be semantically more precise when the target is specifically a test case implementation rather than a general test artifact.

### 2d: Missing `generates` relationship

BASIL uses `hasEvidence` for test case → test run relationships. The `generates` relationship type (available in SPDX 3.0.1) would be semantically more accurate: a test case execution *generates* a test run result.

**Summary of relationship type changes**:

| Current | Should Be | Reason |
|---|---|---|
| `hasDocument` | `hasDocumentation` or `contains` | `hasDocument` is not a valid SPDX 3.0.1 type |
| *(missing)* | `describes` | Standard SpdxDocument → root element practice |
| `hasTest` (for test cases) | `hasTestCase` | More semantically precise |
| `hasEvidence` (for test runs) | `generates` | Test case generates test run results |

## Gap 3: BASIL Identifiers Not Expressed as `ExternalIdentifier`

**Current**: No `ExternalIdentifier` elements are used in the export. BASIL requirement IDs (e.g., `id: 3`, `title: "MEMORY-LSEEK-001"`) and other entity identifiers are buried inside serialized JSON in `Annotation.statement` fields:

```json
{
  "type": "Annotation",
  "annotationType": "other",
  "statement": "{\"id\": 3, \"title\": \"MEMORY-LSEEK-001\", \"status\": \"NEW\", ...}"
}
```

**SPDX 3.0.1 provides**: `ExternalIdentifier` with `externalIdentifierType`. While the safety-specific types (`requirementUID`, `verificationUID`, `evidenceUID`) are SPDX 3.1 additions, the `ExternalIdentifier` mechanism itself exists in 3.0.1 with `externalIdentifierType: "other"`.

**Fix**: Emit `ExternalIdentifier` entries on each requirement element using `externalIdentifierType: "other"` and `identifierLocator` pointing to the BASIL instance. This makes requirement IDs and titles accessible to standard SPDX tooling without parsing opaque JSON annotations.

**Example target**:
```json
{
  "type": "software_File",
  "spdxId": "spdx:file:basil:software-requirement:3",
  "software_primaryPurpose": "requirement",
  "name": "MEMORY-LSEEK-001",
  "description": "The memory_lseek() function shall...",
  "externalIdentifier": [
    {
      "type": "ExternalIdentifier",
      "externalIdentifierType": "other",
      "identifier": "3",
      "comment": "BASIL requirement ID"
    },
    {
      "type": "ExternalIdentifier",
      "externalIdentifierType": "other",
      "identifier": "MEMORY-LSEEK-001",
      "comment": "BASIL requirement title"
    }
  ]
}
```

## Gap 4: Opaque Metadata in Annotations (3.0.1-Addressable Portion)

**Current**: All BASIL-specific metadata is serialized as JSON strings inside `Annotation.statement` fields with `annotationType: "other"`. This makes the data opaque to any SPDX-aware tooling.

While the full resolution of this issue requires SPDX 3.1-dev safety-specific elements (see `spdx_31_functional_safety_gap_analysis.md`), some metadata can be moved to proper SPDX 3.0.1 constructs now:

| Metadata | Current Location | 3.0.1 Target |
|---|---|---|
| Requirement ID / title | Annotation JSON | `ExternalIdentifier` (see Gap 3) |
| Code section / offset / length | Annotation JSON | Already in `software_Snippet` byte range properties (redundant) |
| Coverage percentage | Annotation JSON | `Relationship.completeness` (`complete`/`incomplete`/`noAssertion`) — already emitted |
| BASIL internal fields (`__tablename__`, `version`) | Annotation JSON | Keep as BASIL-extension annotation, but version the annotation schema |

**Fix**:
1. Move identifiers to `ExternalIdentifier` (Gap 3)
2. Remove redundant code section data from annotations (already in snippet properties)
3. Add `"basil:annotationVersion": "2.0"` to remaining annotations so consumers can distinguish old-format from new-format

**Note**: Requirement status, safety-specific coverage metrics, and lifecycle metadata will be addressed in the SPDX 3.1-dev gap analysis, where proper elements (`Requirement` with status property, `EvaluationResult`, `RequirementVerification`) become available.

## Summary: Migration Priority

| Priority | Gap | Effort | Description |
|---|---|---|---|
| **High** | #1: Deduplicate relationships | Low | Dedup `contains`, consolidate `hasRequirement` targets |
| **High** | #2a: Fix invalid `hasDocument` | Low | Replace with `hasDocumentation` or `contains` |
| **Medium** | #2b: Add `describes` relationship | Low | SpdxDocument → root elements |
| **Medium** | #2c: Use `hasTestCase` | Low | More precise than `hasTest` for test cases |
| **Medium** | #2d: Use `generates` | Low | More precise than `hasEvidence` for test runs |
| **Medium** | #3: Add `ExternalIdentifier` | Low | Expose BASIL IDs to SPDX tooling |
| **Low** | #4: Clean up annotations | Medium | Remove redundant data, version the schema |

## Conclusion

These gaps represent improvements that can be made to BASIL's SPDX export within the current 3.0.1 specification — no spec changes or 3.1 features are required. The highest-priority items (relationship deduplication and fixing the invalid `hasDocument` type) are straightforward code fixes in `spdx_manager.py`. Once these are addressed, the export will be a cleaner, more conformant SPDX 3.0.1 document, providing a solid foundation for the SPDX 3.1-dev FunctionalSafety profile migration documented in `spdx_31_functional_safety_gap_analysis.md`.
