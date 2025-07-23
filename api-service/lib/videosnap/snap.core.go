package videosnap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	redis "telegram-service/config"
	"telegram-service/lib/desks"
	"telegram-service/lib/videosnap/setting"
	storeconfig "telegram-service/lib/videosnap/store.config"
	"telegram-service/utils"

	"time"
)

func SnapShootImage(gameNo string, streamkey string) {
	// Get config stream from Redis
	redisDeskKey := fmt.Sprintf("deskSetting:table:%s", streamkey)
	tableSetting, err := redis.Get(redisDeskKey)
	if err != nil {
		utils.ErrorLog(fmt.Errorf("SnapShootImage Failed to get data from Redis: %s", redisDeskKey), err.Error())
	}
	var tableSnapSetting desks.DeskSettingResponseDTO
	if err := json.Unmarshal([]byte(tableSetting), &tableSnapSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("SnapShootImage failed to unmarshal value from Redis: %s", redisDeskKey), err.Error())
	}
	var Server_ID uint // Change to appropriate type if necessary
	if tableSnapSetting.DeskStreamServer != 0 {
		Server_ID = tableSnapSetting.DeskStreamServer
	} else {
		Server_ID = 1 // Default Server ID if StreamServer.ID is not set
	}
	redisServerKey := fmt.Sprintf("videoSetting:server:%d", Server_ID)
	videoConfig, err := redis.Get(redisServerKey)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return
	}
	var videoSnapSetting setting.VideoSnapSettingResponseDTO
	if err := json.Unmarshal([]byte(videoConfig), &videoSnapSetting); err != nil {
		utils.ErrorLog(nil, err.Error())
		return
	}
	currentDate := time.Now().Format("20060102")
	outputDir := fmt.Sprintf("%s/%s/%s", videoSnapSetting.OutputPath, currentDate, streamkey)
	VideoAddr := fmt.Sprintf("%s%s", videoSnapSetting.Rtmp, streamkey)
	ImageOutputPath := fmt.Sprintf("%s/%s/%s/%s.jpeg",
		videoSnapSetting.OutputPath, // Base output path from video settings
		currentDate,                 // Current date in YYYYMMDD format
		streamkey,                   // Use the stream key as folder name
		streamkey+gameNo,            // Use the stream key in the file name
	)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil { // Set permissions to 755
			utils.ErrorLog(nil, err.Error())
			return
		}
		utils.InfoLog(fmt.Sprintf("Created directory: %s with permissions 755", outputDir), string(utils.SuccessMessage))
	}
	currentTime := time.Now()
	snapshotTime := currentTime.Add(-3 * time.Second)
	formattedTime := fmt.Sprintf("%02d:%02d:%02d", snapshotTime.Hour(), snapshotTime.Minute(), snapshotTime.Second())
	// Debug logs for the command execution
	utils.InfoLog(fmt.Sprintf("Preparing to run FFmpeg with the following parameters:\n- Snapshot Time: %s\n- Input Video Address: %s\n- Output Path: %s", formattedTime, VideoAddr, ImageOutputPath), string(utils.SuccessMessage))
	// Stop video Snap
	taskredisKey := fmt.Sprintf("task:snapshoot:%s", gameNo)
	videoSnapTask, err := redis.Get(taskredisKey)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return
	}
	var videotask VideosnapResponseDTO
	if err := json.Unmarshal([]byte(videoSnapTask), &videotask); err != nil {
		utils.ErrorLog(nil, err.Error())
		return
	}
	command := exec.Command("ffmpeg",
		"-i", VideoAddr, // Use RTMP from settings
		"-ss", "00:00:02", // Set the timestamp to capture (2 seconds before the current time)
		"-y",           // Overwrite output files without asking
		"-f", "image2", // Set the output format to image2
		"-frames:v", "1", // Capture only 1 frame
		"-hide_banner",       // Hide the banner
		"-loglevel", "error", // Set log level to error
		ImageOutputPath, // Use the output path from settings
	)
	// Start the command
	if err := command.Start(); err != nil {
		utils.ErrorLog("Failed to start FFmpeg", err.Error())
		return
	}
	utils.InfoLog(fmt.Sprintf("Successfully captured image at %s", ImageOutputPath), string(utils.SuccessMessage))
	VideoOutputPath := fmt.Sprintf("%s/%s/%s/%s.mp4",
		videoSnapSetting.OutputPath, // Base output path from video settings
		currentDate,                 // Current date in YYYYMMDD format
		streamkey,                   // Use the stream key as folder name
		streamkey+gameNo,
	)
	fileupload := []string{
		VideoOutputPath,
		ImageOutputPath,
	}
	UploadPrefix := fmt.Sprintf("%s/%s", currentDate, streamkey)
	pid := command.Process.Pid
	go func() {
		killFfmpegProcess(videotask.ProcessId)
		storeconfig.UploadFilesToGCS(videoSnapSetting.BucketName, videoSnapSetting.ServiceAcccount, UploadPrefix, fileupload, videoSnapSetting.DeleteLocalStore)
		done := make(chan error, 1)
		go func() {
			done <- command.Wait() // Wait for FFmpeg to complete in a separate goroutine
		}()

		select {
		case err := <-done:
			if err != nil {
				utils.WarnLog(fmt.Sprintf("FFmpeg process terminated with PID: %d, Warring: %v", pid, err), string(utils.SuccessMessage))
				killFfmpegProcess(uint(pid))
			} else {
				utils.InfoLog(fmt.Sprintf("FFmpeg process completed with PID: %d", pid), string(utils.SuccessMessage))
			}
		case <-time.After(3 * time.Second): // 2-second timeout
			utils.ErrorLog(fmt.Sprintf("FFmpeg process did not complete on task, Killing PID: %d", pid), string(utils.ErrorMessage))
			killFfmpegProcess(uint(pid))
		}
	}()
}

