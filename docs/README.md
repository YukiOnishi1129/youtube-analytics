# YouTube Viral Theme Detection System - Documentation

Technical specifications for a video trend analysis system designed to support YouTube channel management for engineers.

## 📚 Documentation Structure

### 🎯 Overview
- [Background and Purpose](./overview/background.md) - Project background, objectives, and expected outcomes

### 🏗️ Architecture
- [System Overview](./architecture/system-overview.md) - Overall architecture and microservice composition
- [Backend Design](./architecture/backend-design.md) - Clean Architecture + DDD details
- [Clean Architecture Details](./architecture/clean-architecture.md) - Layer composition and DDD implementation patterns
- [Service Details](./architecture/services.md) - Detailed specifications for auth-svc and yt-svc
- [Testing Strategy](./architecture/testing-strategy.md) - Test pyramid and idempotency testing

### 💾 Database
- [Schema Design](./database/schema.md) - Ingestion table definitions
- [Analytics Tables](./database/analytics-tables.md) - Analysis table definitions
- [Authentication Tables](./database/auth-tables.md) - Account and permission management
- [Metrics Calculation](./database/metrics.md) - Derived metric calculation formulas

### 🔌 API Specifications
- [Analytics API](./api/analytics-proto.md) - Rankings, history, and video details
- [Ingestion API](./api/ingestion-proto.md) - Channel and keyword management
- [Authority API](./api/authority-proto.md) - Authentication and profiles

### 🎨 Frontend
- [Architecture](./frontend/architecture.md) - Next.js design and directory structure
- [Screen Specifications](./frontend/screens.md) - Detailed specifications for each screen
- [Domain Type Definitions](./frontend/domain-types.md) - TypeScript type definitions

### 🚀 Deployment
- [Environment Variables](./deployment/environment.md) - Required environment variable list
- [Schedule Settings](./deployment/schedule.md) - Cloud Scheduler/Tasks configuration
- [Security](./deployment/security.md) - Authentication flow and security settings
- [Infrastructure](./deployment/infrastructure.md) - Cloud Run and resource configuration

### 📋 Use Cases
- [Automated Processes](./use-cases/automated-processes.md) - Batch processing details
- [User Operations](./use-cases/user-operations.md) - Screen-based operations
- [Workflows](./use-cases/workflows.md) - Main processing flows
- [Implementation Tasks](./use-cases/implementation-tasks.md) - MVP implementation sequence

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

For detailed implementation steps, refer to [Implementation Tasks](./use-cases/implementation-tasks.md).

## 📝 License

This documentation is exclusively for Yuki's project.