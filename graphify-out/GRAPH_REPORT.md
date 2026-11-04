# Graph Report - mail-scrapper  (2026-11-04)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 2752 nodes · 8718 edges · 78 communities (69 shown, 9 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 341 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7d93826f`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- .agent/skills/impeccable/scripts/live-browser.js
- .cursor/skills/impeccable/scripts/live-browser.js
- .gemini/skills/impeccable/scripts/live-browser.js
- context.Context
- connectSSE
- connectSSE
- connectSSE
- setLiveState
- go_pkg_strings
- setLiveState
- setLiveState
- worker.go
- NewServer
- el
- .agent/skills/impeccable/scripts/modern-screenshot.umd.js
- el
- .cursor/skills/impeccable/scripts/modern-screenshot.umd.js
- el
- .gemini/skills/impeccable/scripts/modern-screenshot.umd.js
- syncPageChatFocus
- syncPageChatFocus
- syncPageChatFocus
- ApplicationDetailPage.tsx
- package.json
- testing.T
- Worker
- documentRefForElement
- documentRefForElement
- documentRefForElement
- hooks.ts
- showToast
- showToast
- showToast
- initGlobalBar
- initGlobalBar
- initGlobalBar
- handleManualEditActivity
- handleManualEditActivity
- handleManualEditActivity
- dependencies
- captureElementToBlob
- createLiveBrowserDomHelpers
- createLiveBrowserSessionState
- captureElementToBlob
- createLiveBrowserDomHelpers
- createLiveBrowserSessionState
- captureElementToBlob
- createLiveBrowserDomHelpers
- createLiveBrowserSessionState
- compilerOptions
- actOnAgentTarget
- actOnAgentTarget
- actOnAgentTarget
- retry.go
- App.tsx
- compilerOptions
- onAnnotDown
- onAnnotDown
- onAnnotDown
- oauth.go
- mime.go
- TopBar.tsx
- scheduleAcceptCleanup
- scheduleAcceptCleanup
- scheduleAcceptCleanup
- theme.tsx
- .agent/skills/impeccable/scripts/live-browser-ignores.js
- .cursor/skills/impeccable/scripts/live-browser-ignores.js
- .gemini/skills/impeccable/scripts/live-browser-ignores.js
- tsconfig.json
- github.com/dharundds/jobtrack

## God Nodes (most connected - your core abstractions)
1. `NewServer()` - 36 edges
2. `connectSSE()` - 34 edges
3. `connectSSE()` - 34 edges
4. `connectSSE()` - 34 edges
5. `setLiveState()` - 33 edges
6. `resumeSession()` - 33 edges
7. `resumeSession()` - 33 edges
8. `resumeSession()` - 33 edges
9. `setLiveState()` - 33 edges
10. `setLiveState()` - 33 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `NewServer()`  [EXTRACTED]
  cmd/jobtrack/main.go → internal/api/server.go
- `main()` --calls--> `Load()`  [EXTRACTED]
  cmd/jobtrack/main.go → internal/config/config.go
- `SimilarCompanyRole()` --calls--> `isATSBrandName()`  [INFERRED]
  internal/classify/status.go → internal/classify/extract.go
- `TestExtractCompany_greenhouseSharedSender()` --calls--> `SimilarCompanyRole()`  [INFERRED]
  internal/classify/extract_test.go → internal/classify/status.go
- `TestRuleClassifier_thalesScreeningInvite()` --calls--> `testRules()`  [INFERRED]
  internal/classify/extract_test.go → internal/classify/classifier_test.go

## Import Cycles
- 2-file cycle: `internal/store/accounts_test.go -> internal/store/messages_filter_test.go -> internal/store/accounts_test.go`
- 2-file cycle: `internal/store/accounts_test.go -> internal/store/store_test.go -> internal/store/accounts_test.go`
- 2-file cycle: `internal/store/messages_filter_test.go -> internal/store/store_test.go -> internal/store/messages_filter_test.go`
- 3-file cycle: `internal/store/accounts_test.go -> internal/store/messages_filter_test.go -> internal/store/store_test.go -> internal/store/accounts_test.go`
- 3-file cycle: `internal/store/accounts_test.go -> internal/store/store_test.go -> internal/store/messages_filter_test.go -> internal/store/accounts_test.go`

## Communities (78 total, 9 thin omitted)

