package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10                      // 节点数
	seqBits     uint8 = 12                      // 1 毫秒内可生成的 id 序号的二进制位数
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点 ID 的最大值，用于防止溢出
	seqMax      int64 = -1 ^ (-1 << seqBits)    // 同上，用来表示生成 id 序号的最大值
	timeShift   uint8 = workerBits + seqBits    // 时间戳向左的偏移量
	workerShift uint8 = seqBits                 // 节点 ID 向左的偏移量
	epoch       int64 = 1567906170596           // 开始运行时间
)

type Worker struct {
	// 添加互斥锁 确保并发安全
	mu sync.Mutex
	// 记录时间戳
	timestamp int64
	// 该节点的ID
	workerId int64
	// 当前毫秒已经生成的 id 序列号 (从 0 开始累加 ) 1 毫秒内最多生成 4096 个 ID
	seq int64
}

// 实例化对象
func NewWorker(workerId int64) (*Worker, error) {
	// 要先检测 workerId 是否在上面定义的范围内
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New(" Worker ID excess of quantity ")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		seq:       0,
	}, nil
}

// 获取一个新 ID
func (w *Worker) Next() int64 {
	// 获取 id 最关键的一点 加锁 加锁 加锁
	w.mu.Lock()
	defer w.mu.Unlock() // 生成完成后记得 解锁 解锁 解锁
	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.seq = (w.seq + 1) & seqMax
		// 这里要判断，当前工作节点是否在 1 毫秒内已经生成 seqMax 个 ID
		if w.seq == 0 {
			// 如果当前工作节点在 1 毫秒内生成的 ID 已经超过上限 需要等待 1 毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成 ID 的时间不一致 则需要重置工作节点生成 ID 的序号
		w.seq = 0
	}
	w.timestamp = now // 将机器上一次生成 ID 的时间更新为当前时间
	// 第一段 now - epoch 为该算法目前已经运行了多少毫秒
	// 如果在程序跑了一段时间修改了 epoch 这个值 可能会导致生成相同的 ID
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.seq))
	return ID
}