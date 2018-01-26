import * as React from 'react';
import {Link} from 'react-router-dom';

interface Props {

}

export class Nav extends React.Component<Props, {}> {
    render() {
        return (
            <div className={'nav'}>
                <div className={'site_title'}>spamalot-coo</div>
                <ul className={'nav_menu'}>
                    <Link to={'/instances'}>
                        <li>Instances</li>
                    </Link>
                    <Link to={'/configuration'}>
                        <li>Configuration</li>
                    </Link>
                </ul>
            </div>
        );
    }
}