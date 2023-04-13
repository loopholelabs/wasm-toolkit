package wasmfile

var Debug_wasi_snapshot_preview1 = map[string]string{
	"args_get":               "args_get(argv, argv_buf)",
	"args_sizes_get":         "args_sizes_get(argc, argvBufSize)",
	"environ_get":            "environ_get(environ, environBuf)",
	"environ_sizes_get":      "environ_sizes_get(environCount, environBufSize)",
	"clock_res_get":          "clock_res_get(clockId, resolution)",
	"clock_time_get":         "clock_time_get(clockId, precision, time)",
	"fd_advise":              "fd_advise(fd, offset, len, advice)",
	"fd_allocate":            "fd_allocate(fd, offset, len)",
	"fd_close":               "fd_close(fd)",
	"fd_datasync":            "fd_datasync(fd)",
	"fd_fdstat_get":          "fd_fdstat_get(fd, buffer)",
	"fd_fdstat_set_flags":    "fd_fdstat_set_flags(fd, flags)",
	"fd_fdstat_set_rights":   "fd_fdstat_set_rights(fd, fsRightsBase, fdRightsInheriting)",
	"fd_filestat_get":        "fd_filestat_get(fd, bufPtr)",
	"fd_filestat_set_size":   "fd_filestat_set_size(fd, stSize)",
	"fd_filestat_set_times":  "fd_filestat_set_times(fd, stAtim, stMtim, fstFlags)",
	"fd_prestat_get":         "fd_prestat_get(fd, buffer)",
	"fd_prestat_dir_name":    "fd_prestat_dir_name(fd, pathPtr, pathLen)",
	"fd_pwrite":              "fd_pwrite(fd, iovs, iovsLen, offset, nwritten)",
	"fd_write":               "fd_write(fd, iovs, iovsLen, nwritten)",
	"fd_pread":               "fd_pread(fd, iovs, iovsLen, offset nread)",
	"fd_read":                "fd_read(fd, iovs, iovsLen, nread)",
	"fd_readdir":             "fd_readdir(fd, bufPtr, bufLen, cookie, bufusedPtr)",
	"fd_renumber":            "fd_renumber(from, to)",
	"fd_seek":                "fd_seek(fd, offset, whence, newOffsetPtr)",
	"fd_tell":                "fd_tell(fd, offsetPtr)",
	"fd_sync":                "fd_sync(fd)",
	"path_create_directory":  "path_create_directory(fd, pathPtr, pathLen)",
	"path_filestat_get":      "path_filestat_get(fd, flags, pathPtr, pathLen, bufPtr)",
	"path_filestat_set_time": "path_filestat_set_time(fd, fstflags, pathPtr, pathLen, stAtim, stMtim)",
	"path_link":              "path_link(oldFd, oldFlags, oldPath, oldPathLen, newFd, newPath, newPathLen)",
	"path_open":              "path_open(dirfd, dirflags, pathPtr, pathLen, oflags, fsRightsBase, fsRightsInheriting, fsFlags, fd)",
	"path_readlink":          "path_readlink(fd, pathPtr, pathLen, buf, bufLen, bufused)",
	"path_remove_directory":  "path_remove_directory(fd, pathPtr, pathLen)",
	"path_rename":            "path_rename(oldFd, oldPath, oldPathLen, newFd, newPath, newPathLen)",
	"path_symlink":           "path_symlink(oldPath, oldPathLen, fd, newPath, newPathLen)",
	"path_unlink_file":       "path_unlink_file(fd, pathPtr, pathLen)",
	"poll_oneoff":            "poll_oneoff(sin, sout, nsubscriptions, nevents)",
	"proc_exit":              "proc_exit(rval)",
	"proc_raise":             "proc_raise(sig)",
	"random_get":             "random_get(bufPtr, bufLen)",
	"sched_yield":            "sched_yield()",
	"sock_recv":              "sock_recv",
	"sock_send":              "sock_send",
	"sock_shutdown":          "sock_shutdown",
}

