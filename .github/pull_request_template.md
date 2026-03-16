## Summary

<!-- What changed, why, and anything reviewers should focus on. Link issue/ticket if relevant. Do not include secrets or sensitive exploit details. -->

## Breaking changes

- [ ] No breaking changes
- [ ] Breaking change, migration, or operator-visible behavior change is described below

## Testing

- [ ] Relevant automated tests were added/updated or not needed
- [ ] Relevant local checks were run
- [ ] CI checks passed, or failures are explained below

**Commands / checks run:**

## Manual testing

> Integration tests and end-to-end functional tests are not fully implemented yet. Complete the applicable items below.

- [ ] Manual testing not needed (docs/comment-only change)
- [ ] Core MESH flow verified where relevant (control plane, enrollment, connectivity, forensics, ADB pairing)
- [ ] Android client flow verified where relevant (install/upgrade, connect/disconnect, permissions/background behavior)
- [ ] CLI / analyst tooling verified where relevant (affected `meshcli` commands, flags, output, error paths)
- [ ] Cross-platform / backward-compatibility checks completed where relevant

## Security considerations

- [ ] Auth, ACL, network exposure, and transport-security implications were reviewed
- [ ] No secrets, keys, tokens, private certs, or sensitive data were added to the repo, logs, screenshots, or PR text
- [ ] New or updated dependencies were reviewed for necessity and known risk/vulnerabilities
- [ ] Input validation, error handling, and logs do not leak sensitive information
- [ ] Civil society / high-risk environment and F-Droid impacts were reviewed where applicable
- [ ] If security-sensitive, details are being handled per `SECURITY.md` rather than in public discussion

**Security notes:**

## Deployment

- [ ] CI/CD impact was verified
- [ ] Deployment, config, env var, secret, packaging, or release impact is documented

**Notes:**

## Documentation

- [ ] No documentation updates needed
- [ ] README / docs updated
- [ ] User, analyst, or operator-facing changes documented

## Reviewer guidance

<!-- What should reviewers focus on first? List key files, risks, or areas needing careful review. -->