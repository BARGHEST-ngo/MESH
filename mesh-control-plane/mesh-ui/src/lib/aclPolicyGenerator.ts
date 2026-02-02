import type { ACLPolicy, ACLRule } from '../api/usePolicy'
import type { V1User } from '../api/openapi/types.gen'

export function generateNetworkRule(networkName: string): ACLRule {
    return {
        action: 'accept',
        src: [`${networkName}@`],
        dst: [`${networkName}@:*`],
    }
}

/**
 * Generate a complete ACL policy for network isolation
 * 
 * @param networks - Array of network (user) objects from Headscale
 * @returns ACLPolicy object with isolation rules for all networks
 * 
 * @example
 * // For networks: network-1, network-2, network-3
 * // Generates:
 * // {
 * //   "acls": [
 * //     { "action": "accept", "src": ["network-1@"], "dst": ["network-1@:*"] },
 * //     { "action": "accept", "src": ["network-2@"], "dst": ["network-2@:*"] },
 * //     { "action": "accept", "src": ["network-3@"], "dst": ["network-3@:*"] }
 * //   ]
 * // }
 */
export function generateNetworkIsolationPolicy(networks: V1User[]): ACLPolicy {
    // Filter out networks without names and generate rules
    const acls: ACLRule[] = networks
        .filter((network): network is V1User & { name: string } => 
            typeof network.name === 'string' && network.name.length > 0
        )
        .map(network => generateNetworkRule(network.name))

    return {
        acls,
    }
}

export function generateNetworkIsolationPolicyFromNames(networkNames: string[]): ACLPolicy {
    const acls: ACLRule[] = networkNames
        .filter(name => name && name.length > 0)
        .map(name => generateNetworkRule(name))

    return {
        acls,
    }
}

export function addNetworkToPolicy(existingPolicy: ACLPolicy,newNetworkName: string): ACLPolicy {
    const existingAcls = existingPolicy.acls || []
    
    // Check if rule already exists for this network
    const ruleExists = existingAcls.some(rule => 
        rule.src.includes(`${newNetworkName}@`) &&
        rule.dst.includes(`${newNetworkName}@:*`)
    )

    if (ruleExists) {
        return existingPolicy
    }

    return {
        ...existingPolicy,
        acls: [...existingAcls, generateNetworkRule(newNetworkName)],
    }
}

export function removeNetworkFromPolicy(
    existingPolicy: ACLPolicy,
    networkName: string
): ACLPolicy {
    const existingAcls = existingPolicy.acls || []

    return {
        ...existingPolicy,
        acls: existingAcls.filter(rule => 
            !rule.src.includes(`${networkName}@`) ||
            !rule.dst.includes(`${networkName}@:*`)
        ),
    }
}

export function validateNetworkIsolation(policy: ACLPolicy, networks: V1User[]): { valid: boolean; issues: string[] } {
    const issues: string[] = []
    const acls = policy.acls || []

    for (const network of networks) {
        if (!network.name) continue

        const hasRule = acls.some(rule =>
            rule.src.includes(`${network.name}@`) &&
            rule.dst.includes(`${network.name}@:*`)
        )

        if (!hasRule) {
            issues.push(`Network "${network.name}" is missing isolation rule`)
        }
    }

    // Check for cross-network rules (potential isolation breach)
    for (const rule of acls) {
        for (const src of rule.src) {
            for (const dst of rule.dst) {
                // Extract network names from src and dst
                const srcNetwork = src.replace(/@$/, '')
                const dstNetwork = dst.replace(/@:\*$/, '').replace(/@:.*$/, '')
                
                if (srcNetwork && dstNetwork && srcNetwork !== dstNetwork) {
                    // This might be intentional, but flag it
                    issues.push(
                        `Cross-network rule detected: ${src} -> ${dst}`
                    )
                }
            }
        }
    }

    return {
        valid: issues.length === 0,
        issues,
    }
}