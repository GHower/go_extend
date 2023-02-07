package taildir

import (
	"testing"
)

func Test_TailDir(t *testing.T) {
	tailTest := NewTailTest("taildir", t)
	tailTest.CreateFile("1.log", "")
	tailDir := tailTest.StartTailDir(Config{
		Logger: nil,
	})
	tailTest.AppendToFile("1.log", "test hello 2222")

	go tailTest.VerifyOutPut(tailDir, Line{
		Filename: "1.log",
		Keyword:  "test hello",
	})
	tailTest.WaitAndClose(tailDir)
}
