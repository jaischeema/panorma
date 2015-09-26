import React   from 'react';
import Sidebar from './sidebar';

require('../styles/app.scss');

export default class App extends React.Component {
  render() {
    return (
      <div className="container">
        <Sidebar />
        <div className="content">
          {this.props.children}
        </div>
      </div>
    );
  }
}
