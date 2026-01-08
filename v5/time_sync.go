package okx

import (
	"context"
	"time"
)

// TimeSyncResult 表示一次时间同步的结果。
type TimeSyncResult struct {
	ServerTime time.Time
	RoundTrip  time.Duration
	Offset     time.Duration
}

// SyncTime 通过 GET /api/v5/public/time 同步服务器时间，并更新 Client 的 timeOffset。
//
// 使用近似 NTP 的方式：offset = localMid - serverTime，其中 localMid 为本地往返时间的中点。
func (c *Client) SyncTime(ctx context.Context) (TimeSyncResult, error) {
	t0 := c.now()
	st, err := c.NewPublicTimeService().Do(ctx)
	if err != nil {
		return TimeSyncResult{}, err
	}
	t1 := c.now()

	roundTrip := t1.Sub(t0)
	localMid := t0.Add(roundTrip / 2)
	serverTime := st.Time()
	offset := localMid.Sub(serverTime)

	c.timeOffsetNanos.Store(offset.Nanoseconds())

	return TimeSyncResult{
		ServerTime: serverTime,
		RoundTrip:  roundTrip,
		Offset:     offset,
	}, nil
}
