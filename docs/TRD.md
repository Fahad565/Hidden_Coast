# Hidden Coast — Technical Requirements Document

## Architecture Overview

┌─────────────────────────────────────────────────────────────────┐
│ FRONTEND (Next.js) │
│ /chat → Tourist chat UI │
│ /sme → SME dashboard (accept/decline, inventory) │
│ /demo → Full-screen demo mode (for judges) │
└──────────────────────────┬──────────────────────────────────────┘
│ HTTPS
▼
┌─────────────────────────────────────────────────────────────────┐
│ BACKEND (Go — single binary, subcommands) │
│ │
│ cmd/server/main.go │
│ ├── Gin (port 8080) │
│ │ POST /api/chat → proxy to ADK agent │
│ │ GET /api/sme/:id → SME profile │
│ │ POST /api/sme/:id/request → accept/decline │
│ │ GET /api/sme/:id/inventory → live stock │
│ │ GET /api/itinerary/:id → generated plan │
│ │ │
│ └── ADK Agent (port 8081, in-process goroutine) │
│ Concierge Agent │
│ ├── tool: discover_businesses() │
│ ├── tool: query_business_state() │
│ ├── tool: calculate_route() │
│ └── tool: generate_itinerary() │
└──────────────────────────┬──────────────────────────────────────┘
│
▼
┌─────────────────────────────────────────────────────────────────┐
│ DATA (in-memory for hackathon) │
│ │
│ sme_registry (Go map / slice) │
│ ├── 10–15 SME profiles │
│ ├── Each has: products, prices, inventory, hours, location │
│ └── Mutable (SME dashboard decrements stock) │
│ │
│ sessions (Go map) │
│ └── Active tourist conversations + generated itineraries │
└─────────────────────────────────────────────────────────────────┘


## Tech Stack

| Layer | Technology | Why |
|-------|-----------|-----|
| Frontend | Next.js 15 (App Router) + Tailwind | Fast to build, SSR for SEO later |
| Backend | Go 1.27 + Gin | Single binary, fast, ADK-native |
| AI/Agents | Google ADK for Go (v2) | Multi-agent, tool calls, Gemini |
| LLM | Gemini Flash (via API key) | Fast, cheap, good enough for demo |

## Data Model

### SME Profile (the core entity)

