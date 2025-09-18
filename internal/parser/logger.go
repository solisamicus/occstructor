package parser

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type MismatchLogger struct {
	logFile *os.File
}

func NewMismatchLogger() (*MismatchLogger, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("logs/mismatch_%s.log", time.Now().Format("20060102_150405"))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &MismatchLogger{logFile: file}, nil
}

func (l *MismatchLogger) LogMismatch(lineNum int, codes []string, names []string, rawCodes, rawNames string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logEntry := fmt.Sprintf(`
========================================
Time: %s
Line: %d
Code Count: %d
Name Count: %d

CODES:
%s

RAW CODES:
%s

NAMES:
%s

RAW NAMES:
%s

SQL INSERT TEMPLATE:
-- Line %d manual insert
`, timestamp, lineNum, len(codes), len(names),
		formatArray(codes), rawCodes,
		formatArray(names), rawNames, lineNum)

	for i, code := range codes {
		if i < len(names) {
			sqlTemplate := fmt.Sprintf("INSERT INTO occupations (seq, name, level, parent_seq) VALUES ('%s', '%s', 4, '%s');\n",
				code, names[i], getParentSeq(code))
			logEntry += sqlTemplate
		} else {
			sqlTemplate := fmt.Sprintf("-- Missing name for code: %s\n", code)
			logEntry += sqlTemplate
		}
	}

	logEntry += "\n"
	l.logFile.WriteString(logEntry)
}

func (l *MismatchLogger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

func formatArray(arr []string) string {
	result := ""
	for i, item := range arr {
		result += fmt.Sprintf("[%d] %s\n", i, item)
	}
	return result
}

func getParentSeq(code string) string {
	parts := strings.Split(code, "-")
	if len(parts) >= 4 {
		return strings.Join(parts[:3], "-")
	}
	return ""
}
