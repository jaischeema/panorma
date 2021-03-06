import request from 'superagent';

class API {
  sidebarItems(callback) {
    request
      .get('/all_dates')
      .end((err, response) => {
        callback(response);
      })
  }

  fetchMedia(params, callback) {
    request
      .get('/media')
      .query(params)
      .end((err, response) => {
        callback(response);
      })
  }
}

export default API;
