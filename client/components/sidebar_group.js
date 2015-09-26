import React    from 'react';
import { Link } from 'react-router';

export default class SidebarGroup extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      openChild: null
    };
  }

  onSidebarLinkClick(item) {
    this.setState({openChild: item.label});
  }

  componentWillReceiveProps(newProps) {
    if(!newProps.open) {
      this.setState({openChild: null});
    }
  }

  render() {
    return (
      <li>
        <Link
          activeClassName="active"
          to={this.props.url}
          onClick={this.props.onSidebarLinkClick}
        >
          {this.props.label}
        </Link>
        {this.renderChildNodes()}
      </li>
    )
  }

  renderChildNodes() {
    if(this.props.open && !!this.props.children) {
      let childElements = this.props.children.map((child) => {
        let isOpen = (child.label == this.state.openChild);
        return (
          <SidebarGroup
            {...child}
            open={child.label == this.state.openChild}
            onSidebarLinkClick={this.onSidebarLinkClick.bind(this, child)}
            key={child.url}
          />
        );
      });
      return (<ul>{childElements}</ul>);
    }
  }
}
