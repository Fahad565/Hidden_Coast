# Hidden Coast — Sprint Plan

**Stack:** Go (Gin) + ADK + Next.js + Render Blueprint (free tier)
**CI:** `.github/workflows/ci.yml`
**Deploy:** `render.yaml` (auto-deploys on merge to `main`)
**IDE:** GoLand (optional) / GitHub web (primary)

---

## Sprint 0 — Setup & Foundation (Hour 0–3)

| # | Task | Acceptance |
|---|------|-----------|
| 0.1 | `git init`, push to GitHub | Repo exists, `docs/` folder has PRD, TRD, sprints, Tasks |
| 0.2 | Push `render.yaml` to repo root | Render connects, shows 2 services (pending builds) |
| 0.3 | Push `.github/workflows/ci.yml` | CI runs on PR (go test + npm build + docker build) |
| 0.4 | Push `.github/workflows/keep-alive.yml` | Cron pings both URLs every 10 min |
| 0.5 | Create Render project, link GitHub repo | Render detects `render.yaml`, shows deploy status |
| 0.6 | Set `GOOGLE_API_KEY` in Render dashboard (backend service) | Env var visible in Render → Service → Environment |
| 0.7 | `go mod init hidden-coast`, add Gin + ADK deps | `go build ./...` compiles |
| 0.8 | Minimal `main.go`: Gin on :8080, `/api/health` → `{"status":"ok"}` | `curl .../api/health` works after Render deploy |
| 0.9 | Next.js scaffold: `npx create-next-app@latest frontend` | `npm run build` passes in CI |
| 0.10 | Frontend: basic page that fetches `/api/health` | Browser shows "ok" |

**Checkpoint:** Push to `main` → Render deploys both → `https://hidden-coast-backend.onrender.com/api/health` returns `{"status":"ok"}`.

---

## Sprint 1 — Infrastructure + SME Registry (Hour 3–7)

| # | Task | Acceptance |
|---|------|-----------|
| 1.1 | `render.yaml` finalized (both services, correct env vars, free tier) | Render shows both services as "Live" |
| 1.2 | `backend/Dockerfile` (multi-stage, distroless, EXPOSE 8080) | `docker build` passes in CI |
| 1.3 | `frontend/Dockerfile` (multi-stage, node:22-alpine, EXPOSE 3000) | `docker build` passes in CI |
| 1.4 | `backend/Makefile` (run, test, build, docker targets) | `make test` works |
| 1.5 | `backend/internal/models/sme.go` — SME, Product, Location structs | Compiles |
| 1.6 | `backend/internal/models/itinerary.go` — Itinerary, Stop structs | Compiles |
| 1.7 | `backend/internal/models/request.go` — TouristRequest, Item structs | Compiles |
| 1.8 | `backend/internal/registry/registry.go` — in-memory map, Get, SearchByCategory | Unit tests pass |
| 1.9 | `backend/internal/registry/seed.go` — 15 hardcoded Mombasa businesses | `registry.Get("mama-aisha")` returns correct data |
| 1.10 | Gin handler: `GET /api/sme/:id` | `curl .../api/sme/mama-aisha` returns JSON |
| 1.11 | Gin handler: `GET /api/sme/:id/inventory` | Returns products with stock |
| 1.12 | `backend/internal/registry/registry_test.go` | `go test -race ./...` passes in CI |

**Checkpoint:** Push to `main` → Render deploys → `curl .../api/sme/mama-aisha` returns full JSON profile.

---

## Sprint 2 — Agent Tools (Hour 7–12)

| # | Task | Acceptance |
|---|------|-----------|
| 2.1 | `agent/tools.go`: `searchBusinesses(category, area)` | Unit test: returns correct subset |
| 2.2 | `agent/tools.go`: `queryInventory(smeID)` | Unit test: returns products + stock |
| 2.3 | `agent/tools.go`: `checkAvailability(smeID, productID, qty, time)` | Unit test: (true, "Ready by 12:35") for valid case |
| 2.4 | `agent/tools.go`: `calculateWalkTime(from, to)` | Hardcoded map, returns int minutes |
| 2.5 | `agent/tools.go`: `generateItinerary(startTime, duration, budget, stops)` | Unit test: greedy solver, under budget |
| 2.6 | `agent/tools_test.go` — all 5 tools tested | `go test ./agent/ -v` passes in CI |
| 2.7 | Wire tools into a single ADK agent (no sub-agents yet) | Agent responds to "Find food in Old Town" |
| 2.8 | `POST /api/chat` handler: accepts message, calls agent, returns response | `curl -X POST .../api/chat -d '{"message":"hello"}'` works |

