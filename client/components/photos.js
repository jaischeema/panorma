import React from 'react';

export default class Photos extends React.Component {
  render() {
    return (
      <div className="content">
        <h2>Photos</h2>
        {this.props.children}
      </div>
    );
  }
}
