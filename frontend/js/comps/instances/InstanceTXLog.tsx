import * as React from 'react';
import {inject, observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";
import {withRouter} from "react-router";
import {
    AreaChart, Area, ReferenceLine,
    LineChart, ComposedChart, Brush, XAxis, Line, YAxis,
    CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts';
import Divider from "material-ui/Divider";
import {MetricsStore} from "../../stores/MetricsStore";
import {TXData} from "../../models/Instance";


interface Props {
    instanceStore?: InstanceStore;
    metricsStore?: MetricsStore;
    match?: { params: { id: string } };
}

@withRouter
@inject("metricsStore")
@inject("instanceStore")
@observer
export class InstanceTXLog extends React.Component<Props, {}> {
    render() {
        let id = this.props.match.params.id;
        let instance = this.props.instanceStore.instances.get(id);
        let txs = this.props.metricsStore.transactions;
        let entries = [];
        txs.forEach(tx => {
            entries.push(<TX key={tx.hash} tx={tx}></TX>)
        });
        return (
            <div>
                <h3>Transactions ({txs.length})</h3>
                <Divider/>
                <br/>
                <div className={'tx_log'}>
                    {entries}
                </div>
            </div>
        );
    }
}

/*
		hash = giota.Bundle(tx.Transactions).Hash()
		log.Println("Bundle sent to", node,
			"\nhttp://thetangle.org/bundle/"+hash)
	} else {
		hash = tx.Transactions[0].Hash()
		log.Println("Txn sent to", node,
			"\nhttp://thetangle.org/transaction/"+hash)
 */

class TX extends React.Component<{ tx: TXData }, {}> {
    render() {
        let tx = this.props.tx;
        return (
            <span className={'log_entry'}>
              TX
                <a href={`https://thetangle.org/transaction/${tx.hash}`} target="_blank">
                    <i className="fas fa-external-link-alt log_link"></i>
                </a>
              |{' '}{tx.hash}
          </span>
        );
    }
}