package spamalot

import (
	"time"
	"log"
	"github.com/CWarner818/giota"
)

type MetricType byte

const (
	INC_MILESTONE_BRANCH     MetricType = 0
	INC_MILESTONE_TRUNK      MetricType = 1
	INC_BAD_TRUNK            MetricType = 2
	INC_BAD_BRANCH           MetricType = 3
	INC_BAD_TRUNK_AND_BRANCH MetricType = 4
	INC_FAILED_TX            MetricType = 5
	INC_SUCCESSFUL_TX        MetricType = 6
	SUMMARY                  MetricType = 7
)

type Metric struct {
	Kind MetricType  `json:"Kind" bson:"Kind"`
	Data interface{} `json:"Data" bson:"Data"`
}

type Summary struct {
	TXsSucceeded      int     `json:"txs_succeeded"`
	TXsFailed         int     `json:"txs_failed"`
	BadBranch         int     `json:"bad_branch"`
	BadTrunk          int     `json:"bad_trunk"`
	BadTrunkAndBranch int     `json:"bad_trunk_and_branch"`
	MilestoneTrunk    int     `json:"milestone_trunk"`
	MilestoneBranch   int     `json:"milestone_branch"`
	TPS               float64 `json:"tps"`
	ErrorRate         float64 `json:"error_rate"`
}

type TXData struct {
	Hash  giota.Trytes `json:"hash"`
	Count int          `json:"count"`
}

type txandnode struct {
	tx   Transaction
	node Node
}

func newMetricsRouter() *metricsrouter {
	return &metricsrouter{
		metrics:    make(chan Metric),
		stopSignal: make(chan struct{}),
	}
}

type metricsrouter struct {
	metrics    chan Metric
	stopSignal chan struct{}
	relay      chan<- Metric

	startTime time.Time

	txsSucceeded, txsFailed, badBranch, badTrunk, badTrunkAndBranch int
	milestoneTrunk, milestoneBranch                                 int
}

func (mr *metricsrouter) stop() {
	mr.metrics = nil
	mr.stopSignal <- struct{}{}
}

func (mr *metricsrouter) addMetric(kind MetricType, data interface{}) {
	mr.metrics <- Metric{kind, data}
}

func (mr *metricsrouter) addRelay(relay chan<- Metric) {
	mr.relay = relay
}

func (mr *metricsrouter) collect() {
	mr.startTime = time.Now()
exit:
	for {
		select {
		case <-mr.stopSignal:
			break exit
		case metric := <-mr.metrics:
			switch metric.Kind {
			case INC_MILESTONE_BRANCH:
				mr.milestoneBranch++
			case INC_MILESTONE_TRUNK:
				mr.milestoneTrunk++
			case INC_BAD_TRUNK:
				mr.badTrunk++
			case INC_BAD_BRANCH:
				mr.badBranch++
			case INC_BAD_TRUNK_AND_BRANCH:
				mr.badTrunkAndBranch++
			case INC_FAILED_TX:
				mr.txsFailed++
			case INC_SUCCESSFUL_TX:
				mr.txsSucceeded++
				mr.printMetrics(metric.Data.(txandnode))
			}

			if mr.relay != nil && metric.Kind != INC_SUCCESSFUL_TX {
				mr.relay <- metric
			}
		}
	}
}

func (mr *metricsrouter) printMetrics(txAndNode txandnode) {
	tx := txAndNode.tx
	node := txAndNode.node
	var hash giota.Trytes
	if len(tx.Transactions) > 1 {
		hash = giota.Bundle(tx.Transactions).Hash()
		log.Println("Bundle sent to", node,
			"\nhttp://thetangle.org/bundle/"+hash)
	} else {
		hash = tx.Transactions[0].Hash()
		log.Println("Txn sent to", node,
			"\nhttp://thetangle.org/transaction/"+hash)
	}

	// TPS = delta since startup / successful TXs
	dur := time.Since(mr.startTime)
	tps := float64(mr.txsSucceeded) / dur.Seconds()

	// success rate = successful TXs / successful TXs + failed TXs
	successRate := 100 * (float64(mr.txsSucceeded) / (float64(mr.txsSucceeded) + float64(mr.txsFailed)))
	log.Printf("%.2f TPS -- success rate %.0f%% ", tps, successRate)

	log.Printf("Duration: %s Count: %d Milestone Trunk: %d Milestone Branch: %d Bad Trunk: %d Bad Branch: %d Both: %d",
		dur.String(), mr.txsSucceeded, mr.milestoneTrunk,
		mr.milestoneBranch, mr.badTrunk, mr.badBranch, mr.badTrunkAndBranch)

	// send current state of the spammer
	if mr.relay != nil {
		summary := Summary{
			TXsSucceeded:   mr.txsSucceeded, TXsFailed: mr.txsFailed,
			BadBranch:      mr.badBranch, BadTrunk: mr.badBranch, BadTrunkAndBranch: mr.badTrunkAndBranch,
			MilestoneTrunk: mr.milestoneTrunk, MilestoneBranch: mr.milestoneBranch,
			TPS:            tps, ErrorRate: 100 - successRate,
		}
		mr.relay <- Metric{Kind: SUMMARY, Data: summary}

		// send tx
		txData := TXData{Hash: hash, Count: len(tx.Transactions)}
		mr.relay <- Metric{Kind: INC_SUCCESSFUL_TX, Data: txData}
	}

}
