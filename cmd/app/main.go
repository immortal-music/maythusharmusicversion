package main // <--- ဒီစာကြောင်းကို ဒီနေရာမှာ ထည့်ပေးပါ

/*
#cgo CFLAGS: -I../../
#cgo linux LDFLAGS: -L ../../ -lntgcalls -lm -lz
#cgo darwin LDFLAGS: -L ../../ -lntgcalls -lc++ -lz -lbz2 -liconv -framework AVFoundation -framework AudioToolbox -framework CoreAudio -framework QuartzCore -framework CoreMedia -framework VideoToolbox -framework AppKit -framework Metal -framework MetalKit -framework OpenGL -framework IOSurface -framework ScreenCaptureKit

// Currently is supported only dynamically linked library on Windows due to
// https://github.com/golang/go/issues/63903
#cgo windows LDFLAGS: -L../../ -lntgcalls
#include "ntgcalls/ntgcalls.h"
#include "glibc_compatibility.h"
*/
import "C"

import (
	"github.com/Laky-64/gologging"

	"github.com/immortal-music/maythusharmusicversion/config"
	"github.com/immortal-music/maythusharmusicversion/internal/cookies"
	"github.com/immortal-music/maythusharmusicversion/internal/core"
	"github.com/immortal-music/maythusharmusicversion/internal/database"
	"github.com/immortal-music/maythusharmusicversion/internal/modules"
)

func main() {
	gologging.SetLevel(gologging.DebugLevel)
	gologging.GetLogger("webrtc").SetLevel(gologging.WarnLevel)

	l := gologging.GetLogger("Main")

	l.Debug("🔹 Initializing MongoDB...")
	dbCleanup := database.Init(config.MongoURI)
	defer dbCleanup()
	l.Info("✅ Database connected successfully")

	go database.MigrateData(config.MongoURI)

	l.Debug("🔹 Initializing cookies...")
	cookies.Init()

	l.Debug("🔹 Initializing clients...")
	cleanup := core.Init(config.ApiID, config.ApiHash, config.Token, config.StringSession, config.LoggerID)
	defer cleanup()
	modules.Init(core.Bot, core.UBot, core.Ntg)
	l.Info("🚀 Bot is started")
	core.Bot.Idle()
}
