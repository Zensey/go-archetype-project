package app

import (
	"errors"
	"flag"
	"regexp"
	"strings"
	"time"
)

const (
	defaultPgDsn      = "postgres://docker:docker@localhost:5432/app"
	defaultRedisDsn   = "foobared@localhost:6379"
	defaultApiAddress = ":8888"

	QueueWorker1    = "worker#1"
	QueueWorker2    = "worker#2"
	QueueRcvTimeout = time.Second

	defaultBadWords   = "fee,nee,cruul,leent"
	defaultLogLevel   = "info"
	defaultLogBackend = "console"
)

type Config struct {
	PgDsn      string
	RedisPass  string
	RedisAddr  string
	ApiAddr    string
	BadWords   []string
	LogLevel   string
	LogBackend string
}

func GetConfig() (Config, error) {
	dsnPg := flag.String("pg", defaultPgDsn, "e.g.: postgres://lgn:pwd@localhost:5432/app")
	redis := flag.String("redis", defaultRedisDsn, "e.g.: pwd@localhost:5432")
	apiAddr := flag.String("api", defaultApiAddress, "e.g.: :8888")
	words := flag.String("words", defaultBadWords, "bad words e.g.: "+defaultBadWords)
	logLevel := flag.String("ll", defaultLogLevel, "log level e.g.: error, warning, info, debug, trace")
	logBackend := flag.String("lb", defaultLogBackend, "log backend e.g.: console, syslog")
	flag.Parse()

	c := Config{
		PgDsn:      *dsnPg,
		ApiAddr:    *apiAddr,
		LogLevel:   *logLevel,
		LogBackend: *logBackend,
	}
	var redisDsn = regexp.MustCompile(`(.+)@(([\w|\.|\d]+):\d+)`)
	if !redisDsn.MatchString(*redis) {
		return c, errors.New("wrong Redis dsn")
	}
	res := redisDsn.FindAllSubmatch([]byte(*redis), -1)
	c.RedisPass = string(res[0][1])
	c.RedisAddr = string(res[0][2])

	var reWord = regexp.MustCompile(`\W+`)
	for _, word := range reWord.Split(*words, -1) {
		c.BadWords = append(c.BadWords, strings.ToLower(word))
	}
	return c, nil
}
