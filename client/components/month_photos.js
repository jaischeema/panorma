import React from 'react';
import IntervalPhotos from './interval_photos';
import moment from 'moment';

export default class MonthPhotos extends IntervalPhotos {
  fetchParams(props) {
    return {
      year: props.params.year,
      month: props.params.month,
    }
  }

  titleElements() {
    let titleDate = moment({
      month: this.props.params.month - 1,
      year: this.props.params.year
    })
    return <span>{titleDate.format('MMMM YYYY')}</span>
  }
}
