# Watchup Server Agent - API Specification

## Base URL

```
https://watchup.space
```

All endpoints are prefixed with the base URL.

---

## Table of Contents

1. [Authentication](#authentication)
2. [Agent Registration](#agent-registration)
3. [Metrics Submission](#metrics-submission)
4. [Alert Submission](#alert-submission)
5. [Configuration Retrieval](#configuration-retrieval)
6. [Error Responses](#error-responses)
7. [Rate Limiting](#rate-limiting)
8. [Database Schema](#database-schema)

---

## Authentication

### Overview

The API uses two types of authentication:

1. **Master API Key** - Used only for agent registration (one-time)
2. **Server Key** - Used for all ongoing operations (metrics, alerts, config)

### Authentication Headers

```http
Authorization: Bearer {server_key}
Content-Type: application/json
```

---

## Agent Registration

### POST /agent/register

Registers a new agent with a project. Free projects are limited to one agent.

#### Request

**Headers:**
```http
Content-Type: application/json
```

**Body:**
```json
{
  "project_id": "proj_12345",
  "master_api_key": "ma_abcdef123456",
  "server_identifier": "production-api-1"
}
```

**Field Descriptions:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project_id` | string | Yes | Unique project identifier from Watchup dashboard |
| `master_api_key` | string | Yes | Master API key for the project (from project settings) |
| `server_identifier` | string | Yes | Friendly name for the server (e.g., "production-api", "staging-db") |

**Validation Rules:**
- `project_id`: Must match pattern `proj_[a-zA-Z0-9]+`, length 8-50 characters
- `master_api_key`: Must match pattern `ma_[a-zA-Z0-9]+`, length 20-100 characters
- `server_identifier`: Alphanumeric with hyphens/underscores, 3-50 characters

#### Response

**Success (201 Created):**
```json
{
  "success": true,
  "server_key": "srv_89sd0a7f3b2c1d4e",
  "message": "Agent registered successfully",
  "project": {
    "id": "proj_12345",
    "name": "My Production Project",
    "plan": "free"
  },
  "server": {
    "id": "server_abc123",
    "identifier": "production-api-1",
    "registered_at": "2026-03-25T14:30:00Z"
  }
}
```

**Error - Already Registered (409 Conflict):**
```json
{
  "success": false,
  "error": "agent_limit_reached",
  "message": "This project already has an agent installed. Free projects are limited to one agent.",
  "details": {
    "existing_server": "production-api-1",
    "registered_at": "2026-03-20T10:15:00Z",
    "plan": "free",
    "upgrade_url": "https://watchup.space/upgrade"
  }
}
```

**Error - Invalid Master API Key (401 Unauthorized):**
```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Invalid master API key"
}
```

**Error - Project Not Found (404 Not Found):**
```json
{
  "success": false,
  "error": "project_not_found",
  "message": "Project with ID 'proj_12345' does not exist"
}
```

#### Backend Implementation Notes

1. **Validate Master API Key**
   - Check if `master_api_key` matches the project's stored key
   - Verify the key is active and not revoked

2. **Check Agent Limit**
   - Query database for existing agents with this `project_id`
   - If project plan is "free" and agent count >= 1, return 409 error
   - If project plan is "paid", allow multiple agents

3. **Generate Server Key**
   - Generate unique server key: `srv_` + random 16-character alphanumeric
   - Ensure uniqueness in database

4. **Store Agent Record**
   ```sql
   INSERT INTO agents (
     id, project_id, server_key, server_identifier, 
     registered_at, last_seen, status
   ) VALUES (
     uuid(), 'proj_12345', 'srv_89sd0a7f3b2c1d4e', 
     'production-api-1', NOW(), NOW(), 'active'
   );
   ```

5. **Update Project Flag**
   ```sql
   UPDATE projects 
   SET agent_installed = true, 
       agent_count = agent_count + 1
   WHERE id = 'proj_12345';
   ```

---

## Metrics Submission

### POST /server/metrics

Submits system metrics from the agent. Called every 5 seconds.

#### Request

**Headers:**
```http
Authorization: Bearer srv_89sd0a7f3b2c1d4e
Content-Type: application/json
```

**Body:**
```json
{
  "server_key": "srv_89sd0a7f3b2c1d4e",
  "cpu": 45.2,
  "ram": 67.8,
  "timestamp": "2026-03-25T14:30:00Z"
}
```

**Field Descriptions:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `server_key` | string | Yes | Unique server key from registration |
| `cpu` | float | Yes | CPU usage percentage (0-100) |
| `ram` | float | Yes | RAM usage percentage (0-100) |
| `timestamp` | string | Yes | ISO 8601 timestamp (UTC) |

**Validation Rules:**
- `cpu`: 0.0 to 100.0
- `ram`: 0.0 to 100.0
- `timestamp`: Valid ISO 8601 format, not more than 1 minute in the future

#### Response

**Success (200 OK):**
```json
{
  "success": true,
  "message": "Metrics received",
  "metric_id": "metric_xyz789",
  "next_report": "2026-03-25T14:30:05Z"
}
```

**Error - Invalid Server Key (401 Unauthorized):**
```json
{
  "success": false,
  "error": "invalid_server_key",
  "message": "Server key is invalid or has been revoked"
}
```

**Error - Invalid Data (400 Bad Request):**
```json
{
  "success": false,
  "error": "validation_error",
  "message": "Invalid metric values",
  "details": {
    "cpu": "Must be between 0 and 100",
    "ram": "Must be between 0 and 100"
  }
}
```

#### Backend Implementation Notes

1. **Authenticate Server Key**
   ```sql
   SELECT id, project_id, server_identifier, status 
   FROM agents 
   WHERE server_key = 'srv_89sd0a7f3b2c1d4e' 
   AND status = 'active';
   ```

2. **Validate Metrics**
   - Ensure CPU and RAM are within 0-100 range
   - Validate timestamp is recent (within last 5 minutes)

3. **Store Metrics**
   ```sql
   INSERT INTO metrics (
     id, agent_id, project_id, cpu_usage, ram_usage, 
     recorded_at, created_at
   ) VALUES (
     uuid(), 'agent_id', 'proj_12345', 45.2, 67.8,
     '2026-03-25T14:30:00Z', NOW()
   );
   ```

4. **Update Agent Last Seen**
   ```sql
   UPDATE agents 
   SET last_seen = NOW() 
   WHERE server_key = 'srv_89sd0a7f3b2c1d4e';
   ```

5. **Aggregate for Dashboard**
   - Calculate 1-minute, 5-minute, 1-hour averages
   - Store in time-series database (optional: InfluxDB, TimescaleDB)
   - Trigger real-time dashboard updates via WebSocket

6. **Data Retention**
   - Raw metrics: Keep for 7 days
   - 1-minute aggregates: Keep for 30 days
   - 1-hour aggregates: Keep for 1 year

---

## Alert Submission

### POST /server/alerts

Submits an alert when a threshold is exceeded for a sustained period.

#### Request

**Headers:**
```http
Authorization: Bearer srv_89sd0a7f3b2c1d4e
Content-Type: application/json
```

**Body:**
```json
{
  "server_key": "srv_89sd0a7f3b2c1d4e",
  "metric": "cpu",
  "usage": 87.2,
  "duration": 300,
  "top_processes": [
    {
      "pid": 4213,
      "name": "node",
      "cpu": 52.1,
      "memory": 8.3
    },
    {
      "pid": 9021,
      "name": "python",
      "cpu": 18.4,
      "memory": 12.1
    },
    {
      "pid": 1234,
      "name": "chrome",
      "cpu": 15.2,
      "memory": 22.5
    }
  ],
  "timestamp": "2026-03-25T14:35:00Z"
}
```

**Field Descriptions:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `server_key` | string | Yes | Unique server key from registration |
| `metric` | string | Yes | Metric type: "cpu", "ram", or "process_cpu" |
| `usage` | float | Yes | Current usage percentage |
| `duration` | integer | Yes | Duration threshold was exceeded (seconds) |
| `top_processes` | array | Yes | Top 5 processes consuming resources |
| `top_processes[].pid` | integer | Yes | Process ID |
| `top_processes[].name` | string | Yes | Process name |
| `top_processes[].cpu` | float | Yes | Process CPU usage (%) |
| `top_processes[].memory` | float | Yes | Process memory usage (%) |
| `timestamp` | string | Yes | ISO 8601 timestamp (UTC) |

**Validation Rules:**
- `metric`: Must be one of: "cpu", "ram", "process_cpu"
- `usage`: 0.0 to 100.0
- `duration`: Positive integer (typically 120-600 seconds)
- `top_processes`: Array of 1-5 items
- `timestamp`: Valid ISO 8601 format

#### Response

**Success (201 Created):**
```json
{
  "success": true,
  "message": "Alert received and processed",
  "alert_id": "alert_abc123",
  "severity": "warning",
  "notifications_sent": {
    "email": true,
    "sms": false,
    "webhook": true,
    "dashboard": true
  },
  "created_at": "2026-03-25T14:35:00Z"
}
```

**Error - Invalid Server Key (401 Unauthorized):**
```json
{
  "success": false,
  "error": "invalid_server_key",
  "message": "Server key is invalid or has been revoked"
}
```

**Error - Invalid Data (400 Bad Request):**
```json
{
  "success": false,
  "error": "validation_error",
  "message": "Invalid alert data",
  "details": {
    "metric": "Must be one of: cpu, ram, process_cpu",
    "top_processes": "Must contain 1-5 processes"
  }
}
```

#### Backend Implementation Notes

1. **Authenticate Server Key**
   ```sql
   SELECT a.id, a.project_id, a.server_identifier, p.notification_settings
   FROM agents a
   JOIN projects p ON a.project_id = p.id
   WHERE a.server_key = 'srv_89sd0a7f3b2c1d4e' 
   AND a.status = 'active';
   ```

2. **Determine Severity**
   ```javascript
   function determineSeverity(metric, usage, duration) {
     if (metric === 'cpu') {
       if (usage >= 95) return 'critical';
       if (usage >= 85) return 'warning';
       return 'info';
     }
     if (metric === 'ram') {
       if (usage >= 90) return 'critical';
       if (usage >= 80) return 'warning';
       return 'info';
     }
     if (metric === 'process_cpu') {
       if (usage >= 80) return 'critical';
       if (usage >= 60) return 'warning';
       return 'info';
     }
   }
   ```

3. **Store Alert**
   ```sql
   INSERT INTO alerts (
     id, agent_id, project_id, metric_type, usage_value,
     duration_seconds, severity, top_processes, 
     triggered_at, acknowledged, resolved
   ) VALUES (
     uuid(), 'agent_id', 'proj_12345', 'cpu', 87.2,
     300, 'warning', '{"processes": [...]}',
     '2026-03-25T14:35:00Z', false, false
   );
   ```

4. **Send Notifications**

   **Email Notification:**
   ```
   Subject: [Watchup Alert] High CPU Usage on production-api-1
   
   Alert Details:
   - Server: production-api-1
   - Metric: CPU Usage
   - Current: 87.2%
   - Duration: 5 minutes
   - Severity: Warning
   
   Top Processes:
   1. node (PID 4213) - CPU: 52.1%, Memory: 8.3%
   2. python (PID 9021) - CPU: 18.4%, Memory: 12.1%
   
   View in Dashboard: https://watchup.space/dashboard/alerts/alert_abc123
   ```

   **Webhook Notification:**
   ```json
   POST {webhook_url}
   {
     "event": "alert.triggered",
     "alert_id": "alert_abc123",
     "server": "production-api-1",
     "metric": "cpu",
     "usage": 87.2,
     "severity": "warning",
     "timestamp": "2026-03-25T14:35:00Z"
   }
   ```

5. **Update Dashboard**
   - Trigger WebSocket event to connected clients
   - Update alert count badge
   - Show notification toast

6. **Check Alert Rules**
   - Prevent duplicate alerts (same metric within 10 minutes)
   - Respect notification preferences (email, SMS, webhook)
   - Apply quiet hours if configured

---

## Configuration Retrieval

### GET /agent/config

Retrieves current configuration for the agent. Called every 60 seconds.

#### Request

**Headers:**
```http
Authorization: Bearer srv_89sd0a7f3b2c1d4e
```

**Query Parameters:**
```
?server_key=srv_89sd0a7f3b2c1d4e
```

#### Response

**Success (200 OK):**
```json
{
  "success": true,
  "config": {
    "server_key": "srv_89sd0a7f3b2c1d4e",
    "project_id": "proj_12345",
    "server_identifier": "production-api-1",
    "sampling_interval": 5,
    "api_endpoint": "https://watchup.space",
    "alerts": {
      "cpu": {
        "threshold": 80,
        "duration": 300,
        "enabled": true
      },
      "ram": {
        "threshold": 75,
        "duration": 600,
        "enabled": true
      },
      "process_cpu": {
        "threshold": 60,
        "duration": 120,
        "enabled": true
      }
    },
    "features": {
      "disk_monitoring": false,
      "network_monitoring": false,
      "custom_metrics": false
    }
  },
  "updated_at": "2026-03-25T14:00:00Z"
}
```

**Error - Invalid Server Key (401 Unauthorized):**
```json
{
  "success": false,
  "error": "invalid_server_key",
  "message": "Server key is invalid or has been revoked"
}
```

#### Backend Implementation Notes

1. **Authenticate Server Key**
   ```sql
   SELECT a.*, p.plan, p.features
   FROM agents a
   JOIN projects p ON a.project_id = p.id
   WHERE a.server_key = 'srv_89sd0a7f3b2c1d4e' 
   AND a.status = 'active';
   ```

2. **Retrieve Configuration**
   ```sql
   SELECT threshold_cpu, threshold_ram, threshold_process,
          duration_cpu, duration_ram, duration_process,
          sampling_interval, features_enabled
   FROM agent_configs
   WHERE agent_id = 'agent_id';
   ```

3. **Apply Plan Limits**
   ```javascript
   function applyPlanLimits(config, plan) {
     if (plan === 'free') {
       config.features.disk_monitoring = false;
       config.features.network_monitoring = false;
       config.features.custom_metrics = false;
     }
     return config;
   }
   ```

4. **Cache Configuration**
   - Cache in Redis with 60-second TTL
   - Invalidate cache when user updates settings in dashboard

---

## Error Responses

### Standard Error Format

All error responses follow this format:

```json
{
  "success": false,
  "error": "error_code",
  "message": "Human-readable error message",
  "details": {
    "field": "Additional context"
  },
  "timestamp": "2026-03-25T14:35:00Z",
  "request_id": "req_xyz789"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `invalid_credentials` | 401 | Invalid master API key or server key |
| `invalid_server_key` | 401 | Server key is invalid or revoked |
| `agent_limit_reached` | 409 | Free project already has an agent |
| `project_not_found` | 404 | Project ID does not exist |
| `validation_error` | 400 | Request data failed validation |
| `rate_limit_exceeded` | 429 | Too many requests |
| `internal_error` | 500 | Server error |
| `service_unavailable` | 503 | Service temporarily unavailable |

---

## Rate Limiting

### Limits by Endpoint

| Endpoint | Rate Limit | Window |
|----------|------------|--------|
| `/agent/register` | 5 requests | per hour per IP |
| `/server/metrics` | 720 requests | per hour per agent (12/min) |
| `/server/alerts` | 60 requests | per hour per agent |
| `/agent/config` | 120 requests | per hour per agent (2/min) |

### Rate Limit Headers

```http
X-RateLimit-Limit: 720
X-RateLimit-Remaining: 650
X-RateLimit-Reset: 1711382400
```

### Rate Limit Exceeded Response

```json
{
  "success": false,
  "error": "rate_limit_exceeded",
  "message": "Rate limit exceeded. Please try again later.",
  "retry_after": 60,
  "limit": 720,
  "window": "1 hour"
}
```

---

## Database Schema

### Tables

#### agents
```sql
CREATE TABLE agents (
  id VARCHAR(36) PRIMARY KEY,
  project_id VARCHAR(50) NOT NULL,
  server_key VARCHAR(100) UNIQUE NOT NULL,
  server_identifier VARCHAR(50) NOT NULL,
  status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
  registered_at TIMESTAMP NOT NULL,
  last_seen TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_project_id (project_id),
  INDEX idx_server_key (server_key),
  INDEX idx_last_seen (last_seen)
);
```

#### metrics
```sql
CREATE TABLE metrics (
  id VARCHAR(36) PRIMARY KEY,
  agent_id VARCHAR(36) NOT NULL,
  project_id VARCHAR(50) NOT NULL,
  cpu_usage DECIMAL(5,2) NOT NULL,
  ram_usage DECIMAL(5,2) NOT NULL,
  recorded_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_agent_id (agent_id),
  INDEX idx_project_id (project_id),
  INDEX idx_recorded_at (recorded_at),
  FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
);
```

#### alerts
```sql
CREATE TABLE alerts (
  id VARCHAR(36) PRIMARY KEY,
  agent_id VARCHAR(36) NOT NULL,
  project_id VARCHAR(50) NOT NULL,
  metric_type ENUM('cpu', 'ram', 'process_cpu') NOT NULL,
  usage_value DECIMAL(5,2) NOT NULL,
  duration_seconds INT NOT NULL,
  severity ENUM('info', 'warning', 'critical') NOT NULL,
  top_processes JSON NOT NULL,
  triggered_at TIMESTAMP NOT NULL,
  acknowledged BOOLEAN DEFAULT FALSE,
  acknowledged_at TIMESTAMP NULL,
  acknowledged_by VARCHAR(36) NULL,
  resolved BOOLEAN DEFAULT FALSE,
  resolved_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_agent_id (agent_id),
  INDEX idx_project_id (project_id),
  INDEX idx_triggered_at (triggered_at),
  INDEX idx_severity (severity),
  INDEX idx_resolved (resolved),
  FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
);
```

#### agent_configs
```sql
CREATE TABLE agent_configs (
  id VARCHAR(36) PRIMARY KEY,
  agent_id VARCHAR(36) UNIQUE NOT NULL,
  threshold_cpu INT DEFAULT 80,
  threshold_ram INT DEFAULT 75,
  threshold_process INT DEFAULT 60,
  duration_cpu INT DEFAULT 300,
  duration_ram INT DEFAULT 600,
  duration_process INT DEFAULT 120,
  sampling_interval INT DEFAULT 5,
  features_enabled JSON DEFAULT '{}',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE CASCADE
);
```

---

## Testing

### Test Agent Registration

```bash
curl -X POST https://watchup.space/agent/register \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "proj_test123",
    "master_api_key": "ma_test_key_123",
    "server_identifier": "test-server-1"
  }'
```

### Test Metrics Submission

```bash
curl -X POST https://watchup.space/server/metrics \
  -H "Authorization: Bearer srv_test_key" \
  -H "Content-Type: application/json" \
  -d '{
    "server_key": "srv_test_key",
    "cpu": 45.2,
    "ram": 67.8,
    "timestamp": "2026-03-25T14:30:00Z"
  }'
```

### Test Alert Submission

```bash
curl -X POST https://watchup.space/server/alerts \
  -H "Authorization: Bearer srv_test_key" \
  -H "Content-Type: application/json" \
  -d '{
    "server_key": "srv_test_key",
    "metric": "cpu",
    "usage": 87.2,
    "duration": 300,
    "top_processes": [
      {"pid": 4213, "name": "node", "cpu": 52.1, "memory": 8.3}
    ],
    "timestamp": "2026-03-25T14:35:00Z"
  }'
```

### Test Configuration Retrieval

```bash
curl -X GET "https://watchup.space/agent/config?server_key=srv_test_key" \
  -H "Authorization: Bearer srv_test_key"
```

---

## Implementation Checklist

### Backend Development

- [ ] Set up database tables
- [ ] Implement `/agent/register` endpoint
- [ ] Implement `/server/metrics` endpoint
- [ ] Implement `/server/alerts` endpoint
- [ ] Implement `/agent/config` endpoint
- [ ] Add authentication middleware
- [ ] Add rate limiting
- [ ] Add input validation
- [ ] Add error handling
- [ ] Set up notification system (email, webhook)
- [ ] Set up WebSocket for real-time updates
- [ ] Add logging and monitoring
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Deploy to production

### Dashboard Development

- [ ] Create agent management page
- [ ] Create metrics visualization
- [ ] Create alerts page
- [ ] Create configuration editor
- [ ] Add real-time updates via WebSocket
- [ ] Add notification preferences
- [ ] Add alert acknowledgment
- [ ] Add historical data charts

---

## Support

For questions or issues with the API:

- **Documentation**: This file
- **GitHub**: https://github.com/tomurashigaraki22/watchup-agent
- **Email**: support@watchup.space
