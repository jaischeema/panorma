import request from 'superagent';

class API {
  sidebarItems(callback) {
    request
      .get('/all_dates')
      .end((err, response) => {
        callback(response);
      })
  }

  fetchPhotos(params, callback) {
    request
      .get('/photos')
      .query(params)
      .end((err, response) => {
        callback(response);
      })
  }
}

export default API;
