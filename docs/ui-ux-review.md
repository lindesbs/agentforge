# AgentForge UI/UX review and implementation contract

## Source and scope

Applied upstream release **2.1.1**, revision
`74308ad0beac6d880d87aaf5660b5eaf3bd8afdd`:

- [ux-ui-audit](https://github.com/aditya-ariosity/ux-ui-skills/blob/74308ad0beac6d880d87aaf5660b5eaf3bd8afdd/ux-ui-audit/SKILL.md)
- [design-system-review](https://github.com/aditya-ariosity/ux-ui-skills/blob/74308ad0beac6d880d87aaf5660b5eaf3bd8afdd/design-system-review/SKILL.md)

The first skill supplied the journey/evidence method, severity calibration,
recovery requirements and rendered accessibility checks. The second supplied
semantic token roles, consistent control states and explicit component contracts.
The experimental end-to-end skill and AI-execution patterns were not needed for
this existing local configuration editor. No upstream executable code is loaded
into the app, and no global skill installation is required by the application.

Mode: journey audit with local release checks; system scope: this app's controls
and stylesheet, an emerging single-product design system. Evidence consists of
source inspection, a pre-change browser reproduction, rendered screenshots and
automated interactions using a test-only Wails bridge. There is no supplied user
research, analytics, brand brief or Figma library. Improved task completion rates
or user comprehension have not been measured.

## Critical journey

`Open project → choose configuration → inspect/edit → review diff → explicit save → saved state`

| Stage | Observed issue before the change | Implemented behavior | Verification |
| --- | --- | --- | --- |
| Choose configuration | The entire learning-entry form preceded project results and the editor | Project workspace opens first; learning library has separate navigation | Rendered desktop/narrow views and navigation tests |
| Edit or change files | Selecting another file immediately replaced an unsaved native draft | Native dialog defaults to keeping work; Escape cancels; confirmed replacement alone discards | Baseline reproduction, then file-switch/reload/template regression tests |
| Failed project open | Editor state was cleared before inspection completed | Commit new project state only after inspection succeeds | Failing inspection preserves the draft |
| Edit and save | Raw textarea had no associated label; selected file and draft state were not explicit | Visible text label, selected-file state and unsaved badge | DOM/keyboard checks, automated accessibility scan |
| Review and save | An in-flight operation did not consistently lock editing; structured changes could leave a stale preview | Busy controls cannot mutate a submitted draft; edits invalidate review; staged fields must be applied | Validation, staged-field and delayed-save tests |

### S1 · High confidence · Draft replacement

**Evidence:** The baseline browser sequence opened `AGENTS.md`, changed its
textarea, then selected `.codex/config.toml`. The new contents replaced the draft
without confirmation. This was reproduced before editing the implementation.

**Impact:** Unsaved work could be lost during routine navigation. **Cause:**
`clearEditor()` ran before loading the destination, with no dirty-state guard.
**Correction:** Preserve context until loading succeeds and obtain explicit
confirmation before replacing native or structured drafts. **Acceptance:**
Cancel/Escape retain exact values; failed inspection keeps the current document;
confirmed discard changes only UI state, not a saved file.

### S2 · High confidence · Primary task hierarchy

**Evidence:** The pre-change full-page rendering placed the learning library
between the project picker and the file list/editor. **Impact:** Reaching the
primary configuration task required scrolling past a separate workflow; the
effect on users' speed is inferred, not measured. **Cause:** All features shared
one vertical page. **Correction:** Distinct workspace/library views; file list
beside the editor when content fits; templates/details use disclosures.
**Acceptance:** The workspace is the initial view and navigation preserves both
editor and library form state.

### S2 · High confidence · Control semantics and feedback

**Evidence:** The old native textarea had no label, status was far above the edit
area, and CSS defined repeated literal colors without a shared focus contract.
**Impact:** The editor lacked a programmatic field name and consistent task-local
feedback. **Correction:** Associated labels, status/alert regions, selected and
unsaved states, semantic tokens and explicit focus treatment. **Acceptance:**
Keyboard users can open, edit, preview and save; tested views expose names and
states; errors preserve inputs; contrast assertions and axe checks pass.

Preserved strengths: native files remain authoritative; the diff and explicit save
are mandatory; provider-detail parsing stays read-only.

## Design-system decisions

The existing dark, low-glare code-workspace direction is retained. It is a product
choice, not a universal claim about developer preferences. The primary content is
the native document; library forms and reusable templates are secondary tasks.

| Layer | Evidence and decision | Authority / maintenance |
| --- | --- | --- |
| Foundations | Literal repeated colors consolidated into semantic canvas, surface, text, action, focus and status tokens | `frontend/src/style.css`; repository frontend maintainers |
| Controls | Native buttons, inputs, selects, details and dialog; consistent labels, focus, pending and disabled states | Runtime Vue components and UI tests |
| Page structure | Workspace/library navigation plus file/editor layout, with content-driven stacking | `App.vue`; widths verified below |
| Documentation | Contracts and tests travel with behavior changes | This document and `tests/workspace.spec.ts` |
| Design assets | No separate asset library was supplied | Code is the current source of truth; no design/code parity claim |

Shared spacing roles use a small 4/8/12/16/24/32px scale. Control minimum height is
42px (36px for secondary compact controls). These are implementation choices,
not asserted standards. Breakpoints at 1150, 900, 700 and 480px adapt the available
space: navigation becomes horizontal, then file list/editor stack; small forms
become single-column. These choices are tested with actual content rather than
treated as device categories. Long paths wrap, and code remains editable.

Semantic tokens are the supported extension point. Add a new role only for a
distinct shared meaning; avoid arbitrary per-page color overrides. A single
product currently needs no multi-brand token package, governance committee or
separate component framework. Future UI changes should update the contract and
its tests together. The forced-colors override and reduced-motion treatment are
implemented, but not a complete assistive-technology test claim.

## Interaction contracts

| Pattern | Required behavior |
| --- | --- |
| Navigation | Keep drafts while switching workspace/library; communicate current view and selected file |
| Discard dialog | Show only when an action replaces dirty work; default focus on Keep editing; Tab/Shift+Tab remain inside; Escape cancels; retain work until destination load succeeds |
| Native editor | Visible label; provider summary describes loaded content; edits clear stale previews; lock mutation during pending operations |
| Structured fields | Distinguish unapplied fields from the native draft; applying them creates a preview; later changes invalidate it |
| Review/save | Show addition/removal markers as well as colors; re-preview before save; validation errors block saving; pending save cannot submit twice |
| Failure/success | Contextual text and alert/status semantics; failures retain data; successful save restores focus to the editor heading |
| Library | Required labels and native form validation; loading/filter empty states; local saves never change project files |

Browser `beforeunload` is a best-effort additional guard. Native window-manager
close/crash recovery and persistent draft autosave have not been implemented or
claimed. Choosing Browse still uses Wails' native directory chooser; the browser
fixture does not test the OS chooser itself.

## Validation and limitations

Run from `frontend`:

```sh
npm ci
npx playwright install --with-deps chromium
npm run test:ui
```

The checked suite has 18 tests covering the keyboard journey, discard/recovery,
structured drafts, failure preservation, validation, delayed saving, templates,
library state, filtering/empty results, accessibility and responsive layouts.
Playwright uses `tests/fixtures.ts` exclusively in tests; production never loads
its synthetic projects or bridge. Backend behavior remains covered by real Go
tests. CI is configured to run the UI suite and retain failed-test evidence;
remote CI execution is not established by local success.

Local browser validation used Chromium 154 on Debian 13, viewports 1280×900,
1100×900, 768×900 and 320×900, plus doubled base text size and long project paths.
`PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH` optionally selects a preinstalled browser;
otherwise Playwright uses its installed matching browser. Screenshots/traces are
written under ignored `build/ui-tests`. The pre-change capture is retained in
this workspace at `build/design-evidence/before-desktop.png`.

Text tokens were measured against actual declared surfaces (at least 4.5:1 for
the asserted text pairs); focus/control boundaries are checked at 3:1 or higher.
Automated axe scans cover the loaded-file view, discard dialog and learning
view. During validation, the file pane was corrected from an inappropriate
nested complementary landmark to file navigation, and the dialog received
explicit keyboard focus cycling.

No WCAG conformance claim is made. Remaining verification includes screen-reader
testing with a supported native WebKit/AT pair, OS chooser/close behavior, full
native end-to-end saves, native high-contrast mode, browser-level 400% zoom,
non-Latin/RTL content and representative user task testing. These are distinct
from the completed fixture-based interaction and reflow checks.

## Roadmap continuation: templates and definition map

The same two stable skills guided the new views. Other upstream specialists
(analytical dashboards, AI execution experiences, portfolio case studies and
design handoff) do not match this local configuration manager; the experimental
orchestrator is not used implicitly. No new visual identity or component framework
was needed. Existing semantic tokens remain authoritative.

Journeys:

`Open saved file → create tagged snapshot → search/filter → choose compatible target → review replacement diff → explicit save`

`Paste export → read-only review → inspect native content and portability warnings → import new library entry → optional project application`

`Open project → definition map → inspect an unopened node → loaded roles/references → return to retained editor draft`

| Contract | Decision | Evidence |
| --- | --- | --- |
| Template compatibility | Provider and kind must both match; unavailable actions state why | Go mismatch checks and browser disabled-control tests |
| Project overrides | Full replacement, never a claimed automatic merge; keep removal lines visible | Dirty-draft dialog, stale-hash checks and explicit apply/save tests |
| Import trust boundary | Versioned bounded envelope; review has no writes; edits invalidate review; errors retain input; pending imports lock submission | Real-filesystem import tests and browser recovery tests |
| Export privacy | Complete native text, no upload; warning before copying; no local library identity exported | Portable round-trip tests and rendered export view |
| Definitions | Nested-list containment diagram; distinguish unopened paths, loaded snapshots and unverified references | Browser source-read counts and loaded-role navigation tests |
| Keyboard recovery | Focus enters view/review headings; cancel returns to input; close export returns to trigger; import completion returns to saved list | Component contract and browser interactions |

Rendered new views were inspected at 1280, 768 and 320 CSS pixels. They reuse
native controls, visible labels, wrapped paths/tags and status messages. The
definition map uses a simple linked hierarchy with no draggable canvas: users
can navigate the same structure with native buttons and semantic lists. Automated
axe checks cover both new views; they do not establish screen-reader conformance.
Import/export is intentionally text-based and local; native file-picker transfer
and complete vendor portability validation remain outside this implementation.