**Checkpoint:** Push to `main` → Render deploys → POST `/api/chat` with "Find food in Old Town" returns business names.

---

## Sprint 3 — Multi-Agent + SSE + Frontend Chat (Hour 12–18)

| # | Task | Acceptance |
|---|------|-----------|
| 3.1 | `agent/agent.go`: 3-agent topology (Concierge → Discovery + Business) | Agent calls sub-agents |
| 3.2 | `agent/prompts.go`: instruction strings for each agent | Prompts are testable |
| 3.3 | SSE streaming on `/api/chat` (stream tokens as agent thinks) | Browser shows text appearing live |
| 3.4 | Frontend: `app/chat/page.tsx` — chat UI with streaming | Type message → see tokens stream in |
| 3.5 | Frontend: `ItineraryCard.tsx` — structured plan display | Itinerary renders as card, not raw text |
| 3.6 | Frontend: `AgentTrace.tsx` — shows which agent is active | Demo shows "Discovery Agent searching..." |
| 3.7 | Frontend: `lib/api.ts` — fetch wrapper with SSE reader | Clean API layer |
| 3.8 | Full flow test: "3 hours, KSh 3000, food + crafts" | Returns valid itinerary under budget |

**Checkpoint:** Open frontend URL → type demo prompt → watch agents work → see itinerary card.

---

## Sprint 4 — SME Dashboard (Hour 18–24)

| # | Task | Acceptance |
|---|------|-----------|
| 4.1 | `POST /api/sme/:id/request` — accept/decline endpoint | Accept decrements stock |
| 4.2 | `GET /api/sme/:id/requests` — list pending/accepted | Returns request list |
| 4.3 | Wire: tourist confirms itinerary → creates requests for each SME | 2 requests created |
| 4.4 | Frontend: `app/sme/page.tsx` — pick your business | 15 businesses listed |
| 4.5 | Frontend: `app/sme/[id]/page.tsx` — requests + Accept/Decline | Buttons work, stock updates |
| 4.6 | Frontend: inventory view (stock bars) | Visual, updates on accept |
| 4.7 | Full two-sided test: tourist confirms → SME accepts → stock drops → tourist notified | End-to-end works |

**Checkpoint:** Split screen: tourist confirms → SME dashboard shows request → Accept → stock 12→10 → tourist sees confirmation.

---

## Sprint 5 — Demo Polish + Deploy (Hour 24–40)

| # | Task | Acceptance |
|---|------|-----------|
| 5.1 | `app/demo/page.tsx` — full-screen demo mode | Auto-plays the 3-min story |
| 5.2 | Polish agent prompts (reference stories, be warm) | Sounds natural, not robotic |
| 5.3 | Rehearse demo 3×, fix timing | Flows in under 4 minutes |
| 5.4 | Verify Render free tier: both services awake, no 502s | Stable during demo window |
| 5.5 | Record backup demo video (screen capture) | Exists in repo or Loom |
| 5.6 | Pitch deck: 5 slides (problem, solution, demo, architecture, vision) | Ready to present |
| 5.7 | `README.md` — 1-paragraph description + live URL + demo video link | Judges can find it |

**Checkpoint:** Live URL works end-to-end. Demo fits in 3–4 minutes. Backup video exists.

---

## Sprint 6 — Buffer / Stretch (Hour 40–48)

| # | Stretch Goal | Impact |
|---|-------------|--------|
| 6.1 | Voice input (Web Speech API) | "Wow" factor |
| 6.2 | WhatsApp-style UI for SME | More realistic |
| 6.3 | Browser geolocation ("You're at Mombasa Port") | More authentic |
| 6.4 | "Why this plan?" → agent explains reasoning | Transparency |
| 6.5 | Multi-language (Swahili input) | Local authenticity |

### If behind schedule (cut in this order):
1. Cut `AgentTrace` component
2. Cut `/demo` auto-play → just use chat page
3. Cut to 10 SMEs
4. Cut SSE → full response (slower but simpler)
5. Cut SME dashboard → just show "Request sent" in chat   