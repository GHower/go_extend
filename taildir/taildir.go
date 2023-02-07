package taildir

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"go_extend/taildir/matcher"
	"go_extend/taildir/ratelimiter"
	"go_extend/taildir/watch"
	"gopkg.in/tomb.v1"
	"os"
	"path"
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
	LineBufferSize  uint // 行缓冲大小
	ErrorBufferSize uint // 非致命错误缓冲大小
	FilenameMatcher matcher.Matcher
	RateLimiter     *ratelimiter.LeakyBucket
	// 日志记录器,若禁用日志记录，设置为DiscardingLogger
	Logger logger
	Pool   bool // 轮询模式
}

// TailD 对目录的tail
type TailD struct {
	Dirname string
	Lines   chan *Line
	Errors  chan error
	Config

	tails map[string]*Tail

	watch watch.FileWatcher
	dead  error // 致命错误，拉都拉不起来时，err不为nil
	lk    sync.RWMutex
	tomb  tomb.Tomb
}

// 1. 为符合条件的新增文件建立监听
// 2. 为已存在文件建立监听
func (d *TailD) watchDirSync() {
	//defer d.close()

	dir, err := os.ReadDir(d.Dirname)
	if err != nil {
		d.dead = errStop
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		d.dead = errStop
		return
	}
	d.watch = watcher
	// 首次建立目录监听，拉起已经存在且符合条件的文件tailF
	for _, entry := range dir {
		if !entry.IsDir() && d.FilenameMatcher.Match(entry.Name()) {
			_, err := d.TailF(entry.Name())
			if err != nil {
				// todo: 统一发送错误方法
				d.Errors <- err
				continue
			}
		}
	}
	// 触发create文件时，拉起新的tailF
	for {
		select {
		case <-d.tomb.Dying():
			d.dead = errStop
			return
		}
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
	//close(d.Lines)
	d.tomb.Done()
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

func (d *TailD) TailF(filename string) (*Tail, error) {
	pathF := path.Join(d.Dirname, filename)
	tail, err := TailF(pathF)
	if err != nil {
		return nil, err
	}
	d.tails[pathF] = tail
	return tail, nil
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

	t.Lines = make(chan *Line, t.LineBufferSize)
	t.Errors = make(chan error, t.ErrorBufferSize)
	// 自动监听目录下符合文件名规则的文件
	go t.watchDirSync()
	return t, nil
}
