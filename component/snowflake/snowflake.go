package snowflake

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

/*
   雪花算法(snowFlake)的具体实现方案:
*/

type SnowFlake struct {
	mu sync.Mutex

	//每一部分最大的数值
	maxWorkerId int64
	maxSequence int64

	//当前机器的ID号
	workerId int64
	//序列号
	sequence int64
	//上一次生成ID号前41位的毫秒时间戳
	lastTimestamp int64
}

func NewSnowFlake(workerId int64) (*SnowFlake, error) {
	mySnow := new(SnowFlake)
	if workerId >= mySnow.maxWorkerId {
		mySnow.workerId = mySnow.maxWorkerId % workerId
	} else {
		mySnow.workerId = workerId
	}
	mySnow.lastTimestamp = mySnow.timeGen()
	return mySnow, nil
}

func (s *SnowFlake) NextId() (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	nowTimestamp := s.timeGen() //获取当前的毫秒级别的时间戳
	if nowTimestamp < s.lastTimestamp {
		//系统时钟倒退,倒退了s.lastTimestamp-nowTimestamp
		return -1, errors.New(fmt.Sprintf("clock moved backwards, Refusing to generate id for %d milliseconds", s.lastTimestamp-nowTimestamp))
	}
	if nowTimestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) % s.maxSequence
		if s.sequence == 0 {
			//tilNextMills中有一个循环等候当前毫秒时间戳到达lastTimestamp的下一个毫秒时间戳
			nowTimestamp = s.tilNextMills()
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = nowTimestamp
	nextId := fmt.Sprintf("%d%d%d%d", nowTimestamp, s.workerId, s.sequence)
	return strconv.ParseUint(nextId, 0, 0)
}

/*
   获取毫秒的时间戳
*/
func (s *SnowFlake) timeGen() int64 {
	return time.Now().UnixMilli()
}

/*
   获取比lastTimestamp大的当前毫秒时间戳
*/
func (s *SnowFlake) tilNextMills() int64 {
	timeStampMill := s.timeGen()
	for timeStampMill <= s.lastTimestamp {
		timeStampMill = s.timeGen()
	}
	return timeStampMill
}
