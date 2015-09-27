import React from 'react';
import IntervalPhotos from './interval_photos';

export default class YearPhotos extends IntervalPhotos {
  fetchParams(props) {
    return {
      year: props.params.year,
    }
  }

  titleElements() {
    return <span>{this.props.params.year}</span>
  }
}