func SnapVideo(gameNo string, streamkey string) (processId int, err error) {
	// Lookup in Table Setting if set
	redisDeskKey := fmt.Sprintf("deskSetting:table:%s", streamkey)
	tableSetting, err := redis.Get(redisDeskKey)
	if err != nil {
		utils.ErrorLog(fmt.Errorf("Failed to get data from Redis: %s", redisDeskKey), err.Error())
	}
	var tableSnapSetting desks.DeskSettingResponseDTO

	if err := json.Unmarshal([]byte(tableSetting), &tableSnapSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %s", redisDeskKey), err.Error())
	}
	utils.InfoLog(tableSnapSetting, string(utils.SuccessMessage))
	var Server_ID uint // Change to appropriate type if necessary
	if tableSnapSetting.DeskStreamServer != 0 {
		Server_ID = tableSnapSetting.DeskStreamServer
	} else {
		Server_ID = 1 // Default Server ID if StreamServer.ID is not set
	}
	utils.InfoLog(Server_ID, string(utils.SuccessMessage))

	redisServerKey := fmt.Sprintf("videoSetting:server:%d", Server_ID)
	videoInfo, err := redis.Get(redisServerKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get video settings from Redis: %v", err)
	}
	var videoSnapSetting setting.VideoSnapSettingResponseDTO
	if err := json.Unmarshal([]byte(videoInfo), &videoSnapSetting); err != nil {
		return 0, fmt.Errorf("failed to unmarshal value from Redis: %v", err)
	}
	currentDate := time.Now().Format("20060102")
	outputDir := fmt.Sprintf("%s/%s/%s", videoSnapSetting.OutputPath, currentDate, streamkey)
	VideoAddr := fmt.Sprintf("%s%s", videoSnapSetting.Rtmp, streamkey)
	outputPath := fmt.Sprintf("%s/%s/%s/%s.mp4",
		videoSnapSetting.OutputPath,
		currentDate,
		streamkey,
		streamkey+gameNo,
	)
	// Check if the output directory exists, create it if it doesn't
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			utils.ErrorLog(err.Error(), string(utils.ErrorMessage))
			return 0, fmt.Errorf("error creating directory: %v", err)
		}
		utils.InfoLog(fmt.Sprintf("Created directory: %s with permissions 755", outputDir), string(utils.SuccessMessage))
	} else {
		utils.InfoLog(fmt.Sprintf("Directory already exists: %s", outputDir), string(utils.SuccessMessage))
	}
	// Prepare the FFmpeg command
	command := exec.Command("ffmpeg", "-i", VideoAddr,
		"-y",
		"-f", "mp4",
		"-c", "copy",
		"-t", fmt.Sprintf("%d", videoSnapSetting.Duration),
		"-preset", videoSnapSetting.VideoType,
		"-s", videoSnapSetting.VideoSize,
		"-loglevel", "error",
		outputPath,
	) // Use error level to minimize output noise
	var outBuf, errBuf bytes.Buffer
	command.Stdout = &outBuf
	command.Stderr = &errBuf

	// Start the FFmpeg process
	if err := command.Start(); err != nil {
		utils.ErrorLog(fmt.Sprintf("failed to start FFmpeg: %v", err), string(utils.ErrorMessage))
		return 0, fmt.Errorf("failed to start FFmpeg: %v", err)
	}
	pid := command.Process.Pid
	utils.InfoLog(fmt.Sprintf("FFmpeg process started with PID: %d", pid), string(utils.SuccessMessage))

	go func() {
		done := make(chan error, 1)
		go func() {
			done <- command.Wait() // Wait for FFmpeg to complete in a separate goroutine
		}()

		select {
		case err := <-done:
			if err != nil {
				utils.WarnLog(fmt.Sprintf("FFmpeg process terminated with PID: %d, Warring: %v", pid, err), string(utils.SuccessMessage))
				killFfmpegProcess(uint(pid))
			} else {
				utils.InfoLog(fmt.Sprintf("FFmpeg process completed with PID: %d", pid), string(utils.SuccessMessage))
			}
		case <-time.After(time.Duration(videoSnapSetting.Duration) * time.Second): // 1-minute timeout
			utils.ErrorLog(fmt.Sprintf("FFmpeg process did not complete task in :%d Seconds, Os Killing PID: %d", videoSnapSetting.Duration, pid), string(utils.ErrorMessage))
			killFfmpegProcess(uint(pid))
		}
	}()

	return pid, nil
}

func killFfmpegProcess(processId uint) error {
	proc, err := os.FindProcess(int(processId))
	if err != nil {
		return fmt.Errorf("could not find process: %v", err)
	}
	// Try graceful termination first
	if err := proc.Signal(os.Interrupt); err != nil { // SIGINT
		return fmt.Errorf("failed to send interrupt signal: %v", err)
	}
	// Wait for a moment to allow cleanup
	time.Sleep(3 * time.Second)
	// Check if the process is still running
	if err := proc.Signal(os.Kill); err != nil { // SIGKILL
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}
