import React from 'react';

export default class App extends React.Component {
  render() {
    return (
      <div className="container">
        <h1>Panorma</h1>
        {this.props.children}
      </div>
    );
  }
}
