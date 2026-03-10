package chrome

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bogdanfinn/tls-client/profiles"
)

type Profile struct {
	Major       int
	Impersonate string
	Build       int
	PatchMin    int
	PatchMax    int
	SecChUA     string
}

var chromeProfiles = []Profile{
	{131, "chrome131", 6778, 0, 300, "\"Chromium\";v=\"131\", \"Google Chrome\";v=\"131\", \"Not_A Brand\";v=\"24\""},
}

func RandomChromeVersion() (Profile, string, string) {
	rand.Seed(time.Now().UnixNano())
	profile := chromeProfiles[rand.Intn(len(chromeProfiles))]
	patch := rand.Intn(profile.PatchMax-profile.PatchMin+1) + profile.PatchMin
	fullVersion := fmt.Sprintf("%d.0.%d.%d", profile.Major, profile.Build, patch)
	userAgent := fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", fullVersion)
	return profile, fullVersion, userAgent
}

func MapToTLSProfile(impersonate string) profiles.ClientProfile {
	switch impersonate {
	case "chrome131":
		return profiles.Chrome_131
	case "chrome133a":
		return profiles.Chrome_133
	case "chrome136":
		// Fallback to Chrome_133 as Chrome_136 is not available in tls-client v1.14.0
		return profiles.Chrome_133
	case "chrome142":
		// Fallback to Chrome_133 as Chrome_142 is not available in tls-client v1.14.0
		return profiles.Chrome_133
	default:
		return profiles.Chrome_133
	}
}
