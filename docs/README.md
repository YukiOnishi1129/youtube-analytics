# YouTube Viral Theme Detection System - Documentation

Technical specifications for a video trend analysis system designed to support YouTube channel management for engineers.

## 📚 Documentation Structure

### 🎯 01-overview/ - Project Background
- [01-background.md](./01-overview/01-background.md) - Project purpose, objectives, and expected outcomes

### 📋 02-use-cases/ - System Requirements  
- [01-README.md](./02-use-cases/01-README.md) - Use cases overview and categories
- [02-data-collection.md](./02-use-cases/02-data-collection.md) - Automated data collection processes
- [03-admin-management.md](./02-use-cases/03-admin-management.md) - Administrative functions and management
- [04-user-operations.md](./02-use-cases/04-user-operations.md) - End-user features and analytics
- [05-automated-processes.md](./02-use-cases/05-automated-processes.md) - System automation and batch jobs
- [06-workflows.md](./02-use-cases/06-workflows.md) - Business workflows and processes
- [07-management-operations.md](./02-use-cases/07-management-operations.md) - Operational management tasks
- [08-implementation-tasks.md](./02-use-cases/08-implementation-tasks.md) - MVP implementation checklist

### 🧠 03-domain/ - Business Logic
- [01-ubiquitous-language.md](./03-domain/01-ubiquitous-language.md) - Domain terminology and definitions
- [02-bounded-contexts.md](./03-domain/02-bounded-contexts.md) - Context boundaries and responsibilities
- [03-aggregates.md](./03-domain/03-aggregates.md) - Domain aggregates, entities, and invariants
- [04-domain-services.md](./03-domain/04-domain-services.md) - Cross-aggregate logic and domain services
- [05-domain-flows.md](./03-domain/05-domain-flows.md) - Key business process flows
- [06-application-services.md](./03-domain/06-application-services.md) - Use case orchestration layer

### 💾 04-database/ - Data Design
- [01-schema-overview.md](./04-database/01-schema-overview.md) - Database architecture overview
- [02-ingestion-tables.md](./04-database/02-ingestion-tables.md) - Data collection and storage tables
- [03-analytics-tables.md](./04-database/03-analytics-tables.md) - Analytics and metrics tables
- [04-auth-tables.md](./04-database/04-auth-tables.md) - Authentication and authorization tables
- [05-metrics.md](./04-database/05-metrics.md) - Metric calculation formulas
- [06-schema-migration-guide.md](./04-database/06-schema-migration-guide.md) - Migration guide for multi-genre support

### 🔌 05-api/ - Service Interfaces
- [01-authority-proto.md](./05-api/01-authority-proto.md) - Authentication service API
- [02-ingestion-proto.md](./05-api/02-ingestion-proto.md) - Data ingestion service API
- [03-analytics-proto.md](./05-api/03-analytics-proto.md) - Analytics service API

### 🏗️ 06-architecture/ - Technical Design
- [01-system-overview.md](./06-architecture/01-system-overview.md) - High-level system architecture
- [02-clean-architecture.md](./06-architecture/02-clean-architecture.md) - Clean Architecture principles and patterns
- [03-services.md](./06-architecture/03-services.md) - Detailed service specifications
- [04-backend-design.md](./06-architecture/04-backend-design.md) - Backend implementation details
- [05-backend-final-plan.md](./06-architecture/05-backend-final-plan.md) - Final implementation plan
- [06-testing-strategy.md](./06-architecture/06-testing-strategy.md) - Testing approach and strategies

### 🎨 07-frontend/ - UI Implementation
- [01-architecture.md](./07-frontend/01-architecture.md) - Frontend architecture and structure
- [02-domain-types.md](./07-frontend/02-domain-types.md) - TypeScript domain type definitions
- [03-screens.md](./07-frontend/03-screens.md) - Screen specifications and wireframes

### 🚀 08-deployment/ - Operations
- [01-environment.md](./08-deployment/01-environment.md) - Environment variables and configuration
- [02-infrastructure.md](./08-deployment/02-infrastructure.md) - Cloud infrastructure components
- [03-security.md](./08-deployment/03-security.md) - Security considerations and practices
- [04-schedule.md](./08-deployment/04-schedule.md) - Batch job and scheduler configuration

## 🎯 System Purpose

Quantitatively discover "what themes should be used for videos to gain traction" and accelerate YouTube channel growth.

### Key Features
1. **Multi-Genre Support** - Track videos across different regions, languages, and categories
2. **Competitor Channel Tracking** - WebSub real-time detection + periodic snapshots
3. **Smart Video Discovery** - Genre-specific keyword filtering from trending videos
4. **Theme Rankings** - Automatic extraction of successful video patterns
5. **Growth Analysis** - Track view and like count changes over time (0/3/6/12/24/48/72/168h)
6. **Admin Management** - Configure genres, keywords, and categories through admin portal

## 🛠️ Technology Stack

### Backend
- **Language**: Go
- **Architecture**: Clean Architecture + DDD (Hexagonal)
- **Communication**: gRPC (internal) + HTTP (external events)
- **Database**: PostgreSQL (Neon)
- **Infrastructure**: Cloud Run, Cloud Tasks, Cloud Scheduler

### Frontend
- **Framework**: Next.js (App Router)
- **UI**: shadcn/ui
- **State Management**: TanStack Query
- **Authentication**: Auth.js + Identity Platform

## 📖 Reading Guide

1. **New to the project?** Start with [01-overview/01-background.md](./01-overview/01-background.md)
2. **Understanding features?** Read [02-use-cases/01-README.md](./02-use-cases/01-README.md)
3. **Learning the domain?** Begin with [03-domain/01-ubiquitous-language.md](./03-domain/01-ubiquitous-language.md)
4. **Implementing features?** Follow the numbered sequence in each directory

## 🚦 Quick Start

For detailed implementation steps, refer to [02-use-cases/08-implementation-tasks.md](./02-use-cases/08-implementation-tasks.md).

## 📝 License

This documentation is exclusively for Yuki's project.