package logger

import (
	"testing"
	"time"
)

func TestZapLogger(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger)
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithLevel(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithLevel(ErrorLevel))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithAddSource(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithAddSource())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithFileOutput(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithFileOutput("./zap.log"))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithJSONFormat(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithJSONFormat())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithTimeFormat(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithTimeFormat(time.DateTime))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLogger(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger)
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithLevel(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithLevel(ErrorLevel))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithAddSource(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithAddSource())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithFileOutput(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithFileOutput("./logger2.log"))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggeWithJSONFormat(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithJSONFormat())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggrWithTimeFormat(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithTimeFormat(time.DateTime))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggrWithFields(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger)
	if err != nil {
		t.Error(err)
	}
	loggerFiled := logger.WithFields(map[string]any{"key": "value"})
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	loggerFiled.Warn("warn:", "hello world")
	loggerFiled.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLogger(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger)
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:%v", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithLevel(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithAddSource(), WithLevel(ErrorLevel))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithAddSource(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithAddSource())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithFileOutput(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithFileOutput("./logger.log"))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithJSONFormat(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithJSONFormat())
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithTimeFormat(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithTimeFormat(time.Stamp))
	if err != nil {
		t.Error(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithFields(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithTimeFormat(time.Stamp))
	if err != nil {
		t.Error(err)
	}
	loggerFiled := logger.WithFields(map[string]any{"key": "value"})
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	loggerFiled.Warn("warn:", "hello world")
	loggerFiled.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithErrorOutPut(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithErrorOutPut("./loggerout.log"))
	if err != nil {
		t.Error(err)
	}
	loggerFiled := logger.WithFields(map[string]any{"key": "value"})
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	loggerFiled.Warn("warn:", "hello world")
	loggerFiled.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}
