import React from 'react';
import { Router, Route, Link } from 'react-router';

import App         from './components/app';
import Photos      from './components/photos';
import Photo       from './components/photo';
import YearPhotos  from './components/year_photos';
import MonthPhotos from './components/month_photos';
import DayPhotos   from './components/day_photos';

React.render(
  <Router>
    <Route path="/" component={App}>
      <Route path="/photos" component={Photos} />
      <Route path="/photo/:photo_id" component={Photo} />
      <Route path="/from/:year" component={YearPhotos} />
      <Route path="/from/:year/:month" component={MonthPhotos} />
      <Route path="/from/:year/:month/:day" component={DayPhotos} />
    </Route>
  </Router>,
  document.body
);
