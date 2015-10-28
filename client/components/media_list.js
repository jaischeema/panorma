import React from 'react';

export default class MediaList extends React.Component {
  render() {
    return (
      <div className="row">
        {this.props.media.map((media) => { return this.renderMedia(media); })}
      </div>
    );
  }

  renderMedia(media) {
    return (
      <div key={`media-${media.id}`} className="col-xs-6 col-md-4">
        <div className="thumbnail">
          <img src={`/media/${media.id}/medium`} />
          <div className="caption">
            <h5>{media.name}</h5>
            <p>{media.height}</p>
            <p>{media.width}</p>
            <p>{media.ext}</p>
          </div>
        </div>
      </div>
    )
  }
}
