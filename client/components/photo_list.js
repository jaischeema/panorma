import React from 'react';

export default class PhotoList extends React.Component {
  render() {
    return (
      <div className="row">
        {this.props.photos.map((photo) => { return this.renderPhoto(photo); })}
      </div>
    );
  }

  renderPhoto(photo) {
    return (
      <div key={`photo-${photo.id}`} className="col-xs-6 col-md-3">
        <div className="thumbnail">
          <div className="caption">
            <h3>{photo.path}</h3>
            <p>{photo.hash_value}</p>
          </div>
        </div>
      </div>
    )
  }
}
