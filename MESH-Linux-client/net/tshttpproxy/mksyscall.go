// Copyright (c) 2020- 2025 Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause
// Additional contributions by BARGHEST are dedicated to the public domain under CC0 1.0.

package tshttpproxy

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output zsyscall_windows.go mksyscall.go

//sys globalFree(hglobal winHGlobal) (err error) [failretval==0] = kernel32.GlobalFree
//sys winHTTPCloseHandle(whi winHTTPInternet) (err error) [failretval==0] = winhttp.WinHttpCloseHandle
//sys winHTTPGetProxyForURL(whi winHTTPInternet, url *uint16, options *winHTTPAutoProxyOptions, proxyInfo *winHTTPProxyInfo) (err error) [failretval==0] = winhttp.WinHttpGetProxyForUrl
//sys winHTTPOpen(agent *uint16, accessType uint32, proxy *uint16, proxyBypass *uint16, flags uint32) (whi winHTTPInternet, err error) [failretval==0] = winhttp.WinHttpOpen
