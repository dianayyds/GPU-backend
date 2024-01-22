package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func TestCron() {
	c := cron.New()
	i := 1
	EntryID, err := c.AddFunc("*/1 * * * *", func() {
		fmt.Println(time.Now(), "每分钟执行一次", i)
		i++
	})
	fmt.Println(time.Now(), EntryID, err)

	c.Start()
	time.Sleep(time.Minute * 5)
}
