package okx

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
	"strings"
)

// WSOrderBookStore 用于在本地维护 OKX WS 深度数据（合并 snapshot/update，校验 seqId/prevSeqId 与 checksum）。
//
// 说明：
// - books/books-elp/books-l2-tbt/books50-l2-tbt：首次为 snapshot，之后为增量 update（需要合并）
// - books5/bbo-tbt：定量推送（推荐按“全量替换”处理，以避免残留旧档位）
//
// 并发：该结构体非并发安全；请在单一 goroutine 中串行调用 Apply/ApplyMessage/Reset。
// 若跨 goroutine 读取 Snapshot/Ready，请由调用方自行加锁或做串行化。
type WSOrderBookStore struct {
	channel string
	instId  string
	sprdId  string

	verifySequence bool
	verifyChecksum bool

	ready bool

	asks []OrderBookLevel
	bids []OrderBookLevel

	ts       int64
	seqId    int64
	checksum int64
}

type WSOrderBookStoreOption func(*WSOrderBookStore)

// WithWSOrderBookVerifySequence 控制是否校验 prevSeqId 与本地 seqId 连续性（默认开启）。
func WithWSOrderBookVerifySequence(enable bool) WSOrderBookStoreOption {
	return func(s *WSOrderBookStore) {
		s.verifySequence = enable
	}
}

// WithWSOrderBookVerifyChecksum 控制是否校验 checksum（默认开启）。
func WithWSOrderBookVerifyChecksum(enable bool) WSOrderBookStoreOption {
	return func(s *WSOrderBookStore) {
		s.verifyChecksum = enable
	}
}

