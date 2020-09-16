package domain

import (
	"bytes"
	"strconv"
	"time"
)

type Producer struct {
	lastTime      time.Time
	lastTimeCount int
	id            int
}

func (p *Producer) GetNewMsgID(b *bytes.Buffer) {

	postfix := false
	if time.Now().Unix() == p.lastTime.Unix() {
		p.lastTimeCount++
		postfix = true
	} else {
		p.lastTimeCount = 0
		p.lastTime = time.Now()
	}

	//return strconv.FormatInt(p.lastTime.Unix(), 10) + "-" + strconv.Itoa(p.id) + "-" + strconv.Itoa(p.lastTimeCount)
	//return fmt.Sprintf("%d-%d-%d", p.lastTime.Unix(), p.id, p.lastTimeCount)
	b.WriteString(strconv.FormatInt(p.lastTime.Unix(), 16))
	b.WriteString(".")
	b.WriteString(strconv.Itoa(p.id))
	if postfix {
		b.WriteString(".")
		b.WriteString(strconv.Itoa(p.lastTimeCount))
	}
	return
}
