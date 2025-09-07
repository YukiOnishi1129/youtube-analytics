# YouTube Viral Theme Detection System - Documentation

Technical specifications for a video trend analysis system designed to support YouTube channel management for engineers.

## 📚 Documentation Structure

### 🎯 Overview
- [Background and Purpose](./01-overview/background.md) - Project background, objectives, and expected outcomes

### 📋 Use Cases
- [Automated Processes](./02-use-cases/automated-processes.md) - Batch processing details
- [User Operations](./02-use-cases/user-operations.md) - Screen-based operations
- [Data Collection](./02-use-cases/data-collection.md) - Ingestion-related user flows
- [Workflows](./02-use-cases/workflows.md) - Main processing flows
- [Implementation Tasks](./02-use-cases/implementation-tasks.md) - MVP implementation sequence
- [Management & Operations](./02-use-cases/management-operations.md) - Admin/ops flows and tasks

### 🧠 Domain
- [Ubiquitous Language](./03-domain/ubiquitous-language.md) - Shared vocabulary and definitions
- [Bounded Contexts](./03-domain/bounded-contexts.md) - Context boundaries and responsibilities
- [Aggregates](./03-domain/aggregates.md) - Transactional boundaries and invariants
- [Application Services](./03-domain/application-services.md) - Use case orchestration layer
- [Domain Services & Policies](./03-domain/domain-services.md) - Cross-aggregate logic and rules
- [Representative Flows](./03-domain/domain-flows.md) - Key domain flows and sequences

### 💾 Database
- [Schema Design](./04-database/schema.md) - Ingestion table definitions
- [Analytics Tables](./04-database/analytics-tables.md) - Analysis table definitions
- [Authentication Tables](./04-database/auth-tables.md) - Account and permission management
- [Metrics Calculation](./04-database/metrics.md) - Derived metric calculation formulas

### 🔌 API Specifications
- [Analytics API](./05-api/analytics-proto.md) - Rankings, history, and video details
- [Ingestion API](./05-api/ingestion-proto.md) - Channel and keyword management
- [Authority API](./05-api/authority-proto.md) - Authentication and profiles

### 🏗️ Architecture
- [System Overview](./06-architecture/system-overview.md) - Overall architecture and microservice composition
- [Backend Design](./06-architecture/backend-design.md) - Clean Architecture + DDD details
- [Clean Architecture Details](./06-architecture/clean-architecture.md) - Layer composition and DDD implementation patterns
- [Service Details](./06-architecture/services.md) - Detailed specifications for auth-svc and yt-svc
- [Testing Strategy](./06-architecture/testing-strategy.md) - Test pyramid and idempotency testing

### 🎨 Frontend
- [Architecture](./07-frontend/architecture.md) - Next.js design and directory structure
- [Screen Specifications](./07-frontend/screens.md) - Detailed specifications for each screen
- [Domain Type Definitions](./07-frontend/domain-types.md) - TypeScript type definitions

### 🚀 Deployment
- [Environment Variables](./08-deployment/environment.md) - Required environment variable list
- [Schedule Settings](./08-deployment/schedule.md) - Cloud Scheduler/Tasks configuration
- [Security](./08-deployment/security.md) - Authentication flow and security settings
- [Infrastructure](./08-deployment/infrastructure.md) - Cloud Run and resource configuration

### 📋 Use Cases
- [Automated Processes](./08-use-cases/automated-processes.md) - Batch processing details
- [User Operations](./08-use-cases/user-operations.md) - Screen-based operations
- [Data Collection](./08-use-cases/data-collection.md) - Ingestion-related user flows
- [Workflows](./08-use-cases/workflows.md) - Main processing flows
- [Implementation Tasks](./08-use-cases/implementation-tasks.md) - MVP implementation sequence
- [Management & Operations](./08-use-cases/management-operations.md) - Admin/ops flows and tasks

## 🎯 System Purpose

Quantitatively discover "what themes should be used for videos to gain traction" and accelerate YouTube channel growth.

### Key Features
1. **Competitor Channel Tracking** - WebSub real-time detection + periodic snapshots
2. **New Video Discovery** - Keyword filtering from trending videos
3. **Theme Rankings** - Automatic extraction of frequent themes
4. **Growth Analysis** - Tracking view and like count changes over time

## 🛠️ Technology Stack

### Backend
- **Language**: Go
- **Architecture**: Clean Architecture + DDD
- **Communication**: gRPC (internal) + HTTP (external events)
- **Database**: PostgreSQL (Neon)
- **Infrastructure**: Cloud Run, Cloud Tasks, Cloud Scheduler

### Frontend
- **Framework**: Next.js (App Router)
- **UI**: shadcn/ui
- **State Management**: TanStack Query
- **Authentication**: Auth.js + Identity Platform

## 🚦 Quick Start

For detailed implementation steps, refer to [Implementation Tasks](./02-use-cases/implementation-tasks.md).

## 📝 License

This documentation is exclusively for Yuki's project.
