import React from 'react';
import { Router, Route, Link } from 'react-router';

import App          from './components/app';
import MediaIndex   from './components/media_index';
import Media        from './components/media';
import MediaByYear  from './components/media_by_year';
import MediaByMonth from './components/media_by_month';
import MediaByDay   from './components/media_by_day';

React.render(
  <Router>
    <Route path="/" component={App}>
      <Route path="/media" component={MediaIndex} />
      <Route path="/media/:media_id" component={Media} />
      <Route path="/from/:year" component={MediaByYear} />
      <Route path="/from/:year/:month" component={MediaByMonth} />
      <Route path="/from/:year/:month/:day" component={MediaByDay} />
    </Route>
  </Router>,
  document.body
);
