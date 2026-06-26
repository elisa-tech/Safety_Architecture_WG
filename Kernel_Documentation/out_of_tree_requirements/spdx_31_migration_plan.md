# Incremental Plan: BASIL SPDX 3.1-dev FunctionalSafety Compliance

## Goal

Make BASIL import and export requirements, source code, test specifications, and test case implementations encoded as SPDX 3.1-dev, and display the relationships:

- Requirement <-> Requirement (decomposition / derivation)
- Requirement <-> Code (implementation traceability)
- Requirement <-> Test Specification (verification planning)
- Test Specification <-> Test Case Implementation (verification execution)

## Baseline

BASIL's **internal data model** already supports all four relationship types via mapping tables (`SwRequirementSwRequirementModel`, `ApiSwRequirementModel`, `SwRequirementTestSpecificationModel`, `TestSpecificationTestCaseModel`). The work is about aligning the **SPDX serialization layer** (`spdx_manager.py`) and the **import path** (`import_manager.py`) with SPDX 3.1-dev, and ensuring the **UI** surfaces these relationships clearly.

### Key BASIL files involved

| File | Role |
|---|---|
| `api/spdx_manager.py` | SPDX export (JSON-LD generation) |
| `api/import_manager.py` | SPDX import (JSON-LD parsing) |
| `db/models/*.py` | SQLAlchemy data models |
| `api/api.py` | REST API endpoints |
| `app/src/` | React/PatternFly frontend |

### SPDX 3.1-dev FunctionalSafety profile classes available

| Class | Inherits from | Key own properties |
|---|---|---|
| `RequirementVerification` | `Element` | `requirementVerificationMethod` (enum: test, analysis, inspection, demonstration, simulation, other), `verifiedElement` (Element ref), `verifyingElement` (Element ref[]) |
| `EvaluationResult` | `Element` | `evaluatedElement` (Element ref), `evaluationOutcome` (enum: pass, fail, inconclusive), `evaluationOutcomeDescription` (string) |
| `EvidenceRelationship` | `Relationship` | `evidenceType` (enum: assurance, confidence, independence, other), `evidenceSupportType` (enum: supports, refutes) |

### Open question: `Requirement` class in Core

The gap analysis states that a `Requirement` class is being added to SPDX 3.1 **Core** (not FunctionalSafety). The SPDX 3.1-RC1 schema confirms `RequirementVerification`, `EvaluationResult`, and `EvidenceRelationship` in the FunctionalSafety namespace, but the status of the `Requirement` Core class needs verification against the latest `spdx/spdx-3-model` develop branch before Phase 2 begins. If it does not land, requirements remain as `software_File` with `primaryPurpose: "requirement"` — which is already what BASIL does.

---

## Phase 0: Quick wins on the current SPDX 3.0.1 export

**Objective**: Fix issues that don't require SPDX 3.1 and reduce noise in the export.

**Changes to `spdx_manager.py`**:

