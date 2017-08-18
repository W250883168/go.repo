package xcmd

import "os/exec"

type XCommand struct {
}

func exec() {

	cmd := exec.Command("cmd.exe", "/c", "start "+datapath)

	err := cmd.Run()
}

func ffmpeg_cmd() {
	/*go func() {
		filepath := "live555.mp4"
		istream := fmt.Sprintf("rtsp://%s:%s@%s:%d/cam/realmonitor?channel=1&subtype=0", vc.CameraLoginUser, vc.CameraLoginPass, vc.CameraIp, vc.CameraPort)
		ffmpeg_cmd := fmt.Sprintf("start ffmpeg -re -i %s -acodec copy -vcodec libx264 -pix_fmt yuv420p -f mp4 %s", istream, filepath)
		pCmd := exec.Command("cmd.exe ", "/C", ffmpeg_cmd)
		var err error
		if err = command.Start(); err != nil {
			if err = cmd.Wait(); err != nil {
				//TODO 正常录制完毕，写数据库数据
			}
		}

		if err != nil {
			//TODO 录制失败，记录日志/写数据库
			xdebug.PrintStackTrace(err)
		}
	}() */
}