### Community 0 - ".agent/skills/impeccable/scripts/live-browser.js"
Cohesion: 0.03
Nodes (139): applyGlobalBarLabelState(), applyLiveBarPreference(), applyParamDefaults(), applyParamValue(), applyPlaceholderSizingStyles(), applySavedParamValues(), bindEditBadgeProxy(), bufferToBase64() (+131 more)

### Community 1 - ".cursor/skills/impeccable/scripts/live-browser.js"
Cohesion: 0.03
Nodes (139): applyGlobalBarLabelState(), applyLiveBarPreference(), applyParamDefaults(), applyParamValue(), applyPlaceholderSizingStyles(), applySavedParamValues(), bindEditBadgeProxy(), bufferToBase64() (+131 more)

### Community 2 - ".gemini/skills/impeccable/scripts/live-browser.js"
Cohesion: 0.03
Nodes (139): applyGlobalBarLabelState(), applyLiveBarPreference(), applyParamDefaults(), applyParamValue(), applyPlaceholderSizingStyles(), applySavedParamValues(), bindEditBadgeProxy(), bufferToBase64() (+131 more)

### Community 3 - "context.Context"
Cohesion: 0.03
Nodes (28): Client, Message, QuotaLimiter, appliedAtForSplit(), AccountFromContext(), ApplicationFromStore(), EventFromStore(), ApplicationJSON (+20 more)

### Community 4 - "connectSSE"
Cohesion: 0.06
Nodes (83): applyPlaceholderDimensions(), applySavedSessionMeta(), buildParamsPanel(), checkpointPayload(), clampVariantIndex(), closedClipPath(), commitAcceptedVariantToDom(), completeParameterGenerationIfReady() (+75 more)

### Community 5 - "connectSSE"
Cohesion: 0.06
Nodes (83): applyPlaceholderDimensions(), applySavedSessionMeta(), buildParamsPanel(), checkpointPayload(), clampVariantIndex(), closedClipPath(), commitAcceptedVariantToDom(), completeParameterGenerationIfReady() (+75 more)

### Community 6 - "connectSSE"
Cohesion: 0.06
Nodes (83): applyPlaceholderDimensions(), applySavedSessionMeta(), buildParamsPanel(), checkpointPayload(), clampVariantIndex(), closedClipPath(), commitAcceptedVariantToDom(), completeParameterGenerationIfReady() (+75 more)

### Community 7 - "setLiveState"
Cohesion: 0.09
Nodes (67): abortSvelteComponentInjection(), applyEditing(), beginNewLiveConfiguration(), buildLocatorForLeaf(), cancelEditing(), cancelEditingToPicking(), cancelInsertConfigure(), cleanup() (+59 more)

### Community 8 - "go_pkg_strings"
Cohesion: 0.06
Nodes (46): ApplicationActivity, ApplicationCandidate, Classifier, EmailInput, EventDetails, EventRecord, GroupKey, MessageActivity (+38 more)

### Community 9 - "setLiveState"
Cohesion: 0.09
Nodes (67): abortSvelteComponentInjection(), applyEditing(), beginNewLiveConfiguration(), buildLocatorForLeaf(), cancelEditing(), cancelEditingToPicking(), cancelInsertConfigure(), cleanup() (+59 more)

### Community 10 - "setLiveState"
Cohesion: 0.09
Nodes (67): abortSvelteComponentInjection(), applyEditing(), beginNewLiveConfiguration(), buildLocatorForLeaf(), cancelEditing(), cancelEditingToPicking(), cancelInsertConfigure(), cleanup() (+59 more)

### Community 11 - "worker.go"
Cohesion: 0.20
Nodes (7): authStatusResponse, createRuleRequest, mergeApplicationsRequest, patchApplicationRequest, patchRuleRequest, patchSettingsRequest, splitMessageRequest

### Community 12 - "NewServer"
Cohesion: 0.10
Nodes (28): Server, AppError, ApplicationsHandler, AuthHandler, HealthHandler, RulesHandler, SettingsHandler, StatsHandler (+20 more)

### Community 13 - "el"
Cohesion: 0.07
Nodes (57): actionLabel(), applyConfigureBarChrome(), bindConfigureCountPillTooltip(), bindConfigureInlineControlHover(), bindConfigureModifierPillHover(), buildConfigureActionControl(), buildConfigureCountControl(), buildConfigureRow() (+49 more)

### Community 14 - ".agent/skills/impeccable/scripts/modern-screenshot.umd.js"
Cohesion: 0.09
Nodes (53): ae(), be(), bt(), Ce(), Ct(), de(), dt(), _e() (+45 more)

