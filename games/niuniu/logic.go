package niuniu

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	nncard "github.com/Ajenpan/chinese_poker_go/Niuniu"
	"github.com/sirupsen/logrus"
	protobuf "google.golang.org/protobuf/proto"

	pb "hotwave/games/niuniu/proto"
	"hotwave/services/battle"
)

func CreateLogic() battle.GameLogic {
	ret := &Logic{
		log:     logrus.New(),
		players: make(map[int32]*Player),
		info:    &pb.GameInfo{},
		conf:    &Config{},
	}
	return ret
}

type Player struct {
	raw battle.Player
	*pb.GamePlayerInfo
	rawHandCards *nncard.NNHandCards
}

type Config struct {
	Downtime time.Duration
}

func ParseConfig(raw []byte) (*Config, error) {
	ret := &Config{}
	err := json.Unmarshal(raw, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type Logic struct {
	desk    battle.GameDesk
	info    *pb.GameInfo
	conf    *Config
	log     *logrus.Logger
	players map[int32]*Player

	lifeTime  time.Duration
	stageTime time.Duration
}

func (nn *Logic) OnInit(d battle.GameDesk, conf interface{}) error {
	switch v := conf.(type) {
	case []byte:
		var err error
		nn.conf, err = ParseConfig(v)
		if err != nil {
			return err
		}
	case *Config:
		nn.conf = v
	default:
		return fmt.Errorf("unknow config type ")
	}

	nn.desk = d

	nn.info = &pb.GameInfo{
		Status: pb.GameStep_IDLE,
	}
	nn.lifeTime = 0
	return nil
}

func (nn *Logic) OnStart(players []battle.Player) error {
	if len(players) == 0 {
		return fmt.Errorf("player len is 0")
	}
	for _, v := range players {
		nn.addPlayer(v)
	}

	nn.ChangeLogicStatus(pb.GameStep_BEGIN)
	return nil
}

func (nn *Logic) OnMessage(p battle.Player, topic string, msg []byte) {
	nn.log.Infof("recv topic:%s", topic)
	//TODO:
	// msgDes := pb.File_games_niuniu_proto_niuniu_protobufMessages().ByName(protoreflect.Name(topic))
	// if msgDes != nil {
	// 	dynamicpb.NewMessage(msgDes)
	// }
	// New().Unmarshal(msg)
}

func (nn *Logic) OnEvent(topic string, event protobuf.Message) {

}

func (nn *Logic) GameDeskInfoRequest(p battle.Player, req *pb.GameDeskInfoRequest) {
	resp := &pb.GameDeskInfoResponse{
		Info: nn.info,
	}
	nn.desk.SendMessageToPlayer(p, resp)
}

func (nn *Logic) PlayerBankerRequest(player *Player, req *pb.PlayerBanker) {
	if nn.getLogicStatus() != pb.GameStep_BANKER {
		nn.log.Warnf("OnPlayerBankerRequest 游戏状态错误 logic status:%v", nn.getLogicStatus())
		return
	}

	if player.Status != pb.GameStep_BEGIN {
		nn.log.Warnf("OnPlayerBankerRequest 用户状态错误 status:%d", (player.Status))
		return
	}

	notice := &pb.PlayerBankerNotify{
		SeatId: int32(player.raw.GetSeatID()),
		Rob:    req.Rob,
	}
	player.BankerRob = req.Rob
	nn.desk.BroadcastMessage(notice)
}

//
func (nn *Logic) OnPlayerBetRateRequest(p battle.Player, pMsg *pb.PlayerBetRate) {
	nnPlayer := nn.playerConv(p)
	if nnPlayer == nil {
		nn.log.Infof("can't find player uid :%d", p.GetSeatID())
		return
	}

	if nnPlayer.Status != pb.GameStep_BANKER {
		nn.log.Warnf("OnPlayerBetRateRequest 用户状态错误 status:%d", nnPlayer.Status)
		return
	}

	if nn.getLogicStatus() != pb.GameStep_BET {
		nn.log.Warnf("OnPlayerBetRateRequest 游戏状态错误 logic status:%d", nn.getLogicStatus())
		return
	}

	nnPlayer.BetRate = pMsg.Rate
	nnPlayer.Status = pb.GameStep_BET

	notice := &pb.PlayerBetRateNotify{
		SeatId: int32(p.GetSeatID()),
		Rate:   pMsg.Rate,
	}
	nn.desk.BroadcastMessage(notice)
}

func (nn *Logic) OnPlayerOutCardRequest(p battle.Player, pMsg *pb.PlayerOutCard) {
	nnPlayer := nn.playerConv(p)

	if nnPlayer == nil {
		nn.log.Errorf("OnPlayerOutCardRequest player is nil")
		return
	}

	if nnPlayer.Status != pb.GameStep_BET {
		nn.log.Errorf("OnPlayerOutCardRequest 用户状态错误 %d", nnPlayer.Status)
		return
	}

	if nn.getLogicStatus() != pb.GameStep_SHOW_CARDS {
		nn.log.Errorf("OnPlayerOutCardRequest 游戏状态错误 logic status:%v", nn.getLogicStatus())
		return
	}

	nnPlayer.OutCard = &pb.OutCardInfo{
		Cards: nnPlayer.rawHandCards.Bytes(),
		Type:  pb.CardType(nnPlayer.rawHandCards.Type()),
	}

	nnPlayer.Status = pb.GameStep_SHOW_CARDS

	notice := &pb.PlayerOutCardNotify{
		SeatId:  int32(p.GetSeatID()),
		OutCard: nnPlayer.OutCard,
	}

	nn.desk.BroadcastMessage(notice)
}

func (nn *Logic) addPlayer(p battle.Player) *Player {
	ret := &Player{}
	ret.GamePlayerInfo = &pb.GamePlayerInfo{}

	ret.GamePlayerInfo.SeatId = int32(p.GetSeatID())
	ret.raw = p
	nn.players[p.GetSeatID()] = ret
	return ret
}

func (nn *Logic) OnTick(duration time.Duration) {
	nn.lifeTime += duration
	nn.stageTime += duration

	switch nn.getLogicStatus() {
	case pb.GameStep_UNKNOW:
		//do nothing
	case pb.GameStep_BEGIN:
		nn.ChangeLogicStatus(pb.GameStep_BANKER)
	case pb.GameStep_BANKER:
		if nn.checkEndBanker() {
			nn.ChangeLogicStatus(pb.GameStep_BANKER_NOTIFY)
		}
	case pb.GameStep_BANKER_NOTIFY:
		if nn.isStageTimeover() {
			nn.notifyRobBanker()
			nn.ChangeLogicStatus(pb.GameStep_BET)
		}
	case pb.GameStep_BET:
		nn.ChangeLogicStatus(pb.GameStep_DEAL_CARDS)
	case pb.GameStep_DEAL_CARDS:
		nn.sendCardToPlayer()
		nn.ChangeLogicStatus(pb.GameStep_SHOW_CARDS)
	case pb.GameStep_SHOW_CARDS:
		nn.ChangeLogicStatus(pb.GameStep_TALLY)
	case pb.GameStep_OVER:
		nn.desk.ReportGameOver()
	case pb.GameStep_TALLY:
	default:
		//warn
	}
}

func (nn *Logic) OnReset() {

}

func (nn *Logic) getLogicStatus() pb.GameStep {
	return nn.info.Status
}

func (nn *Logic) getStageDowntime(s pb.GameStep) time.Duration {
	//TODO:
	return nn.conf.Downtime
}

func (nn *Logic) ChangeLogicStatus(s pb.GameStep) {
	lastStatus := nn.getLogicStatus()
	nn.info.Status = s

	if lastStatus != s {
		//reset stage time
		nn.stageTime = 0
	}

	donwtime := nn.getStageDowntime(s).Seconds()

	nn.log.Infof("dest status change, before:%v ,now:%v ", lastStatus, s)

	if lastStatus == s {
		nn.log.Errorf("set same status before:%v ,now:%v", lastStatus, s)
	}

	if lastStatus != pb.GameStep_OVER {
		if lastStatus > s {
			nn.log.Errorf("last status is bigger than now before:%v, now:%v", lastStatus, s)
		}
	}

	notice := &pb.GameStatusNotify{
		GameStatus: s,
		TimeDown:   int32(donwtime),
	}

	nn.desk.BroadcastMessage(notice)
}

func (nn *Logic) playerConv(p battle.Player) *Player {
	return nn.getPlayerBySeatId(p.GetSeatID())
}

func (nn *Logic) getPlayerBySeatId(seatid int32) *Player {
	p, ok := nn.players[seatid]
	if ok {
		return p
	}
	return nil
}

func (nn *Logic) isStageTimeover() bool {
	return nn.stageTime >= nn.getStageDowntime(nn.info.Status)
}

func (nn *Logic) checkEndBanker() bool {
	if nn.isStageTimeover() {
		return true
	}

	for _, p := range nn.players {
		if p.BankerRob == 0 {
			return false
		}
	}
	return true
}

func (nn *Logic) notifyRobBanker() {
	for _, p := range nn.players {
		if p.Status != pb.GameStep_BANKER {
			p.Status = pb.GameStep_BANKER
		}
	}

	seats := []int32{}
	var maxRob int32 = -1
	for _, p := range nn.players {
		if (p.BankerRob) > maxRob {
			maxRob = p.BankerRob
			seats = seats[:0]
			seats = append(seats, p.SeatId)
		} else if (p.BankerRob) == maxRob {
			seats = append(seats, p.SeatId)
		}
	}

	if len(seats) == 0 {
		nn.log.Errorf("选庄错误 maxrob:%d", maxRob)
	}

	index := rand.Intn(len(seats))
	bankSeatId := seats[index]
	banker, ok := nn.players[int32(bankSeatId)]

	if !ok {
		nn.log.Errorf("banker seatid error. seatid:%d,index:%d", bankSeatId, index)
		return
	}

	banker.Banker = true
	//庄家不参与下注.提前设置好状态
	banker.Status = pb.GameStep_BET

	notice := &pb.BankerSeatNotify{
		SeatId: bankSeatId,
	}
	nn.desk.BroadcastMessage(notice)
}

func (nn *Logic) sendCardToPlayer() {
	deck := nncard.NewNNDeck()
	deck.Shuffle()

	for _, p := range nn.players {
		p.rawHandCards = deck.DealHandCards()
		p.HandCards = p.rawHandCards.Bytes()
		notice := &pb.PlayerHandCardsNotify{
			SeatId:    p.SeatId,
			HandCards: p.HandCards,
		}
		nn.desk.SendMessageToPlayer(p.raw, notice)
	}

	for _, p := range nn.players {
		p.rawHandCards.Calculate()
	}
}

func (nn *Logic) beginTally() {
	nn.ChangeLogicStatus(pb.GameStep_TALLY)

	var banker *Player = nil

	for _, p := range nn.players {
		if p.Banker {
			banker = p
			break
		}
	}
	if banker == nil {
		nn.log.Errorf("bank is nil")
		return
	}

	notify := &pb.PlayerTallyNotify{}
	// notify.TallInfo = make([]*pb.PlayerTallyNotify_TallyInfo, 0)
	// type tally struct {
	// 	UserId int64
	// 	Coins  int32
	// }

	bankerTally := &pb.PlayerTallyNotify_TallyInfo{
		SeatId: banker.SeatId,
		//Coins:  chips*cardRate*p.BetRate - 100,
	}

	for _, p := range nn.players {
		if p.Banker {
			continue
		}
		var chips int32 = 5
		var cardRate int32 = 1

		if banker.rawHandCards.Compare(p.rawHandCards) {
			//底注*倍率*牌型倍率
			cardRate += int32(banker.rawHandCards.Type())
			cardRate = -cardRate
		} else {
			cardRate += int32(p.rawHandCards.Type())
		}
		temp := &pb.PlayerTallyNotify_TallyInfo{
			SeatId: p.SeatId,
			Coins:  chips * cardRate * p.BetRate,
		}
		// notify.TallInfo = append(notify.TallInfo, temp)
		bankerTally.Coins += temp.Coins
	}

	// notify.TallInfo = append(notify.TallInfo, bankerTally)

	nn.desk.BroadcastMessage(notify)
}

func (nn *Logic) resetDesk() {
	nn.players = make(map[int32]*Player)
	for _, p := range nn.players {
		p.GamePlayerInfo.Reset()
		p.GamePlayerInfo.Status = pb.GameStep_IDLE
		p.GamePlayerInfo.SeatId = p.SeatId
	}
	nn.ChangeLogicStatus(pb.GameStep_IDLE)
}
