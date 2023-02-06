package taildir

import (
	"errors"
	"github.com/nxadm/tail/ratelimiter"
	"go_extend/taildir/matcher"
	"gopkg.in/tomb.v1"
	"os"
	"sync"
	"time"
)

var (
	errStop = errors.New("tail dir stop")
)

type logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Config struct {
	LineBufferSize  uint // 行缓冲大小
	ErrorBufferSize uint // 非致命错误缓冲大小
	FilenameMatcher matcher.Matcher
	RateLimiter     *ratelimiter.LeakyBucket
	// 日志记录器,若禁用日志记录，设置为DiscardingLogger
	Logger logger
}

// TailD 对目录的tail
type TailD struct {
	Dirname string
	Lines   chan *Line
	Errors  chan error
	Config

	tails map[string]*Tail

	dead error // 致命错误，拉都拉不起来时，err不为nil
	lk   sync.RWMutex
	tomb tomb.Tomb
}

// 1. 为符合条件的新增文件建立监听
// 2. 为已存在文件建立监听
func (d *TailD) watchFileSync() {
	defer d.tomb.Done()
	defer d.close()

	dir, err := os.ReadDir(d.Dirname)
	if err != nil {
		d.dead = errStop
		return
	}
	for _, entry := range dir {
		if !entry.IsDir() && d.FilenameMatcher.Match(entry.Name()) {
			f, err := TailF(entry.Name())
			if err != nil {
				// todo: 统一发送错误方法
				d.Errors <- err
				continue
			}
		}
	}
	// 触发create文件时，拉起新的tail
	for {
		select {}
	}
}

func (d *TailD) Err() (reason error) {
	d.lk.Lock()
	reason = d.dead
	d.lk.Unlock()
	return
}

// 对tailDir的关闭
// 1. 关闭lines管道
// 2. 关闭它所打开的文件
func (d *TailD) close() {
	close(d.Lines)
	d.closeFiles()
}

func (d *TailD) closeFiles() {
	for _, tail := range d.tails {
		if tail.file != nil {
			tail.file.Close()
			tail.file = nil
		}
	}
}

type Line struct {
	Filename string    // 文件名
	Keyword  string    // 触发关键字
	SeekInfo SeekInfo  // seek信息
	Time     time.Time // 时间
	Num      int       // 行号
	Text     string    // 行内容
}

func NewLine(text, filename, kw string, lineNum int) *Line {
	return &Line{
		Filename: filename,
		Keyword:  kw,
		Text:     text,
		Time:     time.Now(),
	}
}

// SeekInfo represents arguments to io.Seek. See: https://golang.org/pkg/io/#SectionReader.Seek
type SeekInfo struct {
	Offset int64
	Whence int
}

func TailDir(dirname string, config Config) (*TailD, error) {
	t := &TailD{
		Dirname: dirname,
		Config:  config,
		tails:   make(map[string]*Tail),
	}
	if t.LineBufferSize == 0 {
		t.LineBufferSize = 1000
	}
	if t.ErrorBufferSize == 0 {
		t.ErrorBufferSize = 10
	}
	t.Lines = make(chan *Line, t.LineBufferSize)
	t.Errors = make(chan error, t.ErrorBufferSize)
	// 自动监听目录下符合文件名规则的文件
	go t.watchFileSync()
	return t, nil
}