### Community 15 - "el"
Cohesion: 0.07
Nodes (57): actionLabel(), applyConfigureBarChrome(), bindConfigureCountPillTooltip(), bindConfigureInlineControlHover(), bindConfigureModifierPillHover(), buildConfigureActionControl(), buildConfigureCountControl(), buildConfigureRow() (+49 more)

### Community 16 - ".cursor/skills/impeccable/scripts/modern-screenshot.umd.js"
Cohesion: 0.09
Nodes (53): ae(), be(), bt(), Ce(), Ct(), de(), dt(), _e() (+45 more)

### Community 17 - "el"
Cohesion: 0.07
Nodes (57): actionLabel(), applyConfigureBarChrome(), bindConfigureCountPillTooltip(), bindConfigureInlineControlHover(), bindConfigureModifierPillHover(), buildConfigureActionControl(), buildConfigureCountControl(), buildConfigureRow() (+49 more)

### Community 18 - ".gemini/skills/impeccable/scripts/modern-screenshot.umd.js"
Cohesion: 0.09
Nodes (53): ae(), be(), bt(), Ce(), Ct(), de(), dt(), _e() (+45 more)

### Community 19 - "syncPageChatFocus"
Cohesion: 0.09
Nodes (49): armPageChatForTyping(), attachSteerFocusDebug(), attachSteerFocusGuard(), buildSteerProcessingDots(), buildSteerQueueHint(), clearSteerAwaitTimer(), clearSteerFocusRecoverTimer(), collapsePageChat() (+41 more)

### Community 20 - "syncPageChatFocus"
Cohesion: 0.09
Nodes (49): armPageChatForTyping(), attachSteerFocusDebug(), attachSteerFocusGuard(), buildSteerProcessingDots(), buildSteerQueueHint(), clearSteerAwaitTimer(), clearSteerFocusRecoverTimer(), collapsePageChat() (+41 more)

### Community 21 - "syncPageChatFocus"
Cohesion: 0.09
Nodes (49): armPageChatForTyping(), attachSteerFocusDebug(), attachSteerFocusGuard(), buildSteerProcessingDots(), buildSteerQueueHint(), clearSteerAwaitTimer(), clearSteerFocusRecoverTimer(), collapsePageChat() (+41 more)

### Community 22 - "ApplicationDetailPage.tsx"
Cohesion: 0.13
Nodes (35): invalidateApplicationQueries(), useApplication(), useApplications(), useDeleteApplication(), useIgnoreMessage(), useMergeApplications(), usePatchApplication(), useSplitMessage() (+27 more)

### Community 23 - "package.json"
Cohesion: 0.06
Nodes (41): autoprefixer, eslint, @eslint/js, eslint-plugin-react-hooks, eslint-plugin-react-refresh, globals, postcss, prettier (+33 more)

### Community 24 - "testing.T"
Cohesion: 0.09
Nodes (33): TestExtractCompanyRole(), companyFromDisplayName(), companyFromEmployerDomain(), displayName(), emailAddress(), ExtractCompany(), extractCompanyFromText(), ExtractRole() (+25 more)

### Community 25 - "Worker"
Cohesion: 0.08
Nodes (19): main(), fetchGmailAddress(), Service, NewService(), NewProcessor(), DecryptToken(), EncryptToken(), TestDecryptWrongKey() (+11 more)

### Community 26 - "documentRefForElement"
Cohesion: 0.06
Nodes (39): addManualContextText(), canRestoreManualEditElement(), collectEditableTextRows(), visit(), collectManualContextPieces(), walk(), contextElementForManualEdit(), copyEditContainerContext() (+31 more)

### Community 27 - "documentRefForElement"
Cohesion: 0.06
Nodes (39): addManualContextText(), canRestoreManualEditElement(), collectEditableTextRows(), visit(), collectManualContextPieces(), walk(), contextElementForManualEdit(), copyEditContainerContext() (+31 more)

### Community 28 - "documentRefForElement"
Cohesion: 0.06
Nodes (39): addManualContextText(), canRestoreManualEditElement(), collectEditableTextRows(), visit(), collectManualContextPieces(), walk(), contextElementForManualEdit(), copyEditContainerContext() (+31 more)

### Community 29 - "hooks.ts"
Cohesion: 0.12
Nodes (31): ApiError, apiFetch(), appsParamsKey(), queryKeys, SplitMessageInput, SplitMessageResult, useCreateRule(), useDeleteRule() (+23 more)

