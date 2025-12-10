# Project Overview (Developer Guide)

This document provides a technical overview of the Edmonton Bike Map architecture and data model. For detailed setup, testing, and deployment instructions, see the other docs in this directory.

## System Architecture

Edmonton Bike Map is a full-stack web application providing bike-specific routing and user-generated reviews for Edmonton cyclists. The system has three layers:

1. **Frontend** (Svelte/TypeScript on Vercel) — Single-page application with route planning UI
2. **Backend** (Go on DigitalOcean) — REST API with graph-based routing algorithm
3. **Database** (PostgreSQL) — Nodes (intersections), ways (roads/paths), reviews, and user accounts

```
Frontend (Vercel)                  Backend (DigitalOcean)         Database (PostgreSQL)
┌──────────────────────────┐       ┌──────────────────────┐       ┌─────────────────┐
│ Svelte + SvelteKit       │       │ Go HTTP Server       │       │ Nodes           │
│ TailwindCSS              │       │ ├─ Handlers          │       │ Ways            │
│ Vite + TypeScript        │◄─────►│ ├─ Routing Algorithm │◄─────►│ Reviews         │
└──────────┬───────────────┘       │ ├─ Auth & Middleware │       │ Users           │
           │                       │ └─ Repositories      │       └─────────────────┘
           │ HTTPS                 └──────────────────────┘
           │
        (Browser)
```

## Data Model

The database schema is versioned in `backend/migrations/`. At a high level:

- **nodes** — Intersections and key points with latitude/longitude coordinates
- **ways** — Roads, paths, and bike infrastructure (from OpenStreetMap)
- **way_nodes** — Junction table linking ways to their ordered node sequences
- **reviews** — User ratings and comments on ways
- **users** — User accounts and preferences

See `backend/migrations/` for the authoritative schema definitions.

## Backend Architecture

### Routing Algorithm

Located in `internal/domain/routing/`. The algorithm:

1. Snaps start and end coordinates to the nearest nodes in the network
2. Runs Dijkstra's algorithm to find the lowest-weighted path between those nodes
3. Returns the path as an ordered list of nodes and the total distance in km

Edge weights are pre-computed when the network graph is built.

### Layered Architecture

The backend follows a clean layered architecture:

- **Models** (`internal/models/`) — Data structures representing domain objects (Node, Way, Review, User)
- **Repositories** (`internal/repository/`) — Database access layer. Repositories wrap `*sql.DB` and handle all SQL queries. Named `*Repository` (e.g., `NodeRepository`, `WayRepository`)
- **Services** (`internal/service/`) — Business logic layer. Services orchestrate repositories and domain logic, handling workflows like route finding or user authentication
- **Handlers** (`internal/handler/`) — HTTP request handlers. Handlers parse HTTP requests, call services, and return JSON responses

### Directory Structure

```
backend/
├── cmd/
│   ├── server/           # HTTP server entrypoint
│   ├── update_db/        # Updates the database with latest OpenStreetMap data (nodes / ways)
│   └── jsontodb/         # CSV data import utility
├── internal/
│   ├── domain/           # Core logic (routing, network, geo)
│   ├── handler/          # HTTP endpoint handlers
│   ├── middleware/       # Auth, CORS, etc.
│   ├── models/           # Data structures
│   ├── repository/       # Data access layer
│   ├── server/           # HTTP server setup & routing
│   ├── service/          # Business logic
│   ├── token/            # JWT utilities
│   ├── updater/          # OSM data update service (used by update_db)
│   └── utils/            # Utility functions
└── migrations/           # SQL schema (versioned)
```

## Frontend Architecture

### Structure

```
src/
├── routes/                 # SvelteKit file-based routing
│   ├── +layout.svelte      # Root layout wrapper
│   ├── map/                # Map/route planning feature
│   ├── login/, signup/     # Authentication pages
│   ├── settings/           # User preferences
│   └── api/                # Server-side API endpoints
├── lib/
│   ├── components/         # Reusable UI components (Map, MenuBar, Sidebar, etc.)
│   ├── api/
│   │   ├── client/         # Client-side API functions (reviews.ts, settings.ts)
│   │   └── server/         # Server-side API logic
│   ├── map/                # Map utilities (LeafletMap, mapActions, mapModes)
│   ├── state.svelte.ts     # Selected Way state
│   └── types.ts            # TypeScript type definitions
└── app.html                # Root HTML template
```

### Data Flow and Authentication

Data may pass between the client (browser), frontend server (SvelteKit), and backend server (Go API). Where possible, data is loaded on the frontend server side using SvelteKit's SSR capabilities. Authentication is handled by the frontend server. For dynamic updates, requests may originate on the client. In these cases, the client will send a request to the frontend server, which will attach the authorization header and then send the request on to the backend. The response will then be passed back to the client from the frontend server. The browser does not have the authorization token so it cannot make requests directly to the backend.
