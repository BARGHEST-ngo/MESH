interface EncryptedPayload {
	uri: string
	pin: string
}

interface OnboardingPayload {
	key: string
	url: string
}

export async function encrypt(controlPlaneURL: string, authKey: string): Promise<EncryptedPayload> {
	const digits = 6
	const pinArray = crypto.getRandomValues(new Uint32Array(1))
	const pin = (pinArray[0] % 10 ** digits).toString().padStart(digits, '0')
	const saltArray = crypto.getRandomValues(new Uint8Array(16))
	const ivArray = crypto.getRandomValues(new Uint8Array(12))
	const payload: OnboardingPayload = { key: authKey, url: controlPlaneURL }
	const intent = JSON.stringify(payload)
	const pinBytes = new TextEncoder().encode(pin)
	const pbkdf2 = await crypto.subtle.importKey("raw", pinBytes, "PBKDF2", false, ["deriveBits"])
	const aesKeyBytes = await crypto.subtle.deriveBits({name: "PBKDF2", salt: saltArray, iterations: 600000, hash: "SHA-256" }, pbkdf2, 256)
	const aesKey = await crypto.subtle.importKey("raw", new Uint8Array(aesKeyBytes), "AES-GCM", false, ["encrypt"])
	const plaintext = new TextEncoder().encode(intent)
	const cipher = await crypto.subtle.encrypt({name: "AES-GCM", iv: ivArray}, aesKey, plaintext)
	const blob = new Uint8Array(16 + 12 + cipher.byteLength)
	blob.set(saltArray, 0)
	blob.set(ivArray, 16)
	blob.set(new Uint8Array(cipher), 28)
	return { uri: buildURI(blob), pin}
}

function buildURI(blob: Uint8Array): string {
	return `mesh://onboard?d=${btoa(String.fromCharCode(...blob))
		.replace(/\+/g, "-")
		.replace(/\//g, "_")
		.replace(/=+$/, "")}`
}
