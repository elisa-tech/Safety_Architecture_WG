# Gap Analysis: BASIL SPDX Export vs. SPDX 3.1-dev FunctionalSafety Profile

## Context

This document covers the gaps between BASIL's SPDX export and the **SPDX 3.1-dev** specification, specifically the FunctionalSafety profile. These gaps **cannot be addressed within SPDX 3.0.1** and require the 3.1 spec (currently at Release Candidate stage).

**Prerequisite**: The SPDX 3.0.1 compliance gaps documented in `spdx_functional_safety_gap_analysis.md` should be addressed first, as they provide the foundation for the 3.1-dev migration.

**Critical caveat**: The FunctionalSafety profile is being developed for **SPDX 3.1** in the `spdx/spdx-3-model` repo `develop` branch. The `Requirement` class itself is being added to **Core** (not FunctionalSafety), while supporting classes (`RequirementVerification`, `EvaluationResult`, `EvidenceRelationship`, `Assumption`) are in the FunctionalSafety profile namespace.

The file `mem-2026-05-08T150437.688Z.jsonld` is used as the reference BASIL export for this analysis — a SPDX 3.0.1 JSON-LD export containing 42 software requirements for `/drivers/char/mem.c`, 6 technical annex documents, and associated code snippets.

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

**Note**: The status of the `Requirement` class in SPDX 3.1 Core needs verification against the latest `spdx/spdx-3-model` develop branch before implementation. If it does not land, requirements remain as `software_File` with `primaryPurpose: "requirement"`.

## Gap 2: Requirement Status Property

**Current**: Requirement status (e.g., `NEW`) is buried inside serialized JSON in `Annotation.statement` fields, opaque to SPDX tooling.

**SPDX 3.1-dev target**: The `Requirement` element will carry a `status` property. The agreed status values for SPDX 3.1-dev are:

| Status | Description |
|---|---|
| `Draft` | WIP but not ready for review |
| `Review Ready` | Requirement ready to be reviewed |
| `Active` | Requirement has been approved by stakeholders and can be included in the traceability framework |
| `Retired` | Not applicable anymore |
| `Backlog` | Requirement has been captured and logged, but not yet become a formal requirement draft. May still be incomplete or ambiguous |

BASIL's current internal status values (`NEW`, `IN_PROGRESS`, `APPROVED`, `OBSOLETE`) will need to be mapped to these SPDX 3.1-dev values.

**Important**: The verification of code against applicable requirements using applicable verification methods is a system-level activity whose status is **not** tracked as part of the requirement itself. Verification status is tracked in the SBOM as part of the verification reports (see Gap 4).

## Gap 3: Safety-Specific Relationship Types

**Current**: BASIL uses general-purpose SPDX 3.0.1 relationship types. The FunctionalSafety profile adds safety-specific relationship types that provide proper traceability semantics:

| Relationship Type | Purpose | Currently Used? |
|---|---|---|
| `implementedBy` | Requirement → code artifacts | NO (direction is inverted: snippet → requirement via `hasRequirement`) |
| `verifiedBy` | Requirement → RequirementVerification | NO |
| `tracedToDetail` | Parent → child requirement decomposition | NO |
| `assumes` | Element → Assumption | NO |
| `conformsTo` | Element → Assumption or Specification | NO |
| `evaluatedOn` | Element → evaluated elements | NO |

**Impact**: The current `hasRequirement` relationships go from code snippet to requirement, but the FunctionalSafety profile's `implementedBy` goes from requirement to implementation. The direction matters for traceability: the safety perspective is "this requirement is implemented by that code," not "this code has this requirement."

**Target relationship map**:

```
Requirement --implementedBy--> Snippet (code)
Requirement --tracedToDetail--> Requirement (child)
Requirement --verifiedBy--> RequirementVerification
```

## Gap 4: No `RequirementVerification` Elements

**Current**: No verification data exists in the export.

**SPDX 3.1 provides**: `RequirementVerification` class with:
- `verificationMethod` (`VerificationType`: test, analysis, review, inspection, audit, assessment, demonstration, other)
- `verificationPrecondition`
- `verificationPostcondition`
- `verificationUID` (`ExternalIdentifier`)

BASIL's coverage data (`coverage: 20`, `srs_coverage: 600.0`, `tss_coverage: 0`, `tcs_coverage: 0`) partially addresses verification status but is currently buried in annotation JSON rather than expressed as proper verification elements.

BASIL's test specification to requirement mappings (`SwRequirementTestSpecification`) naturally map to `RequirementVerification` elements, with the verification method derived from the test specification type.

## Gap 5: No `EvaluationResult` Elements

**Current**: No evaluation/assessment data.

**SPDX 3.1 provides**: `EvaluationResult` class with:
- `evaluation` (`EvaluationResultType`: pass, fail, inconclusive) — **required**
- `evaluationBasedOn` (→ `RequirementVerification`) — **required**
- `rationale` — **required**