### Community 30 - "showToast"
Cohesion: 0.09
Nodes (34): abandonForeignSession(), abandonSupersededGo(), applyOriginalAttrsToSvelteAnchor(), clearMountErrorCard(), commitAcceptedSvelteComponentToDom(), componentModuleCandidates(), describeMountFailure(), detectDevServerBase() (+26 more)

### Community 31 - "showToast"
Cohesion: 0.09
Nodes (34): abandonForeignSession(), abandonSupersededGo(), applyOriginalAttrsToSvelteAnchor(), clearMountErrorCard(), commitAcceptedSvelteComponentToDom(), componentModuleCandidates(), describeMountFailure(), detectDevServerBase() (+26 more)

### Community 32 - "showToast"
Cohesion: 0.09
Nodes (34): abandonForeignSession(), abandonSupersededGo(), applyOriginalAttrsToSvelteAnchor(), clearMountErrorCard(), commitAcceptedSvelteComponentToDom(), componentModuleCandidates(), describeMountFailure(), detectDevServerBase() (+26 more)

### Community 33 - "initGlobalBar"
Cohesion: 0.13
Nodes (24): agentHasWorkInFlight(), agentStatusText(), barPaletteForTheme(), brandMarkSvg(), detectPageTheme(), ensureAgentPollTooltip(), ensureSpinKeyframes(), hideAgentPollTooltip() (+16 more)

### Community 34 - "initGlobalBar"
Cohesion: 0.13
Nodes (24): agentHasWorkInFlight(), agentStatusText(), barPaletteForTheme(), brandMarkSvg(), detectPageTheme(), ensureAgentPollTooltip(), ensureSpinKeyframes(), hideAgentPollTooltip() (+16 more)

### Community 35 - "initGlobalBar"
Cohesion: 0.13
Nodes (24): agentHasWorkInFlight(), agentStatusText(), barPaletteForTheme(), brandMarkSvg(), detectPageTheme(), ensureAgentPollTooltip(), ensureSpinKeyframes(), hideAgentPollTooltip() (+16 more)

### Community 36 - "handleManualEditActivity"
Cohesion: 0.19
Nodes (24): clearStoredManualApplyState(), fetchPendingCount(), handleManualEditActivity(), hidePendingApplyDock(), manualApplyLoadingText(), manualApplyStateKey(), manualEditEventForCurrentPage(), numberOrNull() (+16 more)

### Community 37 - "handleManualEditActivity"
Cohesion: 0.19
Nodes (24): clearStoredManualApplyState(), fetchPendingCount(), handleManualEditActivity(), hidePendingApplyDock(), manualApplyLoadingText(), manualApplyStateKey(), manualEditEventForCurrentPage(), numberOrNull() (+16 more)

### Community 38 - "handleManualEditActivity"
Cohesion: 0.19
Nodes (24): clearStoredManualApplyState(), fetchPendingCount(), handleManualEditActivity(), hidePendingApplyDock(), manualApplyLoadingText(), manualApplyStateKey(), manualEditEventForCurrentPage(), numberOrNull() (+16 more)

### Community 39 - "dependencies"
Cohesion: 0.10
Nodes (21): date-fns, @dnd-kit/core, @dnd-kit/sortable, @dnd-kit/utilities, react, react-dom, @tanstack/react-query, dependencies (+13 more)

### Community 40 - "captureElementToBlob"
Cohesion: 0.14
Nodes (19): averageRgb01(), captureAndEmit(), captureElementFromRenderedAncestor(), captureElementToBlob(), compileShader(), cssColorToRgb01(), dominantRgb01(), findBackdropAncestor() (+11 more)

### Community 41 - "createLiveBrowserDomHelpers"
Cohesion: 0.16
Nodes (15): createLiveBrowserDomHelpers(), activeElementDeep(), cloneWithoutChrome(), copyTextViaSelection(), cssId(), liveUiRoot(), makeFrozenAnchor(), raiseHitModal() (+7 more)

### Community 42 - "createLiveBrowserSessionState"
Cohesion: 0.19
Nodes (17): createLiveBrowserSessionState(), clearHandled(), clearScrollY(), clearSession(), isHandled(), loadSession(), markHandled(), nextCheckpointRevision() (+9 more)

### Community 43 - "captureElementToBlob"
Cohesion: 0.14
Nodes (19): averageRgb01(), captureAndEmit(), captureElementFromRenderedAncestor(), captureElementToBlob(), compileShader(), cssColorToRgb01(), dominantRgb01(), findBackdropAncestor() (+11 more)

