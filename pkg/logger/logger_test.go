package logger

import (
	"testing"
	"time"
)

func TestZapLogger(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger)
	if err != nil {

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithTimeFormat(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithTimeFormat(time.DateTime))
	if err != nil {

		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithErrorOutPut(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithErrorOutPut("./sss"))
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestZapLoggerWithColor(t *testing.T) {
	logger, err := NewLoggerWithType(ZapLogger, WithColor(), WithFileOutput("./logger.log"))
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithFields(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger)
	if err != nil {

		t.Fatal(err)
	}
	loggerFiled := logger.WithFields(map[string]any{"key": "value"})
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	loggerFiled.Warn("warn:", "hello world")
	loggerFiled.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithErrorOutPut(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithErrorOutPut("./loggerout.log"), WithFileOutput("./sss"))
	if err != nil {

		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestLogrusLoggerWithColor(t *testing.T) {
	logger, err := NewLoggerWithType(LogrusLogger, WithColor(), WithFileOutput("./logger.log"))
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLogger(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger)
	if err != nil {

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
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

		t.Fatal(err)
	}
	loggerFiled := logger.WithFields(map[string]any{"key": "value"})
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	loggerFiled.Warn("warn:", "hello world")
	loggerFiled.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggerWithLogRotation(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithLogRotation("app.log", 1, 5, 0, true))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10000; i++ {
		logger.Info("Info:hello world")
		logger.Infof("Infof:%v", "hello world")
		logger.Warn("warn:", "hello world")
		logger.Warnf("warn:%v", "hello world")
		logger.Error("error:", "hello world")
		logger.Errorf("error:%v", "hello world")
	}
}

func TestSlogLoggeWithColor(t *testing.T) {
	logger, err := NewLoggerWithType(SlogLogger, WithColor(), WithFileOutput("./logger.log"))
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}

func TestSlogLoggeWithColorTheme(t *testing.T) {
	theme := ColorScheme{
		CodeType: CodeTypeANSI,
		Info: &Color{
			ansi: "\u001B[35m",
		},
	}
	logger, err := NewLoggerWithType(SlogLogger, WithColor(), WithFileOutput("./logger.log"), WithColorScheme(theme))
	if err != nil {
		t.Fatal(err)
	}
	logger.Info("Info:hello world")
	logger.Infof("Infof:%v", "hello world")
	logger.Warn("warn:", "hello world")
	logger.Warnf("warn:%v", "hello world")
	logger.Error("error:", "hello world")
	logger.Errorf("error:%v", "hello world")
}
