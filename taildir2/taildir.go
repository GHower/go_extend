package taildir

import (
	"errors"
	"go_extend/taildir/matcher"
	"go_extend/taildir/ratelimiter"
	"go_extend/taildir/watch"
	"gopkg.in/tomb.v1"
	"path/filepath"
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
	LineBufferSize  uint            // 行缓冲大小，可选
	ErrorBufferSize uint            // 非致命错误缓冲大小,可选
	FilenameMatcher matcher.Matcher // 文件名匹配器,可选
	FileLineMatcher matcher.Matcher // 文件行内容匹配器,可选
	LineProcesSync  LineProcesser   // 行结果异步处理方法,内部自动开协程执行
	RateLimiter     *ratelimiter.LeakyBucket
	// 日志记录器,若禁用日志记录，设置为DiscardingLogger
	Logger logger
	Pool   bool // 轮询模式
}

// TailD 对目录的tail
type TailD struct {
	Dirname string
	Errors  chan error
	Config

	tails map[string]*Tail
	lines chan *Line

	watch watch.FileWatcher
	dead  error // 致命错误，拉都拉不起来时，err不为nil
	lk    sync.RWMutex
	tomb  tomb.Tomb
}

// 1. 为符合条件的新增文件建立监听
// 2. 为已存在文件建立监听
func (d *TailD) watchDirSync() {
}

func (d *TailD) Err() (reason error) {
	return nil
}

// 对tailDir的关闭
// 1. 关闭lines管道
// 2. 关闭它所打开的文件
func (d *TailD) close() {
}

func (d *TailD) closeFiles() {
}

func (d *TailD) TailF(filename string) (*Tail, error) {
	return nil, nil
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
		Dirname: filepath.Clean(dirname),
		tails:   make(map[string]*Tail),
	}

	if config.FilenameMatcher == nil {
		config.FilenameMatcher = matcher.NewFileSufMatcher()
	}
	if config.LineBufferSize == 0 {
		config.LineBufferSize = 1000
	}
	if config.ErrorBufferSize == 0 {
		config.ErrorBufferSize = 10
	}

	t.Config = config

	t.lines = make(chan *Line, t.LineBufferSize)
	t.Errors = make(chan error, t.ErrorBufferSize)
	// 自动监听目录下符合文件名规则的文件
	go t.watchDirSync()
	return t, nil
}
