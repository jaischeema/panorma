import React from 'react';
import MediaByInterval from './media_by_interval';
import moment from 'moment';

export default class MediaByDay extends MediaByInterval {
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