```go
// internal/models/sme.go
package models

type SME struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Category    string    `json:"category"`     // "food", "crafts", "experience"
    Location    Location  `json:"location"`
    Description string    `json:"description"`
    Story       string    `json:"story"`        // the "3 generations" narrative
    Languages   []string  `json:"languages"`
    Opening     string    `json:"opening"`      // "10:00"
    Closing     string    `json:"closing"`      // "18:00"
    Capacity    int       `json:"capacity"`     // max people at once
    Products    []Product `json:"products"`
}

type Location struct {
    Lat       float64 `json:"lat"`
    Lng       float64 `json:"lng"`
    Area      string  `json:"area"`       // "Old Town", "Marikiti"
    WalkFrom  string  `json:"walk_from"`  // "Mombasa Port"
    WalkMins  int     `json:"walk_mins"`
}

type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Price       int     `json:"price"`      // KSh
    PrepMinutes int     `json:"prep_minutes"` // 0 for retail, 35 for food
    Stock       int     `json:"stock"`      // remaining today
    MaxStock    int     `json:"max_stock"`  // capacity today
    Custom      bool    `json:"custom"`     // made-to-order?
}

type TouristRequest struct {
    ID        string   `json:"id"`
    SMEID     string   `json:"sme_id"`
    TouristID string   `json:"tourist_id"`
    Items     []Item   `json:"items"`
    PickupAt  string   `json:"pickup_at"`  // "13:15"
    Total     int      `json:"total"`
    Status    string   `json:"status"`     // "pending", "accepted", "declined"
}

type Item struct {
    ProductID string `json:"product_id"`
    Name      string `json:"name"`
    Qty       int    `json:"qty"`
    Price     int    `json:"price"`
}

Itinerary
type Itinerary struct {
    ID        string      `json:"id"`
    TouristID string      `json:"tourist_id"`
    Stops     []Stop      `json:"stops"`
    TotalCost int         `json:"total_cost"`
    Budget    int         `json:"budget"`
    StartTime string      `json:"start_time"`
    EndTime   string      `json:"end_time"`
}

type Stop struct {
    Time      string  `json:"time"`
    Title     string  `json:"title"`
    Category  string  `json:"category"`
    Cost      int     `json:"cost"`
    Duration  int     `json:"duration_mins"`
    SMEID     string  `json:"sme_id,omitempty"`
    RequestID string  `json:"request_id,omitempty"`
    Note      string  `json:"note"`
}

Agent Design (ADK)
Agent Topology
ConciergeAgent (root)
│
├── Tools (direct):
│   ├── parse_tourist_intent()     → extracts time, budget, interests
│   └── generate_itinerary()       → assembles final plan
│
├── Sub-agents:
│   ├── DiscoveryAgent
│   │   └── Tools:
│   │       ├── search_businesses(category, area, radius)
│   │       └── search_experiences(interest, area)
│   │
│   └── BusinessAgent
│       └── Tools:
│           ├── query_inventory(sme_id)
│           ├── check_availability(sme_id, time, party_size)
│           ├── get_prep_time(sme_id, product_id)
│           └── calculate_walk_time(from, to)

Why only 3 agents? For a hackathon, 3 agents with clear responsibilities
beats 7 agents that overlap. The concierge orchestrates, discovery finds
candidates, business validates feasibility. The planner logic lives in the
concierge's generate_itinerary tool (it's a constraint-satisfaction problem,
not a "reasoning" problem that needs its own agent).

Tool Implementations (Go)
// agent/tools.go

// search_businesses: filters the in-memory registry
func searchBusinesses(_ tool.Context, args struct {
    Category string `json:"category" jsonschema:"food, crafts, or experience"`
    Area     string `json:"area" jsonschema:"Area name, e.g. Old Town"`
}) ([]models.SME, error) {
    return registry.SearchByCategory(args.Category, args.Area)
}

// query_inventory: returns live stock for an SME
func queryInventory(_ tool.Context, args struct {
    SMEID string `json:"sme_id" jsonschema:"The business ID"`
}) ([]models.Product, error) {
    sme, ok := registry.Get(args.SMEID)
    if !ok {
        return nil, fmt.Errorf("SME not found: %s", args.SMEID)
    }
    return sme.Products, nil
}

// check_availability: can this SME fulfill the request at this time?
func checkAvailability(_ tool.Context, args struct {
    SMEID      string `json:"sme_id"`
    ProductID  string `json:"product_id"`
    Qty        int    `json:"qty"`
    RequestedAt string `json:"requested_at" jsonschema:"HH:MM format"`
}) (bool, string, error) {
    sme, _ := registry.Get(args.SMEID)
    product, _ := sme.FindProduct(args.ProductID)
    if product == nil {
        return false, "Product not found", nil
    }
    if product.Stock < args.Qty {
        return false, fmt.Sprintf("Only %d left", product.Stock), nil
    }
    // Check if prep time fits before closing
    readyAt := addMinutes(args.RequestedAt, product.PrepMinutes)
    if readyAt > sme.Closing {
        return false, "Wouldn't be ready before closing", nil
    }
    return true, fmt.Sprintf("Ready by %s", readyAt), nil
}

// generate_itinerary: constraint solver
func generateItinerary(_ tool.Context, args struct {
    StartTime string   `json:"start_time"`
    Duration  int      `json:"duration_mins"`
    Budget    int      `json:"budget"`
    Stops     []StopIn `json:"stops"`
}) (models.Itinerary, error) {
    // Greedy: sort by time, accumulate cost, drop if over budget
    // Return the itinerary
}

Concierge Agent Config
// agent/agent.go

func NewConcierge(ctx context.Context, model *gemini.Model) (*llmagent.Agent, error) {
    discovery, _ := llmagent.New(llmagent.Config{
        Name:        "discovery_agent",
        Model:       model,
        Description: "Finds businesses and experiences matching tourist interests",
        Instruction: `You search for local businesses and experiences.
Use your tools to find candidates. Return a list of promising options
with their IDs, names, categories, and locations.`,
        Tools: []tool.Tool{
            functiontool.Must(searchBusinesses),
            functiontool.Must(searchExperiences),
        },
    })

    business, _ := llmagent.New(llmagent.Config{
        Name:        "business_agent",
        Model:       model,
        Description: "Queries SME live state: inventory, availability, prep time",
        Instruction: `You verify whether a specific business can fulfill a
specific request. Check inventory, prep time, and opening hours.
Return: can they do it? When will it be ready? What's the total cost?`,
        Tools: []tool.Tool{
            functiontool.Must(queryInventory),
            functiontool.Must(checkAvailability),
            functiontool.Must(getPrepTime),
            functiontool.Must(calculateWalkTime),
        },
    })

    concierge, _ := llmagent.New(llmagent.Config{
        Name:        "concierge",
        Model:       model,
        Description: "Main orchestrator: understands tourist intent, delegates, plans",
        Instruction: `You are the Hidden Coast concierge for Mombasa.
A tourist tells you their time, budget, and interests.
1. Parse their constraints.
2. Ask the discovery agent to find candidates.
3. For each candidate, ask the business agent to verify feasibility.
4. Generate an itinerary that fits time + budget.
5. Present it clearly with timestamps, costs, and total.
Be warm, specific, and confident. Reference real business names and
their stories.`,
        SubAgents: []agent.Agent{discovery, business},
        Tools: []tool.Tool{
            functiontool.Must(generateItinerary),
        },
    })

    return concierge, nil
}

API Design (Gin)
POST /api/chat
  Body: { "message": "I'm in Mombasa for 3 hours...", "session_id": "..." }
  Response (SSE stream): agent thinking + final response

GET /api/sme/:id
  Response: full SME profile

