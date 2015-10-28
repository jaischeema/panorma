import React from 'react';

export default class MediaIndex extends React.Component {
  render() {
    return (
      <div className="content">
        <h2>Media</h2>
        {this.props.children}
      </div>
    );
  }
}
