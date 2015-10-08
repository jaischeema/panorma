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
          <img src={`/photos/${photo.id}/large`} />
          <div className="caption">
            <h5>{photo.name}</h5>
            <p>{photo.height}</p>
            <p>{photo.width}</p>
            <p>{photo.ext}</p>
          </div>
        </div>
      </div>
    )
  }
}
