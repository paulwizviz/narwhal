package shared

const (
	osDarwin    = "darwin"
	osDragonFly = "dragonfly"
	osFreeBSD   = "freebsd"
	osLinux     = "linux"
	osNetBSD    = "netbsd"
	osOpenBSD   = "openbsd"
	osPlan9     = "plan9"
	osSolaris   = "solaris"
	osWindows   = "windows"
)

const (
	arch386      = "386"
	archARM      = "arm"
	archARM64    = "arm64"
	archAMD64    = "amd64"
	archPPC64le  = "ppc64le"
	archMIPS64   = "mips64"
	archMIPS64le = "mips64le"
	archS390X    = "s390x"
)

// Source: https://github.com/docker/cli/blob/1f6a1a438c4ae426e446f17848114e58072af2bb/cli/command/manifest/util.go
// of compatible platforms.
//   {os: "darwin", arch: "386"}
//   {os: "darwin", arch: "amd64"}
//   {os: "darwin", arch: "arm"}
//   {os: "darwin", arch: "arm64"}
//   {os: "dragonfly", arch: "amd64"}
//   {os: "freebsd", arch: "386"}
//   {os: "freebsd", arch: "amd64"}
//   {os: "freebsd", arch: "arm"}
//   {os: "linux", arch: "386"}
//   {os: "linux", arch: "amd64"}
//   {os: "linux", arch: "arm"}
//   {os: "linux", arch: "arm64"}
//   {os: "linux", arch: "ppc64le"}
//   {os: "linux", arch: "mips64"}
//   {os: "linux", arch: "mips64le"}
//   {os: "linux", arch: "s390x"}
//   {os: "netbsd", arch: "386"}
//   {os: "netbsd", arch: "amd64"}
//   {os: "netbsd", arch: "arm"}
//   {os: "openbsd", arch: "386"}
//   {os: "openbsd", arch: "amd64"}
//   {os: "openbsd", arch: "arm"}
//   {os: "plan9", arch: "386"}
//   {os: "plan9", arch: "amd64"}
//   {os: "solaris", arch: "amd64"}
//   {os: "windows", arch: "386"}
//   {os: "windows", arch: "amd64"}

// DockerPlatformConfig represents configuration of Docker
// platform
type DockerPlatformConfig struct {
	OS   string
	Arch string
}

// PlatformDarwin386 returns platform configured for
// Darwin and 386
func PlatformDarwin386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osDarwin,
		Arch: arch386,
	}
}

// PlatformDarwinARM returns platform configured for
// Darwin and ARM
func PlatformDarwinARM() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osDarwin,
		Arch: archARM,
	}
}

// PlatformDarwinARM64 returns platform configured for
// Darwin and ARM64
func PlatformDarwinARM64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osDarwin,
		Arch: archARM64,
	}
}

// PlatformDragonFlyAMD64 returns platform configured for
// DragonFly and ARM64
func PlatformDragonFlyAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osDarwin,
		Arch: archAMD64,
	}
}

// PlatformFreeBSD386 returns platform configured for
// FreeBSD and 386
func PlatformFreeBSD386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osFreeBSD,
		Arch: arch386,
	}
}

// PlatformFreeBSDAMD64 returns platform configured for
// FreeBSD and AMD64
func PlatformFreeBSDAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osFreeBSD,
		Arch: archAMD64,
	}
}

// PlatformFreeBSDARM returns platform configured for
// FreeBSD and ARM
func PlatformFreeBSDARM() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osFreeBSD,
		Arch: archARM,
	}
}

// PlatformDarwinArm64 returns platform configured for
// Linux and 386
func PlatformLinux386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: arch386,
	}
}

// PlatformLinuxAMD64 returns platform configured for
// Linux and AMD64
func PlatformLinuxAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archAMD64,
	}
}

// PlatformLinuxARM returns platform configured for
// Linux and ARM
func PlatformLinuxARM() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archARM,
	}
}

// PlatformLinuxARM returns platform configured for
// Linux and ARM
func PlaformLinuxARM64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archARM64,
	}
}

// PlaformLinuxMIPS64 returns platform configured for
// Linux and MIPS64
func PlaformLinuxMIPS64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archMIPS64,
	}
}

// PlatformLinuxPPC64le returns platform configured for
// Linux and PPC64le
func PlatformLinuxPPC64le() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archPPC64le,
	}
}

// PlatformLinuxS390X returns platform configured for
// Linux and s390x
func PlatformLinuxS390X() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osLinux,
		Arch: archS390X,
	}
}

// PlatformNetBSDS386 returns platform configured for
// NetBSD and 386
func PlatformNetBSDS386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osNetBSD,
		Arch: arch386,
	}
}

// PlatformNetBSDAMD64 returns platform configured for
// NetBSD and AMD64
func PlatformNetBSDAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osNetBSD,
		Arch: archAMD64,
	}
}

// PlatformNetBSDARM returns platform configured for
// NetBSD and ARM
func PlatformNetBSDARM() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osNetBSD,
		Arch: archARM,
	}
}

// PlatformOpenBSD386 returns platform configured for
// OpenBSD and ARM
func PlatformOpenBSD386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osOpenBSD,
		Arch: archARM,
	}
}

// PlatformOpenBSDAMD64 returns platform configured for
// OpenBSD and AMD64
func PlatformOpenBSDAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osOpenBSD,
		Arch: archAMD64,
	}
}

// PlatformOpenBSDARM returns platform configured for
// OpenBSD and ARM
func PlatformOpenBSDARM() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osOpenBSD,
		Arch: archARM,
	}
}

// PlatformPlan9386 returns platform configured for
// Plan9 and 380
func PlatformPlan9386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osPlan9,
		Arch: arch386,
	}
}

// PlatformPlan9AMD64 returns platform configured for
// Plan9 and AMD64
func PlatformPlan9AMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osPlan9,
		Arch: archAMD64,
	}
}

// PlatformSolarisAMD64 returns platform configured for
// Solaris and AMD64
func PlatformSolarisAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osPlan9,
		Arch: archAMD64,
	}
}

// Windows386 returns platform configured for
// Windows and 386
func Windows386() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osPlan9,
		Arch: arch386,
	}
}

// WindowsAMD64 returns platform configured for
// Windows and AMD64
func WindowsAMD64() DockerPlatformConfig {
	return DockerPlatformConfig{
		OS:   osPlan9,
		Arch: arch386,
	}
}