var Wasi_errors = map[string]int{
	"WASI_ESUCCESS":        0,
	"WASI_E2BIG":           1,
	"WASI_EACCES":          2,
	"WASI_EADDRINUSE":      3,
	"WASI_EADDRNOTAVAIL":   4,
	"WASI_EAFNOSUPPORT":    5,
	"WASI_EAGAIN":          6,
	"WASI_EALREADY":        7,
	"WASI_EBADF":           8,
	"WASI_EBADMSG":         9,
	"WASI_EBUSY":           10,
	"WASI_ECANCELED":       11,
	"WASI_ECHILD":          12,
	"WASI_ECONNABORTED":    13,
	"WASI_ECONNREFUSED":    14,
	"WASI_ECONNRESET":      15,
	"WASI_EDEADLK":         16,
	"WASI_EDESTADDRREQ":    17,
	"WASI_EDOM":            18,
	"WASI_EDQUOT":          19,
	"WASI_EEXIST":          20,
	"WASI_EFAULT":          21,
	"WASI_EFBIG":           22,
	"WASI_EHOSTUNREACH":    23,
	"WASI_EIDRM":           24,
	"WASI_EILSEQ":          25,
	"WASI_EINPROGRESS":     26,
	"WASI_EINTR":           27,
	"WASI_EINVAL":          28,
	"WASI_EIO":             29,
	"WASI_EISCONN":         30,
	"WASI_EISDIR":          31,
	"WASI_ELOOP":           32,
	"WASI_EMFILE":          33,
	"WASI_EMLINK":          34,
	"WASI_EMSGSIZE":        35,
	"WASI_EMULTIHOP":       36,
	"WASI_ENAMETOOLONG":    37,
	"WASI_ENETDOWN":        38,
	"WASI_ENETRESET":       39,
	"WASI_ENETUNREACH":     40,
	"WASI_ENFILE":          41,
	"WASI_ENOBUFS":         42,
	"WASI_ENODEV":          43,
	"WASI_ENOENT":          44,
	"WASI_ENOEXEC":         45,
	"WASI_ENOLCK":          46,
	"WASI_ENOLINK":         47,
	"WASI_ENOMEM":          48,
	"WASI_ENOMSG":          49,
	"WASI_ENOPROTOOPT":     50,
	"WASI_ENOSPC":          51,
	"WASI_ENOSYS":          52,
	"WASI_ENOTCONN":        53,
	"WASI_ENOTDIR":         54,
	"WASI_ENOTEMPTY":       55,
	"WASI_ENOTRECOVERABLE": 56,
	"WASI_ENOTSOCK":        57,
	"WASI_ENOTSUP":         58,
	"WASI_ENOTTY":          59,
	"WASI_ENXIO":           60,
	"WASI_EOVERFLOW":       61,
	"WASI_EOWNERDEAD":      62,
	"WASI_EPERM":           63,
	"WASI_EPIPE":           64,
	"WASI_EPROTO":          65,
	"WASI_EPROTONOSUPPORT": 66,
	"WASI_EPROTOTYPE":      67,
	"WASI_ERANGE":          68,
	"WASI_EROFS":           69,
	"WASI_ESPIPE":          70,
	"WASI_ESRCH":           71,
	"WASI_ESTALE":          72,
	"WASI_ETIMEDOUT":       73,
	"WASI_ETXTBSY":         74,
	"WASI_EXDEV":           75,
	"WASI_ENOTCAPABLE":     76,
}

func GetWasiParamCodeEnter(wasi_name string) string {

	// Show the arguments for path_open
	if wasi_name == "path_open" {
		// Print out path string
		return `i32.const offset($dd_wasi_var_path)
					i32.const length($dd_wasi_var_path)
					call $debug_func_wasi_context
					local.get 2
					local.get 3
					call $debug_print
					call $debug_func_wasi_done
					`
	} else if wasi_name == "path_create_directory" {
		// Print out path string
		return `i32.const offset($dd_wasi_var_path)
					i32.const length($dd_wasi_var_path)
					call $debug_func_wasi_context
					local.get 1
					local.get 2
					call $debug_print
					call $debug_func_wasi_done
					`
	} else if wasi_name == "path_remove_directory" {
		// Print out path string
		return `i32.const offset($dd_wasi_var_path)
					i32.const length($dd_wasi_var_path)
					call $debug_func_wasi_context
					local.get 1
					local.get 2
					call $debug_print
					call $debug_func_wasi_done
					`
	} else if wasi_name == "path_unlink_file" {
		// Print out path string
		return `i32.const offset($dd_wasi_var_path)
					i32.const length($dd_wasi_var_path)
					call $debug_func_wasi_context
					local.get 1
					local.get 2
					call $debug_print
					call $debug_func_wasi_done
					`
	}
	return ""
}

func GetWasiParamCodeExit(wasi_name string) string {
	if wasi_name == "fd_prestat_dir_name" {
		// Show the dir_name
		return `i32.const offset($dd_wasi_var_path)
					i32.const length($dd_wasi_var_path)
					call $debug_func_wasi_context
					local.get 1
					local.get 2
					call $debug_print
					call $debug_func_wasi_done
					`
	}
	return ""
}
