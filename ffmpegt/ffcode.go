package ffmpegt

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// StartPushVideo2Cloud fpv视频流转到云端
func StartPushVideo2Cloud() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("start ffmpeg push video stream fail", r)
			time.Sleep(1 * time.Second)
			StartPushVideo2Cloud()
		}
	}()

	viper.SetConfigName("url.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %s ", err))
	}
	// // 从配置文件中读取配置项
	// dbHost := viper.GetString("database.host")
	// fmt.Println("------------------->")
	// fmt.Println("dbHost = ", dbHost)

	fpath := viper.GetString("path.ffmpeg")
	localURL := viper.GetString("url.local")
	remoteURL := viper.GetString("url.remote")

	// fpath = "E:/ffmpegYang/ffmpeg/bin/ffmpeg.exe"
	// localURL = "rtsp://192.168.1.4:554/channel=0,stream=0"
	// remoteURL = "rtmp://10.240.34.35:1935/live/1111111111"

	fmt.Println("fpath:", fpath)
	fmt.Println("localURL:", localURL)
	fmt.Println("remoteURL:", remoteURL)

	// ffmpegPath := "E:/ffmpegYang/ffmpeg/bin/ffmpeg.exe"
	// inputURL := "rtsp://192.168.1.4:554/channel=0,stream=0"
	// outputURL := "rtmp://10.240.34.35:1935/live/1111111111"

	// fmt.Println("inputURL : ", localURL)
	// fmt.Println("outputURL : ", remoteURL)

	err := ffmpeg.Input(localURL,
		ffmpeg.KwArgs{"rtsp_transport": "tcp"}).Output(remoteURL, ffmpeg.KwArgs{
		"vcodec": "copy",
		"acodec": "aac",
		"bf":     "0",
		"f":      "flv",
	}).SetFfmpegPath(fpath).Run()
	if err != nil {
		fmt.Println("Error: ", err)
		panic(fmt.Errorf("视频推流异常: %s ", err))
	}
	fmt.Println("starting push video...")
}

// ./ffmpeg -rtsp_transport tcp -i rtsp://192.168.1.4:554/channel=0,stream=0  -vcodec copy -acodec aac -bf 0  -f flv  rtmp://10.240.34.35:1935/live/1111111111

// // StartPushVideo2Cloud fpv视频流转到云端
// func StartPushVideo2Cloud() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			logger.Error("start ffmpeg push video stream fail", r)
// 			return
// 		}
// 	}()

// 	VideoToCloud := config.Global.VideoToCloud
// 	defaultFfmpegPath := "/data/ffmpeg"
// 	var ffmpegPath string
// 	if VideoToCloud == nil {
// 		ffmpegPath = defaultFfmpegPath
// 	} else {
// 		ffmpegPath = VideoToCloud.Ffmpeg
// 		if ffmpegPath == "" {
// 			ffmpegPath = defaultFfmpegPath
// 		}
// 	}
// 	logger.Info("ffmpegPath = ", ffmpegPath)

// 	inputURL := "rtsp://admin:Autel123@192.168.1.6:554/h264/ch33/main/av_stream"
// 	outputURL := "rtmp://43.192.105.32:1935/live/0000000000"

// 	if VideoToCloud.InputURL != "" {
// 		inputURL = VideoToCloud.InputURL
// 	}
// 	if VideoToCloud.OutputURL != "" {
// 		outputURL = VideoToCloud.OutputURL
// 	}

// 	logger.Info("inputURL : ", inputURL)
// 	logger.Info("outputURL : ", outputURL)

// 	err := ffmpeg.Input(inputURL,
// 		ffmpeg.KwArgs{"rtsp_transport": "tcp"}).Output(outputURL, ffmpeg.KwArgs{
// 		"vcodec": "copy",
// 		"acodec": "aac",
// 		"bf":     "0",
// 		"f":      "flv",
// 	}).SetFfmpegPath(ffmpegPath).Run()
// 	if err != nil {
// 		logger.Error("Error: ", err)
// 		return
// 	}
// }
