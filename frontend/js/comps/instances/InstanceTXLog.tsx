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


interface Props {
    instanceStore?: InstanceStore;
    match?: { params: { id: string } };
}

@withRouter
@inject("instanceStore")
@observer
export class InstanceTXLog extends React.Component<Props, {}> {
    render() {
        let id = this.props.match.params.id;
        let instance = this.props.instanceStore.instances.get(id);
        return (
            <div>
                <h3>Transactions</h3>
                <Divider/>
                <br/>
                <div className={'tx_log'}>
                    <span className='log_entry'>TX: BLABLABLABLABLA</span>
                    <span className='log_entry'>TX: BLABLABLABLABLA</span>
                    <span className='log_entry'>TX: BLABLABLABLABLA</span>
                    <span className='log_entry'>TX: BLABLABLABLABLA</span>
                    <span className='log_entry'>TX: BLABLABLABLABLA</span>
                </div>
            </div>
        );
    }
}