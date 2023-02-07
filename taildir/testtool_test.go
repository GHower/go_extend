package taildir

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// TailTest just for test
type TailTest struct {
	Name string // 测试组名，用于隔离各个测试方法所用目录
	path string
	done chan struct{}
	*testing.T
}

func NewTailTest(name string, t *testing.T) TailTest {
	tt := TailTest{name, ".test/" + name, make(chan struct{}), t}
	err := os.MkdirAll(tt.path, os.ModeTemporary|0700)
	if err != nil {
		tt.Fatal(err)
	}

	return tt
}

// CreateFile 创建文件相当于第一次的 echo "xx" >> file
func (t TailTest) CreateFile(name string, contents string) {
	err := ioutil.WriteFile(t.path+"/"+name, []byte(contents), 0600)
	if err != nil {
		t.Fatal(err)
	}
}

// AppendToFile 追加内容到文件,也是echo "xx" >> file
func (t TailTest) AppendToFile(name string, contents string) {
	err := ioutil.WriteFile(t.path+"/"+name, []byte(contents), 0600|os.ModeAppend)
	if err != nil {
		t.Fatal(err)
	}
}

//RemoveFile 相当于rm
func (t TailTest) RemoveFile(name string) {
	err := os.Remove(t.path + "/" + name)
	if err != nil {
		t.Fatal(err)
	}
}

func (t TailTest) RenameFile(oldname string, newname string) {
	oldname = t.path + "/" + oldname
	newname = t.path + "/" + newname
	err := os.Rename(oldname, newname)
	if err != nil {
		t.Fatal(err)
	}
}

// StartTailDir 开启taildir
func (t TailTest) StartTailDir(config Config) *TailD {
	tailD, err := TailDir(t.path, config)
	if err != nil {
		t.Fatal(err)
	}
	return tailD
}

// 验证
func (t TailTest) VerifyOutPut(d *TailD, expectLine Line) bool {
	select {
	case line := <-d.Lines:
		if expectLine.Filename == line.Filename && expectLine.Keyword == expectLine.Keyword {
			return true
		}
	case <-time.After(time.Second):
		return false
	}
	return false
}

func (t TailTest) WaitAndClose(d *TailD) bool {
	<-time.After(time.Second)
	d.tomb.Done()
	<-time.After(time.Second)
	return false
}
