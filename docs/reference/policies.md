# Access Control Lists (ACLs)

ACLs are the **most critical security component** of your MESH deployment. They control which nodes can communicate with each other and on which ports.

!!! danger "Critical security warning"
    **Never use permissive ACLs (`"*"` rules) in production.** This allows any compromised endpoint to access all other devices in your mesh, including other forensic targets and analyst workstations. Always implement the principle of least privilege.

## Where the ACLs are stored

MESH uses `policy.mode: database` in the Headscale config. Therefore the active policy is stored in Headscale's SQLite database, not in a YAML file you edit on disk.

To read or change the policy, open the control plane Web UI and click CUSTOMIZE ACL in the sidebar. Save Changes writes straight to the database and applies the updated ACLs immediately.

## Why ACLs are critical for MESH

In a forensic mesh network, security isolation is paramount:

**Threat model:**

- **Compromised endpoints**: Target devices may be infected with malware
- **Lateral movement**: Malware could attempt to spread to other devices
- **Data exfiltration**: Compromised devices could try to access other forensic targets
- **Analyst protection**: Analyst workstations must be protected from compromised endpoints

**Security goals:**

- Endpoints should **only** be accessible by authorised analysts using a secure collection node
- Endpoints should **never** be able to access each other
- Endpoints should **never** be able to access analyst workstations
- Analysts should have controlled access to specific endpoints
- All access should be logged and auditable

## ACL syntax and structure

ACLs are defined in JSON format (technically HuJSON which can include comments) with the following structure:

```json
{
  // Define groups of users/nodes
  "groups": {
    "group:name": [
      "user1",
      "user2"
    ]
  },

  // Define access control rules
  "acls": [
    {
      "action": "accept",  // or "deny"
      "src": ["source_group_or_user"],
      "dst": ["destination_group_or_user:port"]
    }
  ]
}
```

**Components:**

**`groups`**: Logical groupings of users or nodes

- Format: `group:name`
- Contains list of usernames or node names
- Makes ACL rules more maintainable
- Can be nested or referenced in rules

**`acls`**: List of access control rules

- Evaluated in order (first match wins)
- Each rule has `action`, `src`, and `dst`

**`action`**: What to do with matching traffic

- `accept`: Allow the traffic
- `deny`: Block the traffic

**`src`**: Source of the traffic

- Can be: username, node name, group, or `"*"` (any)
- Multiple sources can be listed

**`dst`**: Destination of the traffic

- Format: `target:port` or `target:*` (all ports)
- Can be: username, node name, group, or `"*"` (any)
- Port can be specific (e.g., `5555` for ADB) or `*` (all ports)

## Testing vs. production ACLs

### Testing ACL (permissive - DO NOT USE IN PRODUCTION)

