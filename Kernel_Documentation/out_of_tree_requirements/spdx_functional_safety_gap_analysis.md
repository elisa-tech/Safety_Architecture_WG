# Gap Analysis: BASIL SPDX Export vs. SPDX 3.1-dev FunctionalSafety Profile

## Context

The file `mem-2026-05-08T150437.688Z.jsonld` is a BASIL-generated SPDX 3.0.1 JSON-LD export containing 42 software requirements for `/drivers/char/mem.c` (functions: `memory_lseek`, `memory_open`, `mmap_mem_prepare`, `open_port`, `read_mem`, `write_mem`), 6 technical annex documents, and associated code snippets.

**Critical caveat**: The FunctionalSafety profile **does not exist** in the released SPDX 3.0 or 3.0.1 specification. It is being developed for **SPDX 3.1** (currently Release Candidate stage) in the `spdx/spdx-3-model` repo `develop` branch. The `Requirement` class itself is being added to **Core** (not FunctionalSafety), while supporting classes (`RequirementVerification`, `EvaluationResult`, `EvidenceRelationship`, `Assumption`) are in the FunctionalSafety profile namespace.

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

## Gap 1: Requirements Modeled as `software_File` Instead of `Requirement`

**Current**: Requirements are `software_File` elements with `software_primaryPurpose: "requirement"`:
```json
{
  "type": "software_File",
  "spdxId": "spdx:file:basil:software-requirement:3",
  "software_primaryPurpose": "requirement",
  "name": "MEMORY-LSEEK-001",
  "description": "The memory_lseek() function shall..."
}
```

**SPDX 3.1 target**: A first-class `Requirement` class (subclass of `Element`) with dedicated properties:

| Property | Type | Required | Currently Mapped? |
|---|---|---|---|
| `requirementStatement` | `xsd:string` | **YES** | NO (content is in `description`) |
| `rationale` | `xsd:string` | No | NO |
| `devLifecycleStage` | `LifecycleScopeType` | No | NO |
| `requirementUID` | `ExternalIdentifier` | No | NO (BASIL ID is in annotation JSON) |

**Impact**: The requirement text lives in the generic `description` field rather than the semantically meaningful `requirementStatement`. There is no way to distinguish a requirement's formal statement from its general description. The BASIL requirement ID (e.g., `id: 3`) and title (e.g., `MEMORY-LSEEK-001`) are not expressed using `ExternalIdentifier` with `externalIdentifierType: "requirementUID"`.

## Gap 2: Metadata Encoded as Serialized JSON in Annotations

**Current**: All BASIL-specific metadata (requirement status, version, coverage, authorship, mapping details, code sections) is serialized as JSON strings inside `Annotation.statement` fields with `annotationType: "other"`:

```json
{
  "type": "Annotation",
  "annotationType": "other",
  "statement": "{\"id\": 3, \"title\": \"MEMORY-LSEEK-001\", \"status\": \"NEW\", ...}"
}
```

**Problem**: This makes the data opaque to any SPDX-aware tooling. A consumer must parse these JSON strings with knowledge of BASIL's internal schema to extract:
- Requirement status (`NEW`, etc.)
- Coverage metrics (`coverage: 20`, `srs_coverage: 600.0`)
- Version tracking (`version: "1"`)
- Code section mappings (`section: "static loff_t memory_lseek..."`)
- BASIL internal table names (`__tablename__: "sw_requirements"`)

**SPDX 3.1 target**: This metadata should be distributed across proper SPDX elements and properties — requirement status could map to lifecycle annotations, coverage data could use `EvaluationResult`, and code traceability should use `implementedBy` relationships.

## Gap 3: No Safety-Specific Relationship Types

**Current**: Only three relationship types are used:
- `contains` (49x) — API file contains snippets
- `hasRequirement` (42x) — snippet has requirement
- `hasDocumentation` (7x) — API has documentation

**SPDX 3.1 FunctionalSafety adds**:

| Relationship Type | Purpose | Currently Used? |
|---|---|---|
| `implementedBy` | Requirement → code artifacts | NO (direction is inverted: snippet → requirement via `hasRequirement`) |
| `verifiedBy` | Requirement → RequirementVerification | NO |
| `tracedToDetail` | Parent → child requirement decomposition | NO |
| `assumes` | Element → Assumption | NO |
| `conformsTo` | Element → Assumption or Specification | NO |
| `evaluatedOn` | Element → evaluated elements | NO |

**Impact**: The current `hasRequirement` relationships go from code snippet to requirement, but the FunctionalSafety profile's `implementedBy` goes from requirement to implementation. The direction matters for traceability: the safety perspective is "this requirement is implemented by that code," not "this code has this requirement."

## Gap 4: No RequirementVerification Elements

**Current**: No verification data exists in the export.

