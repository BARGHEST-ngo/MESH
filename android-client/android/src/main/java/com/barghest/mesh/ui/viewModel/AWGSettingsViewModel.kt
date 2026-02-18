// SPDX-License-Identifier: CC0-1.0
// Authored by BARGHEST. Dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import android.content.Context
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import com.barghest.mesh.App
import com.barghest.mesh.util.TSLog
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext
import java.io.File

data class AWGConfig(
    val enabled: Boolean = false,
    val jc: Int = 0,
    val jmin: Int = 0,
    val jmax: Int = 0,
    val s1: Int = 0,
    val s2: Int = 0,
    val h1: Int = 1,
    val h2: Int = 2,
    val h3: Int = 3,
    val h4: Int = 4
) {
    fun isObfuscationEnabled(): Boolean {
        return jc > 0 || jmin > 0 || jmax > 0 || s1 > 0 || s2 > 0 ||
                h1 != 1 || h2 != 2 || h3 != 3 || h4 != 4
    }

    fun validate(): String? {
        if (jc < 0 || jc > 128) return "Jc must be between 0 and 128"
        if (jmax > 1280) return "Jmax must be ≤ 1280"
        if (jmin > jmax && jmax != 0) return "Jmin must be ≤ Jmax"
        if (s1 in 1..14) return "S1 must be 0 or ≥ 15"
        if (s1 > 150) return "S1 must be ≤ 150"
        if (s2 in 1..14) return "S2 must be 0 or ≥ 15"
        if (s2 > 150) return "S2 must be ≤ 150"
        if (s1 > 0 && s2 > 0 && s1 + 56 == s2) return "S1+56 must not equal S2"
        if (h1 < 1 || h2 < 1 || h3 < 1 || h4 < 1) return "H1-H4 must be ≥ 1"
        val hSet = setOf(h1, h2, h3, h4)
        if (hSet.size != 4) return "H1, H2, H3, H4 must all be different"
        return null
    }
}

class AWGSettingsViewModelFactory : ViewModelProvider.Factory {
    @Suppress("UNCHECKED_CAST")
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return AWGSettingsViewModel() as T
    }
}

class AWGSettingsViewModel : ViewModel() {
    companion object {
        private const val TAG = "AWGSettingsViewModel"
        private const val CONFIG_FILENAME = "amneziawg.conf"
    }

    private val _config = MutableStateFlow(AWGConfig())
    val config: StateFlow<AWGConfig> = _config

    private val _saveStatus = MutableStateFlow<String?>(null)
    val saveStatus: StateFlow<String?> = _saveStatus

    private val _validationError = MutableStateFlow<String?>(null)
    val validationError: StateFlow<String?> = _validationError

    init {
        loadConfig()
    }

    fun updateConfig(newConfig: AWGConfig) {
        _config.value = newConfig
        _validationError.value = newConfig.validate()
    }

    fun loadConfig() {
        viewModelScope.launch {
            withContext(Dispatchers.IO) {
                try {
                    val context = App.get()
                    val file = File(context.filesDir, CONFIG_FILENAME)
                    if (file.exists()) {
                        val config = parseConfigFile(file.readText())
                        _config.value = config
                        TSLog.d(TAG, "Loaded AWG config: $config")
                    } else {
                        TSLog.d(TAG, "No AWG config file found, using defaults")
                    }
                } catch (e: Exception) {
                    TSLog.e(TAG, "Failed to load AWG config: ${e.message}")
                }
            }
        }
    }

    fun saveConfig() {
        val currentConfig = _config.value
        val error = currentConfig.validate()
        if (error != null) {
            _validationError.value = error
            return
        }

        viewModelScope.launch {
            withContext(Dispatchers.IO) {
                try {
                    val context = App.get()
                    val file = File(context.filesDir, CONFIG_FILENAME)
                    file.writeText(generateConfigFile(currentConfig))
                    TSLog.d(TAG, "Saved AWG config to ${file.absolutePath}")
                    _saveStatus.value = "saved"
                } catch (e: Exception) {
                    TSLog.e(TAG, "Failed to save AWG config: ${e.message}")
                    _saveStatus.value = "error"
                }
            }
        }
    }

    fun resetToDefaults() {
        _config.value = AWGConfig()
        _validationError.value = null
    }

    fun clearSaveStatus() {
        _saveStatus.value = null
    }

    private fun parseConfigFile(content: String): AWGConfig {
        var cfg = AWGConfig()
        var inInterface = false

        content.lines().forEach { line ->
            val trimmed = line.trim()
            if (trimmed.isEmpty() || trimmed.startsWith("#") || trimmed.startsWith(";")) {
                return@forEach
            }
            if (trimmed.startsWith("[") && trimmed.endsWith("]")) {
                inInterface = trimmed.lowercase() == "[interface]"
                return@forEach
            }
            if (!inInterface) return@forEach

            val parts = trimmed.split("=", limit = 2)
            if (parts.size != 2) return@forEach

            val key = parts[0].trim().lowercase()
            val value = parts[1].trim().toIntOrNull() ?: return@forEach

            cfg = when (key) {
                "jc" -> cfg.copy(jc = value)
                "jmin" -> cfg.copy(jmin = value)
                "jmax" -> cfg.copy(jmax = value)
                "s1" -> cfg.copy(s1 = value)
                "s2" -> cfg.copy(s2 = value)
                "h1" -> cfg.copy(h1 = value)
                "h2" -> cfg.copy(h2 = value)
                "h3" -> cfg.copy(h3 = value)
                "h4" -> cfg.copy(h4 = value)
                else -> cfg
            }
        }
        return cfg.copy(enabled = cfg.isObfuscationEnabled())
    }

    private fun generateConfigFile(config: AWGConfig): String {
        return buildString {
            appendLine("# AmneziaWG Configuration")
            appendLine("# Generated by MESH Android client")
            appendLine()
            appendLine("[Interface]")
            appendLine("Jc = ${config.jc}")
            appendLine("Jmin = ${config.jmin}")
            appendLine("Jmax = ${config.jmax}")
            appendLine("S1 = ${config.s1}")
            appendLine("S2 = ${config.s2}")
            appendLine("H1 = ${config.h1}")
            appendLine("H2 = ${config.h2}")
            appendLine("H3 = ${config.h3}")
            appendLine("H4 = ${config.h4}")
        }
    }
}