1. **Deduplicate `contains` relationships (Gap #9)**: The same `contains` (API -> snippet) relationship is emitted once per requirement mapped to that snippet. Emit it once per snippet.

2. **Consolidate `hasRequirement` targets**: Instead of N separate `hasRequirement` relationships (one per requirement per snippet), use a single relationship per snippet with multiple targets in the `to` array.

3. **Add `ExternalIdentifier` to requirements**: Emit the BASIL requirement ID and title as `externalIdentifier` entries on each requirement element. Use `externalIdentifierType: "other"` (since `requirementUID` is a 3.1 addition) with `identifierLocator` pointing to the BASIL instance.

**Deliverable**: Cleaner, valid SPDX 3.0.1 export. No schema changes, no UI changes.

**Validation**: Run `spdx3-validate` on the output; compare element/relationship counts before and after.

---

## Phase 1: Upgrade SPDX context to 3.1-dev

**Objective**: Switch the export serialization to SPDX 3.1-dev so new classes and relationship types are valid.

**Steps**:

1. **Update `@context`** in `spdx_manager.py` from `https://spdx.org/rdf/3.0.1/spdx-context.jsonld` to the SPDX 3.1-dev context URL (or embed the 3.1-dev context inline if not yet published).

2. **Update `specVersion`** from `"3.0.1"` to `"3.1.0"` (or `"3.1-dev"`).

3. **Update `spdx3-validate` dependency** — check if a 3.1-compatible version exists, or pin to a dev version. If validation tooling doesn't yet support 3.1, document this and validate manually against the JSON schema from `spdx-3-model`.

4. **Verify the `Requirement` class status**: Check `spdx/spdx-3-model` develop branch. If `Requirement` exists in Core:
   - Add an `SPDXRequirement` class to `spdx_manager.py` with properties: `requirementStatement`, `rationale`, `externalIdentifier` (with `requirementUID` type).
   - Change requirement export from `software_File` to `Requirement`.
   
   If it does not exist, keep `software_File` with `primaryPurpose: "requirement"` and use `description` for the requirement statement.

5. **Update `ExternalIdentifier` types**: Change requirement identifiers from `"other"` to `"requirementUID"` (now valid in 3.1).

**Deliverable**: Export produces valid SPDX 3.1-dev JSON-LD. All existing data still exports correctly.

**Validation**: Validate against 3.1-dev JSON schema; round-trip test (export then re-import).

---

## Phase 2: Fix relationship directions and add safety-specific relationship types

**Objective**: Align relationship semantics with the FunctionalSafety profile's traceability model.

**Changes to `spdx_manager.py`**:

1. **Reverse `hasRequirement` to `implementedBy`**: Currently snippet -> requirement via `hasRequirement`. Change to requirement -> snippet via `implementedBy`. This matches the safety perspective: "this requirement is implemented by that code."

2. **Add `hasSpecification`**: For `SwRequirementTestSpecification` mappings, emit `hasSpecification` from requirement to test specification element.

3. **Add `hasTestCase`**: For `TestSpecificationTestCase` mappings, emit `hasTestCase` from test specification to test case element.

4. **Add `tracedToDetail`**: For `SwRequirementSwRequirement` mappings (parent-child requirement decomposition), emit `tracedToDetail` from parent requirement to child requirement.

5. **Preserve `contains`**: API -> snippet (container relationship, keep as-is).

6. **Preserve `hasDocumentation`**: API -> technical annex documents (keep as-is).

**Resulting relationship map**:

```
Requirement --implementedBy--> Snippet (code)
Requirement --tracedToDetail--> Requirement (child)
Requirement --hasSpecification--> TestSpecification
TestSpecification --hasTestCase--> TestCase
API --contains--> Snippet
API --hasDocumentation--> Document
```

**Changes to `import_manager.py`**:

- Parse `implementedBy` relationships and map back to `ApiSwRequirement` with snippet offsets.
- Parse `tracedToDetail` and map to `SwRequirementSwRequirement`.
- Parse `hasSpecification` and map to `SwRequirementTestSpecification`.
- Parse `hasTestCase` and map to `TestSpecificationTestCase`.

**Deliverable**: Export and import handle all four relationship types. Existing BASIL data round-trips correctly.

**Validation**: Export the `mem.c` dataset, re-import it into a fresh BASIL instance, verify all mappings survive.

---

## Phase 3: Add `RequirementVerification` elements

**Objective**: Express verification planning as first-class SPDX elements, linking requirements to the test specifications and test cases that verify them.

**Changes to `spdx_manager.py`**:

1. **Add `SPDXRequirementVerification` class** with properties:
   - `requirementVerificationMethod`: map from BASIL test spec type. Default to `"test"` for test-based specs, `"analysis"` for analysis-based, `"inspection"` for review-based.
   - `verifiedElement`: reference to the requirement being verified.
   - `verifyingElement`: array of references to test specification and/or test case elements.

2. **Emit one `RequirementVerification` per requirement-to-test-spec mapping**: For each `SwRequirementTestSpecification` row, create a `RequirementVerification` element.

3. **Link via `verifiedBy` relationship**: Emit `verifiedBy` from the requirement to the `RequirementVerification` element.

**Changes to data model** (optional, for richer metadata):

- Add a `verification_method` field to `SwRequirementTestSpecificationModel` (enum: test, analysis, inspection, demonstration, simulation, other). Default: `"test"`.

**Changes to `import_manager.py`**:

- Parse `RequirementVerification` elements and create corresponding `SwRequirementTestSpecification` mappings.
- Extract `requirementVerificationMethod` and store it.

**Deliverable**: Export includes `RequirementVerification` elements linking requirements to their verification artifacts.

---

## Phase 4: Add `EvaluationResult` elements for test runs

**Objective**: Express test execution outcomes as SPDX evaluation results.

**Changes to `spdx_manager.py`**:

1. **Add `SPDXEvaluationResult` class** with properties:
   - `evaluatedElement`: reference to the `RequirementVerification` element being evaluated.
   - `evaluationOutcome`: map from BASIL test run result — `"pass"`, `"fail"`, or `"inconclusive"`.
   - `evaluationOutcomeDescription`: test run summary or log excerpt.

2. **Emit one `EvaluationResult` per test run**: For each `TestRunModel` row (up to the configurable limit), create an `EvaluationResult` element.

3. **Link via relationship**: Emit relationship from the `RequirementVerification` to the `EvaluationResult`.

**Changes to `import_manager.py`**:

- Parse `EvaluationResult` elements and create corresponding `TestRun` records.
- Map `evaluationOutcome` back to BASIL's test result status.

**Deliverable**: Test execution history is represented as proper FunctionalSafety evaluation results in the SPDX export.

---

## Phase 5: Add `EvidenceRelationship` elements

**Objective**: Link evidence artifacts (technical annexes, test reports, justifications) to requirements and verifications using typed evidence relationships.

**Changes to `spdx_manager.py`**:

1. **Add `SPDXEvidenceRelationship` class** with properties:
   - `evidenceType`: `"assurance"` for technical annexes, `"confidence"` for test evidence, `"other"` for justifications.
   - `evidenceSupportType`: `"supports"` (default) or `"refutes"`.
   - Standard relationship fields: `from`, `to`, `relationshipType`.

2. **Replace `hasDocumentation` with `EvidenceRelationship`** for technical annex -> requirement links:
   - `from`: requirement, `to`: technical annex document.
   - `relationshipType`: `hasEvidence`.
   - `evidenceType`: `"assurance"`.
   - `evidenceSupportType`: `"supports"`.

3. **Add evidence links for justifications**: For each `ApiJustification`, create an `EvidenceRelationship` linking the justification to the relevant requirement.

**Changes to `import_manager.py`**:

- Parse `EvidenceRelationship` elements and reconstruct document/justification mappings.

**Deliverable**: Full evidence chain from requirements through verification to evidence artifacts.

---

## Phase 6: Migrate metadata out of opaque JSON annotations

**Objective**: Eliminate the pattern of encoding all BASIL metadata as serialized JSON in `Annotation.statement` fields (Gap #2).

**Changes to `spdx_manager.py`**:

1. **Requirement metadata** — distribute across proper fields:
   - `title` -> `name` (already done)
   - `description` -> `requirementStatement` (if `Requirement` class available) or `description`
   - `status` -> `comment` or a lifecycle annotation
   - `id` / `title` -> `externalIdentifier` (done in Phase 1)
   - `coverage` -> captured by `RequirementVerification` completeness (done in Phase 3)

2. **Snippet metadata** — distribute:
   - `section` / `offset` / `length` -> already in `software_Snippet` byte range properties
   - `coverage` -> relationship completeness field

3. **Keep a minimal annotation** for BASIL-specific data that has no SPDX equivalent (e.g., `__tablename__`, internal `version` counter, `srs_coverage` percentages). Mark these clearly as BASIL-extension annotations.

4. **Version the annotation schema**: Add `"basil:annotationVersion": "2.0"` so consumers can distinguish old-format from new-format annotations.

**Deliverable**: SPDX-aware tooling can read requirement metadata without parsing opaque JSON. BASIL-specific extensions are clearly marked.

---

## Phase 7: BASIL import of external SPDX 3.1-dev documents

**Objective**: BASIL can consume SPDX 3.1-dev documents produced by other tools (e.g., StrictDoc, Polarion, Doors NG) and populate its data model.

**Changes to `import_manager.py`**:

1. **Detect SPDX version**: Check `specVersion` and `@context` to dispatch to 3.0.1 or 3.1 parsing logic.

2. **Parse `Requirement` elements** (or `software_File` with `primaryPurpose: "requirement"`):
   - Extract `requirementStatement` -> BASIL `description`
   - Extract `name` -> BASIL `title`
   - Extract `externalIdentifier` with type `requirementUID` -> BASIL `id`/`title`

3. **Parse relationship graph**:
   - `implementedBy` -> create `ApiSwRequirement` mapping (need to resolve snippet to API + offset)
   - `tracedToDetail` -> create `SwRequirementSwRequirement` mapping
   - `hasSpecification` -> create `SwRequirementTestSpecification` mapping
   - `hasTestCase` -> create `TestSpecificationTestCase` mapping
   - `verifiedBy` -> associate `RequirementVerification` with requirement

4. **Parse FunctionalSafety elements**:
   - `RequirementVerification` -> enrich test spec mapping with verification method
   - `EvaluationResult` -> create `TestRun` record with pass/fail/inconclusive
   - `EvidenceRelationship` -> create document or justification mapping

5. **Conflict resolution**: When importing into an existing BASIL instance, match by `externalIdentifier` / `spdxId` and offer merge/overwrite/skip options.

**Deliverable**: BASIL can ingest SPDX 3.1-dev documents and reconstruct the full traceability model.

**Validation**: Create a hand-crafted SPDX 3.1-dev document with all element types, import it, verify all data appears correctly in BASIL.

---

## Phase 8: UI — Relationship visualization

**Objective**: BASIL's frontend clearly displays all four relationship types.

**UI changes** (`app/src/`):

1. **Requirement decomposition view**: Show parent-child requirement hierarchy (from `tracedToDetail`). Render as an expandable tree or indented list within the existing requirement panel.

2. **Requirement-to-code view**: Already partially exists (snippet display). Enhance to:
   - Show the `implementedBy` direction explicitly.
   - Display completeness (complete/incomplete/noAssertion) from the relationship.
   - Link to source file with line numbers.

3. **Requirement-to-test-specification view**: Show linked test specifications per requirement. Display:
   - Verification method (test/analysis/inspection/...).
   - Verification status (from `EvaluationResult`: pass/fail/inconclusive/not-yet-run).
   - Coverage completeness.

4. **Test-specification-to-test-case view**: Show linked test case implementations per test specification. Display:
   - Repository and file path.
   - Latest test run result.
   - Link to test execution artifacts.

5. **Traceability matrix**: A cross-reference view showing requirements on one axis and their linked artifacts (code, test specs, test cases) on the other. Color-code by verification status. This could extend the existing "dynamic view" in BASIL.

6. **Graph visualization update**: Update the Graphviz DOT export (`generate_diagraph`) to reflect the new relationship types and color-code by element type (requirement, code, test spec, test case, evidence).

**Deliverable**: Users can navigate from any requirement to its code, test specs, test cases, and evidence — and vice versa.

---

## Phase 9: Validation and conformance testing

**Objective**: Ensure BASIL's SPDX 3.1-dev output is conformant and interoperable.

1. **Automated validation in CI**: Run the SPDX 3.1 JSON schema validator on every export as part of BASIL's test suite (extends existing `test_spdx_api_validation.py`).

2. **Round-trip test**: Export from BASIL, import into a second BASIL instance, export again, diff the two exports. They should be semantically equivalent.

3. **Cross-tool interop test**: Import a document from another SPDX 3.1 producer (e.g., StrictDoc) and verify BASIL displays the relationships correctly.

4. **Conformance profile declaration**: Include a `conformsTo` relationship in the `SpdxDocument` pointing to the FunctionalSafety profile IRI, asserting that the document uses the profile.

---

## Dependency graph

```
Phase 0 (dedup/identifiers)
    |
Phase 1 (3.1-dev context + Requirement class check)
    |
Phase 2 (relationship directions + all 4 relationship types)
   / \
  /   \
Phase 3 (RequirementVerification)    Phase 6 (metadata cleanup)
  |                                      |
Phase 4 (EvaluationResult)              |
  |                                      |
Phase 5 (EvidenceRelationship)          |
   \                                   /
    \                                 /
     Phase 7 (SPDX 3.1 import)
         |
     Phase 8 (UI visualization)
         |
     Phase 9 (validation & conformance)
```

Phases 3-5 and Phase 6 can proceed in parallel once Phase 2 is complete.

## Risk: SPDX 3.1 specification stability

The FunctionalSafety profile is at RC stage. Class names, property names, or relationship semantics may change before final release. Mitigation:

- **Track the `spdx/spdx-3-model` develop branch** weekly during implementation.
- **Isolate SPDX-specific code** in `spdx_manager.py` and `import_manager.py` so schema changes are localized.
- **Use feature flags** in BASIL to toggle between 3.0.1 and 3.1-dev export until 3.1 is finalized.
- **Engage with the SPDX FunctionalSafety working group** to influence the spec and get early notice of breaking changes.
