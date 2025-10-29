// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package gp

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go mksyscall.go

//sys enterCriticalPolicySection(machine bool) (handle policyLockHandle, err error) [int32(failretval)==0] = userenv.EnterCriticalPolicySection
//sys impersonateLoggedOnUser(token windows.Token) (err error) [int32(failretval)==0] = advapi32.ImpersonateLoggedOnUser
//sys leaveCriticalPolicySection(handle policyLockHandle) (err error) [int32(failretval)==0] = userenv.LeaveCriticalPolicySection
//sys registerGPNotification(event windows.Handle, machine bool) (err error) [int32(failretval)==0] = userenv.RegisterGPNotification
//sys refreshPolicyEx(machine bool, flags uint32) (err error) [int32(failretval)==0] = userenv.RefreshPolicyEx
//sys unregisterGPNotification(event windows.Handle) (err error) [int32(failretval)==0] = userenv.UnregisterGPNotification
