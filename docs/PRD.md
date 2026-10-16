# Hidden Coast — Product Requirements Document

## One-liner
An AI-powered gateway that connects tourists' time, budget, and intent
with Mombasa's invisible local tourism economy.

## Problem
Tourists in Mombasa discover the same 5–6 attractions (Fort Jesus, Old Town,
Diani, Wasini, beaches). The thousands of small local businesses — food
vendors, artisans, cultural guides, experience providers — are invisible to
them. These SMEs have products, knowledge, and capacity that perfectly match
what tourists want, but no discovery mechanism connects the two.

## Target Users

### Primary: Tourist
- Arriving with limited time (2–6 hours)
- Has a budget (KSh 1,000–10,000)
- Wants "authentic" — not resort-generic
- Speaks English (possibly Swahili)
- On mobile, possibly with spotty data

### Secondary: Local SME
- Operates a small tourism-relevant business (food, crafts, experiences)
- No tech background
- Needs zero AI/agent/API knowledge
- Wants: "tourists find me, I accept or decline"

## Core Hypothesis
> Can agentic AI become the intelligent bridge between tourist demand and
> fragmented local tourism supply?

## The Key Differentiator
**The SME is not a listing. The SME is an AI-accessible node.**

Google Maps knows Mama Aisha exists.
Hidden Coast knows: Mama Aisha makes chicken biryani, has 12 portions left
today, takes 35 minutes to prepare, charges KSh 450/portion, is 12 minutes
from the tourist's current location, and is open until 18:00.

The agent doesn't describe the business. It **queries** it.

## User Journeys

### Tourist Journey (the demo)
1. Tourist arrives at Mombasa port, opens the app
2. Says: "I'm in Mombasa for 3 hours. I want authentic food, something
   local to buy, and I have KSh 3,000."
3. The concierge agent parses: time=3h, budget=KSh 3000, interests=[food,
   crafts, authentic]
4. Discovery agent finds candidate businesses within range
5. Business agent queries each candidate's live state (inventory, hours,
   capacity, prep time)
6. Planner agent assembles a feasible itinerary that fits time + budget
7. Tourist sees: "Your 3-hour Mombasa experience" with timestamps, costs,
   and total
8. Tourist confirms → requests are sent to SMEs
9. SMEs accept/decline
10. Tourist gets confirmed plan with directions

### SME Journey
1. Mama Aisha opens her simple dashboard (mobile web)
2. Sees: "New request: 2× Chicken Biryani, pickup 13:15, KSh 900"
3. Taps **Accept**
4. Sees inventory auto-decrement: 12 → 10 portions
5. End of day: "You served 4 tourists. Revenue: KSh 2,700."

## MVP Boundary (Hackathon Scope)

### IN
- [ ] Tourist chat interface (Next.js)
- [ ] ADK multi-agent orchestration (Concierge → Discovery → Business → Planner)
- [ ] 10–15 seeded SME profiles with structured knowledge
- [ ] Real-time inventory/availability query (in-memory, not a real DB)
- [ ] Itinerary generation with time + budget constraints
- [ ] SME dashboard (accept/decline + inventory view)
- [ ] One continuous 3-minute demo story
- [ ] Deployed on Cloud Run (or Render) with Next.js frontend

### OUT (post-hackathon)
- [ ] Real SME onboarding flow
- [ ] Payment integration (M-Pesa, cards)
- [ ] Real GPS / live location tracking
- [ ] Multi-city expansion
- [ ] Voice interface for SMEs
- [ ] WhatsApp integration
- [ ] Review/rating system
- [ ] Real Postgres (use in-memory for demo)

## Success Criteria (Judging)

| Criterion | What we demonstrate |
|-----------|-------------------|
| **Novelty** | SME as agent-accessible node (not a listing) |
| **Technical depth** | Multi-agent ADK orchestration, real tool calls |
| **Real-world impact** | Two-sided network, local economic empowerment |
| **Demo quality** | One continuous story, no "let me switch screens" |
| **Feasibility** | Clear path to real SMEs (SME dashboard is the proof) |

## Demo Story (3–5 minutes)

[0:00] "I just landed at Mombasa port. I've got 3 hours before my
flight. I want authentic food, something local to buy,
and I have KSh 3,000."

[0:30] Watch the agents work (show agent trace / thinking):
- Concierge parses constraints
- Discovery finds 4 candidates
- Business queries Mama Aisha: "Can you do 2 biryani by 13:15?"
- Business queries Karanja: "Do you have sandals under 2000?"
- Planner assembles timeline

[1:30] "Here's your 3-hour Mombasa experience:"
12:00 🍛 Mama Aisha — 2× Chicken Biryani — KSh 900
13:00 🏛️ Old Town cultural walk (30 min)
14:00 👞 Karanja Leather — Open sandals — KSh 1,500
15:00 🚕 Return to port
TOTAL: KSh 2,650 (under budget ✓)

[2:30] Switch to SME view:
"Mama Aisha's phone lights up: New request — 2× Biryani,
pickup 13:15. [ACCEPT]"
"Karanja's phone: Tourist wants sandals, budget 2000,
stock: 12. [ACCEPT]"

[3:30] Back to tourist: "Both confirmed. Mama Aisha will have your
biryani ready at 13:15. Karanja is 9 minutes from Old Town.
Enjoy Mombasa."

[4:00] Close: "This is what happens when a tourist's intent meets
a local business's reality — in real time, through agents."