### Community 44 - "createLiveBrowserDomHelpers"
Cohesion: 0.16
Nodes (15): createLiveBrowserDomHelpers(), activeElementDeep(), cloneWithoutChrome(), copyTextViaSelection(), cssId(), liveUiRoot(), makeFrozenAnchor(), raiseHitModal() (+7 more)

### Community 45 - "createLiveBrowserSessionState"
Cohesion: 0.19
Nodes (17): createLiveBrowserSessionState(), clearHandled(), clearScrollY(), clearSession(), isHandled(), loadSession(), markHandled(), nextCheckpointRevision() (+9 more)

### Community 46 - "captureElementToBlob"
Cohesion: 0.14
Nodes (19): averageRgb01(), captureAndEmit(), captureElementFromRenderedAncestor(), captureElementToBlob(), compileShader(), cssColorToRgb01(), dominantRgb01(), findBackdropAncestor() (+11 more)

### Community 47 - "createLiveBrowserDomHelpers"
Cohesion: 0.16
Nodes (15): createLiveBrowserDomHelpers(), activeElementDeep(), cloneWithoutChrome(), copyTextViaSelection(), cssId(), liveUiRoot(), makeFrozenAnchor(), raiseHitModal() (+7 more)

### Community 48 - "createLiveBrowserSessionState"
Cohesion: 0.19
Nodes (17): createLiveBrowserSessionState(), clearHandled(), clearScrollY(), clearSession(), isHandled(), loadSession(), markHandled(), nextCheckpointRevision() (+9 more)

### Community 49 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowImportingTsExtensions, baseUrl, isolatedModules, jsx, lib, module, moduleDetection (+11 more)

### Community 50 - "actOnAgentTarget"
Cohesion: 0.29
Nodes (17): actOnAgentTarget(), agentTargetBusyReason(), agentTargetOverlayGone(), agentTargetTaken(), claimAgentTarget(), claimAndActOnAgentTarget(), declineAgentTargetBusy(), declineAgentTargetUnresolvable() (+9 more)

### Community 51 - "actOnAgentTarget"
Cohesion: 0.29
Nodes (17): actOnAgentTarget(), agentTargetBusyReason(), agentTargetOverlayGone(), agentTargetTaken(), claimAgentTarget(), claimAndActOnAgentTarget(), declineAgentTargetBusy(), declineAgentTargetUnresolvable() (+9 more)

### Community 52 - "actOnAgentTarget"
Cohesion: 0.29
Nodes (17): actOnAgentTarget(), agentTargetBusyReason(), agentTargetOverlayGone(), agentTargetTaken(), claimAgentTarget(), claimAndActOnAgentTarget(), declineAgentTargetBusy(), declineAgentTargetUnresolvable() (+9 more)

### Community 53 - "retry.go"
Cohesion: 0.17
Nodes (11): backoffWithJitter(), IsRateLimitError(), IsRetryable(), retryAfter(), TestIsRetryable(), TestIsRetryable_wrapped(), withRetry(), MaskEmail() (+3 more)

### Community 54 - "App.tsx"
Cohesion: 0.26
Nodes (12): react-router-dom, useAuthStatus(), App(), AppLayout(), linkEnd(), links, MobileNav(), Sidebar() (+4 more)

### Community 55 - "compilerOptions"
Cohesion: 0.12
Nodes (15): compilerOptions, allowImportingTsExtensions, isolatedModules, lib, module, moduleDetection, moduleResolution, noEmit (+7 more)

### Community 56 - "onAnnotDown"
Cohesion: 0.24
Nodes (15): beginEditPin(), buildAnnotationsForCapture(), buildPinElement(), cancelEditingPin(), finalizeEditingPin(), initAnnotOverlay(), localCoords(), onAnnotDown() (+7 more)

### Community 57 - "onAnnotDown"
Cohesion: 0.24
Nodes (15): beginEditPin(), buildAnnotationsForCapture(), buildPinElement(), cancelEditingPin(), finalizeEditingPin(), initAnnotOverlay(), localCoords(), onAnnotDown() (+7 more)

### Community 58 - "onAnnotDown"
Cohesion: 0.24
Nodes (15): beginEditPin(), buildAnnotationsForCapture(), buildPinElement(), cancelEditingPin(), finalizeEditingPin(), initAnnotOverlay(), localCoords(), onAnnotDown() (+7 more)

### Community 59 - "oauth.go"
Cohesion: 0.19
Nodes (8): sessionOptions(), AccountIDFromSession(), ClearSession(), NewOAuthState(), SecureCookie(), SetLoggedIn(), SetOAuthState(), ValidateOAuthState()

