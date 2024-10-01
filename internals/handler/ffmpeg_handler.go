package handler

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	StreamPath = "stream"
)

type Quality int

const (
	SHD Quality = iota
	MD
	SD
)

func getFFMPEGCmd(quality Quality) string {
	var resolution string
	var bitrate string

	switch quality {
	case SHD:
		resolution = "1280x720"
		bitrate = "2.5M"
	case MD:
		resolution = "960x540"
		bitrate = "1.5M"
	case SD:
		resolution = "854x480"
		bitrate = "1M"
	default:
		return ""
	}

	return fmt.Sprintf(`ffmpeg -f v4l2 -framerate 30 -video_size %s -i /dev/video0 -c:v libx264 -tune zerolatency -b:v %s -preset ultrafast -pix_fmt yuv420p -g 30 -keyint_min 30 -f dash -window_size 5 -extra_window_size 2 -remove_at_exit 1 -seg_duration 4 -streaming 1 -use_template 1 -use_timeline 1 -live 1 -time_shift_buffer_depth 10 -dash_segment_type mp4 %s/manifest.mpd`, resolution, bitrate, StreamPath)
}

func StartFFMPEGThread() {
	args := strings.Fields(getFFMPEGCmd(MD))
	cmd := exec.Command(args[0], args[1:]...)

	// Run the FFmpeg command
	fmt.Println("Starting FFmpeg for live streaming...")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("FFmpeg command failed: %v", err)
	}
}
