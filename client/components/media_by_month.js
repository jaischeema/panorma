import React from 'react';
import MediaByInterval from './media_by_interval';
import moment from 'moment';

export default class MediaByMonth extends MediaByInterval {
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