This would be the proper place to express BASIL's test run results. Each `TestRun` record in BASIL maps to an `EvaluationResult` element with pass/fail/inconclusive outcome.

## Gap 6: No `EvidenceRelationship` Elements

**Current**: No evidence traceability beyond basic `hasDocumentation` and `hasEvidence` relationships.

**SPDX 3.1 provides**: `EvidenceRelationship` (subclass of `Relationship`) with:
- `evidenceUID` (`ExternalIdentifier`)
- `evidenceCategory` (`EvidenceType`: report, log, recording, observation, other)

The technical annexes (currently `software_File` with `software_primaryPurpose: "documentation"`) could be linked to requirements/verifications via `EvidenceRelationship` with typed evidence categories, providing richer traceability than the current generic `hasDocumentation` relationship.

## Gap 7: No `Assumption` Elements

**Current**: No assumptions are captured.

**SPDX 3.1 provides**: `Assumption` class with:
- `assumptionStatement` — **required**
- `assumptionUID` (`ExternalIdentifier`)

For functional safety, capturing assumptions (e.g., "the kernel runs with `CONFIG_STRICT_DEVMEM` enabled") is essential for safety argumentation. This requires safety analysis input beyond what BASIL currently tracks.

## Gap 8: Missing Safety-Specific `ExternalIdentifier` Types

**Current**: No `ExternalIdentifier` elements (the 3.0.1 gap analysis addresses adding `ExternalIdentifier` with `externalIdentifierType: "other"`).

**SPDX 3.1 adds** three new `ExternalIdentifierType` entries for safety:
- `requirementUID` — linking to requirements management tools
- `verificationUID` — linking to verification tracking
- `evidenceUID` — linking to evidence management

Once the 3.0.1 gap is addressed (using `"other"` type), these identifiers should be upgraded to use the safety-specific types when migrating to 3.1.

## Gap 9: Spec Version and Context Alignment

**Current**: `"@context": "https://spdx.org/rdf/3.0.1/spdx-context.jsonld"` with `"specVersion": "3.0.1"`.

**Impact**: This context does not define `Requirement`, `RequirementVerification`, `EvaluationResult`, `EvidenceRelationship`, `Assumption`, or any of the FunctionalSafety vocabulary terms. When SPDX 3.1 is released, the context must be updated to `https://spdx.org/rdf/3.1/spdx-context.jsonld` and the spec version to `"3.1.0"`.

## Gap 10: Remaining Opaque Metadata in Annotations

After the 3.0.1 cleanup (documented in `spdx_functional_safety_gap_analysis.md`), some BASIL metadata will still remain in opaque JSON annotations because it has no SPDX 3.0.1 equivalent. SPDX 3.1-dev provides proper homes for this data:

| Metadata | Current Location | 3.1-dev Target |
|---|---|---|
| Requirement status | Annotation JSON | `Requirement.status` property (see Gap 2) |
| Coverage metrics (`srs_coverage`, etc.) | Annotation JSON | `EvaluationResult` elements (see Gap 5) |
| Verification method | Not tracked | `RequirementVerification.verificationMethod` (see Gap 4) |
| Test run pass/fail | Annotation JSON | `EvaluationResult.evaluation` (see Gap 5) |

## Summary: Migration Priority

| Priority | Gap | Effort | Description |
|---|---|---|---|
| **High** | #1: Use `Requirement` class | Medium | Replace `software_File` with first-class `Requirement` (blocked on 3.1 Core confirmation) |
| **High** | #3: Use `implementedBy` direction | Low | Reverse relationship direction for safety traceability |
| **High** | #9: Update spec version/context | Trivial | Switch to 3.1 context and version |
| **Medium** | #2: Requirement status property | Low | Map BASIL statuses to SPDX 3.1-dev values |
| **Medium** | #4: Add `RequirementVerification` | Medium | Express verification planning as first-class elements |
| **Medium** | #8: Safety-specific `ExternalIdentifier` types | Low | Upgrade from `"other"` to `requirementUID`/`verificationUID`/`evidenceUID` |
| **Low** | #5: Add `EvaluationResult` | Medium | Express test run outcomes |
| **Low** | #6: Add `EvidenceRelationship` | Medium | Typed evidence traceability |
| **Low** | #7: Add `Assumption` | Medium | Requires safety analysis input |
| **Low** | #10: Migrate remaining annotation metadata | Medium | Distribute across proper 3.1 elements |

## Conclusion

These gaps represent the delta between a conformant SPDX 3.0.1 export (after addressing the gaps in `spdx_functional_safety_gap_analysis.md`) and full compliance with the SPDX 3.1-dev FunctionalSafety profile. The highest-impact changes are adopting the `Requirement` class and reversing the relationship direction to `implementedBy`. All gaps are blocked on the SPDX 3.1 release. The phased migration plan in `spdx_31_migration_plan.md` provides the implementation roadmap.
