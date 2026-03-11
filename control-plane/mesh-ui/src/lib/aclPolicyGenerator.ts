import type { ACLPolicy, ACLRule } from '../api/usePolicy'
import type { V1User } from '../api/openapi/types.gen'

function sanitizeTagComponent(input: unknown): string {
    const raw = String(input ?? '')
    return raw
        .toLowerCase()
        .replace(/[^a-z0-9-]/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')
}

export function networkTagForNetwork(network: Pick<V1User, 'id' | 'name'>): string | null {
    const namePart = sanitizeTagComponent(network.name)
    if (namePart) return `tag:net-${namePart}`

    const idPart = sanitizeTagComponent(network.id)
    if (idPart) return `tag:net-${idPart}`

    return null
}

export function generateNetworkRule(network: Pick<V1User, 'id' | 'name'>): ACLRule {
    const tag = networkTagForNetwork(network)
    if (!tag) {
        return { action: 'accept', src: [], dst: [] }
    }

    return {
        action: 'accept',
        src: [tag],
        dst: [`${tag}:*`],
    }
}

export function generateExitNodeRule(network: Pick<V1User, 'id' | 'name'>): ACLRule | null {
    const tag = networkTagForNetwork(network)
    if (!tag) return null

    return {
        action: 'accept',
        src: [tag],
        dst: ['autogroup:internet:*'],
    }
}

/**
 * Build a full ACL policy that isolates each network from every other network.
 *
 * Each network is allowed to talk only to itself. Optionally, selected
 * networks can be given internet access if they have an approved exit node.
 *
 * @param networks
 *   List of network (user) objects returned by Headscale.
 *
 * @param networksWithExitNodes
 *   Optional set of network names that have at least one approved exit node.
 *   If provided, an `autogroup:internet` rule is added only for those networks.
 *   If omitted or undefined, no exit-node rules are generated.
 *
 * @returns
 *   An ACLPolicy object containing per-network isolation rules and tag owners.
 *
 * @example
 * // Given networks: network-1, network-2, network-3
 * // Produces:
 * // {
 * //   "acls": [
 * //     { "action": "accept", "src": ["tag:net-network-1"], "dst": ["tag:net-network-1:*"] },
 * //     { "action": "accept", "src": ["tag:net-network-2"], "dst": ["tag:net-network-2:*"] },
 * //     { "action": "accept", "src": ["tag:net-network-3"], "dst": ["tag:net-network-3:*"] }
 * //   ],
 * //   "tagOwners": {
 * //     "tag:net-network-1": ["network-1@"],
 * //     "tag:net-network-2": ["network-2@"],
 * //     "tag:net-network-3": ["network-3@"]
 * //   }
 * // }
 */

export function generateNetworkIsolationPolicy(
    networks: V1User[],
    networksWithExitNodes?: Set<string>,
): ACLPolicy {
    const acls: ACLRule[] = []
    const tagOwners: Record<string, string[]> = {}

    for (const network of networks) {
        if (!network) continue
        if (typeof network.name !== 'string' || network.name.length === 0) continue

        const tag = networkTagForNetwork(network)
        if (!tag) continue

        acls.push(generateNetworkRule(network))

        // Only add internet-access rule if this network has an approved exit node.
        if (networksWithExitNodes?.has(network.name)) {
            const exitRule = generateExitNodeRule(network)
            if (exitRule) acls.push(exitRule)
        }

        // Allow the network "user" to own its network tag.
        // Headscale requires user references to include the "@" suffix.
        tagOwners[tag] = [`${network.name}@`]
    }

    return Object.keys(tagOwners).length > 0 ? { acls, tagOwners } : { acls }
}

export function generateNetworkIsolationPolicyFromNames(networkNames: string[]): ACLPolicy {
    const acls: ACLRule[] = networkNames
        .filter(name => typeof name === 'string' && name.length > 0)
        .map(name => {
            const tag = `tag:net-${sanitizeTagComponent(name)}`
            return { action: 'accept', src: [tag], dst: [`${tag}:*`] }
        })

    return { acls }
}

export function addNetworkToPolicy(existingPolicy: ACLPolicy, newNetwork: Pick<V1User, 'id' | 'name'>): ACLPolicy {
    const existingAcls = existingPolicy.acls || []

    const tag = networkTagForNetwork(newNetwork)
    if (!tag) {
        return existingPolicy
    }
    
    // Check if rule already exists for this network
    const ruleExists = existingAcls.some(rule => 
        rule.src.includes(tag) &&
        rule.dst.includes(`${tag}:*`)
    )

    if (ruleExists) {
        return existingPolicy
    }

    return {
        ...existingPolicy,
        acls: [...existingAcls, generateNetworkRule(newNetwork)],
    }
}

export function removeNetworkFromPolicy(existingPolicy: ACLPolicy, network: Pick<V1User, 'id' | 'name'>): ACLPolicy {
    const existingAcls = existingPolicy.acls || []

    const tag = networkTagForNetwork(network)
    if (!tag) {
        return existingPolicy
    }

    return {
        ...existingPolicy,
        acls: existingAcls.filter(rule => 
            !rule.src.includes(tag) ||
            !rule.dst.includes(`${tag}:*`)
        ),
    }
}

export function validateNetworkIsolation(policy: ACLPolicy, networks: V1User[]): { valid: boolean; issues: string[] } {
    const issues: string[] = []
    const acls = policy.acls || []

    for (const network of networks) {
        if (!network?.name) continue

        const tag = networkTagForNetwork(network)
        if (!tag) continue

        const hasRule = acls.some(rule =>
            rule.src.includes(tag) &&
            rule.dst.includes(`${tag}:*`)
        )

        if (!hasRule) {
            issues.push(`Network "${network.name}" is missing isolation rule`)
        }
    }

    // Check for cross-network rules (potential isolation breach)
    for (const rule of acls) {
        for (const src of rule.src) {
            for (const dst of rule.dst) {
                if (!src.startsWith('tag:net-')) continue

                const dstTag = dst.replace(/:\*$/, '')
                if (!dstTag.startsWith('tag:net-')) continue

                if (src !== dstTag) {
                    issues.push(`Cross-network rule detected: ${src} -> ${dst}`)
                }
            }
        }
    }

    return {
        valid: issues.length === 0,
        issues,
    }
}