### Community 60 - "mime.go"
Cohesion: 0.26
Nodes (5): compressSpace(), decodePart(), extractBody(), htmlToText(), stripTagsSimple()

### Community 61 - "TopBar.tsx"
Cohesion: 0.30
Nodes (9): lucide-react, useLogout(), useSyncStatus(), useTriggerSync(), AuthAccountButton(), TopBar(), TopBarProps, formatRelative() (+1 more)

### Community 62 - "scheduleAcceptCleanup"
Cohesion: 0.31
Nodes (11): acceptedDomAlreadyClean(), clearHandledWrapperReloadStamp(), deferredRecoverySuperseded(), ensureAcceptedDomClean(), findAcceptedRuntimeWrappers(), handledWrapperReloadKey(), reloadAfterMissingAcceptedDom(), restoreAcceptedDomFromSnapshot() (+3 more)

### Community 63 - "scheduleAcceptCleanup"
Cohesion: 0.31
Nodes (11): acceptedDomAlreadyClean(), clearHandledWrapperReloadStamp(), deferredRecoverySuperseded(), ensureAcceptedDomClean(), findAcceptedRuntimeWrappers(), handledWrapperReloadKey(), reloadAfterMissingAcceptedDom(), restoreAcceptedDomFromSnapshot() (+3 more)

### Community 64 - "scheduleAcceptCleanup"
Cohesion: 0.31
Nodes (11): acceptedDomAlreadyClean(), clearHandledWrapperReloadStamp(), deferredRecoverySuperseded(), ensureAcceptedDomClean(), findAcceptedRuntimeWrappers(), handledWrapperReloadKey(), reloadAfterMissingAcceptedDom(), restoreAcceptedDomFromSnapshot() (+3 more)

### Community 68 - "theme.tsx"
Cohesion: 0.33
Nodes (8): applyTheme(), getSystemDark(), readStoredMode(), resolveMode(), ThemeContext, ThemeContextValue, ThemeMode, ThemeProvider()

### Community 69 - ".agent/skills/impeccable/scripts/live-browser-ignores.js"
Cohesion: 0.52
Nodes (6): globToRegex(), matchesScope(), normalizeIgnoreRule(), normalizeIgnoreValue(), pageCandidates(), resolveDetectIgnores()

### Community 70 - ".cursor/skills/impeccable/scripts/live-browser-ignores.js"
Cohesion: 0.52
Nodes (6): globToRegex(), matchesScope(), normalizeIgnoreRule(), normalizeIgnoreValue(), pageCandidates(), resolveDetectIgnores()

### Community 71 - ".gemini/skills/impeccable/scripts/live-browser-ignores.js"
Cohesion: 0.52
Nodes (6): globToRegex(), matchesScope(), normalizeIgnoreRule(), normalizeIgnoreValue(), pageCandidates(), resolveDetectIgnores()

## Knowledge Gaps
- **1 isolated node(s):** `github.com/dharundds/jobtrack`
  These have ≤1 connection - possible missing edges. (Counts symbols only; 192 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `createLiveBrowserDomHelpers()` connect `createLiveBrowserDomHelpers` to `setLiveState`, `documentRefForElement`?**
  _High betweenness centrality (0.004) - this node is a cross-community bridge._
- **Are the 22 inferred relationships involving `NewServer()` (e.g. with `.IsDev()` and `.Delete()`) actually correct?**
  _`NewServer()` has 22 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/dharundds/jobtrack` to the rest of the system?**
  _1 weakly-connected nodes found - possible missing edges._
- **Should `.agent/skills/impeccable/scripts/live-browser.js` be split into smaller, more focused modules?**
  _Cohesion score 0.027874902874902876 - nodes in this community are weakly interconnected._
- **Why does `stopFocusOutIntoChrome()` connect `setLiveState` to `documentRefForElement`, `createLiveBrowserDomHelpers`?**
  _High betweenness centrality (0.002) - this node is a cross-community bridge._
- **Should `.cursor/skills/impeccable/scripts/live-browser.js` be split into smaller, more focused modules?**
  _Cohesion score 0.027874902874902876 - nodes in this community are weakly interconnected._
- **Why does `updateGlobalBarState()` connect `initGlobalBar` to `.gemini/skills/impeccable/scripts/live-browser.js`, `setLiveState`?**
  _High betweenness centrality (0.002) - this node is a cross-community bridge._