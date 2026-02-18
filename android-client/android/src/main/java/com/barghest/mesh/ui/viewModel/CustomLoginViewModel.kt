// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package com.barghest.mesh.ui.viewModel

import com.barghest.mesh.ui.util.set
import com.barghest.mesh.ui.view.ErrorDialogType
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow

const val AUTH_KEY_LENGTH = 16

open class CustomLoginViewModel : IpnViewModel() {
  val errorDialog: StateFlow<ErrorDialogType?> = MutableStateFlow(null)
}

class LoginWithAuthKeyViewModel : CustomLoginViewModel() {
  // Sets the auth key and invokes the login flow
  fun setAuthKey(authKey: String, onSuccess: () -> Unit) {
    // The most basic of checks for auth key syntax
    if (authKey.isEmpty()) {
      errorDialog.set(ErrorDialogType.INVALID_AUTH_KEY)
      return
    }
    loginWithAuthKey(authKey) {
      it.onFailure { errorDialog.set(ErrorDialogType.ADD_PROFILE_FAILED) }
      it.onSuccess { onSuccess() }
    }
  }
}

class LoginWithCustomControlURLViewModel : CustomLoginViewModel() {
  private val _successMessage = MutableStateFlow<String?>(null)
  val successMessage: StateFlow<String?> = _successMessage

  fun setControlURL(urlStr: String, onSuccess: () -> Unit) {
    when (urlStr.startsWith("http", ignoreCase = true) &&
            urlStr.contains("://") &&
            urlStr.length > 7) {
      false -> {
        errorDialog.set(ErrorDialogType.INVALID_CUSTOM_URL)
        return
      }
      true -> {
        loginWithCustomControlURL(urlStr) {
          it.onFailure { errorDialog.set(ErrorDialogType.ADD_PROFILE_FAILED) }
          it.onSuccess {
            _successMessage.value = "Control plane successfully added"
          }
        }
      }
    }
  }

  fun clearSuccessMessage() {
    _successMessage.value = null
  }

  fun proceedToAuthKey(onNavigateToAuthKey: () -> Unit) {
    clearSuccessMessage()
    onNavigateToAuthKey()
  }
}
