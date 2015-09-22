import React from 'react';
import { Router, Route, Link } from 'react-router';

import App    from './components/app';
import Photos from './components/photos';
import Photo  from './components/photo';

React.render(
  <Router>
    <Route path="/" component={App}>
      <Route path="/photos" component={Photos}>
        <Route path="/photo/:photo_id" component={Photo} />
      </Route>
    </Route>
  </Router>,
  document.body
);
