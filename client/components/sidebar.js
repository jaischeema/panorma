import React        from 'react';
import { Link }     from 'react-router';
import SidebarGroup from './sidebar_group';
import SidebarItem  from '../sidebar_item';

require('../styles/sidebar.scss');

export default class Sidebar extends React.Component {
  constructor() {
    super();
    this.state = {years: [], openYear: null};
    SidebarItem.loadItems((sidebarItems) => {
      this.setState({years: sidebarItems});
    });
  }

  render() {
    return (
      <div className="sidebar">
        <div className="logo">
          <Link to={"/"}>Panorma</Link>
        </div>
        <div className="menu">
          <h4>Photos</h4>
          <ul>
            {this.state.years.map((year) => { return this.renderSidebarItem(year); })}
          </ul>
        </div>
      </div>
    );
  }

  renderSidebarItem(year) {
    return (
      <SidebarGroup
        {...year}
        open={year.label == this.state.openYear}
        onSidebarLinkClick={this.onSidebarLinkClick.bind(this, year)}
        key={year.url}
      />
    );
  }

  onSidebarLinkClick(year) {
    this.setState({openYear: year.label});
  }
}
