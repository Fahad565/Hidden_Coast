# Hidden Coast — Task Tracker

> One checkbox per task. Feed each task to Jules as a GitHub issue.
> Format: `[ ]` = todo, `[x]` = done, `[~]` = in progress

## Sprint 0 — Setup & Foundation

- [ ] 0.1 Push repo to GitHub with `docs/` folder (PRD, TRD, sprints, Tasks)
- [ ] 0.2 Push `render.yaml` to repo root
- [ ] 0.3 Push `.github/workflows/ci.yml`
- [ ] 0.4 Push `.github/workflows/keep-alive.yml`
- [ ] 0.5 Create Render project, link GitHub repo
- [ ] 0.6 Set `GOOGLE_API_KEY` in Render dashboard
- [ ] 0.7 `go mod init` + deps (Gin, ADK, Gemini)
- [ ] 0.8 Minimal `main.go` with `/api/health`
- [ ] 0.9 Next.js scaffold
- [ ] 0.10 Frontend fetches `/api/health`
- [ ] **Checkpoint:** Render deploys both, `/api/health` returns ok

## Sprint 1 — Infrastructure + SME Registry

- [ ] 1.1 `render.yaml` finalized (env vars, free tier, correct paths)
- [ ] 1.2 `backend/Dockerfile` (distroless)
- [ ] 1.3 `frontend/Dockerfile` (node:22-alpine)
- [ ] 1.4 `backend/Makefile`
- [ ] 1.5 `internal/models/sme.go`
- [ ] 1.6 `internal/models/itinerary.go`
- [ ] 1.7 `internal/models/request.go`
- [ ] 1.8 `internal/registry/registry.go` (Get, SearchByCategory)
- [ ] 1.9 `internal/registry/seed.go` (15 businesses)
- [ ] 1.10 `GET /api/sme/:id` handler
- [ ] 1.11 `GET /api/sme/:id/inventory` handler
- [ ] 1.12 `registry_test.go`
- [ ] **Checkpoint:** `curl .../api/sme/mama-aisha` returns JSON

## Sprint 2 — Agent Tools

- [ ] 2.1 `searchBusinesses` tool
- [ ] 2.2 `queryInventory` tool
- [ ] 2.3 `checkAvailability` tool
- [ ] 2.4 `calculateWalkTime` tool
- [ ] 2.5 `generateItinerary` tool
- [ ] 2.6 `agent/tools_test.go` (all 5)
- [ ] 2.7 Wire tools into single ADK agent
- [ ] 2.8 `POST /api/chat` handler
- [ ] **Checkpoint:** POST `/api/chat` "Find food in Old Town" → returns names

## Sprint 3 — Multi-Agent + SSE + Chat UI

- [ ] 3.1 3-agent topology (Concierge → Discovery + Business)
- [ ] 3.2 `agent/prompts.go`
- [ ] 3.3 SSE streaming on `/api/chat`
- [ ] 3.4 Frontend: `app/chat/page.tsx`
- [ ] 3.5 Frontend: `ItineraryCard.tsx`
- [ ] 3.6 Frontend: `AgentTrace.tsx`
- [ ] 3.7 Frontend: `lib/api.ts` (SSE reader)
- [ ] 3.8 Full flow test
- [ ] **Checkpoint:** Type demo prompt → see streaming + itinerary card

## Sprint 4 — SME Dashboard

- [ ] 4.1 `POST /api/sme/:id/request` (accept/decline)
- [ ] 4.2 `GET /api/sme/:id/requests`
- [ ] 4.3 Tourist confirms → creates requests
- [ ] 4.4 Frontend: `app/sme/page.tsx`
- [ ] 4.5 Frontend: `app/sme/[id]/page.tsx` (Accept/Decline)
- [ ] 4.6 Frontend: inventory stock bars
- [ ] 4.7 Full two-sided test
- [ ] **Checkpoint:** Tourist confirms → SME accepts → stock drops

## Sprint 5 — Demo Polish

- [ ] 5.1 `app/demo/page.tsx` (full-screen)
- [ ] 5.2 Polish agent prompts
- [ ] 5.3 Rehearse 3×
- [ ] 5.4 Verify Render stability
- [ ] 5.5 Record backup video
- [ ] 5.6 Pitch deck (5 slides)
- [ ] 5.7 `README.md` with live URL
- [ ] **Checkpoint:** Live demo works end-to-end in under 4 min

## Sprint 6 — Stretch (if ahead)

- [ ] 6.1 Voice input
- [ ] 6.2 WhatsApp-style SME UI
- [ ] 6.3 Browser geolocation
- [ ] 6.4 "Why this plan?" explanation
- [ ] 6.5 Swahili input

---

## Jules Issue Mapping

| Task(s) | GitHub Issue Title | Priority |
|---------|-------------------|----------|
| 0.7–0.10 | "Sprint 0: Scaffold backend + frontend" | P0 |
| 1.5–1.12 | "Sprint 1: Data models + registry + API handlers" | P0 |
| 2.1–2.8 | "Sprint 2: Agent tools + chat endpoint" | P0 |
| 3.1–3.8 | "Sprint 3: Multi-agent topology + SSE + chat UI" | P0 |
| 4.1–4.7 | "Sprint 4: SME dashboard + request flow" | P1 |
| 5.1–5.7 | "Sprint 5: Demo page + polish" | P1 |

