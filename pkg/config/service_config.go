package config

import (
	"errors"
	"flag"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
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
	PgDsn     string
	RedisPass string
	RedisAddr string
	ApiAddr   string
	BadWords  []string

	LoggerTag string
	logger    *logger.Logger
	Db        *sqlx.DB
	Redis     *redis.Client
}

func (c *Config) ReadConfig() error {
	dsnPg := flag.String("pg", defaultPgDsn, "e.g.: postgres://lgn:pwd@localhost:5432/app")
	redis := flag.String("redis", defaultRedisDsn, "e.g.: pwd@localhost:5432")
	apiAddr := flag.String("api", defaultApiAddress, "e.g.: :8888")
	words := flag.String("words", defaultBadWords, "bad words e.g.: "+defaultBadWords)
	logLevel := flag.String("ll", defaultLogLevel, "log level e.g.: error, warning, info, debug, trace")
	logBackend := flag.String("lb", defaultLogBackend, "log backend e.g.: console, syslog")
	flag.Parse()

	c.PgDsn = *dsnPg
	c.ApiAddr = *apiAddr

	var redisDsn = regexp.MustCompile(`(.+)@(([\w|\.|\d]+):\d+)`)
	if !redisDsn.MatchString(*redis) {
		return errors.New("wrong Redis dsn")
	}
	res := redisDsn.FindAllSubmatch([]byte(*redis), -1)
	c.RedisPass = string(res[0][1])
	c.RedisAddr = string(res[0][2])

	var reWord = regexp.MustCompile(`\W+`)
	for _, word := range reWord.Split(*words, -1) {
		c.BadWords = append(c.BadWords, strings.ToLower(word))
	}
	return c.initLogger(*logLevel, *logBackend)
}

func (c *Config) ConnectDb() (err error) {
	c.Db, err = sqlx.Connect("postgres", c.PgDsn)
	return
}

func (c *Config) ConnectMq() error {
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPass,
		DB:       0,
	})
	_, err := c.Redis.Ping().Result()
	return err
}

func (c *Config) GetBadWords() []string {
	return c.BadWords
}

func (c *Config) initLogger(logLevel, logBackend string) (err error) {
	if c.logger == nil {
		c.logger, err = logger.NewLogger(logger.LookupLogLevel(logLevel), c.LoggerTag, logger.LookupLogBackend(logBackend))
	}
	return
}

func (c *Config) GetLogger() *logger.Logger {
	return c.logger
}
