import React from 'react';
import MediaByInterval from './media_by_interval';

export default class MediaByYear extends MediaByInterval {
  fetchParams(props) {
    return {
      year: props.params.year,
    }
  }

  titleElements() {
    return <span>{this.props.params.year}</span>
  }
}
