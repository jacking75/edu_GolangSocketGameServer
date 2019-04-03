package gohipernetFake

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger, _ = zap.NewProduction()
)

type extendedZapConfig struct{
	MaxSize			int	`json:"maxSize"`
	MaxBackups 		int	`json:"maxBackups"`
	MaxAge 			int	`json:"maxAge"`
	zap.Config
}

func init_Log() {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	configJson, err := ioutil.ReadFile(filepath.FromSlash(currentDir + "/" + "config_logger.json"))
	if err != nil {
		panic(err)
	}

	var myConfig extendedZapConfig

	if err := json.Unmarshal(configJson, &myConfig); err != nil {
		panic(err)
	}

	for index := range myConfig.ErrorOutputPaths{
		if myConfig.ErrorOutputPaths[index] != "stderr"{
			myConfig.ErrorOutputPaths[index], _ = _createFileName(myConfig.ErrorOutputPaths[index])
		}
	}

	enc := zapcore.NewJSONEncoder(myConfig.EncoderConfig)
	Logger = zap.New(zapcore.NewCore(enc, _combineSinkFromConfig(myConfig), myConfig.Level))
}

func _combineSinkFromConfig(myConfig extendedZapConfig) zapcore.WriteSyncer{
	var fileName string
	stdOutLogOn := false
	for index:= range myConfig.OutputPaths{
		if myConfig.OutputPaths[index] != "stdout"{ // 그 외는 텍스트 파일
			fileName = myConfig.OutputPaths[index]
		} else {
			stdOutLogOn = true // 설정파일에 stdout이 있을 경우
		}
	}

	sink := zapcore.AddSync(
		&lumberjack.Logger{
			Filename: fileName,
			MaxSize: 	myConfig.MaxSize, // MB 단위
			MaxBackups: myConfig.MaxBackups,
			MaxAge: 	myConfig.MaxAge, // 28일 단위
		},
	)
	var combineSink zapcore.WriteSyncer
	if stdOutLogOn {
		combineSink = zap.CombineWriteSyncers(sink, os.Stdout)
	} else {
		combineSink = sink
	}
	return combineSink
}

func _createFileName(outputName string) (string, error){
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102_150405")
	fileNameArr := strings.Split(outputName, ".")
	fileName := fileNameArr[0]
	fileExt := "." + fileNameArr[1]
	if len(fileNameArr) > 2{
		return "", errors.New("log ouput name Invalid")
	}
	return fileName + formattedTime + fileExt, nil
}


var IExportLog func(string, string) = _emptyExportLog

func _emptyExportLog(level string, message string) {}