> **Rule:** One issue = one sprint (or half-sprint). Each issue references the
> relevant section of `docs/TRD.md` so Jules has full context.   

Example issues (create these in your GitHub repo)
Issue #1: "Scaffold Go backend with Gin + ADK agent"

## Task
Create the Go backend skeleton for Hidden Coast.

## Context
Read `docs/PRD.md` and `docs/TRD.md` for the full spec.

## Deliverables
- `backend/go.mod` (module: hidden-coast, Go 1.26)
- `backend/cmd/server/main.go` — Gin server on :8080 with:
    - GET /api/health → {"status":"ok"}
    - POST /api/chat → proxy to ADK agent (stub for now, return "hello")
- `backend/internal/models/sme.go` — SME, Product, Location structs (see TRD)
- `backend/internal/registry/registry.go` — in-memory map store
- `backend/Dockerfile` (multi-stage, distroless)
- `backend/Makefile`

## Constraints
- Use `github.com/gin-gonic/gin` for HTTP
- Use `google.golang.org/adk/v2` for the agent (stub the model for now)
- No external DB. In-memory only.
- Must compile: `go build ./...`

## Acceptance
- `docker build -t test .` succeeds
- `curl localhost:8080/api/health` returns `{"status":"ok"}`

Issue #2: "Seed 15 Mombasa SMEs into the registry"

## Task
Populate the in-memory registry with 15 demo businesses.

## Context
Read `docs/TRD.md` → "Seed Data" table for the exact 15 businesses.
The data model is in `backend/internal/models/sme.go`.

## Deliverables
- `backend/internal/registry/seed.go` — a `Seed()` function that
  populates the registry with all 15 businesses from the TRD table.
  Include realistic prices, stock, prep times, opening hours, and
  a "story" field for each.

## Acceptance
- `go test ./internal/registry/` passes
- A test verifies: `registry.Get("mama-aisha").Products[0].Price == 450`

Issue #3: "Implement ADK agent tools"

## Task
Implement the 5 tools the ADK agent calls.

## Context
Read `docs/TRD.md` → "Agent Design" section.
The registry is in `backend/internal/registry/`.

## Deliverables
- `backend/agent/tools.go`:
    - `searchBusinesses(category, area)` → filters registry
    - `queryInventory(smeID)` → returns products + stock
    - `checkAvailability(smeID, productID, qty, time)` → bool + message
    - `calculateWalkTime(from, to)` → int minutes (hardcoded for demo)
    - `generateItinerary(startTime, duration, budget, stops)` → Itinerary
- `backend/agent/tools_test.go` — unit tests for each tool

## Constraints
- Use `google.golang.org/adk/v2/tool/functiontool` for wrapping
- `checkAvailability` must verify: stock >= qty, prep time fits before closing
- `generateItinerary` uses greedy: sort by time, accumulate cost, drop if over budget

## Acceptance
- `go test ./agent/ -v` passes
- `checkAvailability("mama-aisha", "chicken-biryani", 2, "12:00")` → (true, "Ready by 12:35")

Issue #4: "Build Next.js frontend with chat + SME dashboard"

## Task
Create the Next.js frontend.

## Context
Read `docs/PRD.md` → "User Journeys" and `docs/TRD.md` → "Frontend".

## Deliverables
- `frontend/` — Next.js 15 (App Router) + Tailwind CSS
- `frontend/app/chat/page.tsx` — chat UI that POSTs to backend /api/chat
- `frontend/app/sme/page.tsx` — list of 15 businesses (pick one)
- `frontend/app/sme/[id]/page.tsx` — SME dashboard:
    - Shows pending requests with Accept/Decline buttons
    - Shows inventory (stock bars)
- `frontend/lib/api.ts` — fetch wrapper (base URL from env)
- `frontend/Dockerfile` (multi-stage, node:22-alpine)

## Constraints
- Use Tailwind for styling (no CSS modules)
- Chat should show a loading state while waiting for agent response
- SME dashboard should refresh inventory after accept/decline
- Keep it simple. No auth. No routing guards.

## Acceptance
- `npm run build` succeeds
- `docker build -t test .` succeeds

Issue #5: "Wire ADK multi-agent topology"

## Task
Connect the 3-agent ADK topology (Concierge → Discovery + Business).

## Context
Read `docs/TRD.md` → "Agent Design" → "Agent Topology".
Tools are already implemented in `backend/agent/tools.go`.

## Deliverables
- `backend/agent/agent.go` — `NewConcierge()` that creates:
    - DiscoveryAgent (tools: searchBusinesses, searchExperiences)
    - BusinessAgent (tools: queryInventory, checkAvailability, getPrepTime, calculateWalkTime)
    - ConciergeAgent (sub-agents: Discovery + Business, tool: generateItinerary)
- Update `backend/cmd/server/main.go` to wire the concierge into the /api/chat handler
- SSE streaming: stream agent tokens to the client

## Acceptance
- POST /api/chat with "I'm in Mombasa for 3 hours, KSh 3000, food + crafts"
  returns a streaming response with a valid itinerary
- The agent calls at least 3 tools during the response