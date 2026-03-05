# Security best practices

1. **Use HTTPS**: Always use a reverse proxy with SSL/TLS
2. **Restrict API access**: Don't expose port 8080 publicly
3. **Use strong ACLs**: Implement least-privilege access control
4. **Rotate API keys**: Regularly expire and recreate API keys
5. **Monitor logs**: Watch for suspicious activity
6. **Backup regularly**: Automate database backups
7. **Update frequently**: Keep Headscale and Docker updated
8. **Use firewall**: Only allow necessary ports

## Security considerations

### Pre-authentication key rotation

Pre-auth keys should be rotated regularly:

```bash
# List existing keys
docker compose exec headscale headscale preauthkeys list

# Expire old keys
docker compose exec headscale headscale preauthkeys expire --user default KEY_ID

# Create new key
docker compose exec headscale headscale preauthkeys create --user default --reusable --expiration 24h
```

**Best practices:**

- Use short expiration times (24-48 hours)
- Create new keys for each investigation
- Expire keys immediately after use if not reusable
- Never share keys between teams or cases

### Node removal

Remove nodes when investigations complete using the UI or via the backend via:

```bash
# List nodes
docker compose exec headscale headscale nodes list

# Remove a node
docker compose exec headscale headscale nodes delete --identifier NODE_ID
```

**When to remove nodes:**

- Investigation completed
- Device returned to owner
- Node compromised or suspected compromise
- Regular cleanup of old investigations

### Audit logging

Monitor ACL enforcement:

```yaml
# Enable detailed logging in config.yaml
log:
  level: debug  # or trace for maximum detail
  format: json  # easier to parse
```

```bash
# Monitor logs in real-time
docker compose logs -f headscale | grep -i acl
```

**What to monitor:**

- Denied connection attempts (potential compromise)
- Unusual access patterns
- Failed authentication attempts
- ACL rule changes
