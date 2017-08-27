package Parser

import  "go.uber.org/zap"
import "encoding/json"

var MyLogger *zap.Logger

func Init_logger()  {
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": [ "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],

	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}else{
		MyLogger = logger
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}


