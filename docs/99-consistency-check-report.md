# Consistency Check Report

## Summary
This report identifies inconsistencies between domain models, use cases, and table designs. Updated: 2025-09-14

## Fixed Issues ✅

1. **Soft Delete in Domain Model** - Added DeletedAt field and clarified behavior
2. **Target Field in Domain Model** - Added TargetField to Keyword entity
3. **Channel Fields** - Added missing fields to Channel domain entity
4. **Subscription Count Naming** - Standardized to subscription_count
5. **VideoMeta Type** - Added VideoMeta value object definition
6. **Audit Tables** - Added audit_logs table
7. **Job History Tables** - Added batch_jobs table
8. **Channel Subscription Management** - Added to admin use cases
9. **Multi-Genre Registration** - Clarified M:N relationship in use cases
10. **Channel Snapshot Fields** - Added ViewCount and VideoCount

## Critical Issues to Fix

### 1. ❌ Missing Soft Delete in Domain Model
- **Problem**: Keyword domain doesn't specify soft delete behavior
- **Fix**: Update domain model to clarify soft delete behavior

### 2. ❌ Target Field Not in Domain Model  
- **Problem**: `keywords.target_field` exists in table but not in domain
- **Fix**: Add target_field to Keyword entity or remove from table

### 3. ❌ Channel Fields Mismatch
- **Problem**: Channel table has extra fields not in domain (description, country, view_count, video_count)
- **Fix**: Either add to domain model or move to separate tracking table

### 4. ❌ Subscription Count Naming
- **Problem**: Inconsistent naming: SubscriptionCount vs subscriber_count
- **Fix**: Standardize to subscription_count everywhere

### 5. ❌ Missing VideoMeta Type
- **Problem**: RegisterVideoFromTrending uses undefined VideoMeta type
- **Fix**: Define VideoMeta value object in domain model

### 6. ❌ Missing Audit Tables
- **Problem**: Admin use cases reference audit logs but no tables exist
- **Fix**: Add audit_logs table or clarify logging strategy

### 7. ❌ Missing Job History Tables
- **Problem**: Collection job history mentioned but no tables defined
- **Fix**: Add batch_jobs or collection_history table

## Medium Priority Issues

### 8. ⚠️ BuildPattern Domain Service Usage
- **Problem**: Domain service defined but not referenced in use cases
- **Recommendation**: Clarify when BuildPattern is used

### 9. ⚠️ Channel Subscription Management
- **Problem**: Domain commands exist but not in admin use cases
- **Recommendation**: Add subscription management to admin use cases

### 10. ⚠️ Video Multi-Genre Registration
- **Problem**: Domain supports multiple genres but use case describes singular
- **Recommendation**: Update use case to clarify M:N relationship

## Low Priority Issues

### 11. ℹ️ CheckpointHour Type Definition
- **Problem**: Domain uses type alias, table uses CHECK constraint
- **Note**: This is acceptable, just different representations

### 12. ℹ️ Genre CategoryIDs Array
- **Problem**: Using array instead of junction table
- **Note**: PostgreSQL arrays are acceptable for this use case

## Recommendations

1. **Immediate Actions**:
   - Fix naming inconsistencies
   - Add missing type definitions
   - Update domain model to match table fields

2. **Documentation Updates**:
   - Clarify soft delete behavior
   - Document audit logging strategy
   - Add missing use cases for channel subscription

3. **Consider Adding**:
   - audit_logs table
   - batch_jobs table
   - Clear VideoMeta value object definition