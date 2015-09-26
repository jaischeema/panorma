import _      from 'lodash';
import moment from 'moment';
import API    from './api';

export default class SidebarItem {
  constructor(label, url) {
    this.label    = label;
    this.url      = url;
    this.children = []
  }

  static loadItems(callback) {
    let api = new API;
    api.sidebarItems((response) => {

      let sidebarItems = [];

      let itemsGroupedByYears = _.groupBy(response.body.dates, function(item) {
        return item.year
      });

      for(var year of Object.keys(itemsGroupedByYears)) {
        let sidebarItem = new SidebarItem(year, `/from/${year}`)

        let itemsGroupedByMonths = _.groupBy(itemsGroupedByYears[year], function(item) {
          return item.month
        });

        for(var month of Object.keys(itemsGroupedByMonths)) {
          let monthDate = moment({month: month - 1});
          let monthLabel = monthDate.format('MMMM');
          let monthSidebarItem = new SidebarItem(monthLabel, `/from/${year}/${month}`);

          monthSidebarItem.children = _.map(itemsGroupedByMonths[month], function(day) {
            return new SidebarItem(day.day, `/from/${year}/${month}/${day.day}`);
          })

          sidebarItem.children.push(monthSidebarItem);
        }

        sidebarItems.push(sidebarItem);
      }
      callback(sidebarItems);
    });
  }
}