**SPDX 3.1 provides**: `RequirementVerification` class with:
- `verificationMethod` (`VerificationType`: test, analysis, review, inspection, audit, assessment, demonstration, other)
- `verificationPrecondition`
- `verificationPostcondition`
- `verificationUID` (`ExternalIdentifier`)

BASIL's coverage data (`coverage: 20`, `srs_coverage: 600.0`, `tss_coverage: 0`, `tcs_coverage: 0`) partially addresses verification status but is currently buried in annotation JSON rather than expressed as proper verification elements.

## Gap 5: No EvaluationResult Elements

**Current**: No evaluation/assessment data.

**SPDX 3.1 provides**: `EvaluationResult` class with:
- `evaluation` (`EvaluationResultType`: pass, fail, inconclusive) — **required**
- `evaluationBasedOn` (→ `RequirementVerification`) — **required**
- `rationale` — **required**

This would be the proper place to express BASIL's coverage assessment results.

## Gap 6: No EvidenceRelationship Elements

**Current**: No evidence traceability.

**SPDX 3.1 provides**: `EvidenceRelationship` (subclass of `Relationship`) with:
- `evidenceUID` (`ExternalIdentifier`)
- `evidenceCategory` (`EvidenceType`: report, log, recording, observation, other)

The technical annexes (currently `software_File` with `software_primaryPurpose: "documentation"`) could be linked to requirements/verifications via `EvidenceRelationship`.

## Gap 7: No Assumption Elements

**Current**: No assumptions are captured.

**SPDX 3.1 provides**: `Assumption` class with:
- `assumptionStatement` — **required**
- `assumptionUID` (`ExternalIdentifier`)

For functional safety, capturing assumptions (e.g., "the kernel runs with `CONFIG_STRICT_DEVMEM` enabled") is essential for safety argumentation.

## Gap 8: Missing `ExternalIdentifier` with Safety-Specific Types

**Current**: No `ExternalIdentifier` elements.

**SPDX 3.1 adds** three new `ExternalIdentifierType` entries for safety:
- `requirementUID` — linking to requirements management tools
- `verificationUID` — linking to verification tracking
- `evidenceUID` — linking to evidence management

BASIL requirement IDs (e.g., `id: 3`, `title: "MEMORY-LSEEK-001"`) should be expressed as `ExternalIdentifier` objects with type `requirementUID`.

## Gap 9: Duplicate/Redundant Relationships

**Current**: For each snippet-to-requirement mapping, there is a `contains` relationship (API file → snippet) AND a separate `hasRequirement` relationship (snippet → requirement). When a snippet maps to multiple requirements, the `contains` relationship is repeated identically. For example, `spdx:relationship:2`, `4`, `6`, `8`, `10` are all identical `contains` from `spdx:file:basil:api:352` to `spdx:snippet:api:352:1`.

**SPDX best practice**: The `contains` relationship should appear once per snippet, and `hasRequirement` (or `implementedBy` in 3.1) should reference multiple targets in its `to` array rather than creating separate relationships per requirement.

## Gap 10: Spec Version and Context Alignment

**Current**: `"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld"` with `"specVersion": "3.0.1"`.

**Impact**: This context does not define `Requirement`, `RequirementVerification`, `EvaluationResult`, `EvidenceRelationship`, `Assumption`, or any of the FunctionalSafety vocabulary terms. When SPDX 3.1 is released, the context must be updated to `https://spdx.org/rdf/3.1/spdx-context.jsonld`.

## Summary: Migration Priority

| Priority | Gap | Effort | Blocked on SPDX 3.1? |
|---|---|---|---|
| **High** | #1: Use `Requirement` instead of `software_File` | Medium | Yes |
| **High** | #2: Eliminate opaque JSON-in-annotations | High | Partially (some data has no 3.1 home yet) |
| **High** | #3: Use `implementedBy` relationship direction | Low | Yes |
| **Medium** | #8: Use `ExternalIdentifier` with `requirementUID` | Low | Yes |
| **Medium** | #4: Add `RequirementVerification` elements | Medium | Yes |
| **Medium** | #9: Deduplicate relationships | Low | No — can fix now |
| **Low** | #5: Add `EvaluationResult` elements | Medium | Yes |
| **Low** | #6: Add `EvidenceRelationship` elements | Medium | Yes |
| **Low** | #7: Add `Assumption` elements | Medium | Yes (requires safety analysis input) |
| **Low** | #10: Update context/specVersion to 3.1 | Trivial | Yes |

## Conclusion

The current export is a valid SPDX 3.0.1 document that works around the absence of the FunctionalSafety profile by using `software_File` with `software_primaryPurpose: "requirement"` — which is the best available approach today. The main actionable gap *right now* (Gap #9: deduplicating relationships) doesn't require waiting for 3.1. All other significant gaps are blocked on the SPDX 3.1 release. When 3.1 lands, the biggest lift will be restructuring requirements as first-class `Requirement` elements and replacing annotation-embedded metadata with proper FunctionalSafety profile classes and relationships.
