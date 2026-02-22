package okx

func validTradingBotAlgoIDAck(ack *TradingBotAlgoIdAck) bool {
	return ack != nil && ack.AlgoId != ""
}

func validTradingBotSignalCreateAck(ack *TradingBotSignalCreateAck) bool {
	return ack != nil && ack.SignalChanId != "" && ack.SignalChanToken != ""
}

func validTradingBotGridCloseOrderAck(ack *TradingBotGridCloseOrderAck) bool {
	return ack != nil && ack.AlgoId != "" && ack.AlgoClOrdId != "" && ack.OrdId != ""
}

func validTradingBotGridComputeMarginBalanceResult(result *TradingBotGridComputeMarginBalanceResult) bool {
	return result != nil && result.MaxAmt != "" && result.Lever != ""
}

func validTradingBotGridWithdrawIncomeAck(ack *TradingBotGridWithdrawIncomeAck) bool {
	return ack != nil && ack.AlgoId != "" && ack.AlgoClOrdId != "" && ack.Profit != ""
}
