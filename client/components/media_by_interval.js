import React     from 'react';
import MediaList from './media_list';
import API from '../api';

export default class MediaByInterval extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      page:    1,
      media:  [],
      loading: true
    };
    this.loadMedia(props);
  }

  componentWillReceiveProps(props) {
    this.setState({loading: true});
    this.loadMedia(props);
  }

  loadMedia(props) {
    let api = new API;
    let fetchParams = this.fetchParams(props);
    api.fetchMedia(fetchParams, (response) => {
      this.setState({media: response.body, loading: false})
    });
  }

  fetchParams(props) {}
  titleElements() {}

  render() {
    return (
      <div className="container">
        <div className="page-header">
          <h1>{this.titleElements()}</h1>
        </div>
        <div>
          {this.renderContent()}
        </div>
      </div>
    )
  }

  renderContent() {
    if(this.state.loading) {
      return <span>Loading...</span>
    } else {
      return <MediaList media={this.state.media} />
    }
  }
}
