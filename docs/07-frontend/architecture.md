# Frontend Architecture (Next.js)

## Design Principles

- **Domain-driven separation**: app/ features/ shared/ external/
- **Server Components First**, Client Components only when necessary
- **CQRS pattern** adopted in external/handlers (*.command.ts / *.query.ts)
- **Type safety**: All types consolidated in shared/types

## Directory Layout

```
src/
├── app/         # Next.js App Router (pages/layouts only, RSC)
├── features/    # Domain-based UI & logic
├── shared/      # Reusable UI/types/utilities
└── external/    # External service integration
    ├── client/grpc/       # buf-generated TS client code
    ├── services/          # External data → frontend domain transformation
    └── handlers/          # Server functions/actions (CQRS, 'use server')
```

## External Integration Flow

- **client/grpc**: buf generated code (auto-updated, direct import prohibited)
- **services**: Calls generated clients, transforms external DTOs to frontend models
- **handlers**: command/query server functions, exposes Server Actions with 'use server'
- **features**: imports handler actions (UI doesn't directly touch generated code)

## Proto → TS SDK

- **Output**: web/client/src/external/client/grpc
- **Generation**: buf (connect-es or ts-proto)
- **Usage**: Only via handlers/services
- **Generated code handling**: Committed in same repo (if multiple frontends use it in future, private publish to GitHub Packages)

## Test Strategy

- **Unit**: Vitest for features/hooks/components
- **Integration**: features and external.handlers integration
- **Mock**: UI testing with mocked server actions
- **Storybook**: Development and documentation purposes only