// NewWSOrderBookStore 创建一个用于指定频道/产品的本地深度合并器。
func NewWSOrderBookStore(channel, instId string, opts ...WSOrderBookStoreOption) *WSOrderBookStore {
	s := &WSOrderBookStore{
		channel:        channel,
		instId:         instId,
		verifySequence: true,
		verifyChecksum: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// NewWSSprdOrderBookStore 创建一个用于指定频道/Spread 的本地深度合并器。
func NewWSSprdOrderBookStore(channel, sprdId string, opts ...WSOrderBookStoreOption) *WSOrderBookStore {
	s := &WSOrderBookStore{
		channel:        channel,
		sprdId:         sprdId,
		verifySequence: true,
		verifyChecksum: true,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WSOrderBookSnapshot 是本地维护的深度快照（用于读取当前状态）。
type WSOrderBookSnapshot struct {
	Channel  string
	InstId   string
	SprdId   string
	TS       int64
	SeqId    int64
	Checksum int64

	Asks []OrderBookLevel
	Bids []OrderBookLevel
}

// Ready 表示是否已接收并应用过至少一条 snapshot（或 books5/bbo-tbt 的任意推送）。
func (s *WSOrderBookStore) Ready() bool {
	if s == nil {
		return false
	}
	return s.ready
}

// Reset 清空本地状态。
func (s *WSOrderBookStore) Reset() {
	if s == nil {
		return
	}
	s.ready = false
	s.asks = nil
	s.bids = nil
	s.ts = 0
	s.seqId = 0
	s.checksum = 0
}

// Snapshot 返回当前深度快照（深拷贝 asks/bids，便于调用方安全使用）。
func (s *WSOrderBookStore) Snapshot() WSOrderBookSnapshot {
	if s == nil {
		return WSOrderBookSnapshot{}
	}

	out := WSOrderBookSnapshot{
		Channel:  s.channel,
		InstId:   s.instId,
		SprdId:   s.sprdId,
		TS:       s.ts,
		SeqId:    s.seqId,
		Checksum: s.checksum,
	}
	if len(s.asks) > 0 {
		out.Asks = append([]OrderBookLevel(nil), s.asks...)
	}
	if len(s.bids) > 0 {
		out.Bids = append([]OrderBookLevel(nil), s.bids...)
	}
	return out
}

// ApplyMessage 尝试解析并应用一条 WS 原始消息；ok=false 表示不是当前 store 关心的深度消息。
func (s *WSOrderBookStore) ApplyMessage(message []byte) (ok bool, err error) {
	dm, ok, err := WSParseOrderBook(message)
	if err != nil || !ok {
		return ok, err
	}
	if s == nil {
		return true, errors.New("okx: nil ws order book store")
	}
	if s.channel != "" && dm.Arg.Channel != s.channel {
		return false, nil
	}
	if s.instId != "" && dm.Arg.InstId != "" && dm.Arg.InstId != s.instId {
		return false, nil
	}
	if s.sprdId != "" && dm.Arg.SprdId != "" && dm.Arg.SprdId != s.sprdId {
		return false, nil
	}
	return true, s.Apply(dm)
}

// Apply 应用一条已解析的深度推送（通常来自 WSParseOrderBook）。
func (s *WSOrderBookStore) Apply(dm *WSData[WSOrderBook]) error {
	if s == nil {
		return errors.New("okx: nil ws order book store")
	}
	if dm == nil {
		return errors.New("okx: nil ws order book data")
	}
	if dm.Arg.Channel == "" {
		return errors.New("okx: ws order book missing channel")
	}
	if s.channel != "" && dm.Arg.Channel != s.channel {
		return &WSOrderBookChannelMismatchError{Got: dm.Arg.Channel, Want: s.channel}
	}
	if s.channel == "" {
		s.channel = dm.Arg.Channel
	}
	if s.instId != "" && dm.Arg.InstId != "" && dm.Arg.InstId != s.instId {
		return &WSOrderBookInstIdMismatchError{Channel: dm.Arg.Channel, Got: dm.Arg.InstId, Want: s.instId}
	}
	if s.instId == "" {
		s.instId = dm.Arg.InstId
	}
	if s.sprdId != "" && dm.Arg.SprdId != "" && dm.Arg.SprdId != s.sprdId {
		return &WSOrderBookSprdIdMismatchError{Channel: dm.Arg.Channel, Got: dm.Arg.SprdId, Want: s.sprdId}
	}
	if s.sprdId == "" {
		s.sprdId = dm.Arg.SprdId
	}
	if !isOrderBookChannel(dm.Arg.Channel) {
		return fmt.Errorf("okx: ws order book invalid channel %q", dm.Arg.Channel)
	}
	if len(dm.Data) == 0 {
		return errors.New("okx: ws order book empty data")
	}
	if len(dm.Data) != 1 {
		return fmt.Errorf("okx: ws order book expect 1 data item, got %d", len(dm.Data))
	}

	action := dm.Action
	if action == "" {
		if s.ready {
			action = "update"
		} else {
			action = "snapshot"
		}
	}

	upd := dm.Data[0]

	// books5/bbo-tbt 为定量推送，按全量替换处理。
	if isOrderBookFullRefreshChannel(dm.Arg.Channel) {
		if err := s.applySnapshot(upd); err != nil {
			return err
		}
		s.ready = true
		return nil
	}

	switch action {
	case "snapshot":
		if err := s.applySnapshot(upd); err != nil {
			return err
		}
		s.ready = true
		return nil
	case "update":
		if !s.ready {
			return &WSOrderBookNotReadyError{Channel: dm.Arg.Channel, InstId: s.instId, SprdId: s.sprdId}
		}
		if err := s.applyUpdate(dm.Arg.Channel, upd); err != nil {
			s.Reset()
			return err
		}
		return nil
	default:
		return fmt.Errorf("okx: ws order book unknown action %q", action)
	}
}

func (s *WSOrderBookStore) applySnapshot(upd WSOrderBook) error {
	s.asks = append([]OrderBookLevel(nil), upd.Asks...)
	s.bids = append([]OrderBookLevel(nil), upd.Bids...)

	sortOrderBookLevels(s.asks, false)
	sortOrderBookLevels(s.bids, true)

	s.ts = upd.TS
	s.seqId = upd.SeqId

	if s.verifyChecksum {
		expected := wsOrderBookChecksum(s.bids, s.asks)
		if expected != upd.Checksum {
			return &WSOrderBookChecksumError{
				Channel:   s.channel,
				InstId:    s.instId,
				SprdId:    s.sprdId,
				Expected:  expected,
				Got:       upd.Checksum,
				SeqId:     upd.SeqId,
				ChecksumS: wsOrderBookChecksumString(s.bids, s.asks),
			}
		}
		s.checksum = expected
	} else {
		s.checksum = upd.Checksum
	}
	return nil
}

func (s *WSOrderBookStore) applyUpdate(channel string, upd WSOrderBook) error {
	if s.verifySequence && isOrderBookSequencedChannel(channel) {
		// seq 相关字段缺失时，均为默认 0；此时跳过校验以避免误报。
		if !(s.seqId == 0 && upd.PrevSeqId == 0 && upd.SeqId == 0) && upd.PrevSeqId != s.seqId {
			return &WSOrderBookSequenceError{
				Channel:           channel,
				InstId:            s.instId,
				SprdId:            s.sprdId,
				ExpectedPrevSeqId: s.seqId,
				GotPrevSeqId:      upd.PrevSeqId,
				SeqId:             upd.SeqId,
			}
		}
	}

	s.bids = applyOrderBookDelta(s.bids, upd.Bids, true)
	s.asks = applyOrderBookDelta(s.asks, upd.Asks, false)

	s.ts = upd.TS
	s.seqId = upd.SeqId

	if s.verifyChecksum {
		expected := wsOrderBookChecksum(s.bids, s.asks)
		if expected != upd.Checksum {
			return &WSOrderBookChecksumError{
				Channel:   channel,
				InstId:    s.instId,
				SprdId:    s.sprdId,
				Expected:  expected,
				Got:       upd.Checksum,
				SeqId:     upd.SeqId,
				ChecksumS: wsOrderBookChecksumString(s.bids, s.asks),
			}
		}
		s.checksum = expected
	} else {
		s.checksum = upd.Checksum
	}
	return nil
}

// WSOrderBookNotReadyError 表示在未接收 snapshot 的情况下收到 update。
type WSOrderBookNotReadyError struct {
	Channel string
	InstId  string
	SprdId  string
}

func (e *WSOrderBookNotReadyError) Error() string {
	idName := "instId"
	id := e.InstId
	if e.SprdId != "" {
		idName = "sprdId"
		id = e.SprdId
	}
	if id == "" {
		return fmt.Sprintf("okx: ws order book not ready channel=%s", e.Channel)
	}
	return fmt.Sprintf("okx: ws order book not ready channel=%s %s=%s", e.Channel, idName, id)
}

type WSOrderBookChannelMismatchError struct {
	Got  string
	Want string
}

func (e *WSOrderBookChannelMismatchError) Error() string {
	return fmt.Sprintf("okx: ws order book channel mismatch got=%s want=%s", e.Got, e.Want)
}

type WSOrderBookInstIdMismatchError struct {
	Channel string
	Got     string
	Want    string
}

func (e *WSOrderBookInstIdMismatchError) Error() string {
	return fmt.Sprintf("okx: ws order book instId mismatch channel=%s got=%s want=%s", e.Channel, e.Got, e.Want)
}

type WSOrderBookSprdIdMismatchError struct {
	Channel string
	Got     string
	Want    string
}

func (e *WSOrderBookSprdIdMismatchError) Error() string {
	return fmt.Sprintf("okx: ws order book sprdId mismatch channel=%s got=%s want=%s", e.Channel, e.Got, e.Want)
}

// WSOrderBookSequenceError 表示 prevSeqId 与本地 seqId 不连续（通常需要重新订阅获取 snapshot）。
type WSOrderBookSequenceError struct {
	Channel           string
	InstId            string
	SprdId            string
	ExpectedPrevSeqId int64
	GotPrevSeqId      int64
	SeqId             int64
}

func (e *WSOrderBookSequenceError) Error() string {
	idName := "instId"
	id := e.InstId
	if e.SprdId != "" {
		idName = "sprdId"
		id = e.SprdId
	}
	return fmt.Sprintf("okx: ws order book sequence mismatch channel=%s %s=%s prevSeqId=%d want=%d seqId=%d", e.Channel, idName, id, e.GotPrevSeqId, e.ExpectedPrevSeqId, e.SeqId)
}

// WSOrderBookChecksumError 表示 checksum 校验失败（深度可能已不同步）。
type WSOrderBookChecksumError struct {
	Channel  string
	InstId   string
	SprdId   string
	Expected int64
	Got      int64
	SeqId    int64

	// ChecksumS 为用于计算 CRC 的字符串（便于快速排查）。
	ChecksumS string
}

func (e *WSOrderBookChecksumError) Error() string {
	idName := "instId"
	id := e.InstId
	if e.SprdId != "" {
		idName = "sprdId"
		id = e.SprdId
	}
	return fmt.Sprintf("okx: ws order book checksum mismatch channel=%s %s=%s expected=%d got=%d seqId=%d", e.Channel, idName, id, e.Expected, e.Got, e.SeqId)
}

func isOrderBookFullRefreshChannel(channel string) bool {
	switch channel {
	case WSChannelBooks5, WSChannelBboTbt, WSChannelSprdBooks5, WSChannelSprdBboTbt:
		return true
	default:
		return false
	}
}

func isOrderBookSequencedChannel(channel string) bool {
	switch channel {
	case WSChannelBooks, WSChannelBooksL2Tbt, WSChannelBooks50L2Tbt, WSChannelSprdBooksL2Tbt:
		return true
	default:
		return false
	}
}

func sortOrderBookLevels(levels []OrderBookLevel, bids bool) {
	if len(levels) < 2 {
		return
	}
	if bids {
		sort.Slice(levels, func(i, j int) bool {
			return compareDecimalString(levels[i].Px, levels[j].Px) > 0
		})
		return
	}
	sort.Slice(levels, func(i, j int) bool {
		return compareDecimalString(levels[i].Px, levels[j].Px) < 0
	})
}

func applyOrderBookDelta(levels []OrderBookLevel, updates []OrderBookLevel, bids bool) []OrderBookLevel {
	if len(updates) == 0 {
		return levels
	}

	for _, u := range updates {
		idx := searchOrderBookIndex(levels, u.Px, bids)
		if idx < len(levels) && compareDecimalString(levels[idx].Px, u.Px) == 0 {
			if u.Sz == "0" {
				levels = append(levels[:idx], levels[idx+1:]...)
				continue
			}
			levels[idx] = u
			continue
		}
		if u.Sz == "0" {
			continue
		}
		levels = append(levels, OrderBookLevel{})
		copy(levels[idx+1:], levels[idx:])
		levels[idx] = u
	}
	return levels
}

func searchOrderBookIndex(levels []OrderBookLevel, px string, bids bool) int {
	if bids {
		return sort.Search(len(levels), func(i int) bool {
			return compareDecimalString(levels[i].Px, px) <= 0
		})
	}
	return sort.Search(len(levels), func(i int) bool {
		return compareDecimalString(levels[i].Px, px) >= 0
	})
}

func wsOrderBookChecksum(bids, asks []OrderBookLevel) int64 {
	s := wsOrderBookChecksumString(bids, asks)
	sum := crc32.ChecksumIEEE([]byte(s))
	return int64(int32(sum))
}

func wsOrderBookChecksumString(bids, asks []OrderBookLevel) string {
	b := min(25, len(bids))
	a := min(25, len(asks))

	var sb strings.Builder
	sb.Grow((b + a) * 16)

	first := true
	for i := 0; i < b || i < a; i++ {
		if i < b {
			if !first {
				sb.WriteByte(':')
			}
			first = false
			sb.WriteString(bids[i].Px)
			sb.WriteByte(':')
			sb.WriteString(bids[i].Sz)
		}
		if i < a {
			if !first {
				sb.WriteByte(':')
			}
			first = false
			sb.WriteString(asks[i].Px)
			sb.WriteByte(':')
			sb.WriteString(asks[i].Sz)
		}
	}
	return sb.String()
}

func compareDecimalString(a, b string) int {
	aNeg, aInt, aFrac := normalizeDecimalParts(a)
	bNeg, bInt, bFrac := normalizeDecimalParts(b)

	if aNeg != bNeg {
		if aNeg {
			return -1
		}
		return 1
	}
	sign := 1
	if aNeg {
		sign = -1
	}

	if len(aInt) != len(bInt) {
		if len(aInt) < len(bInt) {
			return -1 * sign
		}
		return 1 * sign
	}
	if aInt != bInt {
		if aInt < bInt {
			return -1 * sign
		}
		return 1 * sign
	}

	maxFrac := len(aFrac)
	if len(bFrac) > maxFrac {
		maxFrac = len(bFrac)
	}
	if maxFrac == 0 {
		return 0
	}

	aFrac = aFrac + strings.Repeat("0", maxFrac-len(aFrac))
	bFrac = bFrac + strings.Repeat("0", maxFrac-len(bFrac))
	if aFrac == bFrac {
		return 0
	}
	if aFrac < bFrac {
		return -1 * sign
	}
	return 1 * sign
}

func normalizeDecimalParts(s string) (neg bool, intPart string, fracPart string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return false, "0", ""
	}
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	intPart = s
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		intPart = s[:dot]
		fracPart = s[dot+1:]
	}

	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}
	fracPart = strings.TrimRight(fracPart, "0")
	return neg, intPart, fracPart
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
