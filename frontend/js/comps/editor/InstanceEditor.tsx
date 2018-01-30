import * as React from 'react';
import {observer} from "mobx-react";
import {InstanceStore} from "../../stores/InstanceStore";


interface Props {
    instanceStore: InstanceStore;
}

@observer
export class InterfaceEditor extends React.Component<Props, {}> {
    render() {
        return (
            <div>
                <h2>Instance Editor</h2>
            </div>
        );
    }
}