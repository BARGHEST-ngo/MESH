package com.barghest.mesh

import javax.crypto.Cipher
import javax.crypto.SecretKeyFactory
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.PBEKeySpec
import javax.crypto.spec.SecretKeySpec


object IntentCrypto {

    private fun derive(pin: String, salt: ByteArray): SecretKeySpec {
        if (pin.length == 6 && pin.all { it.isDigit() }) {
            val factory = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256")
            // Convert to UTF-8 bytes first, then to chars, to match Web Crypto's TextEncoder behaviour
            val pinUtf8 = pin.toByteArray(Charsets.UTF_8)
            val pinChars = CharArray(pinUtf8.size) { (pinUtf8[it].toInt() and 0xFF).toChar() }
            val spec = PBEKeySpec(pinChars, salt, 600000, 256)
            val keyBytes = factory.generateSecret(spec).encoded
            return SecretKeySpec(keyBytes, "AES")
        } else {
            throw IllegalArgumentException("PIN must be 6 digits")
        }
    }

    fun decrypt(pin: String, salt: ByteArray, iv: ByteArray, cipherText: ByteArray): String {
        val aesKey = derive(pin, salt)
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        val spec = GCMParameterSpec(128, iv)
        cipher.init(Cipher.DECRYPT_MODE, aesKey, spec)
        val plaintext = cipher.doFinal(cipherText)
        return String(plaintext, Charsets.UTF_8)
    }
}
