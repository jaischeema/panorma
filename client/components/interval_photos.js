import React     from 'react';
import PhotoList from './photo_list';
import API from '../api';

export default class IntervalPhotos extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      page:    1,
      photos:  [],
      loading: true
    };
    this.loadPhotos(props);
  }

  componentWillReceiveProps(props) {
    this.setState({loading: true});
    this.loadPhotos(props);
  }

  loadPhotos(props) {
    let api = new API;
    let fetchParams = this.fetchParams(props);
    api.fetchPhotos(fetchParams, (response) => {
      this.setState({photos: response.body, loading: false})
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
      return <PhotoList photos={this.state.photos} />
    }
  }
}