```json
{
  // WARNING: This is for testing only.
  // Allows all traffic between all nodes.
  "acls": [
    {
      "action": "accept",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

!!! warning "Testing only"
    This permissive ACL is **only** for initial testing and development. It provides **no security** and should **never** be used in production or with real forensic targets.

**When to use:**

- Initial setup and testing
- Verifying connectivity between nodes
- Troubleshooting connection issues

**When to remove:**

- Before connecting any real forensic targets
- Before deploying in production
- Before handling sensitive data

### Production ACL (restrictive - RECOMMENDED)

```json
{
  // Production ACL - Principle of least privilege
  "groups": {
    "group:analysts": ["analyst1", "analyst2"],
    "group:endpoints": ["android-device-1", "android-device-2"]
  },

  "acls": [
    // Analysts can access all endpoints on ADB port
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:5555"]  // ADB port only
    },

    // Analysts can access all endpoints on all ports (for full forensics)
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    },

    // Endpoints CANNOT access each other (critical for isolation)
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["group:endpoints:*"]
    },

    // Endpoints CANNOT access analysts (prevent reverse connections)
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["group:analysts:*"]
    },

    // Analysts can access each other (for collaboration)
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:analysts:*"]
    },

    // Deny all other traffic (default deny)
    {
      "action": "deny",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

## Real-world ACL examples

### Example 1: Single analyst, multiple endpoints

```json
{
  // Simple forensic setup
  "groups": {
    "group:analysts": ["forensic-workstation"],
    "group:endpoints": [
      "suspect-phone-1",
      "suspect-phone-2",
      "suspect-phone-3"
    ]
  },

  "acls": [
    // Analyst can access all endpoints
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    },

    // Endpoints are completely isolated
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["*:*"]
    }
  ]
}
```

### Example 2: Multi-team deployment

```json
{
  // Multiple analyst teams with separate endpoints
  "groups": {
    "group:team-alpha-analysts": ["analyst-alice", "analyst-bob"],
    "group:team-alpha-endpoints": ["case-123-phone-1", "case-123-phone-2"],
    "group:team-bravo-analysts": ["analyst-charlie", "analyst-diana"],
    "group:team-bravo-endpoints": ["case-456-phone-1", "case-456-phone-2"]
  },

  "acls": [
    // Team Alpha: analysts can access their endpoints
    {
      "action": "accept",
      "src": ["group:team-alpha-analysts"],
      "dst": ["group:team-alpha-endpoints:*"]
    },

    // Team Bravo: analysts can access their endpoints
    {
      "action": "accept",
      "src": ["group:team-bravo-analysts"],
      "dst": ["group:team-bravo-endpoints:*"]
    },

    // Team Alpha: analysts can collaborate
    {
      "action": "accept",
      "src": ["group:team-alpha-analysts"],
      "dst": ["group:team-alpha-analysts:*"]
    },

    // Team Bravo: analysts can collaborate
    {
      "action": "accept",
      "src": ["group:team-bravo-analysts"],
      "dst": ["group:team-bravo-analysts:*"]
    },

    // Teams CANNOT access each other's endpoints
    {
      "action": "deny",
      "src": ["group:team-alpha-analysts"],
      "dst": ["group:team-bravo-endpoints:*"]
    },
    {
      "action": "deny",
      "src": ["group:team-bravo-analysts"],
      "dst": ["group:team-alpha-endpoints:*"]
    },

    // Endpoints are completely isolated
    {
      "action": "deny",
      "src": ["group:team-alpha-endpoints"],
      "dst": ["*:*"]
    },
    {
      "action": "deny",
      "src": ["group:team-bravo-endpoints"],
      "dst": ["*:*"]
    },

    // Default deny
    {
      "action": "deny",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

### Example 3: Port-specific access (ADB only)

```json
{
  // Restrict analysts to only ADB access
  "groups": {
    "group:analysts": ["forensic-workstation"],
    "group:endpoints": ["android-device-1", "android-device-2"]
  },

  "acls": [
    // Analysts can ONLY access ADB port (5555)
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:5555"]
    },

    // Deny all other ports
    {
      "action": "deny",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    },

    // Endpoints completely isolated
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["*:*"]
    }
  ]
}
```

## ACL best practices for forensic deployments

### 1. Principle of least privilege

**Always grant the minimum access required:**

- Analysts should only access endpoints they're investigating
- Limit port access to only what's needed (e.g., ADB port 5555)
- Use time-limited access when possible (rotate keys regularly)

### 2. Endpoint isolation

**Endpoints must never access each other:**

```json
// CRITICAL: Always include this rule
{
  "action": "deny",
  "src": ["group:endpoints"],
  "dst": ["group:endpoints:*"]
}
```

**Why this matters:**

- Prevents malware from spreading between forensic targets
- Stops compromised devices from accessing other evidence
- Maintains chain of custody integrity

### 3. Prevent reverse connections

**Endpoints should not initiate connections to analysts:**

```json
// Prevent endpoints from connecting to analysts
{
  "action": "deny",
  "src": ["group:endpoints"],
  "dst": ["group:analysts:*"]
}
```

**Why this matters:**

- Prevents malware from attacking analyst workstations
- Stops data exfiltration from endpoints
- Maintains unidirectional trust model

### 4. Use groups for maintainability

**Bad (hard to maintain):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["analyst1", "analyst2", "analyst3"],
      "dst": ["phone1:*", "phone2:*", "phone3:*"]
    }
  ]
}
```

**Good (easy to maintain):**

```json
{
  "groups": {
    "group:analysts": ["analyst1", "analyst2", "analyst3"],
    "group:endpoints": ["phone1", "phone2", "phone3"]
  },

  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    }
  ]
}
```

### 5. Default deny at the end

**Always end with a default deny rule:**

```json
{
  "acls": [
    // ... your specific rules ...

    // Default deny (catches anything not explicitly allowed)
    {
      "action": "deny",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

## Testing and validating ACLs

After creating or modifying ACLs in the CUSTOMIZE ACL page, always test them:

### 1. Test from analyst workstation

```bash
# Should succeed: Analyst accessing endpoint
ping 100.64.2.1  # Endpoint mesh IP

# Should succeed: ADB connection
meshcli adbpair 100.64.2.1:5555
```

### 2. Test endpoint isolation

```bash
# On endpoint device (via ADB shell, after pairing with `meshcli adbpair`)
adb shell

# Should FAIL: Endpoint trying to access another endpoint
ping 100.64.2.2  # Another endpoint's mesh IP

# Should FAIL: Endpoint trying to access analyst
ping 100.64.1.1  # Analyst mesh IP
```

### 3. Check logs on control plane

```bash
# View ACL enforcement logs
docker compose logs headscale | grep -i acl

# Look for denied connections
docker compose logs headscale | grep -i denied
```

## Common ACL mistakes and how to avoid them

### Mistake 1: Using `"*"` in production

**Wrong:**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

**Right:**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    }
  ]
}
```

### Mistake 2: Forgetting endpoint isolation

**Wrong (endpoints can access each other):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    }
  ]
}
```

**Right (endpoints explicitly denied):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    },
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["group:endpoints:*"]
    }
  ]
}
```

### Mistake 3: Allowing reverse connections

**Wrong (endpoints can connect to analysts):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["*"],
      "dst": ["*:*"]
    }
  ]
}
```

**Right (only analysts initiate connections):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    },
    {
      "action": "deny",
      "src": ["group:endpoints"],
      "dst": ["group:analysts:*"]
    }
  ]
}
```

### Mistake 4: Overly broad port access

**Wrong (all ports open):**

```json
{
  "acls": [
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:*"]
    }
  ]
}
```

**Better (specific ports only):**

```json
{
  "acls": [
    // ADB access only
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:5555"]
    },

    // SSH access only (if needed)
    {
      "action": "accept",
      "src": ["group:analysts"],
      "dst": ["group:endpoints:22"]
    }
  ]
}
```
