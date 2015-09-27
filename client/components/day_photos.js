import React from 'react';
import IntervalPhotos from './interval_photos';
import moment from 'moment';

export default class DayPhotos extends IntervalPhotos {
  fetchParams(props) {
    return {
      year: props.params.year,
      month: props.params.month,
      day: props.params.day
    }
  }

  titleElements() {
    let titleDate = moment({
      month: this.props.params.month - 1,
      year:  this.props.params.year,
      day:   this.props.params.day
    })
    return <span>{titleDate.format('Do MMMM YYYY')}</span>
  }
}