GET /api/sme/:id/inventory
  Response: [{ product, price, stock_remaining }]

POST /api/sme/:id/request
  Body: { "request_id": "...", "action": "accept" | "decline" }
  Response: updated request status + decremented inventory

GET /api/itinerary/:id
  Response: full itinerary with stops, costs, times

GET /api/health
  Response: { "status": "ok" }

Frontend (Next.js)
frontend/
├── app/
│   ├── layout.tsx
│   ├── page.tsx              → redirect to /chat
│   ├── chat/
│   │   └── page.tsx          → tourist chat (the demo)
│   ├── sme/
│   │   ├── page.tsx          → SME list (pick who you are)
│   │   └── [id]/page.tsx     → SME dashboard (requests + inventory)
│   └── demo/
│       └── page.tsx          → full-screen demo (for judges)
├── components/
│   ├── ChatPanel.tsx         → message list + input
│   ├── ItineraryCard.tsx     → the 3-hour plan display
│   ├── SMEDashboard.tsx      → accept/decline + stock
│   ├── AgentTrace.tsx        → shows agent thinking (for demo)
│   └── BudgetBar.tsx         → visual budget tracker
├── lib/
│   ├── api.ts                → fetch wrappers
│   └── types.ts              → shared types
├── tailwind.config.ts
└── package.json

Project Structure (Go backend)
backend/
├── cmd/
│   └── server/
│       └── main.go           ← Gin + ADK in one binary
├── internal/
│   ├── models/
│   │   ├── sme.go
│   │   ├── itinerary.go
│   │   └── request.go
│   ├── registry/
│   │   ├── registry.go       ← in-memory store (map[string]SME)
│   │   └── seed.go           ← 10–15 demo businesses
│   └── handlers/
│       ├── chat.go
│       ├── sme.go
│       └── itinerary.go
├── agent/
│   ├── agent.go              ← ADK agent config
│   ├── tools.go              ← tool implementations
│   └── prompts.go            ← instruction strings
├── Dockerfile
├── go.mod
└── Makefile

Deployment
Single container (hackathon simplicity)

One Dockerfile, one Cloud Run service. Gin on :8080, ADK agent in-process
on :8081 (or embedded in the same Gin router via a handler).

FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]

Dockerfiles
# backend/Dockerfile
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]

# frontend/Dockerfile
FROM node:22-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./
EXPOSE 3000
CMD ["npm", "start"]

render.yaml (Docker Runtime, Free Tier)
# render.yaml (repo root)
services:
  - type: web
    name: hidden-coast-backend
    runtime: docker
    dockerContext: ./backend
    dockerfilePath: ./backend/Dockerfile
    plan: free
    region: oregon
    envVars:
      - key: GOOGLE_API_KEY
        sync: false          # set in Render dashboard
      - key: PORT
        value: "8080"

  - type: web
    name: hidden-coast-frontend
    runtime: docker
    dockerContext: ./frontend
    dockerfilePath: ./frontend/Dockerfile
    plan: free
    region: oregon
    envVars:
      - key: NEXT_PUBLIC_API_URL
        value: https://hidden-coast.onrender.com

Makefile
# backend/Makefile

.PHONY: run test build docker clean

run:
	go run ./cmd/server

test:
	go test -race ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server ./cmd/server

docker:
	docker build -t hidden-coast:latest .

clean:
	rm -rf bin/

docker-compose.yml (root)
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - GOOGLE_API_KEY=${GOOGLE_API_KEY}
    healthcheck:
      test: ["CMD", "/app/server", "-health"]
      interval: 5s
      timeout: 3s

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080
    depends_on:
      - backend


render.yaml (Docker Runtime, Free Tier)
# render.yaml (repo root)
services:
  - type: web
    name: hidden-coast-backend
    runtime: docker
    dockerContext: ./backend
    dockerfilePath: ./backend/Dockerfile
    plan: free
    region: oregon
    envVars:
      - key: GOOGLE_API_KEY
        sync: false          # set in Render dashboard
      - key: PORT
        value: "8080"

  - type: web
    name: hidden-coast-frontend
    runtime: docker
    dockerContext: ./frontend
    dockerfilePath: ./frontend/Dockerfile
    plan: free
    region: oregon
    envVars:
      - key: NEXT_PUBLIC_API_URL
        value: https://hidden-coast-backend.onrender.com

Workflow: You as Architect, Jules as Implementer
YOU (architect)                    JULES (implementer)
─────────────────                  ───────────────────
Write PRD.md, TRD.md, sprints.md
Define the data models
Decide agent topology
        │
        ▼
Create GitHub issues
(one per task, with the spec)
        │
        ▼
Label issue with "jules"
(or assign via jules.google.com)
        │
        │                         Clones repo
        │                         Reads PRD + TRD + existing code
        │                         Implements the task
        │                         Opens PR
        │
        ▼
Review PR in GitHub (or GoLand)
Merge or request changes
        │
        ▼
GitHub Actions runs tests
Render auto-deploys on merge
